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


## Setup
Binary releases are provided for MacOS/Windows. Copy them somewhere in your path. Building from source requires the chunker and toml libraries. Install them and build
```
go get github.com/restic/chunker
go get github.com/BurntSushi/toml
go build dupver.go
```

Copy the executable somewhere in your path. Build scripts for 
MacOS/*nix/Windows are included, see

* `setup` (MacOS, *nix)
* `Install-Dupver.ps1` (Windows)

You will need to edit them to set your desired install folder.

## Usage

### Initialize repository
Initialize a repository with
`dupver -init-repo`

### Initialize working directory
From inside the working directory
`dupver -init -w workdir_name -r repo_path`

Or from the parent directory
`dupver -init -w workdir_name -d working_directory -r repo_path`

### Commit
Stage your files by adding them to aø tar file

`tar cfv tarfile.tar file1 file2 file`

Commit the tarfile with
`dupver -commit -f tarfile.tar -m "commit message" -d working_directory`

Commit (shorhand)
`dupver -ci -f tarfile.tar -d working_directory`

### List commits
List all commits
`dupver -list`

List a specific commit
`dupver -list -c commit_id`

List the last commit
`dupver -list -c commit_id`

List the 2nd to last commit
`dupver -list -c commit_id`

### Restore
Restore another commit
`dupver -checkout -c commit_id`

Restore to a particular file
`dupver -checkout -f filename.tar`
