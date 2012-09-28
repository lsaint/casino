// L'Casino

package main

import (
    "fmt"
    "net"
    "math/rand"
    "time"
    "net/http"

    "casino/game"
    "casino/common"
    "casino/db"
)


var (
    //Op2handle   map[string]func(game.JsonString, net.Conn)
    gameServer          game.GameServer
    httpServer          game.HttpServer 
    authServer          game.AuthServer
    Casino              game.RoundMap // cid:*round
    //Active      map[uint32]map[uint32]int64
    ticker              game.Ticker
)




func init() {
    gameServer.Op2handle = map[string]func(game.JsonString, net.Conn) {
        "server_time"   : OnGetTime,
        "login"         : OnLogin,
        "join"          : OnJoin,
        "user_list"     : OnUserList,
        "rolling_dice"  : OnRollingDice,
        "open_dice"     : OnOpenDice,
        "active"        : OnActive,
        "logout"        : OnLogout,
        "invite"        : OnInvite,
        "send_invite"   : OnSendChlInvite,
        "give_coin"     : OnGiveCoin,
        "get_billboard" : OnGetBillboard,
        "get_winner"    : OnGetWinner,
    }
    Casino = make(map[uint32]*game.Round)
    rand.Seed(int64(time.Now().Nanosecond() * time.Now().Nanosecond()))
    ticker.Active = make(map[uint32]map[uint32]int64)
}


func OnLogin(m game.JsonString, c net.Conn) {
    fmt.Println("OnLogin")
    uid, cid, t, name := m.GetUid(), m.GetCid(), m.GetTime(), m.GetName()
    diff := time.Now().Unix() - t

    if  diff >= common.ENDURE_SEC || diff < 0 {
        rep := game.LoginRep{common.RET_TME, uid, 0, 0, 0, "login"}
        game.SendMsg(c, rep)
        return
    }

    award, lost := TimeCheck(uid)
    fmt.Println("time check", award, lost)
    bal, _ := db.ModifyBalance(uid, int32(award-lost))
    //bal := mod_ret.Balance
    db.SetLoginTime(uid)
    db.SetName(uid, name)

    r, ok := m.GetRound(Casino)
    if ok == false { // 创建
        r = &game.Round{0, make(map[uint32]game.User, common.MAX_USER), make(map[uint32]game.Player, common.MAX_PLAYER), 0}
        Casino[cid] = r
    } 
    user := game.User{c, uid}
    ret := r.Login(user, cid)

    rep := game.LoginRep{ret, uid, award, lost, bal, "login"}
    if ret == common.RET_OK {
        r.Broadcast(rep)
    } else {
        game.SendMsg(c, rep)
    }
    fmt.Println("Login rep:", rep)
}

func TimeCheck(uid uint32) (award, lost uint32){
    intime, outime, _ := db.GetLogTime(uid)
    fmt.Println("intime", intime, "outtime", outime)

    now := time.Now()
    in, _ := time.Parse(time.ANSIC, intime)
    ou, _ := time.Parse(time.ANSIC, outime)

    // check same day
    a1, b1, c1 := now.Date()
    a2, b2, c2 := in.Date()
    if intime == "" || a1 != a2 || b1 != b2 || c1 != c2 {
        award = common.LOGIN_AWARD   // +award
    }

    // check unactive day
    h := now.Sub(ou).Hours()
    if outime == "" || h < 24 {
        lost = 0
    } else {
        lost = uint32(h) / 24 * common.DAY_LOST
    }
    return
}


