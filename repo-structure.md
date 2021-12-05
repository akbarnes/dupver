# Repo Structure

## Directory Organization
- `dupver_settings.toml`
- `snapshots/`
- `files/`
- `trees/`
- `packs/`
- `tags/` (to be implemented later)
- `branches/` (to be implemented later)
- `head.json`

## Repo Settings `dupver_settings.json`
- `repo_version`
- `dupver_version`
- `compression_level`
- `prime`

## Snapshots `snapshots/<snapshot_id>.json`
- `username` (to be added later, this can be queried)
- `hostname` (to be added later, this can be queried)
- `date` (human-readable or iso format)
- `timestamp` (unix epoch, really needed if using iso format?)
- `commit_id`
- `message`

## File Listings `files/<snapshot_id>.json`
Dictionary indexed by filename where each entry has
- `md5hash`?
- `length` 
- `mtime`
- `chunks`

Note: removing permissions because git doesn't have it

## Trees `trees/<snapshot_id.json>`
Dictionary indexed by pack id where each entry is a list of chunk ids

# Preferences `.config/dupver/dupver_prefs.json`

- default prime?
- default compression

## Thoughts on Using Key-Value Store

- Store snapshots (time, message) indexed by snapshot id
- Store pack ids indexed by chunk id

- Store snapshot file list as json file
