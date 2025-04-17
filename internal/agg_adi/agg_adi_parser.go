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
