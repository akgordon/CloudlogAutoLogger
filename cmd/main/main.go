package main

import (
	agg_config_manager "CloudlogAutoLogger/internal/agg_config_manager"
	"CloudlogAutoLogger/internal/agg_listeners"
	agg_logger "CloudlogAutoLogger/internal/agg_logger"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var listeners_list = []*agg_listeners.Listeners{}

func main() {

	args := os.Args
	if len(args) < 2 {
		fmt.Print("Welcome to CloudLog Auto Logger\n")
		fmt.Print("   by the Alan Gordon Group\n")
		fmt.Print("           v1.1\n")
		fmt.Print("\n")
		fmt.Print(" Enter one of the following commands:\n")
		fmt.Print("    S  for Set or update configuration\n")
		fmt.Print("    R  to start listening for UDP packets.\n")
		fmt.Print("    Q  to Quite and close program\n")
		fmt.Print("\n")
		fmt.Print(" To run this program without prompt add these command line options:\n")
		fmt.Print("    run  - to start listening for UDP packets per configuration settings\n")
		fmt.Print("    log  - to log activity. Omit for no logging.\n")
		fmt.Print("    Example: cloudlogautologger.exe run log\n")
		fmt.Print("\n")

		// Start Logger if in interactive mode
		agg_logger.Get().Open("CloundlogAutoLogger.log")

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

				// Some undocumented testers
				if text == "TJ" {
					agg_listeners.TestListener("JS8CALL")
				}
				if text == "TV" {
					agg_listeners.TestListener("VARAC")
				}
				if text == "TW" {
					agg_listeners.TestListener("WSJTX")
				}

				fmt.Print("\nEnter command (S,R or Q):")
			}
		}
	} else {
		autoRun()
	}

ALLDONE:
	stop()
	agg_logger.Get().Close()
	os.Exit(0)

}

func GetUserPromptText(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func set_config() {
	var cd, _ = agg_config_manager.GetConfig()

	fmt.Printf("Enter new values or leave blank to keep current value\n")
	fmt.Printf("API key is from Cloudlog account\n")
	fmt.Printf("Station profile id - This can be found when editing a station profile its a number and displayed in the URL string.\n")
	fmt.Printf("Leave port number as 0 to not enable that listener\n")
	fmt.Printf("\n")

	text := GetUserPromptText("Cloudlog URL  (current value=" + cd.Cloudlog_url + "):")
	if len(text) > 0 {
		cd.Cloudlog_url = text
	}

	text = GetUserPromptText("Cloud log API key:")
	if len(text) > 0 {
		cd.Cloudlog_api_key = text
	}

	text = GetUserPromptText("Station profile ID  (current value=" + cd.Station_profile_id + "):")
	if len(text) > 0 {
		cd.Station_profile_id = text
	}

	text = GetUserPromptText("WSJTX port (current value=" + strconv.Itoa(cd.WSJTX_port) + "):")
	if len(text) > 0 {
		cd.WSJTX_port, _ = strconv.Atoi(text)
	}

	text = GetUserPromptText("JS8CALL port (current value=" + strconv.Itoa(cd.JS8CALL_port) + "):")
	if len(text) > 0 {
		cd.JS8CALL_port, _ = strconv.Atoi(text)
	}

	text = GetUserPromptText("VARAC port (current value=" + strconv.Itoa(cd.VARAC_port) + "):")
	if len(text) > 0 {
		cd.VARAC_port, _ = strconv.Atoi(text)
	}

	agg_config_manager.SaveConfig(cd)
}

func run() {

	listeners_list = nil

	{
		l, stat := agg_listeners.BuildListener("JS8CALL")
		if stat {
			listeners_list = append(listeners_list, &l)
		}
	}

	{
		l, stat := agg_listeners.BuildListener("VARAC")
		if stat {
			listeners_list = append(listeners_list, &l)
		}
	}

	{
		l, stat := agg_listeners.BuildListener("WSJTX")
		if stat {
			listeners_list = append(listeners_list, &l)
		}
	}

	// Start threads
	for _, s := range listeners_list {
		go s.Start()
	}

}

func stop() {
	// Shut down threads
	for _, s := range listeners_list {
		s.Stop()
	}
	listeners_list = nil
}

func autoRun() {
	args := os.Args

	// Check for command line options
	doRun := false
	doLog := false

	for i := 1; i < len(args); i++ {
		if args[i] == "run" {
			doRun = true
		}

		if args[i] == "log" {
			doLog = true
		}
	}

	if doLog {
		agg_logger.Get().Open("CloundlogAutoLogger.log")
	}

	if doRun {
		run()
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n\nPress ENTER to exit\n\n")
		for scanner.Scan() {
			scanner.Text()
			stop()
			agg_logger.Get().Close()
			os.Exit(0)
		}
	}
}
