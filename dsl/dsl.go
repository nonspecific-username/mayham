package dsl


import (
    yaml "gopkg.in/yaml.v2"
)


type DSLConfig struct {
    SpawnNum []SpawnNumMod `yaml:"SpawnNum"`
}


func Load(input string) (*DSLConfig, error) {
    cfg := &DSLConfig{}
    err := yaml.UnmarshalStrict([]byte(input), cfg)
    if err != nil {
        return nil, err
    }
    return cfg, nil
}


func (cfg DSLConfig) String() string {
    output, _ := yaml.Marshal(cfg)
    return string(output)
}

