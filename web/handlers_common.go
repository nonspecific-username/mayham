package web


import (
    "errors"
    "fmt"
    "strconv"
    "reflect"

    apierrors "github.com/nonspecific-username/mayham/web/errors"
    "github.com/nonspecific-username/mayham/web/state"

    "github.com/gin-gonic/gin"
)


type _respFunc func(*gin.Context, int, interface {})
type _bindFunc func(*gin.Context, interface {}) error


type updateFieldRequest struct {
    Value string `yaml:"value" json:"value"`
}


const (
    ctJSON string = "application/json"
    ctYAML string = "application/x-yaml"

)


var (
    respFunc = map[string]_respFunc{ctJSON: jsonResp,
                                    ctYAML: yamlResp}
    bindFunc = map[string]_bindFunc{ctJSON: jsonBind,
                                    ctYAML: yamlBind}
    tagType = map[string]string {ctJSON: "json",
                                 ctYAML: "yaml"}
)


func jsonResp(c *gin.Context, code int, obj interface{}) {
    c.JSON(code, obj)
}


func yamlResp(c *gin.Context, code int, obj interface{}) {
    c.YAML(code, obj)
}


func jsonBind(c *gin.Context, obj interface{}) error {
    err := c.BindJSON(obj)
    if err != nil {
        jsonResp(c, 400, apierrors.ParseError("unmarshal", fmt.Sprintf("%v", err)))
        return err
    }
    return nil
}


func yamlBind(c *gin.Context, obj interface{}) error {
    err := c.BindYAML(obj)
    if err != nil {
        yamlResp(c, 400, apierrors.ParseError("unmarshal", fmt.Sprintf("%v", err)))
        return err
    }
    return nil
}


func checkContentType(c *gin.Context) (string, error) {
    contentType := c.ContentType()
    switch contentType {
    case ctJSON, ctYAML:
        return contentType, nil
    case "":
        return ctJSON, nil
    default:
        c.JSON(400, apierrors.UnsupportedContentType(contentType))
        return "", errors.New("")
    }
}


func checkModPath(c *gin.Context, contentType string, key string) error {
    if _, ok := (*runtimeCfg)[key]; !ok {
        c.JSON(404, apierrors.NotFound("mod", key))
        return errors.New("")
    }
    return nil
}


func checkNumActorsPath(c *gin.Context, contentType string, key string, idx string) (int, error) {
    err := checkModPath(c, contentType, key)
    if err != nil {
        return -1, err
    }

    intIdx, err := strconv.Atoi(idx)
    if err != nil {
        respFunc[contentType](c, 400, apierrors.ParseError("index", idx))
        return -1, errors.New("")
    }

    if intIdx < 0 {
        respFunc[contentType](c, 400, apierrors.InvalidValue("index", idx))
        return -1, errors.New("")
    }

    if intIdx > len((*runtimeCfg)[key].NumActors) - 1 {
        respFunc[contentType](c, 400, apierrors.NotFound("index", idx))
        return -1, errors.New("")
    }

    return intIdx, nil
}


func updateObjectField(c *gin.Context, contentType string, obj interface {}, fieldTag string, newValue string) error {
    // obj MUST be a pointer, otherwise this will crash and burn
    vOrig := reflect.ValueOf(obj).Elem()
    t := vOrig.Type()

    // Create a copy of original
    updatedObj := reflect.New(t)
    updatedObj.Elem().Set(vOrig)

    fieldIdx := -1

    // Find the target field
    for i := 0; i < t.NumField(); i++ {
        f := t.Field(i)
        tag, ok := f.Tag.Lookup(tagType[contentType])
        allowUpdateTag, allowUpdateTagFound := f.Tag.Lookup("singleUpdate")
        // Ignore fields marked with update:no tag
        if ok && tag == fieldTag {
            if allowUpdateTagFound && allowUpdateTag == "no" {
                continue
            }
            fieldIdx = i
            break
        }
    }

    if fieldIdx == -1 {
        respFunc[contentType](c, 400, apierrors.NoSuchField(t.String(), fieldTag))
        return errors.New("")
    }

    // Type checks
    field := updatedObj.Elem().Field(fieldIdx)
    if field.IsValid() && field.CanSet() {
        switch field.Kind() {
        case reflect.Bool:
            boolVal, err := strconv.ParseBool(newValue)
            if err != nil {
                respFunc[contentType](c, 400, apierrors.ParseError("bool", newValue))
                return errors.New("")
            } else {
                field.SetBool(boolVal)
            }
        case reflect.Int:
            intVal, err := strconv.Atoi(newValue)
            if err != nil {
                respFunc[contentType](c, 400, apierrors.ParseError("integer", newValue))
                return errors.New("")
            } else {
                field.SetInt(int64(intVal))
            }
        case reflect.String:
            field.SetString(newValue)
        default:
            respFunc[contentType](c, 400, apierrors.ParseError(fieldTag, newValue))
            return errors.New("")
        }
    }

    // Validate updated object
    ret := updatedObj.MethodByName("Validate").Call([]reflect.Value{})
    var validationErrors *[]error = ret[0].Interface().(*[]error)
    if validationErrors != nil && len(*validationErrors) > 0 {
        msg := "Failed to validate the input data"
        for _, valErr := range(*validationErrors) {
            msg = fmt.Sprintf("%s%v\n", msg, valErr)
        }
        respFunc[contentType](c, 400, apierrors.InvalidValue(fieldTag, msg))
    } else {
        // only replace the original if the updated object passes validation
        vOrig.Set(updatedObj.Elem())
        state.Sync()
        c.Data(204, gin.MIMEHTML, nil)
    }

    return nil
}
