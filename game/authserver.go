package game

import (
    "net"
    "fmt"

    "casino/common"
)

type AuthServer struct {
}

func (aus *AuthServer) Start() {
    ln, err := net.Listen("tcp", ":10133")
    if err != nil {
        fmt.Println("authServer err")
    }
    fmt.Println("authServer runing")
    for {
        conn, err := ln.Accept()
        if err != nil {
            continue
        }
        fmt.Println("auth accept")
        go func(c net.Conn) {
            c.Write([]byte(common.XML_REP))
            c.Write([]byte("\x00"))
        }(conn)
    }
}
