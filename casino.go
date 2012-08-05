// L'Casino

package main

import (
    "fmt"
    "net"
    "encoding/json"
    "bytes"
    "math/rand"
    "time"
    "strconv"
    "net/http"
)


const (
    DELIMITER = byte('\n')
    BUF_SIZE = 128
    DICE_COUNT = 5
    RET_OK = 0
    RET_FL = 1
    RET_ROL = 100 // 未开
    RET_LAT = 101 // 位置被占
    RET_MAX = 200 // 人满
    GS_OPEN = 0
    GS_ROLL = 1
    MAX_PLAYER = 8 // 玩家
    MAX_USER  = 50 // 酱油
    TIME_AC = 12
    OP_INVITE = 111
    URL_INVITE = "http://appstore.yy.com/market/WebServices/AddUserApp?userId="
    XML_REP = `<?xml version="1.0"?><cross-domain-policy><allow-access-from domain="*" to-ports="*"/></cross-domain-policy>`
)

var (
    Op2handle   map[string]func(JsonString, net.Conn)
    Casino      map[uint32]*Round   // cid:*round
    Active      map[uint32]map[uint32]int64
)


type JsonString map[string]interface{}

func (m JsonString) GetRound() (*Round, bool) {
    round, ok :=  Casino[uint32(m["Cid"].(float64))]
    return round, ok
}

func (m JsonString) GetUid() (uint32) {
    return uint32(m["Uid"].(float64))
}

func (m JsonString) GetCid() (uint32) {
    return uint32(m["Cid"].(float64))
}

func (m JsonString) GetPos() (uint32) {
    return uint32(m["Pos"].(float64))
}

func (m JsonString) GetInviteuid() (uint32) {
    return uint32(m["InviteUid"].(float64))
}


type User struct {
    Conn    net.Conn
    Uid     uint32
}

type Player struct {
    User
    Pos     uint32
    Points   string
}

type Round struct {
    Cid         uint32
    Users       map[uint32]User   // uid:user
    Players     map[uint32]Player // pos:player
    Status      uint32
}

func (r *Round) GetAllUsers() (*[]User) {
    au := make([]User, 0, len(r.Users) + len(r.Players))
    for _, user := range r.Users {
       au = append(au, user)
    }
    for _, player := range r.Players {
       au = append(au, player.User)
    }
    return &au
}

func (r *Round) KickPlayer(uid uint32) {
    for pos, player := range r.Players {
        if player.Uid == uid {
           delete(r.Players, pos) 
           return
        }
    }
}

func (r *Round) Broadcast(jn []byte) {
    //fmt.Println("Broadcast", string(jn))
    for _, user := range *(r.GetAllUsers()) {
        sendMsg(user.Conn, jn)
    }
}

func (r *Round) Login(u User, cid uint32) (ret int) {
    r.Cid = cid 
    if len(r.Users) < MAX_USER {
        r.Users[u.Uid] = u
        r.KickPlayer(u.Uid)
    } else { 
        return RET_MAX
    }

    //if len(r.Players) == 0 {
    //    r.Status = GS_OPEN
    //}
    if len(r.Players) == 1 {
        for _, player := range r.Players {
            if player.Uid == u.Uid {
                r.Status = GS_OPEN
            }
        }
    }
    return RET_OK
}

func (r *Round) Join(p Player) (int) {
    ret := RET_FL
    _, exist := r.Players[p.Pos]
    if r.Status == GS_OPEN && 
            len(r.Players) < MAX_PLAYER && exist == false {
        ret = RET_OK
        r.Players[p.Pos] = p
        delete(r.Users, p.Uid)
    } else if exist {
        ret = RET_LAT
        fmt.Println("join false pos exist")
    } else if r.Status != GS_OPEN {
        ret = RET_ROL 
        fmt.Println("join false status roll")
    }

    return ret
}

func (r *Round) Roll() {
    for pos, player := range r.Players {
        if player.Conn == nil || player.Uid == 0 {
            continue
        }
        point := ""
        for i := 0; i < DICE_COUNT; i++ {
            point += strconv.Itoa(rand.Intn(6)+1)
        }

        p := r.Players[pos]
        p.Points = point
        r.Players[pos] = p // wtf
    }
    r.Status = GS_ROLL
}

func (r *Round) Open() {
    r.Status = GS_OPEN
}


func (r *Round) Logout(uid uint32) {
    //delete(r.Players, uid)
    r.KickPlayer(uid)
    if len(r.Players) == 0 {
        r.Status = GS_OPEN
    }
    delete(r.Users, uid)
}


type LoginRep struct {
    Ret     int    
    Uid     uint32
    Op      string
}

type JoinRep struct {
    Ret     int
    Uid     uint32
    Pos     uint32
    Op      string
}

