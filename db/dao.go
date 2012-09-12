package db

import (
    "fmt"
    "net"
    "bytes"
    "strconv"
    "encoding/json"
)

//
type GetLogoutTimeReq struct {
    Op      string
    Uid     uint32
}

type GetLogTimeRep struct {
    Op      string
    Uid     uint32
    Ret     uint32
    Logouttime string
    Logintime  string
}

//
type SetLogoutTimeReq struct {
    Op      string
    Uid     uint32
}

type SetLogoutTimeRep struct {
    Op      string
    Uid     uint32
    Ret     uint32
    Time    string
}

//
type SetLoginTimeReq struct {
    Op      string
    Uid     uint32
}

type SetLoginTimeRep struct {
    Op      string
    Uid     uint32
    Ret     uint32
    Time    string
}

//
type ModifyBalanceReq struct {
    Op      string
    Uid     uint32
    Num     int32
}

type ModifyBalanceRep struct {
    Op      string
    Uid     uint32
    Ret     uint32
    Balance uint32
}

//
type GetBalanceReq struct {
    Op      string
    Uid     []uint32
}

type GetBalanceRep struct {
    Op      string
    Ret     uint32
    Ubl     map[string]uint32
}

type Ublance struct {
    Uid     uint32
    Balance uint32
}


type DbMgr struct {
    Conn    net.Conn
}

var (
    Dao DbMgr
)

func init() {
    Dao = DbMgr{}
    Dao.CheckConnect()
}


func (dao *DbMgr) CheckConnect() {
    if dao.Conn == nil {
        fmt.Println("Connect sqld")
        dao.Conn, _ = net.Dial("tcp", "127.0.0.1:12918")
    }
}


func (dao *DbMgr) Write(inf interface{}) (bool){
    dao.CheckConnect()
    jn, _ := json.Marshal(inf)
    dao.Conn.Write(append(jn, '\n'))
    return true
}


func (dao *DbMgr) ReadRet(inf interface{}) {
    buf := make([]byte, 1280)
    dao.Conn.Read(buf)
    json.Unmarshal(bytes.TrimRight(buf, string('\x00')), inf)
    fmt.Println("ReadRet", inf)
}


func (dao *DbMgr) GetLogTime(uid uint32) (string, string) {
    dao.Write(GetLogoutTimeReq{"getlogtime", uid})

    var rep GetLogTimeRep
    dao.ReadRet(&rep)
    return rep.Logintime, rep.Logouttime
}


func (dao *DbMgr) SetLogoutTime(uid uint32) (b bool) {
    dao.Write(SetLogoutTimeReq{"setlogouttime", uid})

    var rep SetLogoutTimeRep
    dao.ReadRet(&rep)
    return true
}


func (dao *DbMgr) SetLoginTime(uid uint32) (b bool) {
    dao.Write(SetLoginTimeReq{"setlogintime", uid})

    var rep SetLoginTimeRep
    dao.ReadRet(&rep)
    return true
}


func (dao *DbMgr) ModifyBalance(uid uint32, num int32) (*ModifyBalanceRep) {
    dao.Write(ModifyBalanceReq{"modifybalance", uid, num})

    var rep ModifyBalanceRep
    dao.ReadRet(&rep)
    return &rep
}


func (dao *DbMgr) GetBalance(uids []uint32) (map[uint32]uint32) {
    dao.Write(GetBalanceReq{"getbalance", uids})
    
    var rep GetBalanceRep
    dao.ReadRet(&rep)

    ret := make(map[uint32]uint32)
    for k, v := range rep.Ubl {
        i, _ := strconv.Atoi(k)
        ret[uint32(i)] = v
    }

    return ret
}


//func main() {
//    go Dao.SetLogoutTime(10593000)
//    go Dao.GetLogoutTime(10593000)
//    go Dao.ModifyBalance(10593000, 22)
//
//    var line string
//    fmt.Scanln(&line) 
//}

