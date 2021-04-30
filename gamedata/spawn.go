package gamedata


import (
    //"log"
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


type SpawnerList struct {
    Spawners map[string]map[string]map[string]*Spawner `yaml:"Spawners"`
}


var _spawnerList SpawnerList
var unmarshalled bool = false


func (list *SpawnerList) Lookup(selector *dsl.SpawnSelector) []*Spawner {
    var output []*Spawner
    for pkg, mapList := range (*list).Spawners {
        matched, _ := regexp.MatchString(pkg, selector.Package)
        if len(selector.Package) == 0 || matched {
            for mapName, spawners := range mapList {
                matched, _ = regexp.MatchString(mapName, selector.Map)
                if matched {
                    for spawnerName, spawner := range spawners {
                        matched, _ = regexp.MatchString(spawnerName, selector.Spawn)
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


func GetSpawnData() *SpawnerList {
    if unmarshalled {
        return &_spawnerList
    }

    _spawnerList.Spawners = make(map[string]map[string]map[string]*Spawner)
    yaml.UnmarshalStrict([]byte(_spawnDataGenerated), _spawnerList)
    //log.Fatal(err)
    unmarshalled = true

    return &_spawnerList
}
