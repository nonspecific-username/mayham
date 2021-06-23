package web


import (
    "log"
    "fmt"

    "github.com/nonspecific-username/mayham/dsl"
    "github.com/nonspecific-username/mayham/web/state"
    apierrors "github.com/nonspecific-username/mayham/web/errors"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)


type modCreatedResponse struct {
    Id string `yaml:"id" json:"id"`
}


type updateModRequest struct {
    Name string `yaml:"name" json:"name"`
    Description string `yaml: "description" json:"description"`
    Author string `yaml:"author" json:"author"`
    Enabled bool `yaml:"enabled" json:"enabled"`
}


func handleGetModList(c *gin.Context) {
    log.Printf("handleGetModList")

    ct, err := checkContentType(c)
    if err != nil {
        return
    }

    mods := make([]string, 0, len(*runtimeCfg))
    for mod := range *runtimeCfg {
        mods = append(mods, mod)
    }

    respFunc[ct](c, 200, &mods)
}


func handleCreateMod(c *gin.Context) {
    log.Printf("handleCreateMod")
    var err error

    ct, err := checkContentType(c)
    if err != nil {
        return
    }

    cfg := dsl.NewModConfig()

    err = bindFunc[ct](c, cfg)
    if err != nil {
        return
    }

    validationErrors := cfg.Validate()
    if validationErrors != nil && len(*validationErrors) > 0 {
        msg := "Failed to validate the input data"
        for _, valErr := range(*validationErrors) {
            msg = fmt.Sprintf("%s%v\n", msg, valErr)
        }
        respFunc[ct](c, 400, apierrors.InvalidValue("mod", msg))
    } else {
        key := uuid.New().String()
        (*runtimeCfg)[key] = cfg
        state.Sync()
        respFunc[ct](c, 200, &modCreatedResponse{Id: key})
    }
}


func handleGetMod(c *gin.Context) {
    log.Printf("handleGetMod")

    ct, err := checkContentType(c)
    if err != nil {
        return
    }

    key := c.Param("mod")
    err = checkModPath(c, key)
    if err != nil {
        return
    }

    mod := (*runtimeCfg)[key]
    respFunc[ct](c, 200, mod)
}


func handleBulkUpdateMod(c *gin.Context) {
    log.Printf("handleBulkUpdateMod")

    _, err := checkContentType(c)
    if err != nil {
        return
    }

    key := c.Param("mod")
    err = checkModPath(c, key)
    if err != nil {
        return
    }

    mod := dsl.NewModConfig()

    err = updateObjectFields(c, (*runtimeCfg)[key], mod)
    if err != nil {
        return
    }

}


func handleUpdateMod(c *gin.Context) {
    log.Printf("handleUpdateMod")

    ct, err := checkContentType(c)
    if err != nil {
        return
    }

    key := c.Param("mod")
    tgt := c.Param("target")
    err = checkModPath(c, key)
    if err != nil {
        return
    }

    req := &updateFieldRequest{}
    err = bindFunc[ct](c, &req)
    if err != nil {
        return
    }

    err = updateObjectField(c, (*runtimeCfg)[key], tgt, req.Value)
    if err != nil {
        return
    }
}


func handleDeleteMod(c *gin.Context) {
    log.Printf("handleDeleteMod")

    _, err := checkContentType(c)
    if err != nil {
        return
    }

    key := c.Param("mod")
    err = checkModPath(c, key)
    if err != nil {
        return
    }

    delete(*runtimeCfg, key)
    state.Sync()
    c.Data(204, gin.MIMEHTML, nil)

}
