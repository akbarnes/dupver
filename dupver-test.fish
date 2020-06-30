#!/usr/local//bin/fish
set workdir_folder "Electric"
set workdir_path "$HOME/Data/$workdir_folder"

set repo_path "$HOME/temp/.dupver_repo"
echo ''
echo "Initialize repo $repo_path"
echo -------------------------------------
if test -d $repo_path 
    rm -fr $repo_path 
end

mkdir -p $HOME/temp
dupver repo init $repo_path

set workdir_name (echo $workdir_folder | tr '[:upper:]' '[:lower:]')
rm -fr $workdir_path/*

echo ''
echo "Initialize workdir $workdir_name in $workdir_path"
echo -------------------------------------
dupver -r $repo_path init $workdir_folder

if test 1 -eq 1
    set i 0
    echo ''
    echo "Backup set $i"
    echo =====================================

    echo "Copy files to $workdir_path"
    echo -------------------------------------
    # cp Static/Electric/ACTIVSg70k/*.csv $workdir_name
    # cp Static/Electric/ACTIVSg70k/*.PWB $workdir_name
    # cp Static/Electric/ACTIVSg70k/*.pwd $workdir_name
    # cp Static/Electric/ACTIVSg25k/*.pwb $workdir_name
    # cp Static/Electric/ACTIVSg25k/*.pwd $workdir_name
    cp Static/Electric/ACTIVSg10k/*.pwb $workdir_name
    cp Static/Electric/ACTIVSg10k/*.pwd $workdir_name


    set tar_name "$workdir_name$i.tar"
    if test -e $tar_name
        rm $tar_name
    end
    tar cfv $tar_name $workdir_folder

    echo "Checking in $tar_name to $workdir_name"
    echo -------------------------------------
    dupver commit $tar_name
end

# set i 1
# echo "Backup set $i"
# echo =====================================

# echo "Copy files to $workdir_path"
# echo -------------------------------------
# cp Static/Electric/ACTIVSg2000/ACTIVSg2000.PWB $workdir_name
# cp Static/Electric/ACTIVSg2000/ACTIVSg2000.pwd $workdir_name
# cp Static/Electric/ACTIVSg2000/ACTIVSg2000.m $workdir_name
# # cp Static/Electric/ACTIVSg2000/*.csv $workdir_name

# set tar_name "$workdir_name$i.tar"
# if test -e $tar_name
#     rm $tar_name
# end
# tar cfv $tar_name $workdir_folder

# echo "Checking in $tar_name to $workdir_name"
# echo -------------------------------------
# dupver -ci -f $tar_name

if test 1 -eq 0
    set i 2
    echo "Backup set $i"
    echo =====================================

    echo "Copy files to $workdir_path"
    echo -------------------------------------
    cp Static/Electric/ACTIVSg2000/ACTIVSg2000.RAW $workdir_name
    cp ACTIVSg2000_Mod.RAW $workdir_name
    cp ACTIVSg2000_Copy.RAW $workdir_name
    # cp Static/Electric/ACTIVSg2000/*.aux $workdir_name
    # cp Static/Electric/ACTIVSg2000/contab_ACTIVSg2000.m $workdir_name
    # cp Static/Electric/ACTIVSg2000/scenarios_ACTIVSg2000.m $workdir_name

    set tar_name "$workdir_name$i.tar"
    if test -e $tar_name
        rm $tar_name
    end
    tar cfv $tar_name $workdir_folder
end

echo ''
echo "Checking in $tar_name to $workdir_name"
echo -------------------------------------
dupver commit $tar_name

echo ''
echo "List $workdir_name in $repo_path"
echo -------------------------------------
dupver -d $workdir_name log

set snapshots $repo_path/snapshots/$workdir_name/*.json
set snapshot_id (basename $snapshots[1] | cut -c22-61)

echo ''
echo "Check out commit $snapshot_id"
echo -------------------------------------
dupver -d $workdir_name checkout $snapshot_id
