#!/usr/bin/python3

import argparse
import sys


import deepmerge
import json
import pprint
import yaml


parser = argparse.ArgumentParser(description='Generate spawn data')
top_subparsers = parser.add_subparsers(dest='target')
top_subparsers.required=True

spawners_p = top_subparsers.add_parser('spawners', help='Generate spawners data')
spawnopts_p = top_subparsers.add_parser('spawnoptions', help='Generate spawnoptions data')

spwn_subparsers = spawners_p.add_subparsers(dest='spawners_target')

sp_single_p = spwn_subparsers.add_parser('single', help='Parse a single file')
sp_list_p = spwn_subparsers.add_parser('list', help='Parse a file list')
sp_single_p.add_argument('base_path', metavar='base-path', type=str,
                    help='base game path')
sp_single_p.add_argument('json_path', metavar='json-path', type=str,
                    help='path of a parsed json file')
sp_list_p.add_argument('list', metavar='list', type=str,
                    help='path of a file list to parse')

sop_subparsers = spawnopts_p.add_subparsers(dest='spawnopts_target')
sop_uncap_p = sop_subparsers.add_parser('uncap', help='Uncap AliveLimitParam')
sop_uncap_p.add_argument('sop_uncap_list', metavar='list', type=str,
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
        spawn_anchor_data = self.parse_spawners('Spawner_SpawnAnchor_C')
        all_spawners = deepmerge.always_merger.merge(mission_data,
                                                     nonmission_data)
        all_spawners = deepmerge.always_merger.merge(all_spawners,
                                                     spawn_anchor_data)
        output = {
            self.package: {
                self.map_name: all_spawners
            }
        }
        return output

    def build_spawner_info(self, spawner):
        spawner_id = spawner['_jwp_object_name']
        path = [self.spawners_base_path, spawner_id, 'SpawnerComponent']
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

        def parse_simple(obj):
            __nap = 'NumActorsParam'
            __maawp = 'MaxAliveActorsWhenPassive'
            __maawt = 'MaxAliveActorsWhenThreatened'
            __so = 'SpawnOptions'
            output = {}
            spawnopts = obj.get(__so)
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
            output[__nap] = get_bvc(obj, __nap)
            output[__maawp] = get_bvc(obj, __maawp)
            output[__maawt] = get_bvc(obj, __maawt)
            output[__so] = so_str
            return output
        
        def parse_spawnerstyle_den(idx):
            ss = get_export_idx(self.data, idx)
            ss_name= ss['_jwp_object_name']
            __path = list(path)
            __path.append(ss_name)
            __attr = list(attr)
            __attr.append('SpawnerStyle.Object')

            output = parse_simple(ss)
            if output is None:
                return None
            output['Type'] = 'Multiple'
            output['Path'] = '.'.join(__path)
            output['AttrBase'] = '.'.join(__attr)
            return {spawner_id: output}

        def parse_spawnerstyle_bunchlist(idx):
            ss = get_export_idx(self.data, idx)
            ss_name= ss['_jwp_object_name']
            output = {}
            path.append(ss_name)
            attr.append('')
            for i, bunch in enumerate(ss.get('Bunches')):
                attr[-1] = 'SpawnerStyle.Object.Bunches[{}]'.format(i)
                bunch_output = parse_simple(bunch)
                if bunch_output is None:
                    continue
                bunch_output['Type'] = 'Multiple'
                bunch_output['Path'] = '.'.join(path)
                bunch_output['AttrBase'] = '.'.join(attr)
                index_key = '{sid}__bunch_{index}'.format(
                    sid=spawner_id,
                    index=i
                )
                output[index_key] = bunch_output
            return output

        def parse_noop(idx):
            return None

        def get_spawnerstyle_parser(type):
            func_map = {'SpawnerStyle_Den': parse_spawnerstyle_den,
                        'SpawnerStyle_Bunch': parse_spawnerstyle_den,
                        'SpawnerStyle_BunchList': parse_spawnerstyle_bunchlist,
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


def render_bvc(num):
    tpl = '(ValueType=Int,DisabledValueModes=102,ValueFlags=0,ValueMode=AttributeInitializationData,Range=(Value=1.000000,Variance=0.000000),AttributeInitializer=None,AttributeData=None,AttributeInitializationData=(BaseValueConstant={bvc}.000000,DataTableValue=(DataTable=None,RowName="",ValueName=""),BaseValueAttribute=None,AttributeInitializer=None,BaseValueScale=1.000000),BlackboardKey=(KeyName="",bRuntimeKey=False),Condition=None,Actor=None)'
    return tpl.format(bvc=num)


def render_bvc_patch(path, attribute, num,
                     method='SparkEarlyLevelPatchEntry',
                     matcher='MatchAll'):
    output = [method,
              '(1,1,0,{})'.format(matcher),
              path,
              attribute,
              '0',
              '',
              render_bvc(int(num))]
    return ','.join(output)


def render_str_patch(path, attribute, tstring,
                     method='SparkEarlyLevelPatchEntry',
                     matcher='MatchAll'):
    output = [method,
              '(1,1,0,{})'.format(matcher),
              path,
              attribute,
              '0',
              '',
              tstring]
    return ','.join(output)


def uncap_spawnoptions_limit(basepath, filepath):
    output = []

    chunks = filepath.split('/')
    sop_base_path = filepath[len(basepath):].split('.')[0]
    sop_export_name = chunks[-1].split('.')[0]
    sop_full_path = '{}.{}'.format(sop_base_path,
                                   sop_export_name)

    with open(filepath, 'r') as f:
        sop = json.load(f)

    sopd = None
    for exp in sop:
        if 'export_type' not in exp or \
           exp['export_type'] != 'SpawnOptionData':
            continue
        else:
            sopd = exp

    if not sopd or 'Options' not in sopd:
        return None

    for option in sopd['Options']:
        attribute = 'Options.Options[{}].AliveLimitType'.format(
            option['_jwp_arr_idx']
        )
        if 'AliveLimitType' in option and \
           option['AliveLimitType'] != 'ESpawnLimitType::None':
            output.append(render_str_patch(sop_full_path,
                                           attribute,
                                           'None'))
    return output


def uncap_spawnoptions_limit_old(basepath, filepath):
    output = []

    chunks = filepath.split('/')
    sop_base_path = filepath[len(basepath):].split('.')[0]
    sop_export_name = chunks[-1].split('.')[0]
    sop_full_path = '{}.{}'.format(sop_base_path,
                                   sop_export_name)

    with open(filepath, 'r') as f:
        sop = json.load(f)

    sopd = None
    for exp in sop:
        if 'export_type' not in exp or \
           exp['export_type'] != 'SpawnOptionData':
            continue
        else:
            sopd = exp

    if not sopd or 'Options' not in sopd:
        return None

    for option in sopd['Options']:
        attribute = 'Options.Options[{}].AliveLimitParam'.format(
            option['_jwp_arr_idx']
        )
        if 'AliveLimitType' in option and \
           option['AliveLimitType'] != 'ESpawnLimitType::None':
            output.append(render_bvc_patch(sop_full_path,
                                           attribute,
                                           1000))
    return output


if __name__ == '__main__':
    output = {}
    output_txt = ''
    if args.target == 'spawners':
        if args.spawners_target == 'list':
            with open(args.list, 'r') as f:
                data = f.read().splitlines()
            basepath = data[0]
            for filepath in data[1:]:
                print("INFO: processing spawners in {}".format(filepath))
                msp = MapSpawnParser(filepath, basepath)
                output = deepmerge.always_merger.merge(output, msp.parse())
        elif args.spawners_target == 'single':
            msp = MapSpawnParser(args.json_path, args.base_path)
            output = deepmerge.always_merger.merge(output, msp.parse())
            output_txt = yaml.dump(output)
    elif args.target == 'spawnoptions':
        if args.spawnopts_target == 'uncap':
            with open(args.sop_uncap_list, 'r') as f:
                data = f.read().splitlines()
            out_lines = []
            basepath = data[0]
            for filepath in data[1:]:
                print("INFO: processing spawnoptions in {}".format(filepath))
                lines = uncap_spawnoptions_limit(basepath, filepath)
                if lines:
                    out_lines += lines
            output_txt = '\n'.join(out_lines)
    if args.output is not None:
        with open(args.output, 'w+') as f:
            f.write(output_txt)
    else:
        print(output_txt)
