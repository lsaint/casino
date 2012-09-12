
package common

const (
    DELIMITER = byte('\n')
    BUF_SIZE = 128
    DICE_COUNT = 5
    RET_OK = 0
    RET_FL = 1
    RET_ROL = 100 // 未开
    RET_LAT = 101 // 位置被占
    RET_MAX = 200 // 人满
    RET_TME = 300 // 连接超时
    GS_OPEN = 0
    GS_ROLL = 1
    MAX_PLAYER = 8 // 玩家
    MAX_USER  = 50 // 酱油
    TIME_AC = 12
    OP_INVITE = 300
    ENDURE_SEC = 10
    LOGIN_AWARD = 500
    DAY_LOST = 10
    URL_INVITE = "http://appstore.yy.com/market/WebServices/AddUserApp?userId="
    XML_REP = `<?xml version="1.0"?><cross-domain-policy><allow-access-from domain="*" to-ports="*"/></cross-domain-policy>`
)
