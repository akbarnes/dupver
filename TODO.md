# TODO:
- [ ] Split ReadHead into ReadHead & ReadHeadFile, or combine ReadWorkDirConfig & ReadWorkDirConfigFile into a single function
- [ ] Add tags as pointers to commits
- [ ] remove parent-ids parameter
- [ ] Remove branches, though leave this in the repo. Use a mapping of "/" to "_branch_" in the workdir name going to the repo folders. Need to simplify stuff 
- [ ] Change working directory name to project name
- [ ] Warn  if project name isn't unique
- [ ] Add username/fullname/email to commits