package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
)

const DefaultConfigFileName = ".frc"
const CommentCharacter = "#"

var configPaths []string
var allTriplets []Triplet

func init() {
	configPaths = make([]string, 0, 3)
	allTriplets = make([]Triplet, 0, 5)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "https://github.com/amattn/f\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Typical usage is 'f <Num>' or 'f <Shortcut>'.  Just 'f' will print out a menu\n")
		flag.PrintDefaults()
	}
}

func defaultConfigPaths() []string {
	return configPaths
}
func preferedPath() string {
	if len(configPaths) <= 0 {
		panic("call to preferedPath before any paths set")
	}
	return configPaths[0]
}

func prependConfigPath(cfgPth string) {
	if len(cfgPth) == 0 {
		return
	}

	configPaths = append(configPaths, "")
	copy(configPaths[1:], configPaths[0:])
	configPaths[0] = cfgPth
}

func appendConfigPath(cfgPth string) {
	if len(cfgPth) == 0 {
		return
	}
	configPaths = append(configPaths, cfgPth)
}

func joinFrcToDir(dirPath string) string {
	return path.Join(dirPath, DefaultConfigFileName)
}

func appendDefaultPaths() {
	// home directory
	usr, _ := user.Current()
	if usr != nil {
		appendConfigPath(joinFrcToDir(usr.HomeDir))
	}

	// current working directory...
	wd, _ := os.Getwd()
	if wd != "" {
		appendConfigPath(joinFrcToDir(wd))
	}
}

func findConfigFile(pathsToCheck []string) string {
	for _, p := range pathsToCheck {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			// file does not exist, check the next one
		} else {
			// file exists, return it
			return p
		}
	}

	// not found :(
	return ""
}

// convenience, one call load...
func findAndParseConfigFile() error {
	configPath := findConfigFile(configPaths)
	if configPath == "" {
		errStr := "\nCould not find config file in any of the following paths:\n"
		for _, p := range configPaths {
			errStr += fmt.Sprintf("  %s\n", p)
		}
		errStr += "\nrun with the --init flag to create a config file.\n"
		return errors.New(errStr)
	}

	return parseConfig(configPath)
}

func createConfigFile() error {
	configPath := preferedPath()

	content := `# Default .frc content
# https://github.com/amattn/f
# Typical usage is 'f <Num>' or 'f <Shortcut>'.  Just 'f' will print out a menu
# 
# Any characters after a # are considered comments
# The basic format is <Shortcut><WHITESPACE><Command>
# Examples:
#
# db     cd ~/Dropbox                     # requires Dropbox installed
# rmdot  rm -rF .                         # Danger!
# psapps ps -ax | grep .app/
# openf  open https://github.com/amattn/f # Mac OS X only command

ll ls -la
`

	err := ioutil.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		return err
	}

	// assume it worked:
	fmt.Println("Config file written to:", configPath)
	return nil
}

// ######
// #     #   ##   #####   ####  ######
// #     #  #  #  #    # #      #
// ######  #    # #    #  ####  #####
// #       ###### #####       # #
// #       #    # #   #  #    # #
// #       #    # #    #  ####  ######
//

// our config format is simple.
// # is a comment... ignore anything after that until end of line
// for each line:
//   - trim all right and left white space
//   - split by whitespace
//   - lineCompoenent[0] is the shortcut
//   - join.(lineCompoenent[1:], " ") is the command
func parseConfig(pathToFile string) error {
	f, err := os.Open(pathToFile)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if err != nil {
		return err
	}

	i := 0
	for scanner.Scan() {
		i++
		trip := parseLine(i, scanner.Text())
		if trip.IsValid() {
			allTriplets = append(allTriplets, trip)
		}
	}

	if i == 0 {
		return errors.New("Empty config file")
	}

	return nil
}

func parseLine(i int, line string) Triplet {
	cleanedPair, comment := cleanLine(i, line)
	shortcut, command := parsePair(i, cleanedPair)

	return Triplet{
		Shortcut: shortcut,
		Command:  strings.TrimSpace(command),
		Comment:  strings.TrimSpace(comment),
	}
}

func parsePair(i int, pair string) (string, string) {
	pairComponents := strings.Fields(pair)
	switch len(pairComponents) {
	case 0:
		return "", ""
	case 1:
		return pairComponents[0], ""
	default:
		return pairComponents[0], strings.Join(pairComponents[1:], " ")
	}
}

// trim and strip comments
// return before and after first CommentCharacter
func cleanLine(i int, line string) (string, string) {
	trimmed := strings.TrimSpace(line)
	// strip comments
	strippedComponents := strings.SplitN(trimmed, CommentCharacter, 2)

	switch len(strippedComponents) {
	case 0:
		return "", ""
	case 1:
		return strippedComponents[0], ""
	default:
		return strippedComponents[0], strippedComponents[1]
	}
}
