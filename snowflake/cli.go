package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/emakeev/snowflake"
)

var (
	file           string
	make, forceNew bool
)

func init() {
	const (
		fileUsage = "Use specified file to store id"
		makeUsage = "Make a snowflake ID if one doesn't exist. If one already does, this does nothing, so it's always safe to use."
	)
	flag.StringVar(&file, "f", snowflake.DefaultSnowflakeFile, fileUsage)
	flag.StringVar(&file, "file", snowflake.DefaultSnowflakeFile, fileUsage+" (shorthand)")
	flag.BoolVar(&make, "make-snowflake", false, makeUsage+" (shorthand)")
	flag.BoolVar(&make, "m", false, makeUsage+" (shorthand)")
	flag.BoolVar(&forceNew, "force-new-key", false, "Force generation of new ID. WARNING: Deletes existing ID")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr,
			"\nGenerate IDs for a system that are as unique as a snowflake (they're just persistent UUIDs).\n"+
				"Run with no arguments, simply prints the system's snowflake ID.\n"+
				"The snowflake ID should be generated for the system when this package is installed.\n"+
				"You can use this tool to generate an ID if for some reason one doesn't exist.\n")
		fmt.Fprintf(os.Stderr, "\nUsage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	if forceNew {
		if err := snowflake.WriteNew(file); err != nil {
			fmt.Printf("Could not create a new snowflake ID: %v\n", err)
			os.Exit(1)
		}
	}
	if make {
		if _, err := snowflake.Make(file); err != nil {
			fmt.Printf("Could not Make a snowflake ID: %v\n", err)
			os.Exit(2)
		}
	}
	u, err := snowflake.Get(file)
	if err != nil {
		fmt.Printf("Couldn't find a snowflake ID in '%s': %v.\nTry running with the -m option to make one\n", file, err)
		os.Exit(3)
	}
	fmt.Print(u.Encode(), "\n")
}
