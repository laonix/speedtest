package main

import (
	"fmt"
	"log"

	"github.com/laonix/speedtest"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	flagProvider = kingpin.Arg(
		"provider",
		"Provider to perform a speed test with. Should be either 'ookla' or 'netflix'",
	).Required().String()
)

func main() {
	kingpin.Parse()

	result, err := speedtest.Run(*flagProvider)
	if err != nil {
		log.Fatal("speedtest: ", err)
	}
	fmt.Println(result)
}
