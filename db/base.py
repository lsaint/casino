# -*- coding: utf-8 -*-
from datetime import datetime

from sqlalchemy import *


engine = create_engine('mysql://root:111333@localhost/wanted?charset=utf8', echo = False)


if __name__ == "__main__":
    metadata = MetaData(engine)
    table_ltime = Table('ltime', metadata,
            Column('uid', Integer, primary_key=True),
            Column('login_time', DateTime()),
            Column('logout_time', DateTime()),
    )

    table_counter = Table("counter", metadata,
            Column('uid', Integer, primary_key=True),
            Column('balance', Integer(), nullable=False, default=0),
    )

    table_uname = Table("uname", metadata,
            Column('uid', Integer, primary_key=True, autoincrement=False),
            Column('name', String(50), nullable=False, default="L'"),
    )

    table_counter = Table("day_counter", metadata, 
            Column('id', Integer, primary_key=True),
            Column('uid', Integer, nullable=False),
            Column('chip', Integer, nullable=False, default=0),
            Column('date', Date(), nullable=False, default=datetime.now().date),
    )

    metadata.create_all(engine)

