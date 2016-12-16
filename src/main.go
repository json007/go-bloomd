package main

import (
    // "net/http"
    "bloomfilter"
    "log"
    // "html"
    // "io"
    "fmt"
    // "flag"
)
// var bf *bloomfilter.BloomFilter

var (
    listen  = flag.String("l", "127.0.0.1", "Interface to listen on. Default to all addresses.")
    port    = flag.Int("p", 10086, "TCP port number to listen on (default: 11211)")
    //threads = flag.Int("t", runtime.NumCPU(), fmt.Sprintf("number of threads to use (default: %d)", runtime.NumCPU()))
)

func main(){
    flag.Parse()
    address := fmt.Sprintf("%s:%d", *listen, *port)
    server := bloomfilter.NewServer(address)
    log.Fatal(server.ListenAndServer())

    // bf = bloomfilter.New(10)
    // http.HandleFunc("/add", add)
    // log.Fatal(http.ListenAndServe(":80",nil))
}

// func add(w http.ResponseWriter, r *http.Request){
//     bf.Add([]byte(r.URL.RawQuery))
//     fmt.Fprint(w,"Hello, %q", html.EscapeString(r.URL.RawQuery))
//     // io.WriteString(w, "hello, world!\n")
// }

