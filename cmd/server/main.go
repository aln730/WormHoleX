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

func handleRegister(w https.ResponseWrite, r *http.Request){
    name := r.URL.Query().Get("name")
    target := r.URL.Query().Get("target")

    if name == "" || target == ""{
        http.Error(w, "Missing name or target", http.StatusBadRequest)
        return
    }

    targetURL,, err := url.Parse(target)
    if err != nil {
        http.Error(w, "Invalid target URL", http.StatusBadRequest)
        return
    }

    tunnelsMu.Lock()
    tunnels[name] = &Tunnel{
        Target: targetURL, Proxy: httputil.NewSingleHostReverseProxy(targetURL),
    }
    tunnelsMu.Unlock()

    fmt.Fprintf(w,"Tunnel registered: /%s -> %s\n", name, target)
}

func handleProxy(){

}

func main(){
    http.HandleFunc("/register", handleRegister)
    http.HandleFunc("/",handleProxy)

    log.Println("Server listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

