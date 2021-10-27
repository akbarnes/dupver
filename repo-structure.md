# Repo Structure

## Directory Organization
- `dupver_settings.toml`
- `snapshots/`
- `files/`
- `trees/`
- `tags/` (to be implemented later)
- `branches/` (to be implemented later)

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
- `ctime`
- `owner`
- `group`
- `permissions`

## Trees `trees/<snapshot_id.json>`
Dictionary indexed by pack id where each entry is a list of chunk ids

# Preferences `.config/dupver/dupver_prefs.json`
- default prime?
- default compression

