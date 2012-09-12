package game

import (
    "fmt"
    "net"
    "bytes"
    "time"
    "encoding/json"

    "casino/common"
)


type GameServer struct {
    Op2handle   map[string]func(JsonString, net.Conn)
}


func (gas *GameServer) Start() {
    ln, err := net.Listen("tcp", ":30133")
    if err != nil {
        // handle error
        fmt.Println("err")
    }
    fmt.Println("gameServer runing")
    for {
        conn, err := ln.Accept()
        if err != nil {
            // handle error
            fmt.Println("Accept error")
            continue
        }
        go gas.acceptConn(conn)
    }
}

func (gas *GameServer) acceptConn(c net.Conn) {
    fmt.Println("acceptConn")
    buf := make([]byte, common.BUF_SIZE)
    var data bytes.Buffer
    for {
        n, _ := c.Read(buf)
        if n == 0 {
            fmt.Println("close by peer")
            break
        }
        data.Write(buf[:n])

        cn := bytes.Count(data.Bytes(), []byte{common.DELIMITER})
        for ; cn > 0; cn-- {
            jn, err := data.ReadString(common.DELIMITER)
            fmt.Println(time.Now().String()[:19], jn)
            if err != nil {
                fmt.Println("err", err)
                continue
            }

            var unknow interface{}
            err = json.Unmarshal([]byte(jn), &unknow)
            if err != nil {
                fmt.Println("Unmarshal error")
                continue
            }
            
            switch unknow.(type) {
                case map[string]interface{}: //?
                    gas.dispatchOp(unknow, c)
            }
        }
    }
}


func (gas *GameServer) dispatchOp(jn interface{}, c net.Conn) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("dispatchOp recover: ", r)
        }
    }()

    m := JsonString(jn.(map[string]interface{}))
    gas.Op2handle[m["Op"].(string)](m, c)
}


