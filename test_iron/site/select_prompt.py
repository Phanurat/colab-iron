import time
import os
import sqlite3

def get_charactor():
    base_dir = os.path.dirname(os.path.abspath(__file__))
    char_db_path = os.path.join(base_dir, "./promt.db")
    conn = sqlite3.connect(char_db_path)
    cursor = conn.cursor()
    
    try:
        cursor.execute(f"SELECT * FROM charactors")
        row = cursor.fetchone()
        conn.close()
        return row

    except Exception as e:
        print(f"❌ Database error: {e}")
        return None
    finally:
        conn.close()

def main():
    row = get_charactor()
    for i in row:
        print(i)
    time.sleep(5)

if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        print(f"❌ Fatal Error: {e}")
    finally:
        print("✅ Finished execution.")
