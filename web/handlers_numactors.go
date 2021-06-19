package web


import (
    "errors"
    "fmt"
    "log"
    "strconv"

    "github.com/gin-gonic/gin"
)


const (
    templateNumActors404 string = "Error: mod with id: %s doesn't contain NumActors with index: %d"
)


func checkNumActorsPath(c *gin.Context, key string, idx string) (int, error) {
    err := checkModPath(c, key)
    if err != nil {
        return -1, err
    }

    intIdx, err := strconv.Atoi(idx)

    if intIdx < 0 {
        c.Data(404, gin.MIMEHTML, []byte(fmt.Sprintf(templateBadIdx, intIdx)))
        return -1, errors.New("")
    }

    if intIdx > len((*runtimeCfg)[key].NumActors) - 1 {
        c.Data(404, gin.MIMEHTML, []byte(fmt.Sprintf(templateNumActors404, key, intIdx)))
        return -1, errors.New("")
    }

    return intIdx, nil
}


func handleGetNumActorsModList(c *gin.Context) {
    log.Printf("handleGetNumActorsModList")

    key := c.Param("mod")
    err := checkModPath(c, key)
    if err != nil {
        return
    }

    indices := make([]int, 0, len((*runtimeCfg)[key].NumActors))
    for idx, _ := range((*runtimeCfg)[key].NumActors) {
        indices = append(indices, idx)
    }

    c.JSON(200, &indices)
}


func handleCreateNumActorsMod(c *gin.Context) {
    log.Printf("handleCreateNumActorsMod")
}


func handleGetNumActorsMod(c *gin.Context) {
    log.Printf("handleGetNumActorsMod")

    key := c.Param("mod")
    idxStr := c.Param("idx")
    idx, err := checkNumActorsPath(c, key, idxStr)
    if err != nil {
        return
    }

    c.JSON(200, (*runtimeCfg)[key].NumActors[idx])
}


func handleUpdateNumActorsMod(c *gin.Context) {
    log.Printf("handleUpdateNumActorsMod")
}


func handleDeleteNumActorsMod(c *gin.Context) {
    log.Printf("handleDeleteNumActorsMod")
}
