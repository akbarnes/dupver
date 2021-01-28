import json, argparse, os, os.path, toml
from glob import glob
from datetime import datetime 


pwd = os.getcwd()
_, workdir = os.path.split(pwd)

print(f"Workdir: {pwd} -> {workdir}")

cfg_path = ".dupver/config.toml"

with open(cfg_path) as f:
    cfg = toml.load(f)

## Update repos to v0.8 format
# 1. Add prev pointers to snapshot files
# 2. Add branch pointer
# 3. Set branch to main

# glob the snapshots
snapshot_path = os.path.join(cfg['RepoPath'], 'snapshots', cfg['WorkDirName'])
print(f'Snapshot path: {snapshot_path}\n')
snapshots = glob(os.path.join(snapshot_path, '*-*.json'))
snapshots = sorted(snapshots)

if len(snapshots) == 0:
    raise ValueError("No snapshots found")


print(f'Snapshot:\n{snapshots[0]}')

with open(snapshots[0]) as f:
    s = json.load(f)

cid = s['ID']
s['Branch'] = 'main'

if 'ParentIDs' in s and s['ParentIDs'] is not None:
    print(f"Non-empty parent id {s['ParentIDs']} in first file")
else:
    print('No parent ids')
    s['ParentIDs'] = []
    
new_snapshot_path = os.path.join(cfg['RepoPath'], 'snapshots', cfg['WorkDirName'], f"{cid}.json")
print(f"Writing snapshot 1 to {new_snapshot_path}")

with open(new_snapshot_path, 'w') as f:
    json.dump(s, f, indent=2)

pid = cid

for i in range(1,len(snapshots)):
    sp = snapshots[i]

    print(f'\nSnapshot:\n{sp}')

    with open(sp) as f:
        s = json.load(f)

    s['Branch'] = 'main'
    cid = s['ID']

    if 'ParentIDs' in s and s['ParentIDs'] is not None:
        if isinstance(s['ParentIDs'], list) and len(s['ParentIDs']) == 2 and len(s['ParentIDs'][1]) == 0:
            print(f"Removing empty id from {s['ParentIDs']}")
            s['ParentIDs'] = [s['ParentIDs'][0]]
        else:            
            print(f"Non-empty parent id {s['ParentIDs']}, not changing")
    else:
        print('No parent ids')
        s['ParentIDs'] = [pid]

    # save as cid
    new_snapshot_path = os.path.join(cfg['RepoPath'], 'snapshots', cfg['WorkDirName'], f"{cid}.json")
    print(f"Writing snapshot {i+1} to {new_snapshot_path}")

    with open(new_snapshot_path, 'w') as f:
        json.dump(s, f, indent=2)

    pid = cid

branch_folder = os.path.join(cfg['RepoPath'], 'branches', cfg['WorkDirName'])
branch_path = os.path.join(branch_folder, "main.toml")

if not os.path.exists(branch_folder):
    os.makedirs(branch_folder)    

print(f"Writing branch to {branch_path}")
with open(branch_path, 'w') as f:
    toml.dump(b, f)



