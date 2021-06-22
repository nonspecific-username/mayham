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
    Id string `json:"id"`
}


type updateModRequest struct {
    Name string `json:"name"`
    Description string `json:"description"`
    Author string `json:"author"`
    Enabled bool `json:"enabled"`
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
    err = checkModPath(c, ct, key)
    if err != nil {
        return
    }

    mod := (*runtimeCfg)[key]
    respFunc[ct](c, 200, mod)
}


func handleUpdateMod(c *gin.Context) {
    log.Printf("handleUpdateMod")

    ct, err := checkContentType(c)
    if err != nil {
        return
    }

    key := c.Param("mod")
    err = checkModPath(c, ct, key)
    if err != nil {
        return
    }

    mod := (*runtimeCfg)[key]

    req := &updateModRequest{}
    err = bindFunc[ct](c, &req)
    if err != nil {
        return
    }

    if req.Name != "" && mod.Name != req.Name {
        mod.Name = req.Name
    }

    if req.Description != "" && mod.Description != req.Description{
        mod.Description = req.Description
    }

    if req.Author != "" && mod.Author != req.Author {
        mod.Author = req.Author
    }

    if mod.Enabled != req.Enabled {
        mod.Enabled = req.Enabled
    }

    state.Sync()
    c.Data(200, gin.MIMEHTML, nil)
}


func handleDeleteMod(c *gin.Context) {
    log.Printf("handleDeleteMod")

    ct, err := checkContentType(c)
    if err != nil {
        return
    }

    key := c.Param("mod")
    err = checkModPath(c, ct, key)
    if err != nil {
        return
    }

    delete(*runtimeCfg, key)
    state.Sync()
    c.Data(204, gin.MIMEHTML, nil)

}
