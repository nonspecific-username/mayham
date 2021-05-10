package generator


import (
    "math/rand"
    "strconv"
    "strings"
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


func generateSpawnUncapMod(hf *hotfix.Hotfix) {
    for _, line := range(strings.Split(_spawnOptionUncapModGenerated, "\n")) {
        hf.AddRaw(line)
    }
}


func generateSpawnNumMod(hf *hotfix.Hotfix, mod *dsl.SpawnNumMod) {
    rand.Seed(time.Now().UnixNano())
    spawners := gamedata.GetSpawners(mod.Spawn)

    for _, spawner := range(spawners) {
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
            numActors = rand.Intn(mod.Param2 - mod.Param1 + 1) + mod.Param1
        }

        switch mod.MaxActorsMode {
        case dsl.MAScaled:
            prevNumActors, _ := strconv.Atoi(spawner.NumActorsParam)
            prevMaxActors, _ := strconv.Atoi(spawner.MaxAliveActorsWhenThreatened)
            maxActors = prevMaxActors * int(numActors / prevNumActors)
        case dsl.MAFactor:
            prevMaxActors, _ := strconv.Atoi(spawner.MaxAliveActorsWhenThreatened)
            maxActors = prevMaxActors * mod.MaxActorsParam
        case dsl.MAAbsolute:
            maxActors = mod.MaxActorsParam
        case dsl.MAMatch, "":
            maxActors = numActors
        }

        generateSpawnNumHotfix(hf, spawner, numActors, maxActors)
    }

    return
}
