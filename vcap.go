package vcapenv

import (
	"encoding/json"
	"fmt"
)

type VCAPServices struct {
	UserProvided []struct {
		Name        string                     `json:"name"`
		Credentials map[string]json.RawMessage `json:"credentials"`
	} `json:"user-provided"`
}

type Namer interface {
	Name(service, variable string) string
}

type NamerFunc func(service, variable string) string

func (f NamerFunc) Name(service, variable string) string {
	return f(service, variable)
}

func ServiceNamer(service, variable string) string {
	return fmt.Sprintf("%s_%s", service, variable)
}

func (services VCAPServices) Environment(namer Namer) Environment {
	env := make(Environment)

	for _, service := range services.UserProvided {
		for credential, msg := range service.Credentials {

			name := namer.Name(service.Name, credential)
			env[name] = append(env[name], msg)
		}
	}

	return env
}
