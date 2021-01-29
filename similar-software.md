# Software Similar to Dupver

## Binary VCS

## Boar 
https://bitbucket.org/mats_ekberg/boar/wiki/Home
An earlier binary VCS written in Python. This also uses as central repository. It has a simple and well-documented repository format, though only performs file-level deduplication. Block-level deduplication is available, but Linux-only. 

## Dud
https://github.com/kevin-hanselman/dud

## Perforce
Commercial software. AFAIK performs file-level deduplication with compression at the file level.

## Deduplicating Storage

### Zbackup
This allows for a single file such as a tar archive to be stored with block-level deduplication and retrieved with a unique hash.

### Bup
Deduplicating backup software based on the Git packfile format. I'm lumping it under storage because it still retains the ability to store a single file such as a tar archive and retrieve it with a unique hash.

## Deduplicating Backup Software
This includes Borg, Restic and Duplicacy.
