package main

import(
    "flag"
    "fmt"
    "log"
    "net/http"
    "net/url"
)

func main(){

    name := flag.String("name","test", "Subdomain name")
    local := flag.String("local","http://localhost:3000/","Local target URL")
    server := flag.String("server","http://localhost:8080/","Tunnel server URL")
    flag.Parse()

    registerURL := fmt.Sprintf("%s/register?name=%s&target=%s", *server, *name, *local)
    resp, err := http.Get(registerURL)
    if err != nil {
        log.Fatalf("Failed to register tunnel: %v", err)
    }
    defer resp.Body.Close()

    log.Printf("Tunnel registered. Access via %s/%s\n", *server, *name)

    http.HandleFunc("/", func(w http.ResponseWriter, r *httpRequest){
        fmt.Fprintln(w, "Greetings from the local server")
    })

    parsedURL,_ := url.Parse(*local)
    log.Printf("Local server running at %s\n", parsedURL.Host)
    log.Fatal(http.ListenAndServe(parsedURL.Host, nil))
}
