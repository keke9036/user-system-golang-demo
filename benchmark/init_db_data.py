import time

import pymysql.cursors

conn = pymysql.connect(
    host='127.0.0.1',
    user='root',
    password='12345678',
    database='user',
    charset='utf8mb4',
    cursorclass=pymysql.cursors.DictCursor
)

cursor = conn.cursor()

sql = """INSERT INTO user_tab (user_id, user_name, password, nick_name, avatar_url, modify_time, create_time) 
VALUES
("{}",
"testu_{}", 
"$2a$04$frEoozMyVfuxhgM2HQkb1eEp8pEnKFnICOIzt8Z0RY.X3VhOl133u", 
"test_nickname", 
"https://images2.imgbox.com/40/1e/n2bhfC9o_o.jpeg", 
{},
{});
"""

print(sql)
for i in range(1, 10000000):
    try:
        uid = i + 100000
        now = int(round(time.time() * 1000))
        tmp_sql = sql.format(uid, i, now, now)
        cursor.execute(tmp_sql)
        if not i % 10000:
            print(i)
            print(tmp_sql)
            conn.commit()

    except Exception as e:
        print(e)

# Closing the connection
conn.commit()
conn.close()
