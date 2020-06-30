$WorkdirFolder =  "House"
$WorkdirPath = "$HOME\Documents\$WorkdirFolder"
 
$RepoPath = "$HOME\temp\.dupver_repo"
echo "Initialize repo $RepoPath"
echo ----------------------------------------------------
if (test-path $RepoPath) {
    del -Force -Recurse $RepoPath
}
dupver -r $RepoPath repo init 

$WorkdirName = $WorkdirFolder.ToLower()

echo ''
echo "Initialize workdir $WorkdirName ($WorkdirFolder) in $WorkdirPath"
echo ----------------------------------------------------
dupver -r $RepoPath init $WorkdirFolder

$TarName = "${WorkdirName}.tar"
if (test-path $TarName) { del -Force $TarName }
tar cfv $TarName $WorkdirFolder

echo ''
echo "Checking in $TarName to $WorkdirName"
echo ----------------------------------------------------
dupver commit $TarName

echo ''
echo "List $WorkdirName in $RepoPath"
echo ----------------------------------------------------
dupver -d $WorkdirName log

$Snapshots = (dir $RepoPath\snapshots\$WorkdirName\*.json)
$SnapshotId = $Snapshots[0].basename.substring(21,40)

echo ''
echo "Check out commit $SnapshotId"
echo ----------------------------------------------------
$OutTarName = (dupver -d $WorkdirFolder -q checkout $SnapshotId)
echo "Exported to $OutTarName"
echo "Files in tar:"
tar tvf $OutTarName

