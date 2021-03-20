import sys

import json
import pprint
import yaml


if len(sys.argv) < 2:
    print('Error: input path is not specified')
    sys.exit(1)


def get_matching(data, field, value):
    ret = filter(lambda item: item.get(field, '') == value, data)
    return list(ret)


def get_by_export_type(data, export_type):
    return get_matching(data, 'export_type', export_type)


def get_export_idx(data, export_idx):
    ret = get_matching(data, '_jwp_export_idx', export_idx)
    return ret[0] if len(ret) > 0 else None


def get_bvc(obj, param):
    # strip BaseValueConstant from AttributeInitializationData for a param
    aid = 'AttributeInitializationData'
    bvc = 'BaseValueConstant'
    return obj.get(param, {}).get(aid, {}).get(bvc, 1.0)


def get_spawnerstyle_info(obj):
    ss_type = obj.get('SpawnerStyle').get('_jwp_export_dst_type')
    ss_idx = obj.get('SpawnerStyle').get('export')
    return ss_type, ss_idx


def get_spawnerstyle_parser(type):
    func_map = {'SpawnerStyle_Den': parse_spawnerstyle_den,
                'SpawnerStyle_Encounter': parse_spawnerstyle_encounter,
                'SpawnerStyle_Single': parse_spawnerstyle_single}
    return func_map[type]


def parse_spawnerstyle_encounter(data, idx):
    ss = get_export_idx(data, idx)
    output = {'Waves': [],
              'Type': 'Encounter'}
    for wave in ss.get('Waves'):
       wtype, widx = get_spawnerstyle_info(wave)
       func = get_spawnerstyle_parser(wtype)
       wave_output = func(data, widx)
       output['Waves'].append(wave_output)
    return output


def parse_spawnerstyle_single(data, idx):
    __so = 'SpawnOptions'
    ss = get_export_idx(data, idx)
    spawnopts = ss.get(__so)
    so_str = '{path}.{export}'.format(path=spawnopts[1], export=spawnopts[0])

    # Special case: Spawners can be used to place dynamic objects like lootables
    # we need to filter these out
    if 'Enemies' not in so_str:
        return None

    output = {}
    output[__so] = so_str
    output['Type'] = 'Single'
    return output


def parse_spawnerstyle_den(data, idx):
    __nap = 'NumActorsParam'
    __maawp = 'MaxAliveActorsWhenPassive'
    __maawt = 'MaxAliveActorsWhenThreatened'
    __so = 'SpawnOptions'
    ss = get_export_idx(data, idx)
    spawnopts = ss.get(__so)
    so_str = '{path}.{export}'.format(path=spawnopts[1], export=spawnopts[0])

    # Special case: Spawners can be used to place dynamic objects like lootables
    # we need to filter these out
    if 'Enemies' not in so_str:
        return None

    output = {}
    output[__nap] = get_bvc(ss, __nap)
    output[__maawp] = get_bvc(ss, __maawp)
    output[__maawt] = get_bvc(ss, __maawt)
    output[__so] = so_str
    output['Type'] = 'Den'
    return output


def build_spawner_info(data, spawner):
    # Find our SpawnerComponent by it's export index
    sc_idx = spawner.get('SpawnerComponent').get('export')
    sc = get_export_idx(d, sc_idx)
    ss_type, ss_idx = get_spawnerstyle_info(sc)
    ss_func = get_spawnerstyle_parser(ss_type)
    ss = ss_func(data, ss_idx)
    return ss


def parse_spawners(data, type):
    spawners = get_by_export_type(d, type)
    output = []
    for spawner in spawners:
        sinfo = build_spawner_info(data, spawner)
        if sinfo is not None:
            output.append(sinfo)
    return output



input_path = sys.argv[1]
with open(input_path, 'r') as f:
    d = json.load(f)

nonmission_spawners = 'OakSpawner'
mission_spawners = 'OakMissionSpawner'

nonmission_data = parse_spawners(d, nonmission_spawners)
mission_data = parse_spawners(d, mission_spawners)
chunks = input_path.split('/')
map_name = chunks[-2]
idx1 = chunks.index('Game')
idx2 = chunks.index('Maps')
package = '/'.join(chunks[idx1:idx2])

output = {
    package: {
        map_name: {
            'Mission': mission_data,
            'NonMission': nonmission_data
        }
    }
}

print(yaml.dump(output))
