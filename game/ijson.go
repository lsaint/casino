package game

import (
    "fmt"
    "strconv"
)


type JsonString map[string]interface{}

func (m JsonString) GetRound(c RoundMap) (*Round, bool) {
    round, ok :=  c[uint32(m["Cid"].(float64))]
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

func (m JsonString) GetTime() (int64) {
    s1 := m["Time"].(string)[1:]
    p1, p2, p3 := s1[:3], s1[3:7], s1[7:]
    s2 := fmt.Sprintf("%v%v%v", p3, p2, p1)
    t2, _ := strconv.ParseInt(s2, 10, 64)
    return t2 - 20050411
}

func (m JsonString) GetRootChlId() (uint32) {
    return uint32(m["RootChannelId"].(float64)) 
}

func (m JsonString) GetCoin() (int32) {
    return int32(m["Coin"].(float64)) 
}

func (m JsonString) GetTargetUid() (uint32) {
    return uint32(m["TargetUid"].(float64)) 
}

