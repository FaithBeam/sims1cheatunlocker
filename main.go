package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/FaithBeam/patternfinder-go"
	"github.com/gonutz/wui/v2"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	pattern            = "00 56 75 04"
	applicationVersion = "2.0.0"
	applicationName    = "sims1cheatunlocker"
	applicationTitle   = "Sims 1 Cheat Unlocker"
)

func main() {
	// User supplied cli arguments, do cli mode. If not, gui mode
	if len(os.Args[1:]) > 0 {
		fmt.Printf("%s %s\n\n", applicationName, applicationVersion)
		inputPtr := flag.String("i", "", "The path to the Sims.exe")
		flag.Parse()

		// Quit if the user didn't give an input
		if *inputPtr == "" {
			flag.PrintDefaults()
			log.Fatalln("please supply these command line arguments")
		}

		// Quit if the user gave a path to a file that isn't Sims.exe
		if !strings.EqualFold(filepath.Base(*inputPtr), "Sims.exe") {
			log.Fatalln("Please input a Sims.exe")
		}

		err := doPatch(*inputPtr)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		w := wui.NewWindow()
		w.SetWidth(400)
		w.SetHeight(110)
		w.SetResizable(false)
		w.SetTitle(applicationTitle)

		el := wui.NewEditLine()
		browseBtn := wui.NewButton()
		uninstallBtn := wui.NewButton()
		patchBtn := wui.NewButton()

		el.SetWidth(300)
		el.SetHeight(25)
		el.SetPosition(5, 5)

		w.Add(el)

		browseBtn.SetText("Browse")
		browseBtn.SetWidth(70)
		browseBtn.SetHeight(25)
		browseBtn.SetPosition(el.X()+el.Width()+5, el.Y())
		browseBtn.SetOnClick(func() {
			dlg := wui.NewFileOpenDialog()
			dlg.AddFilter("Sims Executable (Sims.exe)", "exe")
			ok, path := dlg.ExecuteSingleSelection(w)
			el.SetText(path)

			if ok && strings.EqualFold(filepath.Base(path), "Sims.exe") {
				found, _ := getPatternOffset(path)
				if found {
					patchBtn.SetEnabled(true)
				}
			}

			backOk, _ := backupExists(path)
			if backOk {
				uninstallBtn.SetEnabled(true)
			}
		})

		w.Add(browseBtn)

		uninstallBtn.SetText("Uninstall")
		uninstallBtn.SetEnabled(false)
		uninstallBtn.SetBounds(browseBtn.X(), browseBtn.Y()+browseBtn.Height()+5, browseBtn.Width(), browseBtn.Height())
		uninstallBtn.SetOnClick(func() {
			ok, bakPath := backupExists(el.Text())
			if ok {
				if err := doUninstall(el.Text(), bakPath); err != nil {
					wui.MessageBoxError(applicationTitle, err.Error())
				} else {
					wui.MessageBoxInfo(applicationTitle, "Uninstalled!")
					uninstallBtn.SetEnabled(false)

					found, _ := getPatternOffset(el.Text())
					if found {
						patchBtn.SetEnabled(true)
					}
				}
			} else {
				wui.MessageBoxInfo(applicationTitle, "Sims backup not found")
			}
		})

		w.Add(uninstallBtn)

		patchBtn.SetText("Patch")
		patchBtn.SetEnabled(false)
		patchBtn.SetBounds(uninstallBtn.X()-uninstallBtn.Width()-5, uninstallBtn.Y(), uninstallBtn.Width(), uninstallBtn.Height())
		patchBtn.SetOnClick(func() {
			err := doPatch(el.Text())
			if err != nil {
				wui.MessageBoxError(applicationTitle, err.Error())
			} else {
				wui.MessageBoxInfo(applicationTitle, "Patched!")

				patchBtn.SetEnabled(false)

				ok, _ := backupExists(el.Text())
				if ok {
					uninstallBtn.SetEnabled(true)
				}
			}
		})

		w.Add(patchBtn)

		if err := w.Show(); err != nil {
			log.Fatalln(err)
		}
	}
}

func getPatternOffset(input string) (bool, int) {
	patternByte := patternfinder.Transform(pattern)
	b, _ := os.ReadFile(input)
	patternFound, offset := patternfinder.Find(b, patternByte)
	return patternFound, offset
}

func backupExists(simsExePath string) (bool, string) {
	bakPath := simsExePath + ".BAK"
	f, err := os.OpenFile(bakPath, os.O_RDONLY, 0666)
	if errors.Is(err, os.ErrNotExist) {
		return false, ""
	}
	if err := f.Close(); err != nil {
		log.Fatalln(err)
	}

	return true, bakPath
}

func doUninstall(simsExePath string, bakPath string) error {
	if err := os.Remove(simsExePath); err != nil {
		return err
	}

	if err := os.Rename(bakPath, simsExePath); err != nil {
		return err
	}

	return nil
}

func doPatch(input string) error {
	patternByte := patternfinder.Transform(pattern)
	b, _ := os.ReadFile(input)
	patternFound, offset := patternfinder.Find(b, patternByte)
	if patternFound {
		backupPath := input + ".BAK"

		// Quit if a previous backup is detected
		if _, err := os.Stat(backupPath); err == nil {
			return errors.New(fmt.Sprintf("the backup %s already exists", backupPath))
		}

		// Create a backup
		if err := os.Rename(input, backupPath); err != nil {
			return err
		}

		// NOP the jnz instruction
		b[offset+2] = 144
		b[offset+2+1] = 144

		// Write the patched Sims.exe
		if err := os.WriteFile(input, b, 0666); err != nil {
			return errors.Join(errors.New("error writing the patched sims.exe"), err)
		}
	} else {
		return errors.New(fmt.Sprintf("couldn't find pattern %s in %s", pattern, input))
	}

	return nil
}
