package report

import (
	"fmt"
	"syscall"
)

// A utility to convert the values to proper strings.
func int8ToStr(arr []int8) string {
	b := make([]byte, 0, len(arr))
	for _, v := range arr {
		if v == 0x00 {
			break
		}
		b = append(b, byte(v))
	}
	return string(b)
}

func GetSystemVersion() (string, error) {
	var uname syscall.Utsname
	if err := syscall.Uname(&uname); err != nil {
		return "", err
	}

	// extract members:
	// type Utsname struct {
	//  Sysname    [65]int8
	//  Nodename   [65]int8
	//  Release    [65]int8
	//  Version    [65]int8
	//  Machine    [65]int8
	//  Domainname [65]int8
	// }

	return fmt.Sprintf("%s %s %s %s",
		int8ToStr(uname.Sysname[:]),
		int8ToStr(uname.Release[:]),
		int8ToStr(uname.Version[:]),
		int8ToStr(uname.Machine[:])), nil
}
