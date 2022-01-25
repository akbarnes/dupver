$ArchiveTypes = (Get-Content .dupver_archive_types)

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
        $ExpandedPath = Join-Path -Path ".dupver_archives" -ChildPath $ArchiveFolder

        if (Test-Path -Path $ExpandedPath) {
            del $ExpandedPath -Recurse
        }

        Write-Output "Creating folder: $ExpandedPath"
        mkdir $ExpandedPath

        Write-Output "Expanding: $ArchivePath"
        Expand-Archive -Path $ArchiveFile -DestinationPath $ExpandedPAth
    }
}


