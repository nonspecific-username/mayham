package gamedata


import (
    "log"
    "regexp"

    yaml "gopkg.in/yaml.v2"

    "github.com/nonspecific-username/mayham/dsl"
)


type Spawner struct {
    AttrBase string `yaml:"AttrBase"`
    MaxAliveActorsWhenPassive string `yaml:"MaxAliveActorsWhenPassive"`
    MaxAliveActorsWhenThreatened string `yaml:"MaxAliveActorsWhenThreatened"`
    NumActorsParam string `yaml:"NumActorsParam"`
    Path string `yaml:"Path"`
    SpawnOptions string `yaml:"SpawnOptions"`
    Type string `yaml:"Type"`
}


type SpawnerList map[string]map[string]map[string]*Spawner


var spawnerList SpawnerList
var unmarshalled bool = false


// TODO: add regex error parsing
func GetSpawners(selector *dsl.SpawnSelector) []*Spawner {
    if !unmarshalled {
        parseSpawnData()
    }
    var output []*Spawner
    for pkg, mapList := range spawnerList {
        matched, _ := regexp.MatchString(selector.Package, pkg)
        if len(selector.Package) == 0 || matched {
            for mapName, spawners := range mapList {
                matched, _ = regexp.MatchString(selector.Map, mapName)
                if matched {
                    for spawnerName, spawner := range spawners {
                        matched, _ = regexp.MatchString(selector.Spawn, spawnerName)
                        if len(selector.Spawn) == 0 || matched {
                            output = append(output, spawner)
                        }
                    }
                }
            }
        }
    }
    return output
}


func parseSpawnData() {
    spawnerList = make(map[string]map[string]map[string]*Spawner)
    err := yaml.UnmarshalStrict([]byte(_spawnDataGenerated), &spawnerList)
    if err != nil {
        log.Fatal(err)
    }
    unmarshalled = true
}
