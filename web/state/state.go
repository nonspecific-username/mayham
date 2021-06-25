package state


import (
    "encoding/json"
    "errors"
    "fmt"

    yaml "gopkg.in/yaml.v2"

    "github.com/nonspecific-username/mayham/dsl"
)


type MultiModConfig map[string]*dsl.ModConfig


func NewMulti() MultiModConfig {
    cfg := make(map[string]*dsl.ModConfig)
    return cfg
}

func LoadMultiYAML(input *[]byte) (MultiModConfig, error, *dsl.ValidationError) {
    cfg := make(map[string]*dsl.ModConfig)
    err := yaml.UnmarshalStrict(*input, cfg)

    if err != nil {
        return nil, err, nil
    }

    for id, subCfg := range cfg {
        validationError := subCfg.Validate()
        if validationError != nil {
            return nil, errors.New(fmt.Sprintf("Failed to validate MultiModConfig[%s]", id)), validationError
        }
    }

    return cfg, nil, nil
}


func (cfg MultiModConfig) YAML() string {
    str, _ := yaml.Marshal(cfg)

    return string(str)
}


func (cfg MultiModConfig) JSON() string {
    str, _ := json.Marshal(cfg)

    return string(str)
}
