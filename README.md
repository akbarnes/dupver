![Go](https://github.com/akbarnes/dupver/workflows/Go/badge.svg)

# Dupver
Dupver is a minimalist deduplicating version control system in Go based on 
the Restic chunking library. It is most similar to the binary
version control system Boar https://bitbucket.org/mats_ekberg/boar/wiki/Home, 
and to a lesser extent the commercial offering Perforce,
though it borrows design conventions from the deduplicating backup
applications Duplicacy, Restic and Borg.
Dupver does not track files, rather it stores snapshots more like
a backup program. Rather than traverse directories itself, Dupver
uses an (uncompressed) tar file as input. Not that *only* tar files
are accepted as input as Dupver relies on the tar container to
 1. provide the list of files in the snapshot
 2. store metadata such as file modification times and permissions
 
Dupver uses a centralized repository to take advantage of deduplication 
between working directories. This means that dupver working 
directories can also be git repositories or subdirectories of git
repositories. 

Dupver has not been tested on repository sizes more than
a few GB, but it is expected to scale up to the low 100's of GB. 

## Setup
[Binary releases](https://github.com/akbarnes/dupver/releases) are provided for MacOS/Linux/Windows. Copy them somewhere in your path. To build from source run `go build` or `go install`.

## Update repos to v0.4 format
1. Add prev pointers to snapshot files
2. Add branch & head pointers
3. Convert snapshot filenames to drop date

## TODO
- [ ] check that branch cannot be called head. this check may not be needed anymore
- [x] fix writing head from outside of wd
- [ ] remove auto-generated tar files after chunking
- [ ] use toml for heads and branches?

## Usage

### Initialize repository
Initialize a repository with
`dupver repo init repo_path`

### Initialize project working directory
From inside the working directory
`dupver -r repo_path init -p project_name`

Or from the parent directory
`dupver -r repo_path init -p project_name working_directory`

### Commit
Stage your files by adding them to aø tar file

`tar cfv tarfile.tar file1 file2 file`

Commit the tarfile with
`dupver commit -m "commit message" tarfile.tar`

Alternatively run 
`dupver commit -m "commit message"`
from within the working directory

### List commits
List all commits
`dupver log`

List a specific commit
`dupver log commit_id`

### Copy
Copy a commit 
`dupver copy dest_repo snapshot_id`

### Check if files are modified/added
`dupver status`

### Restore
Restore another commit
`dupver checkout commit_id`

Restore to a particular file
`dupver checkout -o filename.tar commit_id `
