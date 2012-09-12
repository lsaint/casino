package game


import (
    "fmt"
    "time"
    "encoding/json"

    "casino/common"
    "casino/db"
)


type Ticker struct {
    Active      map[uint32]map[uint32]int64
    Casino      RoundMap
}


func (ticker *Ticker) Start(ca RoundMap) {
    fmt.Println("tickProcess runing")
    ticker.Casino = ca
    c := time.Tick(common.TIME_AC * time.Second)
    for {
        select {
            case <-c:
                ticker.checkActive()
        }
    }
}


func (ticker *Ticker) checkActive() {
    now := time.Now().Unix() 
    for cid, uid2time := range ticker.Active {
        for uid, t := range uid2time {
            if now - t > common.TIME_AC {
                ticker.kickUnactive(cid, uid)    
            }
        }
    }
}


func (ticker *Ticker) kickUnactive(cid, uid uint32) {
    r, ok := ticker.Casino[cid]
    if ok {
        delete(r.Users, uid)
        r.KickPlayer(uid)
        delete(ticker.Active[cid], uid)
        db.Dao.SetLogoutTime(uid)
    }

    rep := LogoutRep{0, uid, "logout"}
    jn, _ := json.Marshal(rep)
    r.Broadcast(jn)
    fmt.Println("kickUnactive", cid, uid)
}

