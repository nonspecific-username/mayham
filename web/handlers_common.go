package web


import (
    "encoding/json"
    "bytes"
    "errors"
    "fmt"
    "io/ioutil"
    "strconv"
    "reflect"

    apierrors "github.com/nonspecific-username/mayham/web/errors"
    "github.com/nonspecific-username/mayham/web/state"

    "github.com/gin-gonic/gin"
    yaml "gopkg.in/yaml.v2"
)


type _respFunc func(*gin.Context, int, interface {})
type _bindFunc func(*gin.Context, interface {}) error
type _unmarshalFunc func([]byte, interface{}) error


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
    unmarshalFunc = map[string]_unmarshalFunc{ctJSON: json.Unmarshal,
                                              ctYAML: yaml.UnmarshalStrict}
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
    ct := c.ContentType()
    switch ct {
    case ctJSON, ctYAML:
        return ct, nil
    case "":
        return ctJSON, nil
    default:
        c.JSON(400, apierrors.UnsupportedContentType(ct))
        return "", errors.New("")
    }
}


func checkModPath(c *gin.Context, key string) error {
    ct := c.ContentType()
    if _, ok := (*runtimeCfg)[key]; !ok {
        respFunc[ct](c, 404, apierrors.NotFound("mod", key))
        return errors.New("")
    }
    return nil
}


func checkNumActorsPath(c *gin.Context, key string, idx string) (int, error) {
    ct := c.ContentType()
    err := checkModPath(c, key)
    if err != nil {
        return -1, err
    }

    intIdx, err := strconv.Atoi(idx)
    if err != nil {
        respFunc[ct](c, 400, apierrors.ParseError("index", idx))
        return -1, errors.New("")
    }

    if intIdx < 0 {
        respFunc[ct](c, 400, apierrors.InvalidValue("index", idx))
        return -1, errors.New("")
    }

    if intIdx > len((*runtimeCfg)[key].NumActors) - 1 {
        respFunc[ct](c, 400, apierrors.NotFound("index", idx))
        return -1, errors.New("")
    }

    return intIdx, nil
}


// Make sure that newValue is valid before passing it here
func updateSingleFieldByIndex(c *gin.Context, obj *reflect.Value, fieldIdx int, newValue *reflect.Value) error {
    field := obj.Field(fieldIdx)
    if field.IsValid() && field.CanSet() {
        field.Set(*newValue)
    }

    return nil
}


func updateObjectField(c *gin.Context, obj interface {}, fieldTag string, newValue string) error {
    ct := c.ContentType()

    orig := reflect.ValueOf(obj).Elem()
    t := orig.Type()
    updated := reflect.New(t)
    updatedDeref := updated.Elem()
    updatedDeref.Set(orig)

    fieldIdx := -1

    // Find the target field
    for i := 0; i < t.NumField(); i++ {
        f := t.Field(i)
        tag, ok := f.Tag.Lookup(tagType[ct])
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

    // Perform type checks
    var parsedValue interface{}
    field := orig.Field(fieldIdx)
    if field.IsValid() && field.CanSet() {
        switch field.Kind() {
        case reflect.Bool:
            boolVal, err := strconv.ParseBool(newValue)
            if err != nil {
                respFunc[ct](c, 400, apierrors.ParseError("bool", newValue))
                return errors.New("")
            } else {
                parsedValue = interface{}(boolVal)
            }
        case reflect.Int:
            intVal, err := strconv.Atoi(newValue)
            if err != nil {
                respFunc[ct](c, 400, apierrors.ParseError("integer", newValue))
                return errors.New("")
            } else {
                parsedValue = interface{}(intVal)
            }
        case reflect.String:
            parsedValue = interface{}(newValue)
        default:
            respFunc[ct](c, 400, apierrors.ParseError(fieldTag, newValue))
            return errors.New("")
        }
    }

    pv := reflect.ValueOf(parsedValue)
    err := updateSingleFieldByIndex(c, &updatedDeref, fieldIdx, &pv)
    if err != nil {
        return err
    }

    if fieldIdx == -1 {
        respFunc[ct](c, 400, apierrors.NoSuchField(t.String(), fieldTag))
        return errors.New("")
    }

    // Validate updated object
    ret := updated.MethodByName("Validate").Call([]reflect.Value{})
    var validationErrors *[]error = ret[0].Interface().(*[]error)
    if validationErrors != nil && len(*validationErrors) > 0 {
        msg := "Failed to validate the input data"
        for _, valErr := range(*validationErrors) {
            msg = fmt.Sprintf("%s%v\n", msg, valErr)
        }
        respFunc[ct](c, 400, apierrors.InvalidValue(fieldTag, msg))
    } else {
        // only replace the original if the updated object passes validation
        orig.Set(updated.Elem())
        state.Sync()
        c.Data(204, gin.MIMEHTML, nil)
    }

    return nil
}


func updateBulkFields(c *gin.Context, obj *reflect.Value, partial *reflect.Value, changesMap *map[string]interface{}) error {
    ct := c.ContentType()
    t := obj.Type()

    for i := 0; i < t.NumField(); i++ {
        f := t.Field(i)
        tag := f.Tag.Get(tagType[ct])
        if sub, ok := (*changesMap)[tag]; ok {
            switch partial.Field(i).Kind() {
            case reflect.Struct:
                objValueField := obj.Field(i)
                partialValueField := partial.Field(i)
                subMap := sub.(map[string]interface{})
                err := updateBulkFields(c, &objValueField, &partialValueField, &subMap)
                if err != nil {
                    return err
                }
            default:
                newValue := partial.Field(i)
                // since we're pulling the interface{} value directly from the field #i,
                // no need for additional type checks. It just works.
                err := updateSingleFieldByIndex(c, obj, i, &newValue)
                if err != nil {
                    return err
                }
            }
        }
    }

    return nil
}


func updateObjectFields(c *gin.Context, obj interface {}, changes interface{}) error{
    ct := c.ContentType()

    body, err := c.GetRawData()
    if err != nil {
        respFunc[ct](c, 400, apierrors.ParseError("unknown", fmt.Sprintf("%v", err)))
    }

    // a hack to allow a binding function to re-read c.Request.Body
    c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

    err = bindFunc[ct](c, changes)
    if err != nil {
        return err
    }

    changesMap := make(map[string]interface{})
    err = unmarshalFunc[ct](body, &changesMap)

    orig := reflect.ValueOf(obj).Elem()
    t := orig.Type()
    updated := reflect.New(t)
    updatedDeref := updated.Elem()
    updatedDeref.Set(orig)
    partial := reflect.ValueOf(changes).Elem()

    err = updateBulkFields(c, &updatedDeref, &partial, &changesMap)
    if err != nil {
        return err
    }

    // Validate updated object
    ret := updated.MethodByName("Validate").Call([]reflect.Value{})
    var validationErrors *[]error = ret[0].Interface().(*[]error)
    if validationErrors != nil && len(*validationErrors) > 0 {
        msg := "Failed to validate the input data"
        for _, valErr := range(*validationErrors) {
            msg = fmt.Sprintf("%s%v\n", msg, valErr)
        }
        respFunc[ct](c, 400, apierrors.InvalidValue("validation", msg))
        return errors.New("")
    } else {
        // only replace the original if the updated object passes validation
        orig.Set(updatedDeref)
        state.Sync()
        c.Data(204, gin.MIMEHTML, nil)
    }

    return nil
}
