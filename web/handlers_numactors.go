package web


import (
    "log"

    //apierrors "github.com/nonspecific-username/mayham/web/errors"

    "github.com/gin-gonic/gin"
)


func handleGetNumActorsModList(c *gin.Context) {
    log.Printf("handleGetNumActorsModList")

    ct, err := checkContentType(c)
    if err != nil {
        return
    }

    key := c.Param("mod")
    err = checkModPath(c, ct, key)
    if err != nil {
        return
    }

    indices := make([]int, 0, len((*runtimeCfg)[key].NumActors))
    for idx, _ := range((*runtimeCfg)[key].NumActors) {
        indices = append(indices, idx)
    }

    respFunc[ct](c, 200, &indices)
}


func handleCreateNumActorsMod(c *gin.Context) {
    log.Printf("handleCreateNumActorsMod")
}


func handleGetNumActorsMod(c *gin.Context) {
    log.Printf("handleGetNumActorsMod")

    ct, err := checkContentType(c)
    if err != nil {
        return
    }

    key := c.Param("mod")
    idxStr := c.Param("idx")
    idx, err := checkNumActorsPath(c, ct, key, idxStr)
    if err != nil {
        return
    }

    respFunc[ct](c, 200, (*runtimeCfg)[key].NumActors[idx])
}


func handleUpdateNumActorsMod(c *gin.Context) {
    log.Printf("handleUpdateNumActorsMod")
}


func handleDeleteNumActorsMod(c *gin.Context) {
    log.Printf("handleDeleteNumActorsMod")
}
