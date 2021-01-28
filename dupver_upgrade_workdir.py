import json, argparse, os, os.path, toml
from glob import glob
from datetime import datetime 


pwd = os.getcwd()
_, workdir = os.path.split(pwd)

print(f"Workdir: {pwd} -> {workdir}")

cfg_path = ".dupver/config.toml"

with open(cfg_path) as f:
    cfg = toml.load(f)

time = datetime.now().strftime("%Y-%m-%dT%H-%M-%S")
cfg_backup_path = f".dupver/config-{time}.toml"
os.rename(cfg_path, cfg_backup_path)

cfg["Version"] = '0.8'
cfg["Branch"] = 'main'
cfg['DefaultRepo'] = 'main'
cfg['Repos'] = {'main': cfg.pop('RepoPath')}

with open(cfg_path, 'w') as f:
    toml.dump(cfg, f)

