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

## TODO
- Add partial checkouts
- Add trees
- Use snapshot structures from Xover
- Add diff?
- Add json output

## Installation
To build:
``` bash
go mod init dupver
go mod tidy
go get github.com/akbarnes/dupver
```

## Usage

### Commit
The `-msg` or `-m` message flag is optional, as is the `-commit` or `-ci` flag as commiting is the default action:
``` bash
gover -commit -msg 'a message' file1 file2 file3
gover -ci -msg 'a message' file1 file2 file3
gover -m  'a message 'file1 file2 file3
gover file1 file2 file3
```

### Log
This takes the optional `-json` or `-j` argument to output json for use with object shells. To list all the snapshots:
``` bash
gover -log
gover -json -log
gover -j -log
gover -l
```

To list the files in a particular snapshot:
``` bash
gover -log snapshot_time
gover -log -json snapshot_time
```

## Checkout
This takes an optional argument to specify an output folder. To checkout a snapshot:
``` bash
gover -checkout snapshot_time
gover -co snapshot_time
gover -out output_folder -co snapshot_time
gover -o output_folder -co snapshot_time
```
