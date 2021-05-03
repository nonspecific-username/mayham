package dsl


type SpawnNumMode string

const (
    Factor SpawnNumMode = "factor"
    Absolute SpawnNumMode = "absolute"
    Random SpawnNumMode = "random"
)


type SpawnNumMAMode string

// TODO: add "absolute" mode
const (
    MAScaled SpawnNumMAMode = "scaled"
    MAMatch SpawnNumMAMode = "match"
    MAFactor SpawnNumMAMode = "factor"
    MAAbsolute SpawnNumMAMode = "absolute"
)


type SpawnSelector struct {
    Package string `yaml:"pkg",omitempty`
    Map string `yaml:"map"`
    Spawn string `yaml:"spawn",omitempty`
}


type SpawnNumMod struct {
    Spawn SpawnSelector `yaml:"spawn"`
    Mode SpawnNumMode `yaml:"mode"`
    Param1 int `yaml: "param1"`
    Param2 int `yaml: "param2",omitempty`
    MaxActorsMode SpawnNumMAMode `yaml:"max_actors_mode",omitempty`
    MaxActorsParam int `yaml:"max_actors_param",omitempty`
}
