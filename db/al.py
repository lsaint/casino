# -*- coding: utf-8 -*-

from sqlalchemy.ext.declarative import *
from sqlalchemy import *


from sqlalchemy.orm import sessionmaker

Base = declarative_base()

class Tb(Base):
    __tablename__='tb2'
    am_appid = Column(Integer,primary_key=True)
    ss = Column(String)

    def __init__(self, idx, s):
        self.am_appid = idx
        self.ss = s


engine = create_engine('mysql://root:111333@localhost/wanted',echo = False)
# 建表
metadata = MetaData(engine)
tb_table = Table('tb2',metadata,
        Column('am_appid', Integer, primary_key=True),
        Column('ss', String(24)),
)
metadata.create_all(engine)



Session = sessionmaker(bind = engine)
session = Session()

# add
#t1 = Tb(38, 'piano')
#session.add(t1)
#session.commit()

# update
ret = session.query(Tb).filter_by(am_appid=37).first()
ret.ss = "piano"
session.commit()

# delete
#ret = session.query(Tb).filter_by(am_appid=37).first()
#session.delete(ret)
#session.commit()
#
