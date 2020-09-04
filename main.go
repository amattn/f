package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	showV             bool
	showVersion       bool
	flagConfigPath    string
	shouldMakeFrcFile bool
	shouldPrint       bool
)

func init() {
	flag.BoolVar(&showV, "v", false, "show version info and exit(0)")
	flag.BoolVar(&showVersion, "version", false, "show version info and exit(0)")
	flag.StringVar(&flagConfigPath, "config", "", "path to config file, defaults to ~/.frc")
	flag.BoolVar(&shouldMakeFrcFile, "init", false, "Attempt to create config file.")
	flag.BoolVar(&shouldPrint, "print", false, "Print command to stdout, do not run in separate process")

	appendDefaultPaths()
}

func main() {
	flag.Parse()

	if showVersion || showV {
		fmt.Println(VersionInfo())
		os.Exit(0)
	}

	prependConfigPath(flagConfigPath)

	if shouldMakeFrcFile {
		err := createConfigFile()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err := findAndParseConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("\n", allTriplets)

	// get the non flag args:
	args := flag.Args()

	if len(args) == 0 {
		PrintMenu(allTriplets)
		return
	}

	primaryArg := args[0]

	// search for the command by shortcut:
	t := FindTripletByShortcut(allTriplets, primaryArg)
	if t != nil {
		t.RunCommand()
		return
	}

	// search for the command by number:
	t = FindTripletByNumber(allTriplets, primaryArg)
	if t != nil {
		t.RunCommand()
		return
	}

	// if not found,
	fmt.Println("No shortcut or number:" + primaryArg)
	PrintMenu(allTriplets)
}
