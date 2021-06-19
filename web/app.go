package web


import (
    "log"

    "github.com/nonspecific-username/mayham/web/state"

    "github.com/gin-gonic/gin"
)

var (
    runtimeCfg *state.MultiModConfig
)


const (
    cfgFilename string = "persistent.yml"
    logFilename string = "mayham.log"
)


func Init() error {
    var err error
    var validationErrors *[]error

    runtimeCfg, err, validationErrors = state.PersistentState(cfgFilename)
    if err != nil {
        log.Printf("Error parsing input file %s: %v:", cfgFilename, err)
        if validationErrors != nil && len(*validationErrors) > 0 {
            for _, valErr := range(*validationErrors) {
                log.Printf("\t%v", valErr)
            }
        }
        return err
    }
    log.Printf("Successfully read configuration with %d mods", len(*runtimeCfg))

    g := gin.New()
    //g.Use(gin.LoggerWithWriter(logFilename))
    g.Use(gin.Recovery())


    g.GET("/mod/", handleGetModList)
    g.POST("/mod/", handleCreateMod)

    g.GET("/mod/:mod", handleGetMod)
    g.PUT("/mod/:mod", handleUpdateMod)
    g.DELETE("/mod/:mod", handleDeleteMod)

    g.GET("/mod/:mod/spawnnum/", handleGetNumActorsModList)
    g.POST("/mod/:mod/spawnnum/", handleCreateNumActorsMod)

    g.GET("/mod/:mod/spawnnum/:idx/", handleGetNumActorsMod)
    g.PUT("/mod/:mod/spawnnum/:idx/", handleUpdateNumActorsMod)
    g.DELETE("/mod/:mod/spawnnum/:idx/", handleDeleteNumActorsMod)

    return g.Run("localhost:8300")
}


func Close() {
    state.ClosePersistentState()
}
