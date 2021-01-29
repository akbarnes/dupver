![Go](https://github.com/akbarnes/dupver/workflows/Go/badge.svg)

# Dupver
Dupver is a minimalist deduplicating version control system in Go based on 
the Restic chunking library. It is most similar to the binary
version control system Boar https://bitbucket.org/mats_ekberg/boar/wiki/Home, 
though it borrows design conventions from the deduplicating backup
applications Duplicacy, Restic and Borg.
Dupver does not track files, rather it stores commits as snapshots more like
a backup program. Rather than traverse directories itself, Dupver
uses an (uncompressed) tar file as input. Not that *only* tar files
are accepted as input as Dupver relies on the tar container to
 1. provide the list of files in the commit
 2. store metadata such as file modification times and permissions
 
Dupver uses a centralized repository to take advantage of deduplication 
between working directories. This means that dupver working 
directories can also be git repositories or subdirectories of git
repositories. 

Dupver has not been tested on repository sizes more than
a few GB, but it is expected to scale up to the low 100's of GB. 

There are a number of [similar software projects](similar-software.md) both
in terms of technical implementation and use case.


## Setup
[Binary releases](https://github.com/akbarnes/dupver/releases) are provided for MacOS/Linux/Windows. Copy them somewhere in your path. To build from source run `go build` or `go install`.

## Update repos to v0.4 format
1. Add prev pointers to snapshot files
2. Add branch & head pointers
3. Convert snapshot filenames to drop date

## Update repos to v0.8  format
1. Remove head pointers
2. Add main branch tag to snapshots

## Usage

### Initialize repository
Initialize a repository with
`dupver repo init repo_path`

Specify a name when initializing a repository
`dupver repo init repo_path repo_name`

### Initialize project working directory
From inside the working directory
`dupver -r repo_path init -p project_name`

Or from the parent directory
`dupver -r repo_path init -p project_name working_directory`

### Commit
Stage your files by adding them to a tar file

`tar cfv tarfile.tar file1 file2 file`

Commit the tarfile with
`dupver commit -m "commit message" tarfile.tar`

Alternatively run 
`dupver commit -m "commit message"`
from within the working directory

To commit to a specific repository run
`dupver commit -n repo_name -m "commit message"`    

### List commits
List all commits
`dupver log`

List a specific commit
`dupver log commit_id`

### Copy
Copy a commit 
`dupver copy dest_repo commit_id`

### Check if files are modified/added
`dupver status`

### Restore
Restore another commit
`dupver checkout commit_id`

Restore to a particular file
`dupver checkout -o filename.tar commit_id `
