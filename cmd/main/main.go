package main

import (
	"CloudlogAutoLogger/internal/agg_config_manager"
	"CloudlogAutoLogger/internal/agg_logger"
	"CloudlogAutoLogger/internal/agg_udp"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Listeners struct {
	Cloudlog_api_key   string
	Station_profile_id string
	Port               int

	client_name string // Name of broadcasting app (e.g. WSJTX, JS8CALL, VARAC)

	// thread control
	verbose    bool
	endFlag    bool
	threadFlag bool
}

var listeners_list = []*Listeners{}

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
				fmt.Print("\nEnter command (S,R or Q):")
			}
		}
	} else {
		// Check for command line options

	}

ALLDONE:
	stop()
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

	text := GetUserPromptText("Cloud log API key:")
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

	text = GetUserPromptText("JS8Call port (current value=" + strconv.Itoa(cd.JS8Call_port) + "):")
	if len(text) > 0 {
		cd.JS8Call_port, _ = strconv.Atoi(text)
	}

	text = GetUserPromptText("VARAC port (current value=" + strconv.Itoa(cd.VARAC_port) + "):")
	if len(text) > 0 {
		cd.VARAC_port, _ = strconv.Atoi(text)
	}

	agg_config_manager.SaveConfig(cd)
}

func run() {

	listeners_list = nil
	cd, stat := agg_config_manager.GetConfig()
	if stat {
		if cd.Cloudlog_api_key != "" {
			s := &Listeners{endFlag: false, threadFlag: true, verbose: true}

			if cd.JS8Call_port != 0 {
				s.Port = cd.JS8Call_port
				s.client_name = "JS8Call"
				listeners_list = append(listeners_list, s)
			}

			if cd.WSJTX_port != 0 {
				s.Port = cd.WSJTX_port
				s.client_name = "WSJTX"
				listeners_list = append(listeners_list, s)
			}

			if cd.VARAC_port != 0 {
				s.Port = cd.VARAC_port
				s.client_name = "VARAC"
				listeners_list = append(listeners_list, s)
			}
		}

		// Start threads
		for _, s := range listeners_list {
			go s.Start()
		}
	}

}

func stop() {
	// Shut down threads
	for _, s := range listeners_list {
		s.Stop()
	}
	listeners_list = nil
}

func (cd *Listeners) Start() {
	var portstg = strconv.Itoa(cd.Port)
	agg_logger.Get().Log("***** Begin listener on port:", portstg)
	for !cd.endFlag {
		stat, udp_pkt := agg_udp.WaitOnUDP(cd.Port, 250, false)
		if stat {
			agg_logger.Get().Log("***** Begin packet ***** Port:", portstg)

			pkt_stg := string(udp_pkt)
			agg_logger.Get().Log(pkt_stg, "")

			//flex_util.BinaryDump(udp_pkt)

			//pkt_stg := string(udp_pkt[:])
			// flds := strings.FieldsFunc(pkt_stg, split)
			// for _, v := range flds {
			// 	agg_logger.Get().Log(portStg, " "+v)
			// }
			agg_logger.Get().Log("***** End packet ***** Port:", portstg)
		}
		//time.Sleep(time.Millisecond)
	}
	cd.threadFlag = false
}

func (cd *Listeners) Stop() {
	cd.endFlag = true
	sleepTime := 2 * time.Second

	for {
		if !cd.threadFlag {
			return
		}

		time.Sleep(sleepTime)
	}
}
