package game

import (
    "net"
    "casino/common"
    "encoding/json"
)

func SendMsg(conn net.Conn, inf interface{}) {
   b, _ := json.Marshal(inf)
   conn.Write(append(b, common.DELIMITER)) 
}
