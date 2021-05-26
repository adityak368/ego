package mdns

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/adityak368/ego/registry"
	"github.com/adityak368/swissknife/logger/v2"
	"github.com/grandcat/zeroconf"
)

const (
	versionKey = "Version"

	defaultServiceName = "_ego._tcp"

	defaultDomain = "local."
)

func listInterfaces() ([]net.Interface, error) {
	interfaces := make([]net.Interface, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, ifi := range ifaces {
		if (ifi.Flags & net.FlagUp) == 0 {
			continue
		}
		// if (ifi.Flags & net.FlagMulticast) > 0 {
		// }
		interfaces = append(interfaces, ifi)
	}

	return interfaces, nil
}

// mdnsRegistry defines the MDNS registry
type mdnsRegistry struct {
	serviceName string
	domain      string
	options     registry.Options
	server      *zeroconf.Server
	services    map[string][]registry.Entry
	mutex       sync.Mutex
	cancelWatch context.CancelFunc
	entriesChan chan *zeroconf.ServiceEntry
	// wg is used to enforce Close() to return after the watcher() goroutine has finished.
	// Otherwise, data race will be possible. [Race Example] in dns_resolver_test we
	// replace the real lookup functions with mocked ones to facilitate testing.
	// If Close() doesn't wait for watcher() goroutine finishes, race detector sometimes
	// will warns lookup (READ the lookup function pointers) inside watcher() goroutine
	// has data race with replaceNetFunc (WRITE the lookup function pointers).
	wg sync.WaitGroup
}

// Init initializes the registry
func (r *mdnsRegistry) Init(opts registry.Options) error {
	r.options = opts
	return nil
}

// Options Returns the client options
func (r *mdnsRegistry) Options() registry.Options {
	return r.options
}

// Register adds the service to the registry
func (r *mdnsRegistry) Register(entry registry.Entry) error {

	ifaces, err := listInterfaces()
	if err != nil {
		return err
	}

	var meta []string
	for k, v := range entry.Metadata {
		meta = append(meta, fmt.Sprintf("%s=%s", k, v))
	}
	meta = append(meta, fmt.Sprintf("%s=%s", versionKey, entry.Version))

	parts := strings.Split(entry.Address, ":")

	if len(parts) != 2 {
		return errors.New("[MDNS-Registry]: Invalid Service Address")
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	server, err := zeroconf.Register(entry.Name, r.serviceName, r.domain, int(port), meta, ifaces)
	if err != nil {
		return err
	}

	logger.Info().Msgf("[MDNS-Registry]: Registered service '%s'", entry.Name)
	r.server = server
	return nil
}

// Deregister removes the service to the registry
func (r *mdnsRegistry) Deregister(serviceName string) error {
	r.server.Shutdown()
	logger.Info().Msgf("[MDNS-Registry]: Deregistered service '%s'", serviceName)
	return nil
}

// GetService Resolves the serviceName and returns the service details
func (r *mdnsRegistry) GetService(serviceName string) ([]registry.Entry, error) {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	endpoints, ok := r.services[serviceName]
	if !ok {
		return make([]registry.Entry, 0), nil
	}

	return endpoints, nil
}

// ListServices returns all the services in the registry
func (r *mdnsRegistry) ListServices() ([]registry.Entry, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var endpoints []registry.Entry
	for _, serviceEndpoints := range r.services {
		endpoints = append(endpoints, serviceEndpoints...)
	}

	return endpoints, nil
}

// Watch sets the registry to watch mode so that it tracks any updates
func (r *mdnsRegistry) Watch() error {

	ifaces, err := listInterfaces()
	if err != nil {
		return err
	}

	resolver, err := zeroconf.NewResolver(zeroconf.SelectIfaces(ifaces))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	err = resolver.Browse(ctx, r.serviceName, r.domain, r.entriesChan)
	if err != nil {
		defer cancel()
		return err
	}

	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case entry := <-r.entriesChan:
				if entry != nil {
					// Received a register broadcast from a service. So save it
					endpoints := make([]registry.Entry, 0)
					for _, ip := range entry.AddrIPv4 {
						endpoint := registry.Entry{
							Name:    entry.ServiceInstanceName(),
							Address: fmt.Sprintf("%s:%d", ip.String(), entry.Port),
						}
						meta := make(map[string]string)
						for _, v := range entry.Text {
							parts := strings.Split(v, "=")
							if len(parts) != 2 {
								continue
							}
							if parts[0] == versionKey {
								endpoint.Version = parts[1]
								continue
							}
							meta[parts[0]] = parts[1]
						}
						endpoint.Metadata = meta
						logger.Info().Msgf("[MDNS-Registry]: Discovered service %s version {%s} at %s", endpoint.Name, endpoint.Version, endpoint.Address)
						endpoints = append(endpoints, endpoint)
					}
					r.mutex.Lock()
					r.services[entry.ServiceInstanceName()] = endpoints
					r.mutex.Unlock()
				}
			}
		}
	}()

	r.cancelWatch = cancel

	return nil
}

// CancelWatch stops the registry watch mode
func (r *mdnsRegistry) CancelWatch() error {
	if r.cancelWatch != nil {
		r.cancelWatch()
		r.wg.Wait()
	}
	return nil
}

// String returns the description of the registry
func (r *mdnsRegistry) String() string {
	return "[MDNS-Registry]: Using MDNS Zeroconf as registry"
}

// New creates a new mdns registry
func New(info ...string) registry.Registry {

	var serviceName = defaultServiceName
	var domain = defaultDomain
	if len(info) == 1 {
		serviceName = info[0]
	}

	if len(info) == 2 {
		serviceName = info[0]
		domain = info[1]
	}

	return &mdnsRegistry{
		services:    make(map[string][]registry.Entry),
		entriesChan: make(chan *zeroconf.ServiceEntry, 10),
		serviceName: serviceName,
		domain:      domain,
		options:     registry.Options{},
	}
}
