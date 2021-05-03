package generator


import (
    //"log"
    "errors"
    "fmt"
    "math/rand"
    "strconv"
    "time"

    "github.com/nonspecific-username/mayham/dsl"
    "github.com/nonspecific-username/mayham/gamedata"
    "github.com/nonspecific-username/mayham/hotfix"
)



func generateSpawnNumHotfix(hf *hotfix.Hotfix, spawner *gamedata.Spawner, numActors int, maxActors int) {
    hf.AddRegular(hotfix.EarlyLevel,
                   0,
                   "MatchAll",
                   spawner.Path,
                   spawner.AttrBase + ".NumActorsParam",
                   0,
                   "",
                   hotfix.RenderBVCOverride(numActors))
    hf.AddRegular(hotfix.EarlyLevel,
                   0,
                   "MatchAll",
                   spawner.Path,
                   spawner.AttrBase + ".MaxAliveActorsWhenPassive",
                   0,
                   "",
                   hotfix.RenderBVCOverride(maxActors))
    hf.AddRegular(hotfix.EarlyLevel,
                   0,
                   "MatchAll",
                   spawner.Path,
                   spawner.AttrBase + ".MaxAliveActorsWhenThreatened",
                   0,
                   "",
                   hotfix.RenderBVCOverride(maxActors))
}


func generateSpawnNumMod(hf *hotfix.Hotfix, mod *dsl.SpawnNumMod) error {
    rand.Seed(time.Now().UnixNano())
    spawners := gamedata.GetSpawners(&(mod.Spawn))

    for i, spawner := range(spawners) {
        if spawner.Type == "Single" {
            continue
        }
        var numActors, maxActors int

        switch mod.Mode {
        case dsl.Factor:
            prevNumActors, _ := strconv.Atoi(spawner.NumActorsParam)
            numActors = prevNumActors * mod.Param1
        case dsl.Absolute:
            numActors = mod.Param1
        case dsl.Random:
            if mod.Param2 == 0 {
                msg := fmt.Sprintf("SpawnNum[%d]: \"param2\" is required when \"mode\" is \"random\"", i)
                return errors.New(msg)
            }
            numActors = rand.Intn(mod.Param2 - mod.Param1 + 1) + mod.Param1
        }

        switch mod.MaxActorsMode {
        case dsl.MAScaled:
            prevNumActors, _ := strconv.Atoi(spawner.NumActorsParam)
            prevMaxActors, _ := strconv.Atoi(spawner.MaxAliveActorsWhenThreatened)
            maxActors = prevMaxActors * int(numActors / prevNumActors)
        case dsl.MAFactor:
            if mod.MaxActorsParam == 0 {
                msg := fmt.Sprintf("SpawnNum[%d]: \"max_actors_param\" is required when \"max_actors_mode\" is \"factor\"", i)
                return errors.New(msg)
            }
            prevMaxActors, _ := strconv.Atoi(spawner.MaxAliveActorsWhenThreatened)
            maxActors = prevMaxActors * mod.MaxActorsParam
        case dsl.MAMatch, "":
            maxActors = numActors
        }

        generateSpawnNumHotfix(hf, spawner, numActors, maxActors)
    }

    return nil
}
