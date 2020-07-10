#!/usr/local/bin/fish

for f in *.json
    set f2 (echo $f | cut -c 22-67)
    echo "cp -p $f $f2"
    echo "rm $f"
    cp -p "$f" "$f2"
    # rm "$f"
end
