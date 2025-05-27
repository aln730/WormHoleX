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

func handleRegister(){

}

func handleProxy(){

}

func main(){
    http.HandleFunc("/register", handleRegister)
    http.HandleFunc("/",handleProxy)

    log.Println("Server listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

