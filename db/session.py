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


class Counter(Base):
    __tablename__ = "counter"
    uid = Column(Integer, primary_key=True)
    balance = Column(Integer, nullable=False)

    def __init__(self, uid, bal):
        self.uid = uid
        self.balance = bal


class Uname(Base):
    __tablename__ = 'uname'
    uid = Column(Integer, primary_key=True)
    name = Column(String, nullable=False)

    def __init__(self, uid, name):
        self.uid = uid
        self.name = name


class DayCounter(Base):
    __tablename__ = 'day_counter'
    id = Column(Integer, primary_key=True)
    uid = Column(Integer, nullable=False)
    chip = Column(Integer, nullable=False)
    date = Column(Date, nullable=False, default=datetime.now().date)

    def __init__(self, uid, chip, date=datetime.now().date()):
        self.uid = uid
        self.chip = chip
        self.date = date


#ret = session.query(Ltime).filter_by(uid=111).first()
#print ret.logout_time.ctime()

