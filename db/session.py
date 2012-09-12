# -*- coding: utf-8 -*-

from datetime import datetime

from sqlalchemy.ext.declarative import *
from sqlalchemy.orm import sessionmaker
from sqlalchemy import *

from base import engine


S = sessionmaker(bind = engine)
session = S()


Base = declarative_base()

class Ltime(Base):
    __tablename__ = 'ltime'
    uid = Column(Integer, primary_key=True)
    login_time = Column(DateTime)
    logout_time = Column(DateTime)

    def __init__(self, uid, intime, outtime):
        self.uid = uid
        self.login_time = intime
        self.logout_time = outtime

    def logoutNow(self):
        self.logout_time = datetime.now()
        session.commit()



class Counter(Base):
    __tablename__ = "counter"
    uid = Column(Integer, primary_key=True)
    balance = Column(Integer, nullable=False)

    def __init__(self, uid, bal):
        self.uid = uid
        self.balance = bal

    def increase(self, num):
        self.balance += num
        session.commit()

    def decrease(self, num):
        self.balance -= num
        if self.balance < 0:
            self.balance = 0
        session.commit()



#ret = session.query(Ltime).filter_by(uid=111).first()
#print ret.logout_time.ctime()

