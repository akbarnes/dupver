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
snapshots = glob(os.path.join(snapshot_path, '*.json'))
snapshots = sorted(snapshots)

if len(snapshots) == 0:
    raise ValueError("No snapshots found")


for i in range(0,len(snapshots)):
    sp = snapshots[i]

    print(f'\nSnapshot:\n{sp}')

    with open(sp) as f:
        s = json.load(f)

    s['Branch'] = 'main'

    with open(sp, 'w') as f:
        json.dump(s, f, indent=2)




