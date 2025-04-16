package main

import (
	"CloudlogAutoLogger/internal/agg_config_manager"
	"CloudlogAutoLogger/internal/agg_wsjtx"
	"bufio"
	"fmt"
	"os"
)

type listeners struct {
	wsjtx_thread *agg_wsjtx.Wsjtx
}

func main() {

	args := os.Args
	if len(args) < 2 {
		fmt.Print("Welcome to CloudLog Auto Logger\n")
		fmt.Print("   by the Alan Gordon Group\n")
		fmt.Print("\n")
		fmt.Print(" Enter one of the following commands:\n")
		fmt.Print("    S  for Set or update configuration\n")
		fmt.Print("    R  to start listening for UDP packets.\n")
		fmt.Print("    Q  to Quite and close program\n")
		fmt.Print("\n")
		fmt.Print(" To run this program without prompt add these command line options:\n")
		fmt.Print("    run  - to start listening for UDP packets per configuration settings\n")
		fmt.Print("    log  - to log activity\n")
		fmt.Print("\n")

		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("Enter command:")
			for scanner.Scan() {
				text := scanner.Text()
				if text == "Q" {
					goto ALLDONE
				}

				if text == "S" {

				}

				if text == "R" {
					run()
				}
			}
		}
	} else {
		// Check for command line options

	}

ALLDONE:
	os.Exit(0)

}

func run() {

	c := agg_config_manager.Get()
	if c.GetConfig() {

	}

}
