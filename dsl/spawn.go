package dsl


type SpawnNumMode string

const (
    Factor SpawnNumMode = "factor"
    Absolute SpawnNumMode = "absolute"
    Random SpawnNumMode = "random"
)


type SpawnNumMAMode string

const (
    MAScaled SpawnNumMAMode = "scaled"
    MAMatch SpawnNumMAMode = "match"
    MAFactor SpawnNumMAMode = "factor"
)


type SpawnSelectorSpec struct {
    Package string `yaml:"pkg",omitempty`
    Map string `yaml:"map"`
    Spawn string `yaml:"spawn",omitempty`
}


type SpawnNumModSpec struct {
    Spawn SpawnSelectorSpec `yaml:"spawn"`
    Mode SpawnNumMode `yaml:"mode"`
    Param1 float32 `yaml: "param1"`
    Param2 float32 `yaml: "param2",omitempty`
    MaxActorsMode SpawnNumMAMode `yaml:"max_actors_mode",omitempty`
    MaxActorsParam float32 `yaml:"max_actors_param",omitempty`
}
