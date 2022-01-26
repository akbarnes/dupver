#!/usr/bin/fish

set arch_types (cat .dupver_archive_types)
mkdir -p .dupver_archives/extracted
set file_list .dupver_archives/archive_files.txt

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
        set extracted_folder ".dupver_archives/extracted/$arch_folder"
        rm -fr "$extracted_folder"
        mkdir -p "$extracted_folder"
        unzip -d "$extracted_folder" $arch_file  

        echo "$arch_file" >> $file_list
        echo "$extracted_folder" >> $file_list
        echo '' >> $file_list
    end
end

