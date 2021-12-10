package main

import (
	"flag"
	"fmt"
	"github.com/FaithBeam/patternfinder-go"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	pattern = "00 56 75 04"
	applicationVersion = "1.0.0"
	applicationName = "sims1cheatunlocker"
)

func main() {
	fmt.Printf("%s %s\n\n", applicationName, applicationVersion)
	inputPtr := flag.String("i", "", "The path to the Sims.exe")
	flag.Parse()

	// Quit if the user didn't give an input
	if *inputPtr == "" {
		fmt.Println("Please supply these command line arguments:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Quit if the user gave a path to a file that isn't Sims.exe
	if !strings.EqualFold(filepath.Base(*inputPtr), "Sims.exe") {
		fmt.Println("Please input a Sims.exe")
		os.Exit(1)
	}

	patternByte := patternfinder.Transform(pattern)
	b, _ := os.ReadFile(*inputPtr)
	patternFound, offset := patternfinder.Find(b, patternByte)
	if patternFound {
		backupPath := *inputPtr + ".BAK"

		// Quit if a previous backup is detected
		if _, err := os.Stat(backupPath); err == nil {
			fmt.Printf("The backup %s already exists.", backupPath)
			os.Exit(1)
		}

		// Create a backup
		err := os.WriteFile(backupPath, b, 0666)
		if err != nil {
			log.Fatalln(err)
		}

		// NOP the jnz instruction
		b[offset + 2] = 144
		b[offset + 2 + 1] = 144

		// Delete the existing Sims.exe
		err = os.Remove(*inputPtr)
		if err != nil {
			log.Fatalln(err)
		}

		// Write the patched Sims.exe
		err = os.WriteFile(*inputPtr, b, 0666)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("%s was patched successfully. A backup was created at %s.", *inputPtr, backupPath)
	} else {
		fmt.Printf("Couldn't find pattern %s in %s.", pattern, *inputPtr)
		os.Exit(1)
	}
}