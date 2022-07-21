package main

import "net/http"
import "fmt"
import "io/ioutil"

func main() {}
    resp, err := http.Get("https://icanhazip.com/")
    if err != nil {
    	// handle error
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
}