type UserListRep struct {
    Ret         int
    Uid         uint32
    Cid         uint32
    Status      uint32
    UserList    []UserRep
    Op          string
}
type UserRep struct {
    Uid     uint32
    Pos     uint32
}

type RollRep struct {
    Ret     int
    Uid     uint32
    Points   string
    Op      string
}

type OpenRep struct {
    Ret         int
    Uid         uint32
    PointsList   []PointRep
    Op          string
}
type PointRep struct {
    Uid         uint32
    Points      string
}

type LogoutRep struct {
    Ret     int
    Uid     uint32
    Op      string
}

func init() {
    Op2handle = map[string]func(JsonString, net.Conn) {
        "login"         : OnLogin,
        "join"          : OnJoin,
        "user_list"     : OnUserList,
        "rolling_dice"  : OnRollingDice,
        "open_dice"     : OnOpenDice,
        "active"        : OnActive,
        "logout"        : OnLogout,
        "invite"        : OnInvite,
    }
    Casino = make(map[uint32]*Round)
    rand.Seed(int64(time.Now().Nanosecond() * time.Now().Nanosecond()))
    Active = make(map[uint32]map[uint32]int64)
}

func gameServer() {
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
        go acceptConn(conn)
    }
}

func authServer() {
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
            c.Write([]byte(XML_REP))
            c.Write([]byte("\x00"))
        }(conn)
    }
}


func httpServer() {
    http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
        status := ""
        for _, round := range Casino {
            status += fmt.Sprintf("%+v\n", *round)
        }
        fmt.Fprint(w, status)
    })
    http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
        p, u := 0, 0
        for _, round := range Casino {
            u += len(round.Users)
            p += len(round.Players)
        }
        fmt.Fprint(w, fmt.Sprintf("user:%v player:%v total:%v", u, p, p+u))
    })
    fmt.Println("httpServer runing")
    http.ListenAndServe(":20133", nil)
}


func acceptConn(c net.Conn) {
    fmt.Println("acceptConn")
    buf := make([]byte, BUF_SIZE)
    var data bytes.Buffer
    for {
        n, _ := c.Read(buf)
        if n == 0 {
            fmt.Println("close by peer")
            break
        }
        data.Write(buf[:n])

        cn := bytes.Count(data.Bytes(), []byte{DELIMITER})
        for ; cn > 0; cn-- {
            jn, err := data.ReadString(DELIMITER)
            fmt.Println("recv", jn)
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
                    dispatchOp(unknow, c)
            }
        }
    }
}


func dispatchOp(jn interface{}, c net.Conn) {
    m := JsonString(jn.(map[string]interface{}))
    Op2handle[m["Op"].(string)](m, c)
}


func sendMsg(conn net.Conn, b []byte) {
   conn.Write(append(b, DELIMITER)) 
}


func checkActive() {
    now := time.Now().Unix() 
    for cid, uid2time := range Active {
        for uid, t := range uid2time {
            if now - t > TIME_AC {
                kickUnactive(cid, uid)    
            }
        }
    }
}


func tickProcess() {
    fmt.Println("tickProcess runing")
    c := time.Tick(TIME_AC * time.Second)
    for {
        select {
            case <-c:
                checkActive()
        }
    }
}


func kickUnactive(cid, uid uint32) {
    r, ok := Casino[cid]
    if ok {
        delete(r.Users, uid)
        r.KickPlayer(uid)
        delete(Active[cid], uid)
    }

    rep := LogoutRep{0, uid, "logout"}
    jn, _ := json.Marshal(rep)
    r.Broadcast(jn)
    fmt.Println("kickUnactive", cid, uid)
}


func getRound(m JsonString) (*Round, bool) {
    round, ok :=  Casino[uint32(m["Cid"].(float64))]
    return round, ok
}

    
func OnLogin(m JsonString, c net.Conn) {
    fmt.Println("OnLogin")
    uid, cid := m.GetUid(), m.GetCid()
    r, ok := m.GetRound()
    if ok == false { // 创建
        r = &Round{0, make(map[uint32]User, MAX_USER), make(map[uint32]Player, MAX_PLAYER), 0}
        Casino[cid] = r
    } 
    user := User{c, uid}
    ret := r.Login(user, cid)

    rep := LoginRep{ret, uid, "login"}
    jn, _ := json.Marshal(rep)
    if ret == RET_OK {
        r.Broadcast(jn)
    } else {
        sendMsg(c, jn)
    }
    fmt.Println("Login rep:", string(jn))
}


