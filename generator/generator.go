package generator


import (
    "github.com/nonspecific-username/mayham/dsl"
    "github.com/nonspecific-username/mayham/hotfix"
)


func Generate(cfg *dsl.ModConfig) string {
    if !cfg.Enabled {
        return ""
    }

    var hf hotfix.Hotfix

    if len(cfg.NumActors) > 0 {
        generateSpawnUncapMod(&hf)
    }

    for _, spawnNumMod := range(cfg.NumActors) {
        if spawnNumMod.Enabled {
            generateNumActorsMod(&hf, &spawnNumMod)
        }
    }

    return hf.Render()
}
