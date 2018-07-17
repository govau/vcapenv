package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/govau/vcapenv"
)

type services map[string]struct{}

var _ flag.Value = (services)(nil)

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

	err := vcapenv.Run(rootServices, vcapenv.BashExporter{os.Stdout})
	if err != nil {
		log.Fatal(err)
	}
}
