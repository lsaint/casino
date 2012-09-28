# -*- coding: utf-8 -*-

import json, time
from datetime import datetime, timedelta, date

from gevent.server import StreamServer

from session import *

DELIMITER = "\n"
WINNER_CACHE_TIME = 10
WINNER_COUNT_T = 10
WINNER_COUNT_Y = 5
BILLBOARD_NUM = 10
END_DATE = date(2012, 10, 11)



g_uid2name = {}

def cacheName():
    global g_uid2name
    ret = session.query(Uname)
    for q in ret:
        g_uid2name[q.uid] = q.name
cacheName()



def accept(socket, address):
    print 'New connection from %s:%s' % address
    buf = ""
    while True:
        buf =  "%s%s" % (buf, socket.recv(1024))
        if not buf:
            break
        if DELIMITER not in buf:
            continue
        la = buf.rindex(DELIMITER)
        buf, lines = buf[la+1:], buf[:la]
        lt = lines.splitlines()
        for line in lt:
            jn = json.loads(line)
            if type(jn) == dict:
                dispatch(jn, socket)


def dispatch(jn, socket):
    print "req", jn
    kwargs = jn["params"][0]
    method = eval(jn["method"])
    reply, err = method(**kwargs)
    jn = {"result":reply, "error":err, "id":jn["id"]}
    print "reply", jn
    socket.send(json.dumps(jn))



def setLogoutTime(**kwargs):
    print "setLogoutTime"
    uid = kwargs["Uid"]
    now = datetime.now()
    ret = session.query(Ltime).filter_by(uid=uid).first()
    if ret:
        ret.logout_time = now
    else:
        t = Ltime(uid, None, now)
        session.add(t)
    session.commit()
    #rep = {"Op":"setlogouttime", "Ret":0, "Uid":uid, "Time":now.ctime()}
    return {"Time":now.ctime()}, None


def setLoginTime(**kwargs):
    print "setLoginTime"
    uid = kwargs["Uid"]
    now = datetime.now()
    ret = session.query(Ltime).filter_by(uid=uid).first()
    if ret:
        ret.login_time = now
    else:
        t = Ltime(uid, now, None)
        session.add(t)
    session.commit()
    #rep = {"Op":"setlogintime", "Ret":0, "Uid":uid, "Time":now.ctime()}
    return {"Time":now.ctime()}, None


def getLogTime(**kwargs):
    print "getLogTime"
    uid = kwargs["Uid"]
    ret = session.query(Ltime).filter_by(uid=uid).first()
    if ret:
        intime = None if not ret.login_time else ret.login_time.ctime()
        outime = None if not ret.logout_time else ret.logout_time.ctime()
    else:
        intime = None
        outime = None
    #rep = {"Op":"getlogtime", "Ret":0, "Uid":uid, "Logintime":intime, "Logouttime":outime}
    return {"Logintime":intime, "Logouttime":outime}, None


def modifyBalance(**kwargs):
    print "modifyBalance"
    uid = kwargs["Uid"]
    n = kwargs["Num"]
    ret = session.query(Counter).filter_by(uid=uid).first()
    if ret:
        if n < 0 and ret.balance < -n:
            return {"Balance":None}, False
        ret.balance += n
        if ret.balance < 0:
            ret.balance = 0
    else:
        if n < 0:
            return {"Balance":None}, False
        ret = Counter(uid, 0)
        if n > 0:
            ret.balance = n
        session.add(ret)

    session.commit()
    #rep = {"Op":"modifybalance", "Ret":0, "Uid":uid, "Balance":ret.balance}
    return {"Balance":ret.balance}, None


def getBalance(**kwargs):
    print "getBalance"
    uids = kwargs["Uid"]
    ret = session.query(Counter).filter(Counter.uid.in_(uids)).all()
    #rep = {"Op":"getbalance", "Ret":0, "Ubl":None}
    bals = {}
    for r in ret:
        bals[str(r.uid)] = r.balance
    return {"Ubl":bals}, None


def setName(**kwargs):
    uid = kwargs["Uid"]
    name = kwargs["Name"]
    g_uid2name[uid] = name
    ret = session.query(Uname).get(uid)
    if ret:
        if name != ret.name:
            ret.name = name
    else:
        t = Uname(uid, name)
        session.add(t)
    session.commit()
    return {"Name":name}, None



def getBillboard(**kwargs):
    print "getBillboard"
    ret = session.query(Counter).order_by("balance desc")[:BILLBOARD_NUM]
    uids = []
    bals = []
    for i in ret:
        uids.append(i.uid)
        bals.append(i.balance)

    #un = session.query(Uname).filter(Uname.uid.in_(uids))
    #dt = {}
    #for i in un:
    #    dt[i.uid] = i.name

    ret = []
    for i in range(len(uids)):
        name =  g_uid2name.get(uids[i])
        name = name if name else ""
        ret.append((name, str(bals[i])))

    print json.dumps(ret)
    return {"Billboard":ret}, None


def setDayCounter(**kwargs):
    uid, chip = kwargs["Uid"], kwargs["Chip"]
    now = datetime.now().date()
    ret = session.query(DayCounter).filter_by(uid=uid, date=now).first()
    if ret:
       ret.chip += chip
    else:
       ret = DayCounter(uid, chip, now) # now!?
       session.add(ret)
    session.commit()

    return {"Chip":ret.chip}, None


g_cacheWinner = {"time":0, "y":[], "t":[]}
def getWinner(**kwargs):
    print "getWinners"
    global g_cacheWinner
    now = time.time()
    if now - g_cacheWinner["time"] < WINNER_CACHE_TIME:
        return {"Today":g_cacheWinner["t"], "Yestoday":g_cacheWinner["y"]}, None

    g_cacheWinner["time"] = now
    ret_t = session.query(DayCounter).filter_by(date=datetime.now().date()).order_by("chip desc")[:WINNER_COUNT_T]
    ret_y = session.query(DayCounter).filter_by(date=datetime.now().date() - timedelta(1)).order_by("chip desc")[:WINNER_COUNT_Y]

    t = []
    for q in ret_t:
        name = g_uid2name.get(q.uid) or ""
        t.append((name, str(q.chip)))
    g_cacheWinner["t"] = t

    y = []
    if datetime.now().date() <= END_DATE:
        for q in ret_y:
            name = g_uid2name.get(q.uid) or ""
            y.append((name, str(q.chip)))
        g_cacheWinner["y"] = y

    return {"Today":t, "Yestoday":y}, None



if __name__ == '__main__':
    server = StreamServer(('127.0.0.1', 12918), accept)
    print 'Starting server on port 12918'
    server.serve_forever()


