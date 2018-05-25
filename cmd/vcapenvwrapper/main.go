package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/govau/vcapenv"
)

type services map[string]struct{}

func (s services) String() string {
	var terms []string
	for service := range s {
		terms = append(terms, service)
	}

	return strings.Join(terms, ",")
}

func (s services) Set(value string) error {
	s[value] = struct{}{}
	return nil
}

func (s services) Name(service, variable string) string {
	if _, ok := s[service]; ok {
		return variable
	}

	return vcapenv.ServiceNamer(service, variable)
}

func main() {
	var rootServices = make(services)

	flag.Var(&rootServices, "root", "set variables in this service at the root level")
	flag.Parse()
	fargs := flag.Args()

	if len(fargs) < 1 {
		log.Fatal("provide a command to run")
	}

	err := vcapenv.Run(rootServices, vcapenv.SetenvExporter)
	if err != nil {
		log.Fatal(err)
	}

	cmd := fargs[0]
	args := fargs[1:]

	command := exec.Command(cmd, args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err = command.Run()
	if err != nil {
		log.Fatal(err)
	}
}
