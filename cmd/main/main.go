package main

import (
	"CloudlogAutoLogger/internal/agg_config_manager"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type listeners struct {
	Cloudlog_api_key   string
	Station_profile_id string
	Port               int

	client_name string // Name of broadcasting app (e.g. WSJTX, JS8CALL, VARAC)

	// thread control
	verbose    bool
	endFlag    bool
	threadFlag bool
}

var listeners_list = []*listeners{}

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
			fmt.Print("\nEnter command (S,R or Q):")
			for scanner.Scan() {
				text := scanner.Text()
				if text == "Q" {
					goto ALLDONE
				}

				if text == "S" {
					set_config()
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
	stop()
	os.Exit(0)

}

func set_config() {
	var cd = agg_config_manager.GetConfig()

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Enter new values or leave blank to keep current value\n")
	fmt.Printf("API key is from Cloudlog account\n")
	fmt.Printf("Station profile id - This can be found when editing a station profile its a number and displayed in the URL string.\n")
	fmt.Printf("Leave port number as 0 to not enable that listener\n")
	fmt.Printf("\n")

	fmt.Printf("Cloud log API key:")
	text := scanner.Text()
	if len(text) > 0 {
		cd.Cloudlog_api_key = text
	}

	fmt.Printf("Station profile ID:")
	text = scanner.Text()
	if len(text) > 0 {
		cd.Station_profile_id = text
	}

	fmt.Printf("WSJTX port (current value=" + strconv.Itoa(cd.WSJTX_port) + "):")
	text = scanner.Text()
	if len(text) > 0 {
		cd.WSJTX_port, _ = strconv.Atoi(text)
	}

	fmt.Printf("JS8Call port (current value=" + strconv.Itoa(cd.JS8Call_port) + "):")
	text = scanner.Text()
	if len(text) > 0 {
		cd.JS8Call_port, _ = strconv.Atoi(text)
	}

	fmt.Printf("VARAC port (current value=" + strconv.Itoa(cd.VARAC_port) + "):")
	text = scanner.Text()
	if len(text) > 0 {
		cd.VARAC_port, _ = strconv.Atoi(text)
	}

	agg_config_manager.SaveConfig(cd)
}

func run() {

	var cd = agg_config_manager.GetConfig()
	if cd.Cloudlog_api_key != "" {
		s := &listeners{endFlag: false, threadFlag: true, Port: 0}
		listeners_list = append(listeners_list, s)
	}

	// Start threads
	for _, s := range listeners_list {
		s.start()
	}

}

func stop() {
	// Shut down threads
	for _, s := range listeners_list {
		s.stop()
	}
	listeners_list = nil
}
