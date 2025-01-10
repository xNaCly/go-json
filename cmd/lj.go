package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xnacly/libjson"
)

func Must[T any](t T, err error) T {
	if err != nil {
		log.Fatalln(err)
	}
	return t
}

func main() {
	args := os.Args
	var file *os.File
	if info, err := os.Stdin.Stat(); err != nil || info.Mode()&os.ModeCharDevice != 0 { // we are in a pipe
		if len(args) == 1 {
			log.Fatalln("Wanted a file as first argument, got nothing, exiting")
		}
		file = Must(os.Open(args[1]))
	} else {
		file = os.Stdin
	}
	query := os.Args[len(os.Args)-1]
	json := Must(libjson.NewReader(file))
	fmt.Printf("%+#v\n", Must(libjson.Get[any](&json, query)))
}
