package main

import (
	"CloudlogAutoLogger/internal/agg_adi"
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

				if text == "T" {
					js8call_sample := "<call:5>ZZ0ZZ <gridsquare:0> <mode:4>MFSK <submode:3>JS8 <rst_sent:3>556 <rst_rcvd:3>556 <qso_date:8>20250416 <time_on:6>211335 <qso_date_off:8>20250416 <time_off:6>211335 <band:3>20m <freq:9>14.106279 <station_callsign:5>N7AKG <my_gridsquare:6>CN85SL <operator:5>N7AKG <eor>"
					js8_fields := agg_adi.Parse_adi(js8call_sample)
					js8adi := agg_adi.Encode_adi(js8_fields)
					fmt.Print(js8adi)

					varac_sample := "<command:3>Log<parameters:245><CALL:5>zz0zz <MODE:7>DYNAMIC <SUBMODE:7>VARA HF <RST_SENT:3>+01 <RST_RCVD:3>+02 <QSO_DATE:8>20250416 <TIME_ON:6>215130 <QSO_DATE_OFF:8>20250416 <TIME_OFF:6>215131 <BAND:3>20m <STATION_CALLSIGN:5>N7AKG <TX_PWR:0> <COMMENT:14>QSO with VarAC <EOR>"
					varac_fields := agg_adi.Parse_adi(varac_sample)
					varaadi := agg_adi.Encode_adi(varac_fields)
					fmt.Print(varaadi)

					wsjtx_sample := "<call:5>zz0zz <gridsquare:4>EM54 <mode:3>FT8 <rst_sent:0> <rst_rcvd:0> <qso_date:8>20250417 <time_on:6>043403 <qso_date_off:8>20250417 <time_off:6>043459 <band:3>20m <freq:9>14.075307 <station_callsign:5>N7AKG <my_gridsquare:6>CN85SL <tx_pwr:2>50 <comment:3>FT8 <operator:5>N7AKG <eor>"
					wsjtx_fields := agg_adi.Parse_adi(wsjtx_sample)
					wsjtxadi := agg_adi.Encode_adi(wsjtx_fields)
					fmt.Print(wsjtxadi)
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

			if cd.JS8Call_port != 0 {
				s := &Listeners{endFlag: false, threadFlag: true, verbose: true}
				s.Port = cd.JS8Call_port
				s.client_name = "JS8Call"
				listeners_list = append(listeners_list, s)
			}

			if cd.WSJTX_port != 0 {
				s := &Listeners{endFlag: false, threadFlag: true, verbose: true}
				s.Port = cd.WSJTX_port
				s.client_name = "WSJTX"
				listeners_list = append(listeners_list, s)
			}

			if cd.VARAC_port != 0 {
				s := &Listeners{endFlag: false, threadFlag: true, verbose: true}
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

			if cd.client_name == "JS8Call" {
				cd.js8call(pkt_stg)
			}

			if cd.client_name == "VARAC" {
				cd.varac(pkt_stg)
			}

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

func (cd *Listeners) js8call(pkt_stg string) {

	//Example UDP packet.
	//<call:5>ZZ0ZZ <gridsquare:0> <mode:4>MFSK <submode:3>JS8 <rst_sent:3>556 <rst_rcvd:3>556 <qso_date:8>20250416 <time_on:6>211335 <qso_date_off:8>20250416 <time_off:6>211335 <band:3>20m <freq:9>14.106279 <station_callsign:5>N7AKG <my_gridsquare:6>CN85SL <operator:5>N7AKG <eor>\x00\
	//	adi_fields := agg_adi.Parse_adi(pkt_stg)

}

func (cd *Listeners) varac(pkt_stg string) {
	//	Example:
	// <command:3>Log<parameters:245><CALL:5>zz0zz <MODE:7>DYNAMIC <SUBMODE:7>VARA HF <RST_SENT:3>+01 <RST_RCVD:3>+02 <QSO_DATE:8>20250416 <TIME_ON:6>215130 <QSO_DATE_OFF:8>20250416 <TIME_OFF:6>215131 <BAND:3>20m <STATION_CALLSIGN:5>N7AKG <TX_PWR:0> <COMMENT:14>QSO with VarAC <EOR>\x00\
	//	adi_fields := agg_adi.Parse_adi(pkt_stg)
}

func (cd *Listeners) wsjtx(pkt_stg string) {
	//	Example:
	// <call:5>zz0zz <gridsquare:4>EM54 <mode:3>FT8 <rst_sent:0> <rst_rcvd:0> <qso_date:8>20250417 <time_on:6>043403 <qso_date_off:8>20250417 <time_off:6>043459 <band:3>20m <freq:9>14.075307 <station_callsign:5>N7AKG <my_gridsquare:6>CN85SL <tx_pwr:2>50 <comment:3>FT8 <operator:5>N7AKG <eor>\x00\
	//	adi_fields := agg_adi.Parse_adi(pkt_stg)
}
