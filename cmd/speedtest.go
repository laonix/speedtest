package main

import (
	"fmt"
	"log"

	"github.com/laonix/speedtest"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	provider = kingpin.Arg(
		"provider",
		"Provider to perform a speed test with. Should be either 'ookla' or 'netflix'.",
	).Required().Enum("ookla", "netflix")
)

func main() {
	kingpin.Parse()

	runner := speedtest.NewRunner()

	result, err := runner.Run(*provider)
	if err != nil {
		log.Fatal("speedtest: ", err)
	}

	fmt.Println(result)
}
