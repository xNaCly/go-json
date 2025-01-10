package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	// "runtime/pprof"

	"github.com/xnacly/libjson"
)

func main() {
	// f, err := os.Create("cpu.pprof")
	// if err != nil {
	// 	panic(err)
	// }
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()
	lj := flag.Bool("libjson", true, "benchmark libjson or gojson")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		log.Fatalln("Wanted a file as first argument, got nothing, exiting")
	}
	file, err := os.Open(args[0])
	if err != nil {
		log.Fatalln(err)
	}
	if *lj {
		_, err := libjson.NewReader(file)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		v := []struct {
			Key1      string
			Array     []any
			Obj       any
			AtomArray []any
		}{}
		d := json.NewDecoder(file)
		err := d.Decode(&v)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
