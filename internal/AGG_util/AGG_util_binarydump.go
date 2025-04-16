package AGG_util

import (
	"CloudlogAutoLogger/internal/AGG_logger"
	"fmt"
)

func BinaryDump(data []byte, length int, useLogger bool) {
	const bytesPerLine = 16

	if useLogger {
		AGG_logger.Get().Log("           0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F\n", "")
	} else {
		fmt.Printf("           0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F\n")
	}

	for i := 0; i < len(data); i += bytesPerLine {

		if length < 0 {
			break
		}
		length -= bytesPerLine

		outStg := ""

		// Print the offset
		outStg = fmt.Sprintf("%08x  ", i)
		//fmt.Printf("%08x  ", i)

		// Print the binary representation
		for j := 0; j < bytesPerLine; j++ {
			if i+j < len(data) {
				//fmt.Printf("%08b ", data[i+j])
				outStg += fmt.Sprintf("%02x ", data[i+j])
			} else {
				outStg += fmt.Sprintf("         ")
			}
		}

		// Print the ASCII representation
		outStg += fmt.Sprintf(" |")
		for j := 0; j < bytesPerLine; j++ {
			if i+j < len(data) {
				b := data[i+j]
				if b >= 32 && b <= 126 {
					outStg += fmt.Sprintf("%c", b)
				} else {
					outStg += fmt.Sprintf(".")
				}
			} else {
				outStg += fmt.Sprintf(" ")
			}
		}
		outStg += fmt.Sprintf("|")

		// Print
		if useLogger {
			AGG_logger.Get().Log(outStg, "")
		} else {
			fmt.Println(outStg)
		}
	}
}
