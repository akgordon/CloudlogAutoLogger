package agg_config_manager

import (
	"CloudlogAutoLogger/internal/agg_logger"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type AGG_config_manager struct {
	Cloudlog_api_key   string
	Station_profile_id string
	WSJTX_port         int
	VARAC_port         int
	JS8Call_port       int
}

const (
	filename  string = "cloudlog_auto_logger.ini"
	crypt_key string = "n5M7rBYZvO+2Oq6SeZIyIeoV44AY3hlrG/u/ouTu8lQ6ZY71We9XGJsb97Ud3XyI"
)

var theConfig_manager *AGG_config_manager = nil

func Get() *AGG_config_manager {
	if theConfig_manager == nil {
		theConfig_manager = &AGG_config_manager{}
		theConfig_manager.init()
	}
	return theConfig_manager
}

func (cd *AGG_config_manager) init() {
	cd.Cloudlog_api_key = ""
	cd.Station_profile_id = ""
	cd.WSJTX_port = 0
	cd.VARAC_port = 0
	cd.JS8Call_port = 0
}

func (cd *AGG_config_manager) GetConfig() bool {

	// Clear config
	cd.init()

	// Open file
	var filePtr *os.File
	var err error
	filePtr, err = os.Open(filename)
	if err != nil {
		agg_logger.Get().Log(err.Error(), "")
		return false
	}
	defer filePtr.Close()

	// Read in structure
	scanner := bufio.NewScanner(filePtr)
	for scanner.Scan() {
		ln := scanner.Text()
		if strings.Index(ln, "Cloudlog_api_key =") == 0 {
			ekey := ln[18:]
			eba := []byte(ekey)
			key := []byte(crypt_key)
			var dba []byte
			dba, err = decrypt(key, eba)
			if err != nil {
				agg_logger.Get().Log(err.Error(), "")
				return false
			}
			cd.Cloudlog_api_key = string(dba)
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
		if strings.Index(ln, "JS8Call_port =") == 0 {
			port := ln[14:]
			cd.JS8Call_port, _ = strconv.Atoi(port)
		}

		return true
	}

	return true
}

func (cd *AGG_config_manager) SetConfig() bool {

	// Open file
	var filePtr *os.File
	var err error
	filePtr, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		agg_logger.Get().Log(err.Error(), "")
		return false
	}

	// First encrypt the api key
	key := []byte(crypt_key)
	ekey, err := encrypt(cd.Cloudlog_api_key, key)

	filePtr.WriteString("Cloudlog_api_key =" + ekey + "\n")
	filePtr.WriteString("Station_profile_id =" + cd.Station_profile_id + "\n")
	filePtr.WriteString("WSJTX_port =" + fmt.Sprint(cd.WSJTX_port) + "\n")
	filePtr.WriteString("VARAC_port =" + fmt.Sprint(cd.VARAC_port) + "\n")
	filePtr.WriteString("JS8Call_port =" + fmt.Sprint(cd.JS8Call_port) + "\n")

	return true
}
