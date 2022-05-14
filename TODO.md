# TODO:
- [x] Add diff
- [x] Print warnings/info to stderr
- [x] Add repack/prune
- [x] Add global preferences (for diff tool or editor)
- [x] Fix timestamp issue on checkout - not storing milliseconds?
- [ ] Add test for partial checkouts
- [ ] Add unit test for snapshot creation
- [ ] Add unit test for file hash
- [ ] Add unit test for default prefs creation
- [ ] Add unit test for correct prefs version
- [ ] Add unit test for correct repo version
- [ ] Don't create new pack until new chunks are found

# Maybe TODO:
- [x] Re-add WorkDirPath to work well with Xojo
- [ ] Use -m flag for commit message, promit for message if absent
- [x] Add username/fullname/email to commits
- [x] Add quiet mode to just print snapshot ids on log
- [x] Preserve date stamps on checkout?
- [x] Add option to generate random polynomial on init
- [ ] Command to print repo stats
- [ ] Add codecov
- [x] Decompress archives before deduplicating
