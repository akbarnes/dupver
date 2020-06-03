# go get github.com/restic/chunker
# go get github.com/BurntSushi/toml
# go build dupver.go commit.go pack.go config.go randstring.go

# $InstallFolder = "$HOME\AppData\Local\Executables"
# if (-not (test-path $InstallFolder)) { mkdir $Installfolder }
# copy dupver.exe $InstallFolder
 
$repo_path = "$HOME\.dupver_repo"
if (test-path $repo_path) { del -force -recurse $repo_path }
dupver -init-repo -r $repo_path

$workdir_name = "property"
$workdir_path = "$HOME\Documents\Admin\Property"
if (test-path $workdir_path\.dupver) { del -force -recurse $workdir_path\.dupver }

cd $workdir_path
dupver -init -w $workdir_name -r $repo_path
cd ..

$tar_name = "${workdir_name}.tar"
if (test-path $tar_name) { del -force $tar_name }
tar cfv $tar_name $workdir_path

dupver -r $repo_path -w $workdir_name -ci -f $tar_name
dupver -r $repo_path -w $workdir_name -list
