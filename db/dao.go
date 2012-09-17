package db

import (
    //"fmt"
    "strconv"
    "net/rpc"
    "net/rpc/jsonrpc"
)


/*
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

*/

type argUid struct {
    Uid     uint32
}

type getLogTimeRep struct {
    Logouttime string                                                      
    Logintime  string
}

type setTimeRep struct {
    Time        string 
}


type argModifyBalance struct {
    Uid     uint32
    Num     int32
}

type modifyBalanceRep struct {
    Balance     uint32
}

type argGetBalance struct {
    Uid     []uint32
}

type getBalanceRep struct {
    Ubl     map[string]uint32
}

func (rep *getBalanceRep) keyToUint32() (map[uint32]uint32) {
    ret := make(map[uint32]uint32)
    for k, v := range rep.Ubl {
        uid, _ := strconv.Atoi(k)
        ret[uint32(uid)] = v
    }
    return ret
}

const (
    PY_RPC_ADDR =  "127.0.0.1:12918"
)


func GetLogTime(uid uint32) (intime, outtime string, err error) {
    var cli *rpc.Client
    if cli, err = jsonrpc.Dial("tcp", PY_RPC_ADDR); err != nil {
        return 
    }
    defer cli.Close()
    
    args := &argUid{uid}
    reply :=  new(getLogTimeRep)
    if err = cli.Call("getLogTime", args, reply); err != nil {
        return 
    }
    intime, outtime = reply.Logintime, reply.Logouttime
    return
}


func setTime(uid uint32, isSetLogout bool) (t string, err error){
    var cli *rpc.Client
    if cli, err = jsonrpc.Dial("tcp", PY_RPC_ADDR); err != nil {
        return 
    }
    defer cli.Close()

    args := &argUid{uid}
    reply :=  new(setTimeRep)
    funcName  := "setLoginTime"
    if isSetLogout {
        funcName = "setLogoutTime"
    }
    if err = cli.Call(funcName, args, reply); err != nil {
        return 
    }
    t = reply.Time
    return
}

func SetLogoutTime(uid uint32) (string, error) {
    return setTime(uid, true)
}

func SetLoginTime(uid uint32) (string, error) {
    return setTime(uid, false)
}


func ModifyBalance(uid uint32, num int32) (balance uint32, err error) {
    var cli *rpc.Client
    if cli, err = jsonrpc.Dial("tcp", PY_RPC_ADDR); err != nil {
        return 
    }
    defer cli.Close()

    args := &argModifyBalance{uid, num}
    reply :=  new(modifyBalanceRep)
    if err = cli.Call("modifyBalance", args, reply); err != nil {
        return 
    }
    balance = reply.Balance
    return
}


func GetBalance(uids []uint32) (ubl map[uint32]uint32, err error) {
    var cli *rpc.Client
    if cli, err = jsonrpc.Dial("tcp", PY_RPC_ADDR); err != nil {
        return 
    }
    defer cli.Close()

    args := &argGetBalance{uids}
    reply :=  new(getBalanceRep)
    if err = cli.Call("getBalance", args, reply); err != nil {
        return 
    }
    ubl = reply.keyToUint32()
    return
}


