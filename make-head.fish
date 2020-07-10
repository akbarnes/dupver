#!/usr/local/bin/fish

set files *-*.json

for f in $files 
    echo $f
end

set last_file $files[-1]
echo "Last file: $last_file"

set cid (echo $files[-1] | cut -c 22-61)
echo "Last commit id: $cid"
set branch main

set f head.json
echo '{' >$f
echo "  \"CommitID\": \"$cid\"" >>$f
echo "  \"BranchName\": \"$branch\"" >>$f
echo '}' >>$f

set f main.json
echo '{' >$f
echo "  \"CommitID\": \"$cid\"" >>$f
echo '}' >>$f

bat head.json
bat main.json
