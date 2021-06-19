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

func LoadMultiYAML(input *[]byte) (MultiModConfig, error, *[]error) {
    var errorsOutput []error
    cfg := make(map[string]*dsl.ModConfig)
    err := yaml.UnmarshalStrict(*input, cfg)

    if err != nil {
        return nil, err, nil
    }

    for id, subCfg := range cfg {
        errs := subCfg.Validate()
        if len(*errs) > 0 {
            for _, e := range(*errs) {
                msg := fmt.Sprintf("ModConfig[%s]: %v", id, e)
                errorsOutput = append(errorsOutput, errors.New(msg))
            }

        }
    }
    if len(errorsOutput) > 0 {
        return nil, errors.New("Failed to validate MultiModConfig"), &errorsOutput
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
