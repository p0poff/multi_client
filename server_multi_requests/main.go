// hello world
package main

import (
	"github.com/jessevdk/go-flags"
	"log"
	"os"
)

type opts struct {
	Port    string `short:"p" long:"port" env:"PORT" default:"8080" description:"port to listen on"`
	Dbg     bool   `long:"dbg" env:"DEBUG" description:"debug mode"`
	Timeout int    `short:"t" long:"timeout" env:"TIMEOUT" default:"30" description:"timeout in seconds"`
}

func setupLog(dbg bool) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)
}

func main() {
	//go run -mod=vendor ./ -p 8000

	var opts opts

	p := flags.NewParser(&opts, flags.Default)
	_, err := p.Parse()

	if err != nil {
		log.Println("[CRITICAL] Error parsing flags")
		return
	}

	setupLog(opts.Dbg)

	log.Println("[INFO] App start!")

	if opts.Dbg {
		log.Println("[INFO] Debug mode enabled")
	}

	s := NewServer(&opts)

	if err = s.Start(); err != nil {
		log.Printf("[ERROR] %v", err)
	}

	log.Println("[INFO] App stop!")
}
