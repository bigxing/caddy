// Package config contains utilities and types necessary for
// launching specially-configured server instances.
package config

import (
	"log"
	"net"
	"os"

	"github.com/mholt/caddy/middleware"
)

const (
	defaultHost = "localhost"
	defaultPort = "2015"
	defaultRoot = "."

	// The default configuration file to load if none is specified
	DefaultConfigFile = "Caddyfile"
)

// config represents a server configuration. It
// is populated by parsing a config file (via the
// Load function).
type Config struct {
	// The hostname or IP on which to serve
	Host string

	// The port to listen on
	Port string

	// The directory from which to serve files
	Root string

	// HTTPS configuration
	TLS TLSConfig

	// Middleware stack
	Middleware map[string][]middleware.Middleware

	// Functions (or methods) to execute at server start; these
	// are executed before any parts of the server are configured,
	// and the functions are blocking
	Startup []func() error

	// Functions (or methods) to execute when the server quits;
	// these are executed in response to SIGINT and are blocking
	Shutdown []func() error

	// The path to the configuration file from which this was loaded
	ConfigFile string
}

// Address returns the host:port of c as a string.
func (c Config) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

// TLSConfig describes how TLS should be configured and used,
// if at all. A certificate and key are both required.
type TLSConfig struct {
	Enabled     bool
	Certificate string
	Key         string
}

// Load loads a configuration file, parses it,
// and returns a slice of Config structs which
// can be used to create and configure server
// instances.
func Load(filename string) ([]Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// turn off timestamp for parsing
	flags := log.Flags()
	log.SetFlags(0)

	p, err := newParser(file)
	if err != nil {
		return nil, err
	}

	cfgs, err := p.parse()
	if err != nil {
		return []Config{}, err
	}

	for i := 0; i < len(cfgs); i++ {
		cfgs[i].ConfigFile = filename
	}

	// restore logging settings
	log.SetFlags(flags)

	return cfgs, nil
}

// IsNotFound returns whether or not the error is
// one which indicates that the configuration file
// was not found. (Useful for checking the error
// returned from Load).
func IsNotFound(err error) bool {
	return os.IsNotExist(err)
}

// Default makes a default configuration
// that's empty except for root, host, and port,
// which are essential for serving the cwd.
func Default() []Config {
	cfg := []Config{
		Config{
			Root: defaultRoot,
			Host: defaultHost,
			Port: defaultPort,
		},
	}
	return cfg
}
