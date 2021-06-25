package dsl


import (
    "encoding/json"
    "errors"
    "fmt"
    "regexp"

    yaml "gopkg.in/yaml.v2"
)


type NumActorsMode = string

const (
    Factor NumActorsMode = "factor"
    Absolute NumActorsMode = "absolute"
    Random NumActorsMode = "random"
    RandomFactor NumActorsMode = "randomfactor"
)


type NumActorsMAMode = string

const (
    MAScaled NumActorsMAMode = "scaled"
    MAMatch NumActorsMAMode = "match"
    MAFactor NumActorsMAMode = "factor"
    MAAbsolute NumActorsMAMode = "absolute"
)


type SpawnSelector struct {
    Package string `yaml:"pkg" json:"pkg"`
    Map string `yaml:"map" json:"map"`
    Spawn string `yaml:"spawn" json:"spawn"`
}


type NumActorsMod struct {
    Enabled bool `yaml:"enabled" json:"enabled"`
    Spawn SpawnSelector `yaml:"spawn" json:"spawn"`
    Mode NumActorsMode `yaml:"mode" json:"mode"`
    Param1 int `yaml:"param1" json:"param1"`
    Param2 int `yaml:"param2" json:"param2"`
    MaxActorsMode NumActorsMAMode `yaml:"max_actors_mode" json:"max_actors_mode"`
    MaxActorsParam int `yaml:"max_actors_param" json:"max_actors_param"`
}


func NewSpawnSelector() *SpawnSelector {
    return &SpawnSelector{}
}


func (selector *SpawnSelector) FromYAML(input *[]byte) (*SpawnSelector, error, *ValidationError) {
    err := yaml.UnmarshalStrict(*input, selector)

    if err != nil {
        return nil, err, nil
    }

    validationError := selector.Validate()
    if validationError != nil {
        return nil, errors.New("Failed to validate SpawnSelector"), validationError
    }

    return selector, nil, nil
}


func (selector *SpawnSelector) FromJSON(input *[]byte) (*SpawnSelector, error, *ValidationError) {
    err := json.Unmarshal(*input, selector)

    if err != nil {
        return nil, err, nil
    }

    validationError := selector.Validate()
    if validationError != nil {
        return nil, errors.New("Failed to validate SpawnSelector"), validationError
    }

    return selector, nil, nil
}


func (selector SpawnSelector) YAML() string {
    str, _ := yaml.Marshal(selector)

    return string(str)
}


func (selector SpawnSelector) JSON() string {
    str, _ := json.Marshal(selector)

    return string(str)
}


func (spawn *SpawnSelector) Validate() *ValidationError {
    var errorsOutput []*ValidationError

    if _, err := regexp.Compile(spawn.Package); err != nil {
        desc := fmt.Sprintf("invalid regexp: %v", err)
        errorsOutput = append(errorsOutput, NewValidationError("pkg", desc, spawn.Package))
    }

    if _, err := regexp.Compile(spawn.Map); err != nil {
        desc := fmt.Sprintf("invalid regexp: %v", err)
        errorsOutput = append(errorsOutput, NewValidationError("map", desc, spawn.Map))
    }

    if _, err := regexp.Compile(spawn.Spawn); err != nil {
        desc := fmt.Sprintf("invalid regexp: %v", err)
        errorsOutput = append(errorsOutput, NewValidationError("spawn", desc, spawn.Spawn))
    }

    if len(errorsOutput) == 0 {
        return nil
    }

    validationError := NewValidationError("", "SpawnSelector validation failed", errorsOutput)
    return validationError
}


func NewNumActorsMod() *NumActorsMod {
    return &NumActorsMod{}
}


func (mod *NumActorsMod) FromYAML(input *[]byte) (*NumActorsMod, error, *ValidationError) {
    err := yaml.UnmarshalStrict(*input, mod)

    if err != nil {
        return nil, err, nil
    }

    validationError := mod.Validate()
    if validationError != nil {
        return nil, errors.New("Failed to validate NumActorsMod"), validationError
    }

    return mod, nil, nil
}


func (mod *NumActorsMod) FromJSON(input *[]byte) (*NumActorsMod, error, *ValidationError) {
    err := json.Unmarshal(*input, mod)

    if err != nil {
        return nil, err, nil
    }

    validationError := mod.Validate()
    if validationError != nil {
        return nil, errors.New("Failed to validate NumActorsMod"), validationError
    }

    return mod, nil, nil
}


func (mod NumActorsMod) YAML() string {
    str, _ := yaml.Marshal(mod)

    return string(str)
}


func (mod NumActorsMod) JSON() string {
    str, _ := json.Marshal(mod)

    return string(str)
}


func (mod *NumActorsMod) Validate() *ValidationError {
    var errorsOutput []*ValidationError


    switch mod.Mode {
    case Factor, Absolute, Random, RandomFactor:
        break
    default:
        errorsOutput = append(errorsOutput, NewValidationError("mode", "invalid value", mod.Mode))
    }

    switch mod.MaxActorsMode {
    case MAScaled, MAMatch, MAFactor, MAAbsolute, "":
        break
    default:
        errorsOutput = append(errorsOutput, NewValidationError("max_actors_mode", "invalid value", mod.MaxActorsMode))
    }

    emptySS := SpawnSelector{}
    if mod.Spawn != emptySS {
        if spawnSelectorError := mod.Spawn.Validate(); spawnSelectorError != nil {
            spawnSelectorError.Field = "spawn"
            errorsOutput = append(errorsOutput, spawnSelectorError)
        }
    }

    if mod.Mode == "" {
        errorsOutput = append(errorsOutput, NewValidationError("mode", "missing param", nil))
    }

    if mod.Param1 == 0 {
        errorsOutput = append(errorsOutput, NewValidationError("param1", "missing param", nil))
    }

    if (mod.Mode == Random || mod.Mode == RandomFactor) && mod.Param2 == 0 {
        errorsOutput = append(errorsOutput, NewValidationError("param2", "missing param", nil))
    }

    if (mod.Mode == Random || mod.Mode == RandomFactor) && mod.Param1 >= mod.Param2 {
        errorsOutput = append(errorsOutput, NewValidationError("param1", "value is lower than param2", mod.Param1))
        errorsOutput = append(errorsOutput, NewValidationError("param2", "value is higher than param1", mod.Param2))
    }

    if mod.MaxActorsMode == MAFactor && mod.MaxActorsParam == 0 {
        errorsOutput = append(errorsOutput, NewValidationError("max_actors_param", "missing param (max_actors_mode=factor)", nil))
    }

    if mod.MaxActorsMode == MAAbsolute && mod.MaxActorsParam == 0 {
        errorsOutput = append(errorsOutput, NewValidationError("max_actors_param", "missing param (max_actors_mode=absolute)", nil))
    }

    if len(errorsOutput) == 0 {
        return nil
    }

    validationError := NewValidationError("", "NumActorsMod validation failed", errorsOutput)
    return validationError
}
