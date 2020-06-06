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
* [ ] combine small chunks into single files like restic
* [ ] store packs in gotiny/msgpack/protobuf
* [ ] decompress archives before deduplicating
* [ ] use buffered file io for speed
* [ ] print when deduplication occurs
* [x] identify revisions with hashes rather than integers so repositories can be merged
* [ ] support deletions of snapshots
* [ ] support copy between repos
* [x] file metadata: mtime, ctime, hash, permissions?
* [ ] support `status` and `diff` commands
* [x] move repository out of working directory
* [ ] check to make sure I don't overwrite workdir config
* [ ] use json for snapshots. compression probably isn't needed as files are about 100 kB for 1 GB data 
* [x] don't write TOML by hand
* [ ] don't hardcode chunker polynomial
* [ ] read repo path & workdir name from config files
* [ ] don't ignore errors
* [ ] multiple repos and names for repos
* [ ] unit tests


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
`dupver -init`

### Backup
Stage your files by adding them to a gzipped tar file

`tar cfvz tarfile.tgz file1 file2 file`

Commit the tarfile with
`dupver -backup -f tarfile.tgz -m "commit message"`

### List commits
List all commits
`dupver -list`

List a specific commit
`dupver -list -r 1`

List the last commit
`dupver -list -r -1`

List the 2nd to last commit
`dupver -list -r -2`

### Restore
Restore the last commit to snapshot<n>.tgz
`dupver -restore`

Restore another commit
`dupver -restore -r 1`

Restore to a particular file
`dupver -restore -f filename.tgz`
