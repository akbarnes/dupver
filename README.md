![Go](https://github.com/akbarnes/dupver/workflows/Go/badge.svg)

# Dupver
Dupver is a minimalist deduplicating version control system in Go based on 
the Restic chunking library. It is most similar to the binary
version control system Boar https://bitbucket.org/mats_ekberg/boar/wiki/Home.
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

## Setup
[Binary releases](https://github.com/akbarnes/dupver/releases) are provided for MacOS/Linux/Windows. Copy them somewhere in your path. Building from source requires the chunker and toml libraries. Install them and build
```
go get github.com/restic/chunker
go get github.com/BurntSushi/toml
go build dupver.go
```

Copy the executable somewhere in your path. Build scripts for 
MacOS/Linux/Windows are included, see

* `setup` (MacOS, Linux)
* `Install-Dupver.ps1` (Windows)

You will need to edit them to set your desired install folder.

## Usage

### Initialize repository
Initialize a repository with
`dupver init-repo repo_path`

### Initialize project working directory
From inside the working directory
`dupver init -r repo_path -p project_name`

Or from the parent directory
`dupver init -r repo_path -p project_name working_directory`

### Commit
Stage your files by adding them to aø tar file

`tar cfv tarfile.tar file1 file2 file`

Commit the tarfile with
`dupver -m "commit message" commit tarfile.tar`

### List commits
List all commits
`dupver log`

List a specific commit
`dupver log commit_id`

### Check if files are modified/added
`dupver status`

### Restore
Restore another commit
`dupver checkout commit_id`

Restore to a particular file
`dupver -f filename.tar checkout commit_id`
