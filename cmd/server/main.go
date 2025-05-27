package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Tunnel struct {
	Target *url.URL
	Proxy  *httputil.ReverseProxy
}

var (
	tunnels   = make(map[string]*Tunnel)
	tunnelsMu sync.RWMutex

	tcpListenAddr = flag.String("tcp", "", "TCP listen address for forwarding, e.g. :9000")
	tcpTargetAddr = flag.String("tcp-target", "", "TCP target address for forwarding, e.g. localhost:22")
)

func handleRegister(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	target := r.URL.Query().Get("target")

	if name == "" || target == "" {
		http.Error(w, "Missing name or target", http.StatusBadRequest)
		return
	}

	targetURL, err := url.Parse(target)
	if err != nil {
		http.Error(w, "Invalid target URL", http.StatusBadRequest)
		return
	}

	tunnelsMu.Lock()
	tunnels[name] = &Tunnel{
		Target: targetURL,
		Proxy:  httputil.NewSingleHostReverseProxy(targetURL),
	}
	tunnelsMu.Unlock()

	fmt.Fprintf(w, "Tunnel registered: /%s -> %s\n", name, target)
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[1:]

	tunnelsMu.RLock()
	tunnel, ok := tunnels[name]
	tunnelsMu.RUnlock()

	if !ok {
		http.Error(w, "Tunnel not found", http.StatusNotFound)
		return
	}

	log.Printf("Proxying request for %s to %s\n", name, tunnel.Target)
	tunnel.Proxy.ServeHTTP(w, r)
}

func startTCPForwarder(listenAddr, targetAddr string) {
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("TCP listen error on %s: %v", listenAddr, err)
	}
	log.Printf("TCP forwarding started on %s -> %s", listenAddr, targetAddr)

	go func() {
		for {
			clientConn, err := ln.Accept()
			if err != nil {
				log.Printf("TCP accept error: %v", err)
				continue
			}
			go handleTCPConn(clientConn, targetAddr)
		}
	}()
}

func handleTCPConn(clientConn net.Conn, targetAddr string) {
	defer clientConn.Close()

	serverConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Printf("TCP dial target error: %v", err)
		return
	}
	defer serverConn.Close()

	// Bidirectional copy
	go func() {
		_, _ = io.Copy(serverConn, clientConn)
	}()
	_, _ = io.Copy(clientConn, serverConn)
}

func main() {
	flag.Parse()

	if *tcpListenAddr != "" && *tcpTargetAddr != "" {
		startTCPForwarder(*tcpListenAddr, *tcpTargetAddr)
	} else {
		log.Println("TCP forwarding disabled (missing --tcp and --tcp-target flags)")
	}

	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/", handleProxy)

	log.Println("HTTP Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
