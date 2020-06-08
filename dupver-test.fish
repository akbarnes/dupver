#!/usr/local//bin/fish
# go get github.com/restic/chunker
# go get github.com/BurntSushi/toml
# go build dupver.go commit.go pack.go config.go randstring.go

# set InstallFolder "$HOME/AppData/Local/Executables"
# if (-not (test-path $InstallFolder)) { mkdir $Installfolder }
# copy dupver.exe $InstallFolder

set workdir_folder "Electric"
set workdir_path "$HOME/Data/$workdir_folder"
  
set repo_path "$HOME/.dupver_repo"
echo "Initialize repo $repo_path"
echo -------------------------------------
if test -d $repo_path 
    rm -fr $repo_path 
end
dupver -init-repo -r $repo_path

set workdir_name (echo $workdir_folder | tr '[:upper:]' '[:lower:]')
# if test -d $workdir_path/.dupver
#     rm -fr $workdir_path/.dupver 
# end
rm -fr $workdir_path/*

echo "Initialize workdir $workdir_name in $workdir_path"
echo -------------------------------------
dupver -init -d $workdir_folder -r $repo_path

# -rw-rw-r--@ 1 art  staff   251K May 16  2017 ACTIVSg200.EPC
# -rw-rw-r--@ 1 art  staff   100K May 16  2017 ACTIVSg200.RAW
# -rw-rw-r--@ 1 art  staff   492K Jun  4  2018 ACTIVSg200.aux
# -rw-rw-r--@ 1 art  staff   1.0M Jun  4  2018 ACTIVSg200.pwb
# -rw-rw-r--@ 1 art  staff   279K Jun  4  2018 ACTIVSg200.pwd
# -rw-rw-r--@ 1 art  staff   3.7M Mar 20  2017 ACTIVSg200.tsb
# -rw-rw-r--@ 1 art  staff   121K Jan 16  2018 ACTIVSg2000.pwd
# -rw-rw-r--@ 1 art  staff    17K Aug  6  2018 ACTIVSg200_GIC_data.gic
# -rw-rw-r--@ 1 art  staff    16K Jul 19  2017 ACTIVSg200_dynamics.dyd
# -rw-rw-r--@ 1 art  staff    13K Jul 19  2017 ACTIVSg200_dynamics.dyr
# -rw-rw-r--@ 1 art  staff    55K Oct 30  2017 case_ACTIVSg200.m
# -rw-rw-r--@ 1 art  staff    10K Oct 30  2017 contab_ACTIVSg200.m
# -rw-rw-r--@ 1 art  staff   2.5M Oct 30  2017 scenarios_ACTIVSg200.m

# -rw-rw-r--@ 1 art  staff   2.6M Sep 11  2019 ACTIVSg2000.EPC
# -rw-rw-r--@ 1 art  staff    10M Sep 11  2019 ACTIVSg2000.PWB
# -rw-rw-r--@ 1 art  staff   1.1M Sep 11  2019 ACTIVSg2000.RAW
# -rw-rw-r--@ 1 art  staff   6.0M Sep 11  2019 ACTIVSg2000.aux
# -rw-rw-r--@ 1 art  staff   306K Sep 11  2019 ACTIVSg2000.con
# -rw-rw-r--@ 1 art  staff   652K Sep 11  2019 ACTIVSg2000.m
# -rw-rw-r--@ 1 art  staff   3.1M Sep 11  2019 ACTIVSg2000.pwd
# -rw-rw-r--@ 1 art  staff   225M Sep 11  2019 ACTIVSg2000.tsb
# -rw-rw-r--@ 1 art  staff   220K Sep 11  2019 ACTIVSg2000_GIC_data.gic
# -rw-rw-r--@ 1 art  staff   670K Sep 11  2019 ACTIVSg2000_contingencies.aux
# -rw-rw-r--@ 1 art  staff   323K Sep 11  2019 ACTIVSg2000_dynamics.AUX
# -rw-rw-r--@ 1 art  staff   245K Sep 11  2019 ACTIVSg2000_dynamics.dyd
# -rw-rw-r--@ 1 art  staff   230K Sep 11  2019 ACTIVSg2000_dynamics.dyr
# -rw-rw-r--@ 1 art  staff    92M Sep 11  2019 ACTIVSg2000_load_time_series_MVAR.csv
# -rw-rw-r--@ 1 art  staff   102M Sep 11  2019 ACTIVSg2000_load_time_series_MW.csv
# -rw-rw-r--@ 1 art  staff   644K Sep 11  2019 case_ACTIVSg2000.m
# -rw-rw-r--@ 1 art  staff   151K Sep 11  2019 contab_ACTIVSg2000.m
# -rw-rw-r--@ 1 art  staff   3.5M Sep 11  2019 scenarios_ACTIVSg2000.m

set i 1
echo "Backup set $i"
echo =====================================

echo "Copy files to $workdir_path"
echo -------------------------------------
cp Static/Electric/ACTIVSg2000/ACTIVSg2000.PWB $workdir_name
cp Static/Electric/ACTIVSg2000/ACTIVSg2000.pwd $workdir_name
cp Static/Electric/ACTIVSg2000/ACTIVSg2000.m $workdir_name
# cp Static/Electric/ACTIVSg2000/*.csv $workdir_name

set tar_name "$workdir_name$i.tar"
if test -e $tar_name
    rm $tar_name
end
tar cfv $tar_name $workdir_folder

echo "Checking in $tar_name to $workdir_name"
echo -------------------------------------
dupver -ci -f $tar_name

set i 2
echo "Backup set $i"
echo =====================================

echo "Copy files to $workdir_path"
echo -------------------------------------
cp Static/Electric/ACTIVSg2000/ACTIVSg2000.RAW $workdir_name
cp ACTIVSg2000_Mod.RAW $workdir_name
# cp ACTIVSg2000_Copy.RAW $workdir_name
# cp Static/Electric/ACTIVSg2000/*.aux $workdir_name
# cp Static/Electric/ACTIVSg2000/contab_ACTIVSg2000.m $workdir_name
# cp Static/Electric/ACTIVSg2000/scenarios_ACTIVSg2000.m $workdir_name

set tar_name "$workdir_name$i.tar"
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
