import json, glob, argparse, os

parser = argparse.ArgumentParser()
parser.add_argument('workdir')
args = parser.parse_args()

workdir = args.workdir

if workdir is None:
    pwd = os.getcwd()
    _, workdir = os.path.split(pwd)

print(workdir)