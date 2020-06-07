# Dupver
Dupver is a minimalist deduplicating version control system in Go based on 
the Restic chunking library.
Dupver does not track files, rather it stores snapshots more like
a backup program. Rather than traverse directories itself, Dupver
uses a gzipped tar file as input. Not that *only* gzipped tar files
are accepted as input as Dupver relies on the tar container to
provide the list of files in the snapshot.

Dupver stores deduplicated chunks as individual gzipped files
under the .dupver folder of a repository. The commit history
is stored as a plaintext .toml file also under .dupver.

## TODO
* [x] zip/gzip compression
* [x] handle multiple files via stage file or zip/tar
* [x] restore
* [ ] support `status` and `diff` commands
* [ ] multiple repos and names for repos
* [ ] merge workdirs - check what Hg does as it didn't do lightweight branching
* [ ] combine small chunks into single files like restic
* [ ] store packs in gotiny/msgpack/protobuf
* [ ] decompress archives before deduplicating
* [ ] print when deduplication occurs
* [ ] merge workdirs - check what Hg does as it didn't do lightweight branching
* [ ] support deletions of snapshots
* [ ] support copy between repos
* [x] identify revisions with hashes rather than integers so repositories can be merged
* [x] file metadata: mtime, ctime, hash, permissions?
* [ ] check that I have all the metadata I need
* [x] move repository out of working directory
* [ ] check to make sure I don't overwrite workdir config
* [x] use json for snapshots. compression probably isn't needed as files are about 100 kB for 1 GB data 
* [x] don't write TOML by hand
* [x] read repo path & workdir name from config files
* [ ] don't ignore errors
* [ ] unit tests
* [ ] move main() functionality into functions
* [ ] record last commit (head)



### Maybe TODO
* [ ] users/emails
* [ ] commit tags
* [ ] workdir tag: default=main 
* [ ] commit graph
* [ ] branches



Notes
* Removed buffered io todo, I'm not doing small writes


## Setup
This requires the chunker and toml libraries. Install them and build
```
go get github.com/restic/chunker
go get github.com/BurntSushi/toml
go build dupver.go
```

Copy the executable somewhere in your path. Build scripts for 
Windows/*nix are included, see

* setup.sh (*nix)
* Install-Dupver.ps1 (Windows)

You will need to edit them to set your desired install folder.

## Usage

### Initialize repository
Initialize a repository with
`dupver -init-repo`

### Initialize working directory
From inside the working directory
`dupver -init -w workdir_name`

Or from the parent directory
`dupver -init -w workdir_name -d working_directory`

### Commit
Stage your files by adding them to a tar file

`tar cfv tarfile.tar file1 file2 file`

Commit the tarfile with
`dupver -commit -f tarfile.tar -m "commit message" -d working_directory`

Commit (shorhand)
`dupver -ci -f tarfile.tar -d working_directory`

### List commits
List all commits
`dupver -list`

List a specific commit
`dupver -list -s snapshot_id`

List the last commit
`dupver -list -s snapshot_id`

List the 2nd to last commit
`dupver -list -s snapshot_id`

### Restore
Restore another commit
`dupver -checkout -s snapshot_id`

Restore to a particular file
`dupver -checkout -f filename.tar`
