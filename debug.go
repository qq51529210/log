package log

import "log"

type debugLogger struct {
	log.Logger
}

func f() {
	log.Println()
}
