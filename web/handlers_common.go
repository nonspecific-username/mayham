package web


import (
    "errors"
    "fmt"
    "strconv"

    apierrors "github.com/nonspecific-username/mayham/web/errors"

    "github.com/gin-gonic/gin"
)


type _respFunc func(*gin.Context, int, interface {})
type _bindFunc func(*gin.Context, interface {}) error


const (
    ctJSON string = "application/json"
    ctYAML string = "application/x-yaml"

)


var (
    respFunc = map[string]_respFunc{ctJSON: jsonResp,
                                    ctYAML: yamlResp}
    bindFunc = map[string]_bindFunc{ctJSON: jsonBind,
                                    ctYAML: yamlBind}
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
        jsonResp(c, 400, apierrors.ParseError("unmarshalError", fmt.Sprintf("%v", err)))
        return err
    }
    return nil
}


func yamlBind(c *gin.Context, obj interface{}) error {
    err := c.BindYAML(obj)
    if err != nil {
        yamlResp(c, 400, apierrors.ParseError("unmarshalError", fmt.Sprintf("%v", err)))
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
