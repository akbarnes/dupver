## Global Preferences Structure
Global preferences for Dupver are stored in a `.dupver/` subdirectory
in the user's home directory. This holds the preferences file
`prefs.toml` which contains the paths forthe selected diff tool,
editor, and the default repository path to use when creating 
project working directories.

```
DiffTool = "bcompare.exe"
Editor = "notepad++"
DefaultRepo = "/Users/KumarSW/.dupver_repo"
```

## Project Working Directory Structure
Project-specific settings are stored in a `.dupver/` subdirectory
in the project working directory. This contains the project settings
file `config.toml`, which contains the project name, branch, 
associated repository, and the default repository to use:

```
WorkDirName = "electronics"
Branch = "main"
DefaultRepo = "store"

[Repos]
  main = "/Users/KumarSW/.dupver_repo"
  store = "/Users/KumarSW/Store"
```


## Repository Structure
The repository version 2 format consists of a folder containing:

```
config.toml
branches/
tags/
snapshots/
trees/
packs/
```

### Configuration File
The configuration file contains
```
Version = 2
ChunkerPolynomial = 14583769067428845
CompressionLevel = 0
```

- `Version` is an integer indicating the repository format version
- `ChunkerPolynomial` is another integer indicating the current chunker polynomial
- `CompressionLevel` is the zip compression level, where 0 is store-on and 8 is the DEFLATE algorithm. Currently only levels 0 and 8 are supported.
### Snapshots Structure
The snapshots directory contains a set of subdirectories named after each project, e.g.

```
arduino/
electronics/
fdtd-results/
ham-radio/
kicad/
``` 

Within each project subdirectory is a set of TOML files, one for each branch, where 
the toml file basename is that of the branch name, eg: `main.toml`. Branch names are 
therefore required to be valid filenames. The branch files store the most recent 
snapshot ID for that branch, eg. `CommitID = "e0d25d16fd071dfd3de4e4b17ef076e152be289b"`

### Snapshots Structure
The snapshots directory contains a set of subdirectories named after each project, e.g.

```
arduino/
electronics/
fdtd-results/
ham-radio/
kicad/
``` 

Within each project subdirectory is a set of TOML files, one for each tag, where 
the toml file basename is that of the tag name, eg: `v0.2.2.toml`. Tag names are 
therefore required to be valid filenames. The tag files store the most recent 
snapshot ID for that branch, eg. `CommitID = "e0d25d16fd071dfd3de4e4b17ef076e152be289b"`

### Snapshots Structure
The snapshots directory contains a set of subdirectories named after each project, e.g.

```
arduino/
electronics/
fdtd-results/
ham-radio/
kicad/
``` 

Within each project subdirectory is a set of JSON files, one for each snapshot which
are named with random 40-character strings:

```
5bd64d62b42a8c15f3241b040bb8564b89cb9cf7.json
8babe30756618c4dcc5f7720774344918dea8ad9.json
8c2cbc7fe87d2de06d6e4dbefdbade282eeef344.json
94f792bc1128766c1791aa0e41a7262bbd99523d.json
625c003452e8363209998efccebfbfb054aa621f.json
```

Here is an example file which contains:
- Snapshot ID
- Snapshot branch
- Commit message
- Commit time
- Snapshot ID of the previous commit (ParentIDs)
- Files and metadata
- A list of chunk IDs for each chunk in the snapshot

Unlike Git, empty directories are included in the snapshot. For each file or directory:
- Path relative to the working directory **parent** directory
- Modification Time
- Size
- Hash

