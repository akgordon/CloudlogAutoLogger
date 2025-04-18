package agg_listeners

import (
	"CloudlogAutoLogger/internal/agg_adi"
	"CloudlogAutoLogger/internal/agg_config_manager"
	"CloudlogAutoLogger/internal/agg_logger"
	"CloudlogAutoLogger/internal/agg_tcp"
	"CloudlogAutoLogger/internal/agg_udp"
	"strconv"
	"time"
)

type Listeners struct {
	Cloudlog_url       string
	Cloudlog_api_key   string
	Station_profile_id string
	Port               int

	client_name string // Name of broadcasting app (e.g. WSJTX, JS8CALL, VARAC)

	// thread control
	verbose    bool
	endFlag    bool
	threadFlag bool
}

func BuildListener(mode string) (Listeners, bool) {
	cd, stat := agg_config_manager.GetConfig()
	var rtnListener = Listeners{
		Cloudlog_url:       cd.Cloudlog_url,
		Cloudlog_api_key:   cd.Cloudlog_api_key,
		Station_profile_id: cd.Station_profile_id,
		endFlag:            false,
		threadFlag:         true,
		verbose:            true}

	if stat {
		if cd.Cloudlog_api_key != "" {
			if mode == "JS8CALL" {
				if cd.JS8CALL_port != 0 {
					rtnListener.Port = cd.JS8CALL_port
					rtnListener.client_name = "JS8CALL"
					return rtnListener, true
				}
			}

			if mode == "VARAC" {
				if cd.JS8CALL_port != 0 {
					rtnListener.Port = cd.VARAC_port
					rtnListener.client_name = "VARAC"
					return rtnListener, true
				}
			}

			if mode == "WSJTX" {
				if cd.JS8CALL_port != 0 {
					rtnListener.Port = cd.WSJTX_port
					rtnListener.client_name = "WSJTX"
					return rtnListener, true
				}
			}
		}
	}

	return rtnListener, false
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

			if cd.client_name == "JS8CALL" {
				cd.Js8call(pkt_stg)
			}

			if cd.client_name == "VARAC" {
				cd.VARAC(pkt_stg)
			}

			if cd.client_name == "WSJTX" {
				cd.WSJTX(pkt_stg)
			}

			agg_logger.Get().Log("***** End packet ***** Port:", portstg)
		}
		//time.Sleep(time.Millisecond)
	}
	cd.threadFlag = false
	agg_logger.Get().Log("***** Listener stopped on port:", portstg)
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

func TestListener(mode string) {

	if mode == "JS8CALL" {
		l, stat := BuildListener("JS8CALL")
		if stat {
			js8call_sample := "<call:5>ZZ0ZZ <gridsquare:0> <mode:4>MFSK <submode:3>JS8 <rst_sent:3>556 <rst_rcvd:3>556 <qso_date:8>20250416 <time_on:6>211335 <qso_date_off:8>20250416 <time_off:6>211335 <band:3>20m <freq:9>14.106279 <station_callsign:5>N7AKG <my_gridsquare:6>CN85SL <operator:5>N7AKG <eor>"
			l.Js8call(js8call_sample)
		}
	}

	if mode == "VARAC" {
		l, stat := BuildListener("VARAC")
		if stat {
			VARAC_sample := "<command:3>Log<parameters:245><CALL:5>zz0zz <MODE:7>DYNAMIC <SUBMODE:7>VARA HF <RST_SENT:3>+01 <RST_RCVD:3>+02 <QSO_DATE:8>20250416 <TIME_ON:6>215130 <QSO_DATE_OFF:8>20250416 <TIME_OFF:6>215131 <BAND:3>20m <STATION_CALLSIGN:5>N7AKG <TX_PWR:0> <COMMENT:14>QSO with VARAC <EOR>"
			l.VARAC(VARAC_sample)
		}
	}

	if mode == "WSJTX" {
		l, stat := BuildListener("WSJTX")
		if stat {
			WSJTX_sample := "<call:5>zz0zz <gridsquare:4>EM54 <mode:3>FT8 <rst_sent:0> <rst_rcvd:0> <qso_date:8>20250417 <time_on:6>043403 <qso_date_off:8>20250417 <time_off:6>043459 <band:3>20m <freq:9>14.075307 <station_callsign:5>N7AKG <my_gridsquare:6>CN85SL <tx_pwr:2>50 <comment:3>FT8 <operator:5>N7AKG <eor>"
			l.WSJTX(WSJTX_sample)
		}
	}
}

func (cd *Listeners) Js8call(pkt_stg string) {

	//Example UDP packet.
	//<call:5>ZZ0ZZ <gridsquare:0> <mode:4>MFSK <submode:3>JS8 <rst_sent:3>556 <rst_rcvd:3>556 <qso_date:8>20250416 <time_on:6>211335 <qso_date_off:8>20250416 <time_off:6>211335 <band:3>20m <freq:9>14.106279 <station_callsign:5>N7AKG <my_gridsquare:6>CN85SL <operator:5>N7AKG <eor>\x00\

	// Parse the ADI to clean it up then re-strigify it
	adi_fields := agg_adi.ParseADIRecord(pkt_stg)
	adi_stg := agg_adi.Encode_adi(adi_fields)
	send2Cloudlog(cd, adi_stg)

}

func (cd *Listeners) VARAC(pkt_stg string) {
	//	Example:
	// <command:3>Log<parameters:245><CALL:5>zz0zz <MODE:7>DYNAMIC <SUBMODE:7>VARA HF <RST_SENT:3>+01 <RST_RCVD:3>+02 <QSO_DATE:8>20250416 <TIME_ON:6>215130 <QSO_DATE_OFF:8>20250416 <TIME_OFF:6>215131 <BAND:3>20m <STATION_CALLSIGN:5>N7AKG <TX_PWR:0> <COMMENT:14>QSO with VARAC <EOR>\x00\
	//	adi_fields := agg_adi.ParseADIRecord(pkt_stg)

	// Parse the ADI to clean it up then re-strigify it
	adi_fields := agg_adi.ParseADIRecord(pkt_stg)
	adi_stg := agg_adi.Encode_adi(adi_fields)
	send2Cloudlog(cd, adi_stg)
}

func (cd *Listeners) WSJTX(pkt_stg string) {
	//	Example:
	// <call:5>zz0zz <gridsquare:4>EM54 <mode:3>FT8 <rst_sent:0> <rst_rcvd:0> <qso_date:8>20250417 <time_on:6>043403 <qso_date_off:8>20250417 <time_off:6>043459 <band:3>20m <freq:9>14.075307 <station_callsign:5>N7AKG <my_gridsquare:6>CN85SL <tx_pwr:2>50 <comment:3>FT8 <operator:5>N7AKG <eor>\x00\

	// Parse the ADI to clean it up then re-strigify it
	adi_fields := agg_adi.ParseADIRecord(pkt_stg)
	adi_stg := agg_adi.Encode_adi(adi_fields)
	send2Cloudlog(cd, adi_stg)
}

func send2Cloudlog(cp *Listeners, pkt_stg string) bool {

	_payload := make(map[string]string)
	_payload["key"] = cp.Cloudlog_api_key
	_payload["station_profile_id"] = cp.Station_profile_id
	_payload["type"] = "adif"
	_payload["string"] = pkt_stg

	url := cp.Cloudlog_url + "index.php/api/qso"

	return agg_tcp.SendADI2Cloudlog(url, _payload)
}