func OnJoin(m JsonString, c net.Conn) {
    fmt.Println("OnJoin")
    uid := m.GetUid()
    rep := JoinRep{RET_FL, uid, 0, "join"}

    r, ok := m.GetRound()
    if ok == false {
        jn, _ := json.Marshal(rep)
        sendMsg(c, jn)
        fmt.Println("OnJoin did not found round")
        return
    }

    pos := m.GetPos()
    rep.Pos, rep.Ret = pos, r.Join(Player{User{c, uid}, pos, ""})
    jn, _ := json.Marshal(rep)
    if rep.Ret == RET_OK {
        r.Broadcast(jn)
    } else {
        sendMsg(c, jn)
    }
    fmt.Println("OnJoin rep:", rep)
}


func OnUserList(m JsonString, c net.Conn) {
    fmt.Println("OnUserList")
    uid, cid := m.GetUid(), m.GetCid()
    
    rep := UserListRep{}
    rep.Ret = RET_OK
    rep.Uid = uid 
    rep.Cid = cid

    r, ok := m.GetRound()
    if ok == false {
        rep.Ret = RET_FL
        jn, _ := json.Marshal(rep)
        sendMsg(c, jn)
        fmt.Println("OnUserList did not found round")
        return
    }

    rep.Status = r.Status
    rep.Op  = "user_list"
    rep.UserList = make([]UserRep, 0) //+
    for _, player := range r.Players {
        ur := UserRep{player.Uid, player.Pos}
        rep.UserList = append(rep.UserList, ur)
    }
    for _, user  := range r.Users {
        ur :=  UserRep{user.Uid, 0}
        rep.UserList = append(rep.UserList, ur)
    }

    jn, _ := json.Marshal(rep)
    sendMsg(c, jn)
    fmt.Println("OnUserList rep:", rep)
}


func OnRollingDice(m JsonString, c net.Conn) {
    fmt.Println("OnRollingice")

    r, ok := m.GetRound()

    if ok == false || r.Status != GS_OPEN {
        uid := m.GetUid()
        rep := RollRep{RET_FL, uid, "", "rolling_dice"}
        jn, _ := json.Marshal(rep)
        sendMsg(c, jn)
        fmt.Println("OnRollingDice fail", ok, r.Status)
        return
    }

    r.Roll()

    for _, player := range r.Players {
        rep := RollRep{RET_OK, player.Uid, player.Points, "rolling_dice"}
        jn, _ := json.Marshal(rep)
        sendMsg(player.Conn, jn)
        fmt.Println("OnRollingDice rep:", rep)
    }
}


func OnOpenDice(m JsonString, c net.Conn) {
    fmt.Println("OnOpenDice")
    uid := m.GetUid()
    
    rep := OpenRep{}
    rep.Op = "open_dice"
    rep.Uid = uid

    r, ok := m.GetRound()
    if ok == false || r.Status != GS_ROLL {
        rep.Ret = RET_FL
        jn, _ := json.Marshal(rep)
        sendMsg(c, jn)
        fmt.Println("OnOpenDice fail", ok, r.Status)
        return
    }
    rep.Ret = RET_OK
    r.Open()

    rep.PointsList = make([]PointRep, 0)
    for _, player := range r.Players {
        pr := PointRep{player.Uid, player.Points}  
        rep.PointsList = append(rep.PointsList, pr)
    }

    jn, _ := json.Marshal(rep)
    r.Broadcast(jn)
    fmt.Println("OnOpenDice rep:", rep)
}


func OnLogout(m JsonString, c net.Conn) {
    fmt.Println("OnLogout")
    uid, cid := m.GetUid(), m.GetCid()
    r, ok := m.GetRound()
    if ok {
        uid := m.GetUid()
        r.Logout(uid)
    }
    _, ok = Active[cid]
    if ok {
        delete(Active[cid], uid)
    }

    rep := LogoutRep{0, uid, "logout"}
    jn, _ := json.Marshal(rep)
    r.Broadcast(jn)
    fmt.Println("OnLogout rep:", rep)
}


func OnActive(m JsonString, c net.Conn) {
    uid, cid := m.GetUid(), m.GetCid()
    _, exist := Active[cid]
    if exist {
        Active[cid][uid] = time.Now().Unix()
    } else {
        Active[cid] = map[uint32]int64{uid : time.Now().Unix()}
    }
}

func OnInvite(m JsonString, c net.Conn) {
    fmt.Println("OnInvite")
    uid, iuid := m.GetUid(), m.GetInviteuid()
    if uid != OP_INVITE {
        return
    }
    go func (iuid uint32) {
        url := fmt.Sprintf("%v%v", URL_INVITE, iuid)
        rep, _ := http.Get(url)
        defer rep.Body.Close()
        fmt.Println("invite url:", url)
    }(iuid)
}


func main() {
    go authServer()
    go tickProcess()
    go httpServer()
    gameServer()
}

