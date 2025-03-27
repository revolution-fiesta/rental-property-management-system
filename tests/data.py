import psycopg2

def Connect():
  return psycopg2.connect("dbname=postgres user=postgres password=270153 host=localhost port=5433")

def InsertData(conn):
  cur = conn.cursor()
  with open("data.pg.sql", "r", encoding='utf-8') as f:
      cur.execute(f.read())
  conn.commit()
  cur.close()

if __name__ == "__main__":
  conn = Connect()
  InsertData(conn)
  conn.close()