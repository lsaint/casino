# -*- coding: utf-8 -*-

import json

from gevent.server import StreamServer

from session import *

import sys
sys.stdout = sys.stderr

DELIMITER = "\n"



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
    print jn
    op  = jn.get("Op")
    if op:
        handles[op](jn, jn["Uid"],socket)


def setLogoutTime(jn, uid, s):
    print "setLogoutTime"
    now = datetime.now()
    ret = session.query(Ltime).filter_by(uid=uid).first()
    if ret:
        ret.logout_time = now
    else:
        t = Ltime(uid, None, now)
        session.add(t)
    session.commit()
    rep = {"Op":"setlogouttime", "Ret":0, "Uid":uid, "Time":now.ctime()}
    sendMsg(s, rep)


def setLoginTime(jn, uid, s):
    print "setLoginTime"
    now = datetime.now()
    ret = session.query(Ltime).filter_by(uid=uid).first()
    if ret:
        ret.login_time = now
    else:
        t = Ltime(uid, now, None)
        session.add(t)
    session.commit()
    rep = {"Op":"setlogintime", "Ret":0, "Uid":uid, "Time":now.ctime()}
    sendMsg(s, rep)


def getLogTime(jn, uid, s):
    print "getLogTime"
    ret = session.query(Ltime).filter_by(uid=uid).first()
    if ret:
        intime = None if not ret.login_time else ret.login_time.ctime()
        outime = None if not ret.logout_time else ret.logout_time.ctime()
    else:
        intime = None
        outime = None
    rep = {"Op":"getlogtime", "Ret":0, "Uid":uid, "Logintime":intime, "Logouttime":outime}
    sendMsg(s, rep)



def modifyBalance(jn, uid, s):
    print "modifyBalance"
    n = jn["Num"]
    ret = session.query(Counter).filter_by(uid=uid).first()
    if ret:
        ret.balance += n
        if ret.balance < 0:
            ret.balance = 0
    else:
        ret = Counter(uid, 0)
        if n > 0:
            ret.balance = n
        session.add(ret)
    session.commit()

    rep = {"Op":"modifybalance", "Ret":0, "Uid":uid, "Balance":ret.balance}
    sendMsg(s, rep)


def getBalance(jn, uids, s):
    print "getBalance"
    ret = session.query(Counter).filter(Counter.uid.in_(uids)).all()
    rep = {"Op":"getbalance", "Ret":0, "Ubl":None}
    bals = {}
    for r in ret:
        bals[str(r.uid)] = r.balance
    rep["Ubl"] = bals
    sendMsg(s, rep)


handles = {
        "setlogouttime" : setLogoutTime,
        "setlogintime"  : setLoginTime,
        "getlogtime"    : getLogTime,
        "modifybalance" : modifyBalance,
        "getbalance"    : getBalance,
}

def sendMsg(socket, dt):
    socket.send(json.dumps(dt) + "\n")


if __name__ == '__main__':
    server = StreamServer(('127.0.0.1', 12918), accept)
    print 'Starting server on port 12918'
    server.serve_forever()