```
{
  "ID": "dfb50c2694ac595dc34db02e93623841f5a531a9",
  "Branch": "main",
  "Message": "add repeaters",
  "Time": "2021/03/06 21:12:04",
  "ParentIDs": [
    "94f792bc1128766c1791aa0e41a7262bbd99523d"
  ],
  "Files": [
    {
      "Path": "Ham Radio/",
      "ModTime": "2021/03/06 21:08:57",
      "Size": 0,
      "Hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    },
    {
      "Path": "Ham Radio/.dupver/",
      "ModTime": "2021/03/06 21:08:57",
      "Size": 0,
      "Hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    },
    {
      "Path": "Ham Radio/Repeaters.xlsx",
      "ModTime": "2021/03/03 21:15:58",
      "Size": 33887,
      "Hash": "aaca55f733e875dd9e7c9930424acbde46238104d48ab803e437469ceae5746a"
    },
    {
      "Path": "Ham Radio/TH_D74A_Manual_ReducedSize.pdf",
      "ModTime": "2021/03/04 06:10:21",
      "Size": 1255522,
      "Hash": "b72e00ef58fa0264932d8febf01189e8a6f8153e0cf03975399323bd3a685970"
    },
    {
      "Path": "Ham Radio/.dupver/config.toml",
      "ModTime": "2021/03/06 21:10:07",
      "Size": 146,
      "Hash": "56002d2b987d93ad5b10207328ab93ee1544d3799b054dc385cc950c57802d8d"
    }
  ],
  "ChunkIDs": [
    "b5acd270ff83d571a8e3d15347b4b4136ea05d132fcdc9157d96dd25f771d122",
    "f35c3c2a2ef8daebd0ffaabf736d1f91afd100c24bddf595ab148db3d7d82e9d"
  ]
}
```

### Trees Structure
The trees directory has no subdirectories and consists of a collection of json files
names, each named with a random 40-character file name:

```
5bd64d62b42a8c15f3241b040bb8564b89cb9cf7.json
8babe30756618c4dcc5f7720774344918dea8ad9.json
8c2cbc7fe87d2de06d6e4dbefdbade282eeef344.json
09e63e8c71982a56b02368a9e06419fe0063125c.json
33c5cf0042911a68a1424f104653e51e52a72ea4.json
```


Each file consists of a directionary mapping chunk hashes to pack file IDs:

```
{
  "0156ed74d9324c5931482531dfb96e83e919e2829fa85be963b53bc13e423da0": "efac7d09e71e9685839d952f534bcde0d0dbc32fa37793f7d08c37fff92979c6",
  "4658225a1debad6800433e1d92c0f951daa21d5d5215a86835453abe06ee2437": "efac7d09e71e9685839d952f534bcde0d0dbc32fa37793f7d08c37fff92979c6",
  "b045c6a0ca53a4553ded9a5b39fe464dbae11a096fd228ab7b4d866cffa4170d": "efac7d09e71e9685839d952f534bcde0d0dbc32fa37793f7d08c37fff92979c6",
  "dd6d79d6d54f7e53ac619e4cbea0494c8f9dff59ed580aa5c2346ab46ddf25ac": "efac7d09e71e9685839d952f534bcde0d0dbc32fa37793f7d08c37fff92979c6",
  "e90141fd453446398a01ebe5e7e49e2de3c66fef092b7c1de015f1d5a601738a": "efac7d09e71e9685839d952f534bcde0d0dbc32fa37793f7d08c37fff92979c6"
}
```

### Packs Structure
The packs directory contains the actual data. It consists of a set of subdirectories
with two-character file names:

```
00
0f
2c
3b
3d
4a
9c
11
16
```

Each subdirectory contains those pack files where the first two characters of the 
pack file name matches the subdirectory name. Pack files are in .zip format and
are given random 64-character base names, eg. `115c2d137ea0571747037cc03c0c2103fa3380d28ebe153494c52a42323c4d6b.zip`
Pack files contain a set of chunks, where the chunk name is it's SHA hash, eg.

```
84ebe552402dcf7bbdf3657cd8de1eabd39f36f915e64844da220060299c1faa
29941845ac6ce3499124d2a800dd5d3293ceb8cc0790aeffb9c6fab34844fb65
d8508d7483b83234266b94b8f09c3c14bd169ec2b04d464e1c4dc1826a6f4a15
```

