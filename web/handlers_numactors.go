package web


import (
    "log"
    "fmt"

    "github.com/nonspecific-username/mayham/dsl"
    apierrors "github.com/nonspecific-username/mayham/web/errors"
    "github.com/nonspecific-username/mayham/web/state"

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

    ct, err := checkContentType(c)
    if err != nil {
        return
    }

    key := c.Param("mod")
    err = checkModPath(c, ct, key)
    if err != nil {
        return
    }

    mod := dsl.NewNumActorsMod()
    err = bindFunc[ct](c, mod)
    if err != nil {
        return
    }

    validationErrors := mod.Validate()
    if validationErrors != nil && len(*validationErrors) > 0 {
        msg := "Failed to validate the input data"
        for _, valErr := range(*validationErrors) {
            msg = fmt.Sprintf("%s%v\n", msg, valErr)
        }
        respFunc[ct](c, 400, apierrors.InvalidValue("numactors", msg))
    } else {
        (*runtimeCfg)[key].NumActors = append((*runtimeCfg)[key].NumActors, *mod)
        state.Sync()
        idxStr := fmt.Sprintf("%d", len((*runtimeCfg)[key].NumActors) - 1)
        respFunc[ct](c, 200, &modCreatedResponse{Id: idxStr})
    }
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

    (*runtimeCfg)[key].NumActors = append(
        (*runtimeCfg)[key].NumActors[:idx],
        (*runtimeCfg)[key].NumActors[idx+1:]...
    )
    state.Sync()
    c.Data(204, gin.MIMEHTML, nil)
}
