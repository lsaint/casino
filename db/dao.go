package db

import (
    //"fmt"
    "strconv"
    "net/rpc"
    "net/rpc/jsonrpc"
)


type nullArg struct {
}

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

type argSetName struct {
    Uid     uint32
    Name    string
}

type setNameRep struct {
    Name    string
}

type argGetBillboard struct {
    Uid     uint32
}

type getBillboardRep struct {
    Billboard   [][]string
}


type getWinnerRep struct {
    Today       [][]string
    Yestoday    [][]string
}

type argSetDayCounter struct {
    Uid         uint32 
    Chip        int32
}

type setDayCounterRep struct {
    Chip        int32 
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

func SetName(uid uint32, name string) (Rname string, err error) {
    var cli *rpc.Client
    if cli, err = jsonrpc.Dial("tcp", PY_RPC_ADDR); err != nil {
        return 
    }
    defer cli.Close()

    args := &argSetName{uid, name}
    reply := new(setNameRep)
    if cli.Call("setName", args, reply); err != nil {
        return
    }
    Rname = reply.Name
    return
}

func GetBillboard(uid uint32) (ret [][]string, err error) {
    var cli *rpc.Client
    if cli, err = jsonrpc.Dial("tcp", PY_RPC_ADDR); err != nil {
        return 
    }
    defer cli.Close()

    args := &argGetBillboard{uid}
    reply := new(getBillboardRep)
    if cli.Call("getBillboard", args, reply); err != nil {
        return
    }
    ret = reply.Billboard
    return
}

func SetDayCounter(uid uint32, chip int32) (ret int32, err error) {
    var cli *rpc.Client
    if cli, err = jsonrpc.Dial("tcp", PY_RPC_ADDR); err != nil {
        return
    }
    defer cli.Close()

    args := &argSetDayCounter{uid, chip}
    reply := new(setDayCounterRep)
    if cli.Call("setDayCounter", args, reply); err != nil {
        return
    }
    ret = reply.Chip
    return 
}

func GetWinner() (t [][]string, y [][]string, err error) {
    var cli *rpc.Client
    if cli, err = jsonrpc.Dial("tcp", PY_RPC_ADDR); err != nil {
        return
    }
    defer cli.Close()

    args := &nullArg{}
    reply := new(getWinnerRep)
    if cli.Call("getWinner", args, reply); err != nil {
        return
    }
    t = reply.Today
    y = reply.Yestoday
    return 
}




