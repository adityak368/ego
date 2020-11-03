package client

import "github.com/adityak368/ego/registry"

// Options is the config for the client
type Options struct {
	Version  string
	Name     string
	Target   string
	Metadata map[string]string
	Registry registry.Registry
}
