package registry

import (
	"crypto/tls"
	"time"
)

// Entry defines the entry in the registry
type Entry struct {
	Name     string
	Version  string
	Metadata map[string]string
	Address  string
}

// Options is the config for the registry
type Options struct {
	Host      string
	Port      uint16
	Timeout   time.Duration
	Secure    bool
	TLSConfig *tls.Config
}
