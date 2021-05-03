package dsl


import (
    "errors"
    "fmt"
)


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
    Spawn *SpawnSelector `yaml:"spawn"`
    Mode SpawnNumMode `yaml:"mode"`
    Param1 int `yaml: "param1"`
    Param2 int `yaml: "param2",omitempty`
    MaxActorsMode SpawnNumMAMode `yaml:"max_actors_mode",omitempty`
    MaxActorsParam int `yaml:"max_actors_param",omitempty`
}


func (mod *SpawnNumMod) Validate() *[]error {
    var errorsOutput []error

    switch mod.Mode {
    case Factor, Absolute, Random:
        break
    default:
        msg := fmt.Sprintf("\"mode\" \"%s\" is invalid", mod.Mode)
        errorsOutput = append(errorsOutput, errors.New(msg))
    }

    switch mod.MaxActorsMode {
    case MAScaled, MAMatch, MAFactor, MAAbsolute, "":
        break
    default:
        msg := fmt.Sprintf("\"max_actors_mode\" \"%s\" is invalid", mod.MaxActorsMode)
        errorsOutput = append(errorsOutput, errors.New(msg))
    }

    if mod.Spawn == nil {
        msg := "\"spawn\" is required"
        errorsOutput = append(errorsOutput, errors.New(msg))
    }

    if mod.Mode == "" {
        msg  := "\"mode\" is required"
        errorsOutput = append(errorsOutput, errors.New(msg))
    }

    if mod.Param1 == 0 {
        msg  := "\"param1\" is required"
        errorsOutput = append(errorsOutput, errors.New(msg))
    }

    if mod.Mode == Random && mod.Param2 == 0 {
        msg  := "\"param2\" is required when \"mode\" is \"random\""
        errorsOutput = append(errorsOutput, errors.New(msg))
    }

    if mod.MaxActorsMode == MAFactor && mod.MaxActorsParam == 0 {
        msg  := "\"max_actors_param\" is required when \"max_actors_mode\" is \"factor\""
        errorsOutput = append(errorsOutput, errors.New(msg))
    }

    if mod.MaxActorsMode == MAAbsolute && mod.MaxActorsParam == 0 {
        msg  := "\"max_actors_param\" is required when \"max_actors_mode\" is \"absolute\""
        errorsOutput = append(errorsOutput, errors.New(msg))
    }

    if len(errorsOutput) == 0 {
        return nil
    }

    return &errorsOutput
}
