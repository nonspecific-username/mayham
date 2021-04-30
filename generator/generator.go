package generator


import (
    "log"

    "github.com/nonspecific-username/mayham/dsl"
    "github.com/nonspecific-username/mayham/hotfix"
)


func Generate(cfg *dsl.DSLConfig) (string, error) {
    var hf hotfix.Hotfix

    for _, spawnNumMod := range(cfg.SpawnNum) {
        err := generateSpawnNumMod(&hf, &spawnNumMod)
        if err != nil {
            log.Fatal(err)
            return "", err
        }
    }

    return hf.Render(), nil
}
