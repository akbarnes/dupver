# go get github.com/restic/chunker
# go get github.com/BurntSushi/toml
# go build dupver.go commit.go pack.go config.go randstring.go

# $InstallFolder = "$HOME\AppData\Local\Executables"
# if (-not (test-path $InstallFolder)) { mkdir $Installfolder }
# copy dupver.exe $InstallFolder
 
$RepoPath = "$HOME\.dupver_repo"
echo "Initialize repo $RepoPath"
echo -------------------------------------
if (test-path $RepoPath) { del -force -recurse $RepoPath }
dupver -init-repo -r $RepoPath

$WorkdirName = "property"
$WorkdirFolder = "Property"
$WorkdirPath = "$HOME\Documents\Admin\${WorkdirFolder}"
if (test-path $WorkdirPath\.dupver) { del -force -recurse $WorkdirPath\.dupver }


echo "Initialize workdir $WorkdirName in $WorkdirPath"
echo -------------------------------------
dupver -init -d $WorkdirFolder -w $WorkdirName -r $RepoPath

$TarName = "${WorkdirName}.tar"
if (test-path $TarName) { del -force $TarName }
tar cfv $TarName $WorkdirFolder

echo "Checking in $TarName to $WorkdirName"
echo -------------------------------------
dupver -ci -f $TarName


echo "List $WorkdirName in $RepoPath"
echo -------------------------------------
dupver -list -d Property

$Snapshots = (dir $RepoPath\snapshots\$WorkdirName\*.toml)
$SnapshotId = $Snapshots[0].basename.substring(21,40)

echo "Check out commit $SnapshotId"
echo -------------------------------------
dupver -co -s $SnapshotId