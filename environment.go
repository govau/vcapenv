package vcapenv

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type Environment map[string][]json.RawMessage

// EnvironmentValue extracts a suitable variable value from an environment
// json message.
//
// Specifically, unquote strings but leaves most other types as-is.
func EnvironmentValue(msg json.RawMessage) string {
	var v interface{}

	err := json.Unmarshal(msg, &v)
	if err != nil {
		log.Println(err)
		return string(msg)
	}

	switch value := v.(type) {
	case string:
		return value
	}

	return string(msg)
}

// Iter calls some function for each variable in the environment
//
// Variable names will be uppercased
// When multiple values exist for a variable, the 'first' is exposed.
func (e Environment) Iter(exporter VariableExporter) {
	for variable, values := range e {
		if len(values) > 0 {
			variable = strings.ToUpper(variable)
			value := EnvironmentValue(values[0])

			exporter.Export(variable, value)
		}
	}
}

type Source interface {
	Environment(Namer) Environment
}

func GetVCAPServices() (VCAPServices, error) {
	var services VCAPServices
	vcap := os.Getenv("VCAP_SERVICES")
	err := json.NewDecoder(strings.NewReader(vcap)).Decode(&services)

	return services, err
}

func GetVCAPApplication() (VCAPApplication, error) {
	var application VCAPApplication
	vcap := os.Getenv("VCAP_APPLICATION")
	err := json.NewDecoder(strings.NewReader(vcap)).Decode(&application)

	return application, err
}

func Run(namer Namer, exporter VariableExporter) error {
	services, err := GetVCAPServices()
	if err != nil {
		return err
	}

	services.Environment(namer).Iter(exporter)
	return nil
}
