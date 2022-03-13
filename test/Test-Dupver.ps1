cd ..
go test
go get github.com/akbarnes/dupver/cmd/dupver

If (!$?) {
    cd test
    Throw "Unit tests failed"
}

cd test

$RepoPath = "test-repo"

If (!(Test-Path data)) {
    New-Item -Path data -Itemtype directory -Force
}

cd ./data

$BaseUrl = "https://egriddata.org/sites/default/files"

$FileUrls = "uiuc-150bus.pwb", "uiuc-150bus.pwd", "ACTIVSg200_0.PWB", "ACTIVSg200.pwd", "ACTIVSg_2000_0.PWB", "ACTIVSg2000.pwd"
$FileNames = "uiuc-150bus.pwb", "uiuc-150bus.pwd", "ACTIVSg200.pwb", "ACTIVSg200.pwd", "ACTIVSg_2000.pwb", "ACTIVSg2000.pwd"

Write-Host "Downloading test data if needed..."

for ($i = 0; $i -lt $FileUrls.Length; $i++) {
    $FileUrl = $FileUrls[$i]
    $FileName = $FileNames[$i]

    If (!(Test-Path $FileName)) {
        Invoke-Webrequest -Uri "$BaseUrl/$FileUrl" -Outfile $FileName
    }
}

Write-Host "Done"
cd ..

If (Test-Path $RepoPath) {
    Remove-Item -Force -Recurse $RepoPath/*
} Else {
    New-Item -Path $RepoPath -Itemtype directory -Force
}

#Write-Host ''
Write-Host "Initializing repo $RepoPath..."
#Write-Host '----------------------------------------------------'
cd $RepoPath
dupver init
Write-Host "Done"
Write-Host ''

Write-Host "Checking repo configuration..."
$cfg = Get-Content .dupver/repo_config.json | ConvertFrom-Json

If ($cfg.CompressionLevel -ne 0) {
    Throw "Dupver compression level in repo configuration isn't equal to 0"
} 

If ($cfg.PackSize -ne (500 * 1024 * 1024)) {
    Throw "Dupver pack size isn't equal to 500 MB"
} 

$poly0 = "3abc9bff07d9e5"
If ($cfg.ChunkerPoly -ne "$poly0") {
    $poly = $cfg.ChunkerPoly
    Throw "Dupver chunker polynomial of $poly in repo configuration isn't equal to $poly0"
} 

Write-Host "Done"
Write-Host "Passed unit tests"
Write-Host ''
del .dupver -Recurse -Force

copy ../data/uiuc-150bus.* .

Write-Host "Commiting files:"
dupver commit "First commit"
Write-Host ''

Write-Host "Dupver log output:"
dupver log
Write-Host ''

$SnapshotId = (dupver log -q).Substring(0, 8)
Write-Host "Dupver log output for first snapshot ${SnapshotId}:"
dupver log "$SnapshotId"

If (Test-Path "../test-repo-second-commit") {
    Remove-Item -Recurse -Force "../test-repo-second-commit" 
}

Write-Host "Checkout out files"
dupver checkout -out "../test-repo-first-commit" "$SnapshotId"
cd ..

Write-Host "Checking that files have been restored correctly..."
$RepoFiles = Get-ChildItem -Recurse test-repo

Foreach ($File in $RepoFiles) {
    if (Test-Path -PathType Container -Path $File) { 
        continue
    }

    if ($File.FullName.Contains(".dupver")) {
        continue
    }
    
    $FileName = $File.Name
    $OriginalFile = Get-Content -Raw "test-repo/$FileName"
    $RestoredFile = Get-Content -Raw "test-repo-first-commit/$FileName"
     
    If ($OriginalFile -ne $RestoredFile) {
         Throw "Binary content for file $FileName doesn't match"
    }
}

Write-Host "Done"
Write-Host "Passed unit tests"
Write-Host ''
	
cd test-repo
$NewFiles = "ACTIVSg200.pwb", "ACTIVSg200.pwd"
$NewFiles | % { copy "../data/$_" . }

Write-Host "Dupver status output after adding files:"
dupver status

Write-Host "Checking that status is correct after adding files..."
$StatusOutput = dupver status

If ($StatusOutput.Length -ne 2) { 
    $n = $StatusOutput.Length
	Throw "Length of status $n not equal to two for status $StatusOutput"
}

$FileStatus = @{}

Foreach ($Line in $StatusOutput) {
	$Status = $Line.Substring(0,1)
	$File = $Line.SubString(2,$Line.Length - 2)
	$FileStatus[$File] = $Status
}

Foreach ($File in $NewFiles) { 
    If (!$FileStatus.ContainsKey($File)) {
        Thow "Status output doesn't contain $File"
    }

    If ($FileStatus[$File] -ne "+") {
        Thow "Status output for $File not equal to +"
    }
}

Write-Host "Done"
Write-Host ''

Write-Host "Checking that repack works correctly..."
dupver commit "Second commit"
dupver repack

Write-Host "Checkout out files after repack"

If (Test-Path "../test-repo-second-commit") {
    Remove-Item -Recurse -Force "../test-repo-second-commit" 
}

dupver checkout -out "../test-repo-second-commit" last
cd ..

Write-Host "Checking that files have been restored correctly"
$RepoFiles = Get-ChildItem -Recurse test-repo

Foreach ($File in $RepoFiles) {
   $FileName = $File.Name
   $OriginalFile = Get-Content -Raw "test-repo/$FileName"
   $RestoredFile = Get-Content -Raw "test-repo-second-commit/$FileName"
    
   If ($OriginalFile -ne $RestoredFile) {
        Throw "Binary content for file $FileName doesn't match"
   }
}

Write-Host "Done"
Write-Host "Passed unit tests"
