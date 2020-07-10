#!/usr/local/bin/fish

for f in *.json
    set f2 (echo $f | cut -c 22-67)
    echo "mv '$f' '$f2'"
    mv "$f" "$f2"
end
