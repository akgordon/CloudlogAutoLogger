package agg_adi

import (
	"strconv"
	"strings"
)

func split(r rune) bool {
	return r == '<'
}

// Pattern:  "<aaaa:bbb>cccc "
func Parse_adi(pkt_stg string) map[string]string {

	adi_fields := make(map[string]string)

	flds := strings.FieldsFunc(pkt_stg, split)
	for _, v := range flds {
		i := strings.Index(v, ">")
		a := v[:i]
		a = strings.ToLower(a)
		if a != "eor" {
			j := strings.Index(a, ":")
			b := a[j+1:]
			a = a[:j]
			c := v[i+1:]

			l, _ := strconv.Atoi(b)
			if l > 0 {
				// limit len
				if l > len(c) {
					l = len(c)
				}
				c = c[:l]
			} else {
				c = ""
			}

			adi_fields[a] = c
		} else {
			break
		}
	}

	return adi_fields
}

// <call:5>N9EAT
// <band:4>70cm
// <mode:3>SSB
// <freq:10>432.166976
// <qso_date:8>20190616
// <time_on:6>170600
// <time_off:6>170600
// <rst_rcvd:2>59
// <rst_sent:2>55
// <qsl_rcvd:1>N
// <qsl_sent:1>N
// <country:24>United States Of America
// <gridsquare:4>EN42
// <sat_mode:3>U/V
// <sat_name:4>AO-7
// <prop_mode:3>SAT
// <name:5>Marty
// <eor>"

func Encode_adi(fields map[string]string) string {
	adi := ""

	for key, value := range fields {
		p := "<"
		p += key
		p += ":"
		p += strconv.Itoa(len(value))
		p += ">"
		p += value
		p += " "

		adi += p
	}

	adi += "<eor>"

	return adi

}
