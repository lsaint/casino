package game

import (
    "fmt"
    "net"
    "strconv"
    "math/rand"
    "casino/common"
)

type RoundMap map[uint32]*Round


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

func (r *Round) Broadcast(inf interface{}) {
    //fmt.Println("Broadcast", string(jn))
    for _, user := range *(r.GetAllUsers()) {
        SendMsg(user.Conn, inf)
    }
}

func (r *Round) Login(u User, cid uint32) (ret int) {
    r.Cid = cid 
    if len(r.Users) < common.MAX_USER {
        r.Users[u.Uid] = u
        r.KickPlayer(u.Uid)
    } else { 
        return common.RET_MAX
    }

    if len(r.Players) == 1 {
        for _, player := range r.Players {
            if player.Uid == u.Uid {
                r.Status = common.GS_OPEN
            }
        }
    }
    return common.RET_OK
}

func (r *Round) Join(p Player) (int) {
    ret := common.RET_FL
    _, exist := r.Players[p.Pos]

    for pos, player := range r.Players {
        if p.Uid == player.Uid {
            delete(r.Players, pos)
        }
    }

    if r.Status == common.GS_OPEN && 
            len(r.Players) < common.MAX_PLAYER && exist == false {
        ret = common.RET_OK
        r.Players[p.Pos] = p
        delete(r.Users, p.Uid)
    } else if exist {
        ret = common.RET_LAT
        fmt.Println("join false pos exist")
    } else if r.Status != common.GS_OPEN {
        ret = common.RET_ROL 
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
        for i := 0; i < common.DICE_COUNT; i++ {
            point += strconv.Itoa(rand.Intn(6)+1)
        }

        p := r.Players[pos]
        p.Points = point
        r.Players[pos] = p // wtf
    }
    r.Status = common.GS_ROLL
}

func (r *Round) Open() {
    r.Status = common.GS_OPEN
}


func (r *Round) Logout(uid uint32) {
    //delete(r.Players, uid)
    r.KickPlayer(uid)
    if len(r.Players) == 0 {
        r.Status = common.GS_OPEN
    }
    delete(r.Users, uid)
}
