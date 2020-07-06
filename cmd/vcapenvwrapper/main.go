package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"log/syslog"
	"os"
	"os/exec"
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

func getOutputStreams() (stdout, stderr io.Writer) {
	stdout = os.Stdout
	stderr = os.Stderr

	syslogURL, ok := os.LookupEnv("SYSLOG_URL")
	if !ok {
		return
	}

	applicationLabel := "UNKNOWN_APP"
	application, err := vcapenv.GetVCAPApplication()
	if err == nil {
		applicationLabel = fmt.Sprintf(
			"%s.%s.%s",
			application.OrganizationName,
			application.SpaceName,
			application.ApplicationName,
		)
	}

	syslogout, err := syslog.Dial("udp", syslogURL, syslog.LOG_INFO|syslog.LOG_DAEMON, applicationLabel)
	if err != nil {
		return
	}

	syslogerr, err := syslog.Dial("udp", syslogURL, syslog.LOG_ERR|syslog.LOG_DAEMON, applicationLabel)
	if err != nil {
		return
	}

	stdout = io.MultiWriter(os.Stdout, syslogout)
	stderr = io.MultiWriter(os.Stderr, syslogerr)
	return stdout, stderr
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

	stdout, stderr := getOutputStreams()

	cmd := fargs[0]
	args := fargs[1:]

	command := exec.Command(cmd, args...)
	command.Stdin = os.Stdin
	command.Stdout = stdout
	command.Stderr = stderr
	err = command.Run()
	if err != nil {
		log.Fatal(err)
	}
}
