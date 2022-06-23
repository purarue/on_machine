package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/seanbreckenridge/on_machine"
	"os"
)

const (
	PRINT       = 1
	MATCH_PATHS = 2
)

type Command = int

type MatchConfig struct {
	base      string
	delimiter string
	json      bool
	skiplast  bool
}

type OnMachineConfig struct {
	pattern   string
	command   Command
	matchConf *MatchConfig
}

func parseFlags() (*OnMachineConfig, error) {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: on_machine [-h] [-cmd <print|match>] [OPTIONS] [PATTERN]

Tool to determine which operating system/machine you're on.

Commands:
print [default]: prints the resulting pattern after interpolating the pattern
match: does directory/path matching based on the pattern, changes the default pattern to '%o/%d/%h'

print
---
PATTERN is a printf-styled format string, supporting the following sequences:

%o - Operating System (using uname)
%d - Distro (using lsb_release)
%h - Hostname (name of the computer)
%a - Arch (detected by golang)
%O - Golang OS (unmodified golang detected operating system)

By default, this uses '%o_%d_%h'

match
---
Directory/path matching, Uses the pattern to match directory structures.
Can provide the base path to use with -base, that replaces '/' with
OS-specific path separator in the pattern. For more information, see the docs:
https://github.com/seanbreckenridge/on_machine

Options:
`)
		flag.PrintDefaults()
	}
	cmd := flag.String("cmd", "print", "on_machine command to run")
	base := flag.String("base", "", "Base directory to use to match paths")
	printJson := flag.Bool("json", false, "print results as a JSON array")
	delimiter := flag.String("delimiter", "\n", "delimiter to print between matches")
	// this is false by default because including a new line as the last delimiter
	// works better for processing lines in the shell
	skiplast := flag.Bool("skip-last-delim", false, "dont print the delimiter after the last match")
	nullchar := flag.Bool("print0", false, "use the null character as the delimiter")

	flag.Parse()
	var pattern string
	// parse command
	var command Command
	switch *cmd {
	case "print":
		command = PRINT
	case "match":
		command = MATCH_PATHS
	default:
		fmt.Printf("Unknown command '%s'. Provide either 'print' or 'match'\n", *cmd)
		os.Exit(1)
	}
	// set pattern
	switch flag.NArg() {
	case 1:
		pattern = flag.Arg(0)
	default:
		// set default pattern
		switch command {
		case PRINT:
			pattern = "%o_%d_%h"
		case MATCH_PATHS:
			pattern = "%o/%d/%h"
		}
	}
	// match based parsing

	var matchConfig *MatchConfig = nil
	if command == MATCH_PATHS {
		matchBase := string(*base)
		if matchBase != "" {
			if !on_machine.DirExists(matchBase) {
				fmt.Fprintf(os.Stderr, "Directory doesnt exist: '%s'\n", matchBase)
				os.Exit(1)
			}
		}
		// handle delimiter flag
		delim := *delimiter
		if *nullchar {
			delim = "\000"
		}
		matchConfig = &MatchConfig{
			base:      matchBase,
			delimiter: delim,
			skiplast:  *skiplast,
			json:      *printJson,
		}
	}
	return &OnMachineConfig{
		pattern:   pattern,
		command:   command,
		matchConf: matchConfig,
	}, nil
}

func run() error {
	conf, err := parseFlags()
	if err != nil {
		return err
	}
	switch conf.command {
	case PRINT:
		res := on_machine.ReplaceFields(conf.pattern)
		fmt.Println(res)
	case MATCH_PATHS:
		matched, _ := on_machine.MatchPaths(conf.pattern, conf.matchConf.base)
		// print to STDOUT
		if conf.matchConf.json {
			jsonBytes, err := json.Marshal(matched)
			if err != nil {
				return err
			}
			fmt.Print(string(jsonBytes))
		} else {
			for i, p := range matched {
				fmt.Print(p)
				if i != len(matched)-1 {
					fmt.Print(conf.matchConf.delimiter)
				} else {
					if !conf.matchConf.skiplast {
						fmt.Print(conf.matchConf.delimiter)
					}
				}
			}
		}
	}
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
