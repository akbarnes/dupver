#!/usr/local/bin/fish

set files *.json

for f in $files 
    echo $f
end

set cid (echo $files[-1] | cut -c 22-61)
echo $cid
set branch main

set f head.json
echo '{' >$f
echo "  \"CommitID\": \"$cid\"" >>$f
echo "  \"BranchName\": \"$branch\"" >>$f
echo '}' >>$f

