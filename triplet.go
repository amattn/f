package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"text/tabwriter"
)

type Triplet struct {
	Shortcut string
	Command  string
	Comment  string
}

func (t Triplet) IsValid() bool {
	return (t.Command != "")
}

func (t Triplet) IsEqual(other Triplet) bool {
	if t.Shortcut == other.Shortcut && t.Command == other.Command && t.Comment == other.Comment {
		return true
	}
	return false
}

func (t Triplet) String() string {
	if len(t.Comment) == 0 {
		return fmt.Sprintln("shortcut:", t.Shortcut, " command:", t.Command)
	}
	return fmt.Sprintln("shortcut:", t.Shortcut, " command:", t.Command, "#", t.Comment)
}

func (t Triplet) RunCommand() {
	commandComponents := strings.Fields(t.Command)
	if shouldPrint {
		fmt.Printf("%s", strings.Join(commandComponents, " "))
	} else {
		// shCommandComponents := append([]string{"sh"}, commandComponents...)
		cmd := exec.Command(commandComponents[0], commandComponents[1:]...)
		log.Println(cmd.Args)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Println(err)
			if exitErr, ok := err.(*exec.ExitError); ok {
				if waitStatus, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					os.Exit(waitStatus.ExitStatus())
				}
			} else {
				os.Exit(1)
			}
		}
	}
}

func FindTripletByNumber(triplets []Triplet, arg string) *Triplet {
	n, err := strconv.ParseUint(arg, 10, 32)
	num := int(n)
	if err != nil {
		return nil
	}
	if num <= 0 {
		return nil
	}
	if num > len(triplets) {
		return nil
	}

	index := num - 1

	return &(triplets[index])
}

func FindTripletByShortcut(triplets []Triplet, arg string) *Triplet {
	for _, t := range triplets {
		if t.Shortcut == arg {
			return &t
		}
	}
	return nil
}

func PrintMenu(triplets []Triplet) {
	fmt.Println("Enter the shortcut or number to run a command\n")
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 4, 0, 1, ' ', 0)
	fmt.Fprintf(w, "Num\tShort\tCommand\t# Comment\t\n")
	for i, t := range triplets {
		fmt.Fprintf(w, "%d\t%s\t%s\t# %s\t\n", i+1, t.Shortcut, t.Command, t.Comment)
	}
	fmt.Fprintln(w)
	w.Flush()
}
