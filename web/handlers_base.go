package web


import (
    "log"
    "fmt"

    "github.com/nonspecific-username/mayham/web/state"

    "github.com/gin-gonic/gin"
)


const (
    templateMod404 string = "Error: mod with id: %s doesn't exist"
)


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
}


func handleGetMod(c *gin.Context) {
    log.Printf("handleGetMod")

    key := c.Param("mod")

    if mod, ok := (*runtimeCfg)[key]; ok {
        c.JSON(200, mod)
    } else {
        c.Data(404, gin.MIMEHTML, []byte(fmt.Sprintf(templateMod404, key)))
    }
}


func handleUpdateMod(c *gin.Context) {
    log.Printf("handleUpdateMod")
}


func handleDeleteMod(c *gin.Context) {
    log.Printf("handleDeleteMod")

    key := c.Param("mod")

    if _, ok := (*runtimeCfg)[key]; ok {
        delete(*runtimeCfg, key)
        state.Sync()
        c.Data(204, gin.MIMEHTML, nil)
    } else {
        c.Data(404, gin.MIMEHTML, []byte(fmt.Sprintf(templateMod404, key)))
    }

}
