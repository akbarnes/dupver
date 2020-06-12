#!/usr/local//bin/fish
# go get github.com/restic/chunker
# go get github.com/BurntSushi/toml
# go build dupver.go commit.go pack.go config.go randstring.go

# set InstallFolder "$HOME/AppData/Local/Executables"
# if (-not (test-path $InstallFolder)) { mkdir $Installfolder }
# copy dupver.exe $InstallFolder

set workdir_folder "Electric"
set workdir_path "$HOME/Data/$workdir_folder"
  
set repo_path "$HOME/.dupver_repo"
echo "Initialize repo $repo_path"
echo -------------------------------------
if test -d $repo_path 
    rm -fr $repo_path 
end
dupver -init-repo -r $repo_path

set workdir_name (echo $workdir_folder | tr '[:upper:]' '[:lower:]')
rm -fr $workdir_path/*

echo "Initialize workdir $workdir_name in $workdir_path"
echo -------------------------------------
dupver -init -d $workdir_folder -r $repo_path

# set i 0
# echo "Backup set $i"
# echo =====================================

# echo "Copy files to $workdir_path"
# echo -------------------------------------
# # cp Static/Electric/ACTIVSg70k/*.csv $workdir_name
# cp Static/Electric/ACTIVSg70k/*.PWB $workdir_name
# cp Static/Electric/ACTIVSg70k/*.pwd $workdir_name
# cp Static/Electric/ACTIVSg25k/*.pwb $workdir_name
# cp Static/Electric/ACTIVSg25k/*.pwd $workdir_name
# cp Static/Electric/ACTIVSg10k/*.pwb $workdir_name
# cp Static/Electric/ACTIVSg10k/*.pwd $workdir_name


# set tar_name "$workdir_name$i.tar"
# if test -e $tar_name
#     rm $tar_name
# end
# tar cfv $tar_name $workdir_folder

# echo "Checking in $tar_name to $workdir_name"
# echo -------------------------------------
# dupver -ci -f $tar_name

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

echo "Checking in $tar_name to $workdir_name"
echo -------------------------------------
dupver -ci -f $tar_name

echo "List $workdir_name in $repo_path"
echo -------------------------------------
dupver -list -d $workdir_name

# set Snapshots (dir $repo_path/snapshots/$workdir_name/*.toml)
# set SnapshotId $Snapshots[0].basename.substring(21,40)

# echo "Check out commit $SnapshotId"
# echo -------------------------------------
# dupver -co -d $workdir_name -s $SnapshotId
