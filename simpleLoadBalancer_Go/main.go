package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	Attempts int = iota
	Retry
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

// SetAlive uses mux because multiple goroutines may try to call this function
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

func (b *Backend) IsAlive() (alive bool) {
	b.mux.RLock()
	alive = b.Alive
	b.mux.RLock()
	return alive
}

type ServerPool struct {
	backends []*Backend
	current  uint64
}

func (s *ServerPool) AddBackend(backend *Backend) {
	s.backends = append(s.backends, backend)
}

func (s *ServerPool) MarkBackendStatus(backendUrl *url.URL, alive bool) {
	for _, b := range s.backends {
		if b.URL.String() == backendUrl.String() {
			b.SetAlive(alive)
			break
		}
	}
}

// NextIndex finds an index in backends[] for the waiting clients
func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}

// GetNextPeer returns next active peer to take a connection
func (s *ServerPool) GetNextPeer() *Backend {
	next := s.NextIndex()
	l := len(s.backends) + next
	for i := next; i < l; i++ {
		idx := i % len(s.backends)
		if s.backends[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx))
			}
			return s.backends[idx]
		}
	}
	return nil
}

func (s *ServerPool) HealthCheck() {
	for _, b := range s.backends {
		status := "up"
		alive := isBackendAlive(b.URL)
		b.SetAlive(alive)
		if !alive {
			status = "down"
		}
		log.Printf("%s [%s]\n", b.URL, status)
	}
}

func GetRetryFromContext(request *http.Request) int {
	if retry, ok := request.Context().Value(Retry).(int); ok {
		return retry
	}
	return 0
}

// Balances the incoming request
func (s *ServerPool) lb(writer http.ResponseWriter, request *http.Request) {
	attempts := GetRetryFromContext(request)
	if attempts > 3 {
		log.Printf("%s(%s) Max attempts reached, terminating\n", request.RemoteAddr, request.URL.Path)
		http.Error(writer, "Service not available", http.StatusServiceUnavailable)
		return
	}

	peer := s.GetNextPeer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(writer, request)
		return
	}
	http.Error(writer, "Service not available", http.StatusServiceUnavailable)
}

// Passive health checks
func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false
	}
	// Close the connection if confirmed alive
	_ = conn.Close()
	return true
}

// Runs a routine for check status of the backends every 20 secs
func passiveHealthCheck() {
	t := time.NewTicker(time.Second * 20)
	for {
		select {
		case <-t.C:
			log.Println("Starting health check...")
			serverPool.HealthCheck()
			log.Println("Health check completed")
		}
	}
}

var serverPool ServerPool

func main() {
	var serverList string
	var port int
	flag.StringVar(&serverList, "backends", "", "Load balanced backends, use commas to separate")
	flag.IntVar(&port, "port", 3030, "Port to serve")
	flag.Parse()

	if len(serverList) == 0 {
		log.Fatal("Please provide one or more backends to load balance")
	}

	// parse servers
	tokens := strings.Split(serverList, ",")
	for _, tok := range tokens {
		serverUrl, err := url.Parse(tok)
		if err != nil {
			log.Fatal(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(serverUrl)
		// Actively check if backend is alive when error happens
		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			log.Printf("[%s] %s\n", serverUrl.Host, e.Error())
			retries := GetRetryFromContext(request)
			if retries < 3 {
				select {
				case <-time.After(10 * time.Millisecond):
					ctx := context.WithValue(request.Context(), Retry, retries+1)
					proxy.ServeHTTP(writer, request.WithContext(ctx))
				}
				return
			}

			serverPool.MarkBackendStatus(serverUrl, false)

			// Move to next peer
			attempts := GetRetryFromContext(request)
			log.Printf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
			ctx := context.WithValue(request.Context(), Attempts, attempts+1)
			serverPool.lb(writer, request.WithContext(ctx))
		}
		serverPool.AddBackend(&Backend{
			URL:          serverUrl,
			Alive:        true,
			ReverseProxy: proxy,
		})
		log.Printf("Configured server: %s\n", serverUrl)
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(serverPool.lb),
	}

	go passiveHealthCheck()

	log.Printf("Load Balancer started at :%d\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
