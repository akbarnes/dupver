#!/usr/local//bin/fish
# go get github.com/restic/chunker
# go get github.com/BurntSushi/toml
# go build dupver.go commit.go pack.go config.go randstring.go

# set InstallFolder "$HOME/AppData/Local/Executables"
# if (-not (test-path $InstallFolder)) { mkdir $Installfolder }
# copy dupver.exe $InstallFolder
 
set repo_path "$HOME/.dupver_repo"
echo "Initialize repo $repo_path"
echo -------------------------------------
if test -d $repo_path 
    rm -fr $repo_path 
end
dupver -init-repo -r $repo_path

set workdir_folder "FDTD"
set workdir_name (echo $workdir_folder | tr '[:upper:]' '[:lower:]')
set workdir_path "$HOME/Results/$workdir_folder"
if test -d $workdir_path/.dupver
    rm -fr $workdir_path/.dupver 
end


echo "Initialize workdir $workdir_name in $workdir_path"
echo -------------------------------------
dupver -init -d $workdir_folder -r $repo_path

set tar_name "$workdir_name.tar"
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
