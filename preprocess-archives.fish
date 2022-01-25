#!/usr/bin/fish

set arch_types (cat .dupver_archive_types)

for arch_type in $arch_types

    if test (string length $arch_type) -eq 0
        continue
    end

    echo $arch_type

    for arch_file in **.$arch_type
        echo $arch_file

        if string match '.dupver*' $arch_file
            continue
        end

        set arch_folder (string replace ".$arch_type" "_$arch_type" $arch_file)
        set extracted_folder ".dupver_archives/$arch_folder"
        rm -fr "$extracted_folder"
        mkdir -p "$extracted_folder"
        unzip -d "$extracted_folder" $arch_file  
    end
end
