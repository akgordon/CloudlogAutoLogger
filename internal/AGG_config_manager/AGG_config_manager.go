package agg_config_manager

import (
	agg_logger "CloudlogAutoLogger/internal/AGG_logger"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type AGG_config_manager struct {
	Cloudlog_url       string
	Cloudlog_api_key   string
	Station_profile_id string
	WSJTX_port         int
	VARAC_port         int
	JS8CALL_port       int
}

const (
	filename  string = "cloudlog_auto_logger.ini"
	crypt_key string = "2e899fe6ffb07f2e1a63f2f619a7f4daddee80eb92a9ed5e328a0c1e2a1c0c58"
)

func (cd *AGG_config_manager) init() {
	cd.Cloudlog_url = ""
	cd.Cloudlog_api_key = ""
	cd.Station_profile_id = ""
	cd.WSJTX_port = 0
	cd.VARAC_port = 0
	cd.JS8CALL_port = 0
}

func GetConfig() (AGG_config_manager, bool) {

	//newKey := Genkey()
	//fmt.Println(newKey)

	// Setup return structure
	var cd AGG_config_manager
	cd.init()

	// Open file
	var filePtr *os.File
	var err error
	filePtr, err = os.Open(filename)
	if err != nil {
		agg_logger.Get().Log(err.Error(), "")
		return cd, false
	}
	defer filePtr.Close()

	// Read in structure
	scanner := bufio.NewScanner(filePtr)
	for scanner.Scan() {
		ln := scanner.Text()
		if strings.Index(ln, "Cloudlog_api_key =") == 0 {
			ekey := ln[18:]
			if len(ekey) > 0 {
				dkey := decrypt(ekey, crypt_key)
				if err != nil {
					agg_logger.Get().Log(err.Error(), "")
					return cd, false
				}
				cd.Cloudlog_api_key = string(dkey)
			}
		}

		if strings.Index(ln, "Cloudlog_url =") == 0 {
			cd.Cloudlog_url = ln[14:]
		}

		if strings.Index(ln, "Station_profile_id =") == 0 {
			cd.Station_profile_id = ln[20:]
		}

		if strings.Index(ln, "WSJTX_port =") == 0 {
			port := ln[12:]
			cd.WSJTX_port, _ = strconv.Atoi(port)
		}

		if strings.Index(ln, "VARAC_port =") == 0 {
			port := ln[12:]
			cd.VARAC_port, _ = strconv.Atoi(port)
		}

		if strings.Index(ln, "JS8CALL_port =") == 0 {
			port := ln[14:]
			cd.JS8CALL_port, _ = strconv.Atoi(port)
		}
	}

	return cd, true
}

func SaveConfig(cd AGG_config_manager) bool {

	// Open file
	var filePtr *os.File
	var err error
	filePtr, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		agg_logger.Get().Log(err.Error(), "")
		return false
	}

	// First encrypt the api key
	ekey := encrypt(cd.Cloudlog_api_key, crypt_key)

	filePtr.WriteString("Cloudlog_api_key =" + ekey + "\n")
	filePtr.WriteString("Cloudlog_url =" + cd.Cloudlog_url + "\n")
	filePtr.WriteString("Station_profile_id =" + cd.Station_profile_id + "\n")
	filePtr.WriteString("WSJTX_port =" + fmt.Sprint(cd.WSJTX_port) + "\n")
	filePtr.WriteString("VARAC_port =" + fmt.Sprint(cd.VARAC_port) + "\n")
	filePtr.WriteString("JS8CALL_port =" + fmt.Sprint(cd.JS8CALL_port) + "\n")

	return true
}
