package game

import (
    "net"
    "fmt"
    "casino/common"
    "encoding/json"
)

func SendMsg(conn net.Conn, inf interface{}) {
   b, _ := json.Marshal(inf)
   conn.Write(append(b, common.DELIMITER)) 
}

func GetNameByUid(uid) {
    pre := "http://222.88.95.242:3737/batchgetuserinfo/uids/"
    url := fmt.Sprintf("%v%v/", pre, uid)

    res, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    robots, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Fatal(err)
    }
    res.Body.Close()
    fmt.Printf("%s", robots)
}
