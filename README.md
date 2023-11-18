**This feature has been merged into my [Widescreen Patcher](https://github.com/FaithBeam/Sims-1-Complete-Collection-Widescreen-Patcher) which is easier to use than this program. Check the Optional tab.**

# Sims 1 Cheat Unlocker

This program unlocks all cheats in The Sims 1. It was tested against the base Sims 1 and the Complete
Collection, so it probably works for the other expansions as long as you are using a NoCD/crack that has been
decompressed.

I was inspired to make this after finding a Russian application that does the same thing on ModTheSims and wanted to
recreate it.

![](https://i.imgur.com/x98PqGl.jpg)

There is a csv file in the download with explanations to what some commands do.

## Usage

**GUI**

[Video Guide](https://www.youtube.com/watch?v=eocyNaSjlUk)

1. Open sims1cheatunlocker.exe
2. Click Browse
3. Select your Sims.exe
4. Click Patch

**CLI**

1. Open cmd or PowerShell as administrator in the directory containing sims1cheatunlocker
2. Enter this command. Change the path to your Sims.exe if it is different:
    ```
    .\sims1cheatunlocker.exe -i "C:\Program Files (x86)\Maxis\The Sims\Sims.exe"
    ```
3. Run The Sims
4. Press ctrl + shift + c
5. Enter help 

## Uninstall

**Via sims1cheatunlocker.exe**

1. Open sims1cheatunlocker.exe
2. Click Browse
3. Select your Sims.exe
4. Click Uninstall

**Manual**

1. Delete your Sims.exe
2. Rename Sims.exe.BAK to Sims.exe

## Manually Compile

1. Install Golang 1.20.1 or newer
2. ```git clone https://github.com/FaithBeam/sims1cheatunlocker```
3. ```cd sims1cheatunlocker```
4. ```go build```
5. ```go get github.com/akavel/rsrc``` 
   1. Used for embedding the sims1cheatunlocker.exe.manifest into the exe
6. ```rsrc.exe -manifest .\sims1cheatunlocker.exe.manifest```
   
