# Dupver
A simple deduplicating version control system based on restic

## TODO
* [x] zip/gzip compression
* [x] handle multiple files via stage file or zip/tar
* [x] restore
* [ ] decompress archives before deduplicating
* [ ] use buffered file io for speed
* [ ] print when deduplication occurs

## Setup
This requires the chunker and toml libraries. Install them with
```
go get github.com/restic/chunker
go get github.com/BurntSushi/toml
```

## Usage

### Backup
Stage your files by adding them to a gzipped tar file

`tar cfvz tarfile.tgz file1 file2 file`

Commit the tarfile with
`go run dupver.go *backup *file tarfile.tgz *message "commit message"`

### List commits
List all commits
`go run dupver.go *list`

List a specific commit
`go run dupver.go *list *revision 1`

List the last commit
`go run dupver.go *list *revision *1`

List the 2nd to last commit
`go run dupver.go *list *revision *2`

# Restore
Restore the last commit to snapshot<n>.tgz
`go run dupver.go *restore`

Restore another commit
`go run dupver.go *restore *revision 1`

Restore to a particular file
`go run dupver.go *restore *file filename.tgz`
