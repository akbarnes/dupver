#/usr/bin/fish
# go get github.com/restic/chunker
# go get github.com/BurntSushi/toml
# go build dupver.go commit.go pack.go config.go randstring.go

# $InstallFolder = "$HOME/AppData/Local/Executables"
# if (-not (test-path $InstallFolder)) { mkdir $Installfolder }
# copy dupver.exe $InstallFolder
 
$repo_path = "$HOME/.dupver_repo"
echo "Initialize repo $repo_path"
echo -------------------------------------
if (test-path $repo_path) { del -force -recurse $repo_path }
dupver -init-repo -r $repo_path

$workdir_name = "fdtd"
$workdir_folder = "FDTD"
$workdir_path = "$HOME/Results/${workdir_folder}"
if (test-path $workdir_path/.dupver) { del -force -recurse $workdir_path/.dupver }


echo "Initialize workdir $workdir_name in $workdir_path"
echo -------------------------------------
dupver -init -d $workdir_folder -r $repo_path

$tar_name = "${workdir_name}.tar"
if (test-path $tar_name) { del -force $tar_name }
tar cfv $tar_name $workdir_folder

echo "Checking in $tar_name to $workdir_name"
echo -------------------------------------
dupver -ci -f $tar_name


echo "List $workdir_name in $repo_path"
echo -------------------------------------
dupver -list -d Property

$Snapshots = (dir $repo_path/snapshots/$workdir_name/*.toml)
$SnapshotId = $Snapshots[0].basename.substring(21,40)

echo "Check out commit $SnapshotId"
echo -------------------------------------
dupver -co -d Property -s $SnapshotId