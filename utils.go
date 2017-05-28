package tmsh

import (
	"strings"
)

func removeCarriageReturn(str []byte) (retStr []byte) {
	retStr = make([]byte, 0)
	for _, c := range str {
		if c != 13 {
			retStr = append(retStr, c)
		}
	}
	return
}

func removeSpaceAndBackspace(str string) string {
	return strings.Replace(str, " \b", "", -1)
}