func OnJoin(m game.JsonString, c net.Conn) {
    fmt.Println("OnJoin")
    uid := m.GetUid()
    rep := game.JoinRep{common.RET_FL, uid, 0, 0, "join"}

    r, ok := m.GetRound(Casino)
    if ok == false {
        game.SendMsg(c, rep)
        fmt.Println("OnJoin did not found round")
        return
    }

    pos := m.GetPos()
    rep.Pos, rep.Ret = pos, r.Join(game.Player{game.User{c, uid}, pos, ""})
    if rep.Ret == common.RET_OK {
        bal, _ := db.GetBalance([]uint32{uid})
        rep.Coin = bal[uid]
        r.Broadcast(rep)
    } else {
        game.SendMsg(c, rep)
    }
    fmt.Println("OnJoin rep:", rep)
}


func OnUserList(m game.JsonString, c net.Conn) {
    fmt.Println("OnUserList")
    uid, cid := m.GetUid(), m.GetCid()
    
    rep := game.UserListRep{}
    rep.Ret = common.RET_OK
    rep.Uid = uid 
    rep.Cid = cid

    r, ok := m.GetRound(Casino)
    if ok == false {
        rep.Ret = common.RET_FL
        game.SendMsg(c, rep)
        fmt.Println("OnUserList did not found round")
        return
    }

    rep.Status = r.Status
    rep.Op  = "user_list"
    rep.UserList = make([]game.UserRep, 0) //+

    ubl := GetPlayerBalance(r.Players)
    for _, player := range r.Players {
        ur := game.UserRep{player.Uid, player.Pos, ubl[player.Uid]}
        rep.UserList = append(rep.UserList, ur)
    }
    for _, user  := range r.Users {
        ur :=  game.UserRep{user.Uid, 0, 0}
        rep.UserList = append(rep.UserList, ur)
    }

    game.SendMsg(c, rep)
    fmt.Println("OnUserList rep:", rep)
}


func GetPlayerBalance(players map[uint32]game.Player) (ret map[uint32]uint32) {
    l := len(players)
    if l == 0 {
        ret = make(map[uint32]uint32)
    } else {
        uids := make([]uint32, l)
        i := 0
        for _, player := range players {
             uids[i] = player.Uid   
             i++
        }
        ret, _ = db.GetBalance(uids)
    }
    return
}


func OnRollingDice(m game.JsonString, c net.Conn) {
    fmt.Println("OnRollingice")

    r, ok := m.GetRound(Casino)

    if ok == false || r.Status != common.GS_OPEN {
        uid := m.GetUid()
        rep := game.RollRep{common.RET_FL, uid, "", "rolling_dice"}
        game.SendMsg(c, rep)
        fmt.Println("OnRollingDice fail", ok, r.Status)
        return
    }

    r.Roll()

    for _, player := range r.Players {
        rep := game.RollRep{common.RET_OK, player.Uid, player.Points, "rolling_dice"}
        game.SendMsg(player.Conn, rep)
        fmt.Println("OnRollingDice rep:", rep)
    }

    for _, user := range r.Users {
        rep := game.RollRep{common.RET_OK, user.Uid, "", "rolling_dice"}
        game.SendMsg(user.Conn, rep)
    }
}


func OnOpenDice(m game.JsonString, c net.Conn) {
    fmt.Println("OnOpenDice")
    uid := m.GetUid()
    
    rep := game.OpenRep{}
    rep.Op = "open_dice"
    rep.Uid = uid

    r, ok := m.GetRound(Casino)
    if ok == false || r.Status != common.GS_ROLL {
        rep.Ret = common.RET_FL
        game.SendMsg(c, rep)
        fmt.Println("OnOpenDice fail", ok, r.Status)
        return
    }
    rep.Ret = common.RET_OK
    r.Open()

    ubl := GetPlayerBalance(r.Players)
    rep.PointsList = make([]game.PointRep, 0)
    for _, player := range r.Players {
        pr := game.PointRep{player.Uid, player.Points, ubl[player.Uid]}  
        rep.PointsList = append(rep.PointsList, pr)
    }

    r.Broadcast(rep)
    fmt.Println("OnOpenDice rep:", rep)
}


