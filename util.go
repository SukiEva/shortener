package shortener

import "log"

func DropError(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}
