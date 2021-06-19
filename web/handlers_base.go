package web


import (
    "errors"
    "log"
    "fmt"

    "github.com/nonspecific-username/mayham/dsl"
    "github.com/nonspecific-username/mayham/web/state"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)


const (
    templateBadIdx string = "Error: bad index: %d"
    templateMod404 string = "Error: mod with id: %s doesn't exist"
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


func checkModPath(c *gin.Context, key string) error {
    if _, ok := (*runtimeCfg)[key]; !ok {
        c.Data(404, gin.MIMEHTML, []byte(fmt.Sprintf(templateMod404, key)))
        return errors.New("")
    }
    return nil
}


func handleGetModList(c *gin.Context) {
    log.Printf("handleGetModList")

    mods := make([]string, 0, len(*runtimeCfg))
    for mod := range *runtimeCfg {
        mods = append(mods, mod)
    }

    c.JSON(200, &mods)
}


func handleCreateMod(c *gin.Context) {
    log.Printf("handleCreateMod")
    var err error

    cfg := dsl.NewModConfig()

    /* The one and only case when we're willing to accept a yaml
       via the API request: create a new mod. This is only useful
       for importing new mods. */
    contentType := c.ContentType()
    switch contentType {
    case "application/json":
        err = c.BindJSON(cfg)
    case "application/x-yaml":
        err = c.BindYAML(cfg)
    default:
        msg := fmt.Sprintf("Unsupported content-type %s", contentType)
        c.Data(400, gin.MIMEHTML, []byte(msg))
    }
    if err != nil {
        msg := fmt.Sprintf("%v", err)
        c.Data(400, gin.MIMEHTML, []byte(msg))
        return
    }

    validationErrors := cfg.Validate()
    if validationErrors != nil && len(*validationErrors) > 0 {
        msg := "Failed to validate the input data"
        for _, valErr := range(*validationErrors) {
            msg = fmt.Sprintf("%s%v\n", msg, valErr)
        }
        c.Data(400, gin.MIMEHTML, []byte(msg))
    } else {
        key := uuid.New().String()
        (*runtimeCfg)[key] = cfg
        state.Sync()
        resp := &modCreatedResponse{Id: key}
        c.JSON(200, resp)
    }
}


func handleGetMod(c *gin.Context) {
    log.Printf("handleGetMod")

    key := c.Param("mod")
    err := checkModPath(c, key)
    if err != nil {
        return
    }

    mod := (*runtimeCfg)[key]
    c.JSON(200, mod)
}


func handleUpdateMod(c *gin.Context) {
    log.Printf("handleUpdateMod")

    key := c.Param("mod")
    err := checkModPath(c, key)
    if err != nil {
        return
    }

    mod := (*runtimeCfg)[key]

    req := &updateModRequest{}
    err = c.BindJSON(req)
    if err != nil {
        msg := fmt.Sprintf("%v", err)
        c.Data(400, gin.MIMEHTML, []byte(msg))
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

    key := c.Param("mod")
    err := checkModPath(c, key)
    if err != nil {
        return
    }

    delete(*runtimeCfg, key)
    state.Sync()
    c.Data(204, gin.MIMEHTML, nil)

}
