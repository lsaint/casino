
package game

import (
    "net/http"
    "fmt"
)

type HttpServer struct {
}


func (hts *HttpServer) Start(c RoundMap) {
    http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
        status := ""
        for _, round := range c {
            if len(round.Players) > 0 {
                status += fmt.Sprintf("%+v\n", *round)
            }
        }
        fmt.Fprint(w, status)
    })

    http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
        p, u := 0, 0
        for _, round := range c {
            u += len(round.Users)
            p += len(round.Players)
        }
        fmt.Fprint(w, fmt.Sprintf("user:%v player:%v total:%v", u, p, p+u))
    })

    fmt.Println("httpServer runing")
    http.ListenAndServe(":20133", nil)
}
