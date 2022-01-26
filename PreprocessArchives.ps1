$ArchiveTypes = (Get-Content .dupver_archive_types)


$ArchiveMaps = @()
$ExpandedFolder = Join-Path -Path ".dupver_archives" -ChildPath "extracted"

if (!(Test-Path $ExpandedFolder)) {
    New-Item -ItemType Container -Path $ExpandedFolder
}

foreach ($ArchiveType in $ArchiveTypes) {
    if ($ArchiveType.Length -eq 0) { 
        continue
    }

    $ArchiveFiles = Get-ChildItem "*.$ArchiveType" -Recurse

    foreach ($ArchiveFile in $ArchiveFiles) {
        if (Test-Path -Path $ArchiveFile -PathType Container) {
            continue
        }

        if ($ArchiveFile.Name.StartsWith('~')) {
            continue
        }

        $ArchivePath = $ArchiveFile | Resolve-Path -Relative

        if ($ArchivePath.StartsWith("./") -or $ArchivePath.StartsWith(".\")) {
            $ArchivePath = $ArchivePath.Substring(2)
        }

        if ($ArchivePath.StartsWith(".dupver") -or $ArchivePath.StartsWith(".dupver_archive")) {
            continue
        }

        $ArchiveFolder = $ArchivePath.Replace(".$ArchiveType","_$ArchiveType")
        $ExpandedPath = Join-Path -Path $ExpandedFolder -ChildPath $ArchiveFolder
        $ArchiveMap = @{}
        $ArchiveMap["ArchiveFile"] = $ArchivePath
        $ArchiveMap["ExtractedFolder"] = Join-Path -Path "expanded" -ChildPath $ArchiveFolder 
        $ArchiveMaps += $ArchiveMap

        if (Test-Path -Path $ExpandedPath) {
            del $ExpandedPath -Recurse
        }

        Write-Host "Creating folder: $ExpandedPath"
        New-Item -ItemType Container -Path $ExpandedPath

        Write-Host "Expanding: $ArchivePath"
        Expand-Archive -Path $ArchiveFile -DestinationPath $ExpandedPath
    }
}

$ArchiveMaps | ConvertTo-Json | Out-File ".dupver_archives/archive_files.json"
