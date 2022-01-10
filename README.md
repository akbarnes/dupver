![Go](https://github.com/akbarnes/dupver/workflows/Go/badge.svg)

# Dupver
Dupver is a minimalist deduplicating version control system in Go based on 
the Restic chunking library. It is most similar to the binary
version control system [Boar](https://github.com/mekberg/boar).
though it borrows design conventions from the deduplicating backup
applications Duplicacy, Restic and Borg.
Dupver does not track files, rather it stores commits as snapshots more like
a backup program. 

Dupver has not been tested on repository sizes more than
a few GB, but it is expected to scale up to the low 100's of GB. 

There are a number of [similar software projects](similar-software.md) both
in terms of technical implementation and use case.

## Installation
To build:
``` bash
go mod init dupver
go mod tidy
go get github.com/akbarnes/dupver
```

## Usage

### Commit
To commit from within a directory
``` bash
dupver commit 'a message' 
dupver ci 'a message' 
```

### Log
To list all the snapshots:
``` bash
dupver log
```

This takes an optional "quiet" argument, which when enabled causes log to only print the snapshot ids in chronological order.
``` bash
dupver log -quiet
dupver log -q
```

To list the files in a particular snapshot:
``` bash
dupver log snapshot_id
```
This takes an optional "quiet" argument, which when enabled causes log to only print the relative file paths
``` bash
dupver log -quiet snapshot_id
dupver log -q snapshot_id
```

## Checkout
This takes an optional argument to specify an output folder. To checkout a snapshot:
``` bash
dupver checkout snapshot_id
dupver co snapshot_id
dupver -out output_folder co snapshot_id
dupver -o output_folder co snapshot_id
```

## Repack
This consolidates small packs from multiple commits. It will also skip chunks that are not associated with a snapshot, allowing for deletion of snapshots.
``` bash
dupver repack
```
