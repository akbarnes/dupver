$WorkdirFolder =  "KiCAD"
$WorkdirPath = "$HOME\Documents\Books\$WorkdirFolder"
 
$RepoPath = "$HOME\.dupver_repo"
echo "Initialize repo $RepoPath"
echo ----------------------------------------------------
if (test-path $RepoPath) {
    del -Force -Recurse $RepoPath
}
dupver -init-repo 
# mkdir $RepoPath\snapshots

$WorkdirName = $WorkdirFolder.ToLower()

echo ''
echo "Initialize workdir $WorkdirName ($WorkdirFolder) in $WorkdirPath"
echo ----------------------------------------------------
dupver -init -d $WorkdirFolder

$TarName = "${WorkdirName}.tar"
if (test-path $TarName) { del -Force $TarName }
tar cfv $TarName $WorkdirFolder

echo ''
echo "Checking in $TarName to $WorkdirName"
echo ----------------------------------------------------
dupver -ci -f $TarName

echo ''
echo "List $WorkdirName in $RepoPath"
echo ----------------------------------------------------
dupver -list -d $WorkdirName

$Snapshots = (dir $RepoPath\snapshots\$WorkdirName\*.json)
$SnapshotId = $Snapshots[0].basename.substring(21,40)

echo ''
echo "Check out commit $SnapshotId"
echo ----------------------------------------------------
dupver -co -d $WorkdirFolder -s $SnapshotId

echo ''
echo "Check out commit $SnapshotId"
echo ----------------------------------------------------
dupver -co -d $WorkdirFolder -s $SnapshotId