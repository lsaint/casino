# -*- coding: utf-8 -*-

from sqlalchemy import *


engine = create_engine('mysql://root:111333@localhost/wanted', echo = False)


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

    metadata.create_all(engine)
