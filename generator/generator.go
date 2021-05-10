package generator


import (
    "github.com/nonspecific-username/mayham/dsl"
    "github.com/nonspecific-username/mayham/hotfix"
)


func Generate(cfg *dsl.DSLConfig) string {
    var hf hotfix.Hotfix

    if len(cfg.SpawnNum) > 0 {
        generateSpawnUncapMod(&hf)
    }

    for _, spawnNumMod := range(cfg.SpawnNum) {
        generateSpawnNumMod(&hf, &spawnNumMod)
    }

    return hf.Render()
}
