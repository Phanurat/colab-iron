import os
import sqlite3
import threading

BASE_DIR = os.path.dirname(os.path.abspath(__file__))
BASE_DIR = os.path.abspath(os.path.join(BASE_DIR, ".."))
db_files = {}

# การแมป prompt_id ตามช่วงของโปรเจกต์
def get_charactor():
    base_dir = os.path.dirname(os.path.abspath(__file__))
    char_db_path = os.path.join(base_dir, "db/promt.db")
    conn = sqlite3.connect(char_db_path)
    cursor = conn.cursor()

    try:
        # ดึงข้อมูลทั้งหมดจากตาราง charactors
        cursor.execute("SELECT * FROM charactors")
        rows = cursor.fetchall()

        # แปลงผลลัพธ์ให้เป็น list ของ id_charactor และข้อมูลที่เกี่ยวข้อง
        id_charactors = [row[0] for row in rows]
        title_charactor = [row[1] for row in rows]
        prompt_content = [row[2] for row in rows]
        log = [row[3] for row in rows]

        return id_charactors, title_charactor, prompt_content, log
    except Exception as e:
        print(f"❌ Database error: {e}")
        return None
    finally:
        conn.close()

def clear_db():
    base_dir = os.path.dirname(os.path.abspath(__file__))
    char_db_path = os.path.join(base_dir, "db/promt.db")
    conn = sqlite3.connect(char_db_path)
    cursor = conn.cursor()

    try:
        cursor.execute("DELETE FROM mapp_prompt")
        conn.commit()
    except Exception as e:
        print(f"❌ Database error: {e}")
    finally:
        conn.close()

def mapp_project_dashboard(folder, char_id, title, content):
    base_dir = os.path.dirname(os.path.abspath(__file__))
    char_db_path = os.path.join(base_dir, "db/promt.db")
    conn = sqlite3.connect(char_db_path)
    cursor = conn.cursor()
    project_id = folder
    prompt_id = char_id
    prompt_rolplay = title
    prompt_content = content

    try:
        cursor.execute("INSERT INTO mapp_prompt (project_id, prompt_id, prompt_rolplay, prompt_content) VALUES (?, ?, ?, ?)", (project_id, prompt_id, prompt_rolplay, prompt_content))
        print("Save Mapping Project!")
        conn.commit()
    except Exception as e:
        print(f"❌ Database error: {e}")
    finally:
        conn.close()

# ฟังก์ชัน scan_dbs
def scan_dbs(id_charactors, title_charactor, prompt_content, log):
    clear_db()
    for idx, folder in enumerate(os.listdir(BASE_DIR)):
        path = os.path.join(BASE_DIR, folder, "fb_comment_system.db")
        if os.path.isfile(path):
            db_files[folder] = path
            # หาช่วงโปรเจกต์และแมป id_charactors
            project_id = int(folder.replace("acc", ""))  # แยกเลขโปรเจกต์
            prompt_idx = (project_id - 1) // 10  # แบ่งโปรเจกต์เป็นช่วงๆ ละ 10
            
            # ดึงข้อมูลที่ตรงกับ prompt_idx จากแต่ละลิสต์
            char_id = id_charactors[prompt_idx % len(id_charactors)]
            title = title_charactor[prompt_idx % len(title_charactor)]
            content = prompt_content[prompt_idx % len(prompt_content)]
            # content = "Test"
            project_log = log[prompt_idx % len(log)]

            # พิมพ์ข้อมูลที่เชื่อมโยงกับโปรเจกต์
            print("======================================================================================================================================")
            print(f"Project : {folder} -> id_charactor : {char_id} -> Title: {title} -> Prompt Content: {content} -> Log: {project_log}")
            # title_charactor
            #db/promp.db | id_charactor | title_charactor | prompt_content | log | mapp_project
            #
            mapp_project_dashboard(folder, char_id, title, content)
            print("======================================================================================================================================")

# ดึงข้อมูลจากฐานข้อมูล
id_charactors, title_charactor, prompt_content, log = get_charactor()

# เรียกใช้ฟังก์ชัน scan_dbs
scan_dbs(id_charactors, title_charactor, prompt_content, log)
