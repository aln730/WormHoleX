package main

import (
        "fmt"
        "log"
        "net/http"
        "net/http/httputil"
        "net/url"
        "sync"           
)

type Tunnel struct {
        Target *url.URL               // The destination URL to proxy to
        Proxy  *httputil.ReverseProxy
}

var (
        tunnels = make(map[string]*Tunnel)
        tunnelsMu sync.RWMutex
)



