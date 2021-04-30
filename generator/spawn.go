package generator


import (
    //"log"
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

    for _, spawner := range(spawners) {
        if spawner.Type == "Single" {
            continue
        }
        var numActors, maxActors int

        // TODO: add a check for param2 if mode is "random"
        switch mod.Mode {
        case dsl.Factor:
            prevNumActors, _ := strconv.Atoi(spawner.NumActorsParam)
            numActors = prevNumActors * mod.Param1
        case dsl.Absolute:
            numActors = mod.Param1
        case dsl.Random:
            numActors = rand.Intn(mod.Param2 - mod.Param1 + 1) + mod.Param1
        }

        // TODO: make "match" the default option
        // TODO: check for maxactorsparam in mode is "factor"
        switch mod.MaxActorsMode {
        case dsl.MAScaled:
            prevNumActors, _ := strconv.Atoi(spawner.NumActorsParam)
            prevMaxActors, _ := strconv.Atoi(spawner.MaxAliveActorsWhenThreatened)
            maxActors = prevMaxActors * int(numActors / prevNumActors)
        case dsl.MAMatch:
            maxActors = numActors
        case dsl.MAFactor:
            prevMaxActors, _ := strconv.Atoi(spawner.MaxAliveActorsWhenThreatened)
            maxActors = prevMaxActors * mod.MaxActorsParam
        }

        generateSpawnNumHotfix(hf, spawner, numActors, maxActors)
    }

    return nil
}
