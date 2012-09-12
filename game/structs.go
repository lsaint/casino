package game


type TimeRep struct {
    Ret         int
    Op          string
    Time        int64
}

type LoginRep struct {
    Ret         int    
    Uid         uint32
    RewardCoin  uint32
    LostCoin    uint32
    Coin        uint32
    Op          string
}

type JoinRep struct {
    Ret     int
    Uid     uint32
    Pos     uint32
    Coin    uint32
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
    Coin    uint32
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
    Coin        uint32
}

type LogoutRep struct {
    Ret     int
    Uid     uint32
    Op      string
}

type ProxyInvite struct {
    Op      string
    Sid     uint32
    Subsid  uint32
}

type GiveCoinRep struct {
    Ret         int
    Uid         uint32
    TargetUid   uint32
    Coin        int32
    Op          string
}

