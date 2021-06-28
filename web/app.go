package web


import (
    "fmt"
    "log"
    "net/http"
    "path"

    "github.com/nonspecific-username/mayham/web/state"
    "github.com/nonspecific-username/mayham/web/ui"

    "github.com/gin-gonic/gin"
)

var (
    runtimeCfg *state.MultiModConfig
)


const (
    cfgFilename string = "persistent.yml"
    logFilename string = "mayham.log"
)


func serveStaticFile(c *gin.Context) {
    c.FileFromFS(path.Join("react-app/build/", c.Request.URL.Path), http.FS(ui.Assets))
    log.Printf("%s ---> %s", path.Join("react-app/build/", c.Request.URL.Path), c.Request.URL.Path)
}


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

    contents, err := ui.Assets.ReadDir("react-app/build")
    if err != nil {
        log.Fatal("Could not load embedded UI: %v", err)
    }
    for _, item := range(contents) {
        g.GET(fmt.Sprintf("/%s", item.Name()), serveStaticFile)
    }
    g.GET("/static/*filepath", serveStaticFile)
    g.GET("/", func(c *gin.Context) {c.FileFromFS("react-app/build/index.htm", http.FS(ui.Assets))})
    g.GET("/index.html", func(c *gin.Context) {c.FileFromFS("react-app/build/index.htm", http.FS(ui.Assets))})

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