func OnLogout(m game.JsonString, c net.Conn) {
    fmt.Println("OnLogout")
    uid, cid := m.GetUid(), m.GetCid()
    r, ok := m.GetRound(Casino)
    if ok {
        uid := m.GetUid()
        r.Logout(uid)
    }
    _, ok = ticker.Active[cid]
    if ok {
        delete(ticker.Active[cid], uid)
    }
    db.SetLogoutTime(uid)

    rep := game.LogoutRep{0, uid, "logout"}
    r.Broadcast(rep)
    fmt.Println("OnLogout rep:", rep)
}


func OnActive(m game.JsonString, c net.Conn) {
    //fmt.Println("OnActive")
    uid, cid := m.GetUid(), m.GetCid()
    _, exist := ticker.Active[cid]
    if exist {
        ticker.Active[cid][uid] = time.Now().Unix()
    } else {
        ticker.Active[cid] = map[uint32]int64{uid : time.Now().Unix()}
    }
}

func OnInvite(m game.JsonString, c net.Conn) {
    fmt.Println("OnInvite")
    uid, iuid := m.GetUid(), m.GetInviteuid()
    if iuid != common.OP_INVITE {
        return
    }
    go func (uid uint32) {
        url := fmt.Sprintf("%v%v", common.URL_INVITE, uid)
        rep, _ := http.Get(url)
        defer rep.Body.Close()
        //fmt.Println("invite url:", url)
    }(uid)
}

func OnGetTime(m game.JsonString, c net.Conn) {
    fmt.Println("OnGetTime")
    rep := game.TimeRep{0, "server_time", time.Now().Unix()}
    game.SendMsg(c, rep)
}


func OnSendChlInvite(m game.JsonString, c net.Conn) {
    fmt.Println("OnSendChlInvite")
    go func () {
        conn, err := net.Dial("tcp", "127.0.0.1:37771")
        if err != nil {
            fmt.Println("conn invite proxy err")
        } else {
            msg := game.ProxyInvite{"Doinvite", m.GetRootChlId(), m.GetCid()}
            game.SendMsg(conn, msg)
        }
    }()
}


func OnGiveCoin(m game.JsonString, c net.Conn) {
    fmt.Println("OnGiveCoin")
    coin, uid, tuid := m.GetCoin(), m.GetUid(), m.GetTargetUid()
    r, ok := m.GetRound(Casino)
    if ok == false {
        rep := game.GiveCoinRep{common.RET_FL, uid, tuid, coin, "give_coin"}    
        game.SendMsg(c, rep)
        return
    }

    ret := common.RET_OK
    _, err1 := db.ModifyBalance(uid, -coin)
    _, err2 := db.ModifyBalance(tuid, coin)
    db.SetDayCounter(uid, -coin)
    db.SetDayCounter(tuid, coin)
    if err1 != nil || err2 != nil {
        ret = common.RET_FL 
    }
    rep := game.GiveCoinRep{ret, uid, tuid, coin, "give_coin"}    
    r.Broadcast(rep)
}


func OnGetBillboard(m game.JsonString, c net.Conn) {
    fmt.Println("OnGetBillboard")
    uid := m.GetUid()

    rep := game.GetBillboardRep{Ret:common.RET_FL, Uid:uid, Op:"get_billboard"}
    bb, err := db.GetBillboard(uid)
    if err == nil {
        rep.Ret = common.RET_OK
        rep.Billboard = bb
    }
    game.SendMsg(c, rep)
}


func OnGetWinner(m game.JsonString, c net.Conn) {
    fmt.Println("OnGetWinner")

    rep := game.GetWinnerRep{Ret:common.RET_FL, Op:"get_winner"} 
    t, y, err := db.GetWinner()
    if err == nil {
        rep.Ret = common.RET_OK
        rep.Today = t
        rep.Yestoday = y
    }
    game.SendMsg(c, rep)
}


func main() {
    go authServer.Start()
    go ticker.Start(Casino)
    go httpServer.Start(Casino)
    gameServer.Start()
}

