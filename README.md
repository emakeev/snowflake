# Golang implementation of  Snowflake: Simple Persistent System-wide Unique IDs

Snowflake is a simple Go implementation of API for creation, management and usage of per-machine UUIDs (called a 'snowflake'). Usage is simple:

		import (
			"fmt"
			"github.com/emakeev/snowflake"
		)
    
		if err := snowflake.WriteNew(file); err != nil {
			fmt.Printf("Could not create a new snowflake ID: %v\n", err)
			return ...
		}
		...
		u, err := snowflake.Get()
		if err != nil {
			fmt.Printf("Couldn't find a snowflake ID in '%s': %v.\nTry running with the -m option to make one\n", file, err)
			return ...
		}
		fmt.Print(u.Encode(), "\n")  // would print something similar to: 7232c1c3-f6d1-4aec-bedd-c7e4c10dc8d3
		...
    
There is also a CLI tool that can be run from the command line:

$ snowflake
008b6f86-0f90-4d71-5814-b73eeb87c1ac

Use -h flag to see all available options:

$ snowflake -h

Generate IDs for a system that are as unique as a snowflake (they're just persistent UUIDs).
Run with no arguments, simply prints the system's snowflake ID.
The snowflake ID should be generated for the system when this package is installed.
You can use this tool to generate an ID if for some reason one doesn't exist.

Usage of snowflake:
  -f string
    	Use specified file to store id (default "/etc/snowflake")
  -file string
    	Use specified file to store id (shorthand) (default "/etc/snowflake")
  -force-new-key
    	Force generation of new ID. WARNING: Deletes existing ID
  -m	Make a snowflake ID if one doesn't exist. If one already does, this does nothing, so it's always safe to use. (shorthand)
  -make-snowflake
    	Make a snowflake ID if one doesn't exist. If one already does, this does nothing, so it's always safe to use. (shorthand)
 
 
For more, see original snoflake python documentation at: https://github.com/shaddi/snowflake/blob/master/README.md
