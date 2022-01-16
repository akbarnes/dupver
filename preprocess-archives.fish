#!/usr/bin/fish

set arch_types docx xlsx pptx vsdx slx zip qgz

for arch_type in $arch_types
    echo $arch_type

    for arch_file in **.$arch_type
        echo $arch_file

        if string match '.dupver*' $arch_file
            continue
        end

        set extracted_folder ".dupver_archives/$arch_file"
        rm -fr "$extracted_folder"
        mkdir -p "$extracted_folder"
        unzip -d "$extracted_folder" $arch_file  
    end
end
