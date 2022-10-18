package helpers

import "fmt"

func Log(message string, params ...interface{}) {
	if DEBUG_POLICIES_ALG {
		fmt.Printf(message+"\n", params...)
	}
}
