package main


import (
    "fmt"

    "github.com/nonspecific-username/mayham/hotfix"
    "github.com/nonspecific-username/mayham/dsl"
    "github.com/nonspecific-username/mayham/gamedata"
)


var test_data string = `
SpawnNum:
- spawn:
    amap: ".*"
  mode: factor
  param1: 5
  max_actors_mode: match
`


func main() {
    var hf hotfix.Hotfix
    hf.AddRegular(hotfix.EarlyLevel,
                   0,
                   "MatchAll",
                   "ObjectPath",
                   "Attr",
                   0,
                   "",
                   hotfix.RenderBVCOverride(100))
    hf.AddRegular(hotfix.EarlyLevel,
                   0,
                   "MatchAll",
                   "ObjectPath",
                   "Attr",
                   0,
                   "",
                   hotfix.RenderBVCOverride(300))

    fmt.Printf(hf.Render() + "\n")
    cfg, e := dsl.Load(test_data)
    if e != nil {
        fmt.Println(e)
    } else {
        fmt.Println(*cfg)
        fmt.Println(gamedata.GetSpawners(&(cfg.SpawnNum[0].Spawn)))
    }

}
