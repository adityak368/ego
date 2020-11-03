package server

import "github.com/adityak368/ego/registry"

// Options is the config for the Server
type Options struct {
	Version  string
	Name     string
	Address  string
	Metadata map[string]string
	Registry registry.Registry
}
