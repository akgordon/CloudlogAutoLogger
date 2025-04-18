package agg_adi

import (
	"strconv"
	"strings"
)

// List of ADI paramters that we need to ignore and not process
// Example "parameters"  some ADI file wrape the ADI in this ADI wrapper
var adiIgnoreList = []string{"parameters"}

func contains(slice []string, item string) bool {
	return strings.Contains(strings.Join(slice, ","), item)
}

// Pattern:  "<aaaa:bbb>cccc "
// Can NOT use a simple string splitter because "<" may be in the data field
func ParseADIRecord(adi string) map[string]string {
	adi_fields := make(map[string]string)

	beginkeyIdx := 0
	endkeyIdx := 0
	colonIdx := 0
	maxIdx := len(adi)

	// get to beginning of key segment
	for ; beginkeyIdx < maxIdx; beginkeyIdx++ {
		if adi[beginkeyIdx] == '<' {
			endkeyIdx = beginkeyIdx
			colonIdx = beginkeyIdx

			// get to end of key segment
			for ; endkeyIdx < maxIdx; endkeyIdx++ {
				if adi[endkeyIdx] == '>' {
					break
				}
			}

			// Break up key into name and length
			if adi[endkeyIdx] == '>' {
				// Find the colon
				for colonIdx = beginkeyIdx; colonIdx < endkeyIdx; colonIdx++ {
					if adi[colonIdx] == ':' {
						break
					}
				}
				if adi[colonIdx] == ':' {
					keyName := adi[beginkeyIdx+1 : colonIdx]
					keylens := adi[colonIdx+1 : endkeyIdx]
					keylen, _ := strconv.Atoi(keylens)

					// See if parameter is in ignore list
					if !contains(adiIgnoreList, keyName) {
						if keylen > 0 {
							// Get data portion
							data := adi[endkeyIdx+1 : endkeyIdx+keylen+1]
							adi_fields[keyName] = data
							beginkeyIdx = endkeyIdx + keylen
						} else {
							adi_fields[keyName] = ""
						}
					} else {
						beginkeyIdx = endkeyIdx - 1
					}
				} else {
					// No colon - so no data
				}
			}
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
