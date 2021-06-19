package dsl


import (
    "errors"
    "fmt"
    "regexp"
)


type NumActorsMode string

const (
    Factor NumActorsMode = "factor"
    Absolute NumActorsMode = "absolute"
    Random NumActorsMode = "random"
    RandomFactor NumActorsMode = "randomfactor"
)


type NumActorsMAMode string

const (
    MAScaled NumActorsMAMode = "scaled"
    MAMatch NumActorsMAMode = "match"
    MAFactor NumActorsMAMode = "factor"
    MAAbsolute NumActorsMAMode = "absolute"
)


type SpawnSelector struct {
    Package string `yaml:"pkg",omitempty json:"pkg"`
    Map string `yaml:"map" json:"map"`
    Spawn string `yaml:"spawn",omitempty json:"spawn"`
}


type NumActorsMod struct {
    Enabled bool `yaml:"enabled" json:"enabled"`
    Spawn *SpawnSelector `yaml:"spawn" json:"spawn"`
    Mode NumActorsMode `yaml:"mode" json:"mode"`
    Param1 int `yaml: "param1" json:"param1"`
    Param2 int `yaml: "param2",omitempty json:"param2"`
    MaxActorsMode NumActorsMAMode `yaml:"max_actors_mode",omitempty json:"max_actors_mode"`
    MaxActorsParam int `yaml:"max_actors_param",omitempty json:"max_actors_param"`
}


func (spawn *SpawnSelector) Validate() *[]error {
    var errorsOutput []error

    if spawn.Map == "" {
        msg := "\"map\" is required"
        errorsOutput = append(errorsOutput, errors.New(msg))
    }

    if _, err := regexp.Compile(spawn.Package); err != nil {
        msg := fmt.Sprintf("\"package\" is an invalid regexp: %v", err)
        errorsOutput = append(errorsOutput, errors.New(msg))
    }

    if _, err := regexp.Compile(spawn.Map); err != nil {
        msg := fmt.Sprintf("\"map\" is an invalid regexp: %v", err)
        errorsOutput = append(errorsOutput, errors.New(msg))
    }

    if _, err := regexp.Compile(spawn.Spawn); err != nil {
        msg := fmt.Sprintf("\"spawn\" is an invalid regexp: %v", err)
        errorsOutput = append(errorsOutput, errors.New(msg))
    }

    if len(errorsOutput) == 0 {
        return nil
    }

    return &errorsOutput
}


func (mod *NumActorsMod) Validate() *[]error {
    var errorsOutput []error


    switch mod.Mode {
    case Factor, Absolute, Random, RandomFactor:
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
    } else {
        if spawnSelectorErrors := mod.Spawn.Validate(); spawnSelectorErrors != nil {
            for _, e := range(*spawnSelectorErrors) {
                msg := fmt.Sprintf("\"spawn\": %v", e)
                errorsOutput = append(errorsOutput, errors.New(msg))
            }
        }
    }

    if mod.Mode == "" {
        msg  := "\"mode\" is required"
        errorsOutput = append(errorsOutput, errors.New(msg))
    }

    if mod.Param1 == 0 {
        msg  := "\"param1\" is required"
        errorsOutput = append(errorsOutput, errors.New(msg))
    }

    if (mod.Mode == Random || mod.Mode == RandomFactor) && mod.Param2 == 0 {
        msg  := "\"param2\" is required when \"mode\" is either \"random\" or \"randomfactor\""
        errorsOutput = append(errorsOutput, errors.New(msg))
    }

    if (mod.Mode == Random || mod.Mode == RandomFactor) && mod.Param1 >= mod.Param2 {
        msg  := "\"param1\" must be less than \"param2\" when \"mode\" is either \"random\" or \"randomfactor\""
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
