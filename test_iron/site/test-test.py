import time
import os
import sqlite3

def get_charactor(id_prompt):
    base_dir = os.path.dirname(os.path.abspath(__file__))
    char_db_path = os.path.join(base_dir, "./promt.db")
    conn = sqlite3.connect(char_db_path)
    cursor = conn.cursor()
    
    try:
        cursor.execute("SELECT * FROM charactor ORDER BY RANDOM() LIMIT 1")
        row = cursor.fetchone()
        if row and id_prompt < len(row):
            prompt_text = row[id_prompt]
            return [id_prompt, prompt_text]
        else:
            print(f"⚠️ ไม่พบ column index {id_prompt} ใน row นี้")
            return None
    except Exception as e:
        print(f"❌ Database error: {e}")
        return None
    finally:
        conn.close()

def main():
    try:
        for i in range(1, 10):  # ✅ วน index 1 ถึง 5
            data_info = get_charactor(i)
            if data_info:
                id_prompt, prompt_text = data_info
                print(f"🧠 Prompt ID => {id_prompt} = {prompt_text}")
            time.sleep(1)
        id_prompt, prompt_text = data_info
        return [id_prompt, prompt_text]
    
    except Exception as e:
        print(f"❌ Error in loop: {e}")
        time.sleep(5)

if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        print(f"❌ Fatal Error: {e}")
    finally:
        print("✅ Finished execution.")
