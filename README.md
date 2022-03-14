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

Binary [releases](github://github.com/akbarnes/dupver/releases) are provided. 
Copy the executable into your path and set permissions appropriately. Otherwise
follow the instructions below to build from source. This assumes that Go 1.17
is present on your system.

### Building

``` bash
go mod init dupver
go mod tidy
go get github.com/akbarnes/dupver
```

### Configuration

When `dupver` is run for the first time, it will create a default preferences
file in `$HOME/.dupver.json` 

`` json
{
  "PrefsMajorVersion": 2,
  "PrefsMinorVersion": 0,
  "DupverMajorVersion": 4,
  "DupverMinorVersion": 0,
  "Editor": "vi",
  "DiffTool": "kdiff3",
  "ArchiveTool": "7z"
}
```

At present only the `DiffTool` entry can be modified. The `Editor` entry is
reserved for future usage and currently only `7z` is supported as an archive
handler. If `7z` is not in your path `ArchiveTool` can be modified to specify
an absolute path to the `7z` executable.  

Ignored files are specified with `.dupver_ignore` within the repository. Both
single- and double-star globs are supported, where single-stars match any 
character other than path separators and double-stars match any character
*incuding* path separators.

``` fish
.git/**
*.log
```

## Usage

### Initialization

To initialize the repository

``` bash
dupver init [--random-poly|-r]
```

This is not required as calling `dupver commit` within an unitialized 
repository will initialize the repository with default parameters prior
to making the initial comit. By default, dupver will initialize the 
repository with a known good polynomical (used to determine 
when chunk boundaries occur). However, previding the `--random-poly`
flag will generate a random polynomial.

### Commit

To commit from within a directory

``` bash
dupver {commit|ci} [MESSAGE]
```

### Log

To list all the snapshots:

``` bash
dupver log
```

This takes an optional "quiet" argument, which when enabled causes log to only print the snapshot ids in chronological order.

``` bash
dupver log [-quiet|-q]
```

To list the files in a particular snapshot:

``` bash
dupver log <snapshot_id>
```

When the "quiet" argument is specified, this causes log to only print the relative file paths

``` bash
dupver log [-quiet|q] <snapshot_id>
```

### Checkout

This takes an optional argument to specify an output folder. To checkout a snapshot:

``` bash
dupver {checkout|co} [-out|-o <output_folder>] <snapshot_id>
```

To check out the last snapshot, there are the shortcuts

``` bash
dupver co {latest|last}
```

### Repack

This consolidates small packs from multiple commits. It will also skip chunks that are not associated 
with a snapshot, allowing for deletion of snapshots.

``` bash
dupver repack
```

## Handling Archives

Special handling for archive files requires the 7-zip `7z` executable to be present on the system. 
The archive types should be specified in `.dupver_archive_types` as below. Archive files are converted
via `7z` to a store-only `.zip` file in a temporary folder which is then added to the archive, allowing
for deduplication to be performed. Because the stored archive is chunked, it is likely that some
chunks will span two or more files within the archive, so effective deduplication will depend on the
file ordering within the archive not changing. By default this is by name. 

Warning! The `7z` implementation on Linux will not preserve user/group permissions and currently
`dupver` does not support other archive tools, so use this option with caution. See the following 
section on how to use `gzip` to create compressed archives that do not require special handling if 
this affects you.

Example `.dupver_archive_types`:

```
zip
7z
tgz
tbz
txz
docx
xlsx
pptx
vsdx
mlx
slx
qgz
```

### Gzip

Gzip includes the `--rsyncable` option that resets the dictionary periodically so small changes
in the raw file will only cause local changes in the compressed file. If this is used then
there is not a need for special handling of gzip files. Regrettably, AFAIK the `--rsyncable` option
is only present in `gzip` itself and not in programs or libraries that call it, including both
`tar` and the Python gzip library.

Usage of `gzip` with `--rsyncable`:

``` bash 
gzip [--keep|k] --rsyncable <raw_file_or_folder>
```

(Warning! This deletes the original file! Use `--keep` to preserve it.)
