import argparse
import sys


import deepmerge
import json
import pprint
import yaml


parser = argparse.ArgumentParser(description='Generate spawn data')
subparsers = parser.add_subparsers(dest='target')
subparsers.required=True
single_p = subparsers.add_parser('single', help='Parse a single file')
list_p = subparsers.add_parser('list', help='Parse a file list')
single_p.add_argument('base_path', metavar='base-path', type=str,
                    help='base game path')
single_p.add_argument('json_path', metavar='json-path', type=str,
                    help='path of a parsed json file')
list_p.add_argument('list', metavar='list', type=str,
                    help='path of a file list to parse')
parser.add_argument('--output', type=str, help='file to write output to')
args = parser.parse_args()


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


class MapSpawnParser(object):
    def __init__(self, filepath, basepath):
        with open(filepath, 'r') as f:
            self.data = json.load(f)

        chunks = filepath.split('/')
        self.map_name = chunks[-1].split('_')[0]
        self.map_base_path = filepath[len(basepath):].split('.')[0]
        self.map_export_name = chunks[-1].split('.')[0]
        idx = chunks.index('Maps')
        self.package = '/'.join(chunks[idx-1:idx])
        self.spawners_base_path = "{base}.{map}:PersistentLevel".format(
            base=self.map_base_path,
            map=self.map_export_name
        )

    def parse(self):
        nonmission_data = self.parse_spawners('OakSpawner')
        mission_data = self.parse_spawners('OakMissionSpawner')
        all_spawners = deepmerge.always_merger.merge(mission_data,
                                                     nonmission_data)
        output = {
            self.package: {
                self.map_name: all_spawners
            }
        }
        return output

    def build_spawner_info(self, spawner):
        spawner_id = spawner['_jwp_object_name']
        path = [self.spawners_base_path, spawner_id]
        attr = ['SpawnerComponent.Object.']

        def parse_spawnerstyle_encounter(idx):
            ss = get_export_idx(self.data, idx)
            ss_name= ss['_jwp_object_name']
            output = {}
            path.append('')
            attr.append('')
            for i, wave in enumerate(ss.get('Waves')):
               path[-1] = ss_name
               attr[-1] = 'SpawnerStyle.Object.Waves[{}]'.format(i)
               wtype, widx = get_spawnerstyle_info(wave)
               func = get_spawnerstyle_parser(wtype)
               wave_output = func(widx)
               if wave_output is None:
                   continue
               index_key = '{sid}__wave_{index}'.format(
                   sid=spawner_id,
                   index=i
               )
               output[index_key] = list(wave_output.values())[0]
            return output
        
        def parse_spawnerstyle_single(idx):
            __so = 'SpawnOptions'
            ss = get_export_idx(self.data, idx)
            ss_name= ss['_jwp_object_name']
            __path = list(path)
            __path.append(ss_name)
            __attr = list(attr)
            __attr.append('SpawnerStyle.Object')
            spawnopts = ss.get(__so)
            if spawnopts is None:
                msg = ("WARNING: spawner {} doesn't have SpawnOptions "
                       "defined, skipping...")
                print(msg.format(spawner_id))
                return None
            so_str = '{path}.{export}'.format(path=spawnopts[1], export=spawnopts[0])
        
            # Special case: Spawners can be used to place dynamic objects like lootables
            # we need to filter these out
            if 'Enemies' not in so_str:
                msg = "WARNING: spawner {} is not an enemy spawner, skipping..."
                print(msg.format(spawner_id))
                return None
        
            output = {}
            output[__so] = so_str
            output['Type'] = 'Single'
            output['Path'] = '.'.join(__path)
            output['AttrBase'] = '.'.join(__attr)
            return {spawner_id: output}
        
        def parse_spawnerstyle_den(idx):
            __nap = 'NumActorsParam'
            __maawp = 'MaxAliveActorsWhenPassive'
            __maawt = 'MaxAliveActorsWhenThreatened'
            __so = 'SpawnOptions'
            ss = get_export_idx(self.data, idx)
            ss_name= ss['_jwp_object_name']
            __path = list(path)
            __path.append(ss_name)
            __attr = list(attr)
            __attr.append('SpawnerStyle.Object')
            spawnopts = ss.get(__so)
            if spawnopts is None:
                msg = ("WARNING: spawner {} doesn't have SpawnOptions "
                       "defined, skipping...")
                print(msg.format(spawner_id))
                return None
            so_str = '{path}.{export}'.format(path=spawnopts[1], export=spawnopts[0])
        
            # Special case: Spawners can be used to place dynamic objects like lootables
            # we need to filter these out
            if 'Enemies' not in so_str:
                msg = "WARNING: spawner {} is not an enemy spawner, skipping..."
                print(msg.format(spawner_id))
                return None

            output = {}
            output[__nap] = get_bvc(ss, __nap)
            output[__maawp] = get_bvc(ss, __maawp)
            output[__maawt] = get_bvc(ss, __maawt)
            output[__so] = so_str
            output['Type'] = 'Multiple'
            output['Path'] = '.'.join(__path)
            output['AttrBase'] = '.'.join(__attr)
            return {spawner_id: output}

        def parse_noop(idx):
            return None

        def get_spawnerstyle_parser(type):
            func_map = {'SpawnerStyle_Den': parse_spawnerstyle_den,
                        'SpawnerStyle_Bunch': parse_spawnerstyle_den,
                        'SpawnerStyle_Encounter': parse_spawnerstyle_encounter,
                        'SpawnerStyle_Single': parse_spawnerstyle_single,
                        'OakSpawnerStyle_PlayerInstanced': parse_noop}
            return func_map[type]

        # Find our SpawnerComponent by it's export index
        sc_idx = spawner.get('SpawnerComponent').get('export')
        sc = get_export_idx(self.data, sc_idx)
        try:
            ss_type, ss_idx = get_spawnerstyle_info(sc)
        except AttributeError as e:
            msg = ("WARNING: skipping Spawner {} since it's "
                   "SpawnerComponent has no SpawnerStyle attached")
            print(msg.format(spawner.get('_jwp_object_name')))
            return None
        ss_func = get_spawnerstyle_parser(ss_type)
        ss = ss_func(ss_idx)
        return ss

    def parse_spawners(self, type):
        spawners = get_by_export_type(self.data, type)
        output = {}
        for spawner in spawners:
            sinfo = self.build_spawner_info(spawner)
            if sinfo is not None:
                output = deepmerge.always_merger.merge(output, sinfo)
        return output

if __name__ == '__main__':
    output = {}
    if args.target == 'list':
        with open(args.list, 'r') as f:
            data = f.read().splitlines()
        basepath = data[0]
        for filepath in data[1:]:
            print("INFO: processing {}".format(filepath))
            msp = MapSpawnParser(filepath, basepath)
            output = deepmerge.always_merger.merge(output, msp.parse())
    elif args.target == 'single':
        msp = MapSpawnParser(args.json_path, args.base_path)
        output = deepmerge.always_merger.merge(output, msp.parse())
    output_txt = yaml.dump(output)
    if args.output is not None:
        with open(args.output, 'w+') as f:
            f.write(output_txt)
    else:
        print(output_txt)
