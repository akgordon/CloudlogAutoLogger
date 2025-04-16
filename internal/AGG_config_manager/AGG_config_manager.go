package AGG_config_manager

import (
	"CloudlogAutoLogger/internal/AGG_logger"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type AGG_config_manager struct {
	cloudlog_api_key   string
	station_profile_id string
	WSJTX_port         int
	VARAC_port         int
	JS8Call_port       int
}

const (
	filename  string = "cloudlog_auto_logger.ini"
	crypt_key string = "n5M7rBYZvO+2Oq6SeZIyIeoV44AY3hlrG/u/ouTu8lQ6ZY71We9XGJsb97Ud3XyI"
)

func (cd *AGG_config_manager) getConfig() bool {

	// Clear config
	cd.cloudlog_api_key = ""
	cd.station_profile_id = ""
	cd.WSJTX_port = 0
	cd.VARAC_port = 0
	cd.JS8Call_port = 0

	// Open file
	var filePtr *os.File
	var err error
	filePtr, err = os.Open(filename)
	if err != nil {
		AGG_logger.Get().Log(err.Error(), "")
		return false
	}
	defer filePtr.Close()

	// Read in structure
	scanner := bufio.NewScanner(filePtr)
	for scanner.Scan() {
		ln := scanner.Text()
		if strings.Index(ln, "cloudlog_api_key =") == 0 {
			ekey := ln[18:]
			eba := []byte(ekey)
			key := []byte(crypt_key)
			var dba []byte
			dba, err = decrypt(key, eba)
			if err != nil {
				AGG_logger.Get().Log(err.Error(), "")
				return false
			}
			cd.cloudlog_api_key = string(dba)
		}
		if strings.Index(ln, "station_profile_id =") == 0 {
			cd.station_profile_id = ln[20:]
		}
		if strings.Index(ln, "WSJTX_port =") == 0 {
			port := ln[12:]
			cd.WSJTX_port, _ = strconv.Atoi(port)
		}
		if strings.Index(ln, "VARAC_port =") == 0 {
			port := ln[12:]
			cd.VARAC_port, _ = strconv.Atoi(port)
		}
		if strings.Index(ln, "JS8Call_port =") == 0 {
			port := ln[14:]
			cd.JS8Call_port, _ = strconv.Atoi(port)
		}

		return true
	}

	return true
}

func (cd *AGG_config_manager) setConfig() bool {

	// Open file
	var filePtr *os.File
	var err error
	filePtr, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		AGG_logger.Get().Log(err.Error(), "")
		return false
	}

	// First encrypt the api key
	key := []byte(crypt_key)
	ekey, err := encrypt(cd.cloudlog_api_key, key)

	filePtr.WriteString("cloudlog_api_key =" + ekey + "\n")
	filePtr.WriteString("station_profile_id =" + cd.station_profile_id + "\n")
	filePtr.WriteString("WSJTX_port =" + fmt.Sprint(cd.WSJTX_port) + "\n")
	filePtr.WriteString("VARAC_port =" + fmt.Sprint(cd.VARAC_port) + "\n")
	filePtr.WriteString("JS8Call_port =" + fmt.Sprint(cd.JS8Call_port) + "\n")

	return true
}
