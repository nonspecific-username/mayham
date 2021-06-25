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
    runtimeCfg, err, validationError := state.PersistentState(cfgFilename)
    if err != nil {
        log.Printf("Error parsing input file %s: %v:", cfgFilename, err)
        if validationError != nil {
            log.Printf(validationError.String())
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
    g.PUT("/mod/:mod", handleBulkUpdateMod)
    g.PUT("/mod/:mod/:target", handleUpdateMod)
    g.DELETE("/mod/:mod", handleDeleteMod)

    g.GET("/mod/:mod/numactors/", handleGetNumActorsModList)
    g.POST("/mod/:mod/numactors/", handleCreateNumActorsMod)

    g.GET("/mod/:mod/numactors/:idx", handleGetNumActorsMod)
    g.PUT("/mod/:mod/numactors/:idx", handleBulkUpdateNumActorsMod)
    g.PUT("/mod/:mod/numactors/:idx/:target", handleUpdateNumActorsMod)
    g.DELETE("/mod/:mod/numactors/:idx", handleDeleteNumActorsMod)

    return g.Run("localhost:8300")
}


func Close() {
    state.ClosePersistentState()
}
