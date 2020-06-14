$WorkdirFolder =  "Computer"
$WorkdirPath = "$HOME\OneDrive\Documents\$WorkdirFolder"
 
$RepoPath = "$HOME\.dupver_repo"
echo "Initialize repo $RepoPath"
echo -------------------------------------
if (test-path $RepoPath) {
    del -Force -Recurse $RepoPath
}
dupver -init-repo -r $RepoPath

$WorkdirName = $WorkdirFolder.ToLower()

echo ''
echo "Initialize workdir $WorkdirName in $WorkdirPath"
echo -------------------------------------
dupver -init -d $WorkdirFolder -r $RepoPath

$TarName = "${WorkdirName}.tar"
if (test-path $TarName) { del -Force $TarName }
tar cfv $TarName $WorkdirFolder

echo ''
echo "Checking in $TarName to $WorkdirName"
echo -------------------------------------
dupver -ci -f $TarName

echo ''
echo "List $WorkdirName in $RepoPath"
echo -------------------------------------
dupver -list -d $WorkdirName

$Snapshots = (dir $RepoPath\snapshots\$WorkdirName\*.json)
$SnapshotId = $Snapshots[0].basename.substring(21,40)
$Snapshots = (dir $RepoPath\snapshots\$WorkdirName\*.json)
$SnapshotId = $Snapshots[0].basename.substring(21,40)

echo "Check out commit $SnapshotId"
echo -------------------------------------
dupver -co -d $WorkdirFolder -s $SnapshotId
echo "Check out commit $SnapshotId"
echo -------------------------------------
dupver -co -d $WorkdirFolder -s $SnapshotId