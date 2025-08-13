from flask import Flask, request, jsonify
import sqlite3
import os
import threading
from flask import send_from_directory

app = Flask(__name__)
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
BASE_DIR = os.path.abspath(os.path.join(BASE_DIR, ".."))
db_files = {}

# ✅ Load DBs จากหลายโฟลเดอร์
def scan_dbs():
    db_files.clear()
    for folder in os.listdir(BASE_DIR):
        path = os.path.join(BASE_DIR, folder, "fb_comment_system.db")
        if os.path.isfile(path):
            db_files[folder] = path
            
@app.route('/')
def index():
    return send_from_directory('.','index.html')

@app.route('/check-acc')
def check_acc():
    return send_from_directory('.','check_acc.html')

@app.route('/log')
def log_page():
    return send_from_directory('.','log.html')

@app.route('/bio-intro')
def bio_intro():
    return send_from_directory('.','bio_intro.html')

@app.route('/change-city')
def change_city():
    return send_from_directory('.','change_city.html')

@app.route('/change-name')
def change_name():
    return send_from_directory('.', 'change_name.html')

@app.route("/proxy-info")
def proxy_info_page():
    return send_from_directory(".", "proxy_info.html")

@app.route('/switch-bio')
def switch_bio():
    return send_from_directory('.', 'switch_bio.html')

@app.route('/switch-lock')
def switch_lock():
    return send_from_directory('.', 'switch_lock.html')

@app.route('/switch-unlock')
def switch_unlock():
    return send_from_directory('.', 'switch_unlock.html')

#====================== Upload File
from werkzeug.utils import secure_filename

UPLOAD_ALLOWED_TYPES = {'caption_photo', 'cover_photo', 'profile_photo', 'story_photo'}
ALLOWED_EXTENSIONS = {'png', 'jpg', 'jpeg', 'gif'}

def allowed_file(filename):
    return '.' in filename and filename.rsplit('.', 1)[1].lower() in ALLOWED_EXTENSIONS

@app.route('/upload/<project>/<folder_name>', methods=['POST'])
def upload_project_photo(project, folder_name):
    if folder_name not in UPLOAD_ALLOWED_TYPES:
        return jsonify({'error': 'folder_name ไม่ถูกต้อง'}), 400

    project_path = os.path.join(BASE_DIR, project)
    if not os.path.exists(project_path):
        return jsonify({'error': f'ไม่พบโปรเจกต์ {project}'}), 404

    if 'file' not in request.files:
        return jsonify({'error': 'ไม่พบไฟล์'}), 400

    file = request.files['file']
    if file.filename == '':
        return jsonify({'error': 'ไม่ได้เลือกไฟล์'}), 400

    if file and allowed_file(file.filename):
        filename = secure_filename(file.filename)
        folder_path = os.path.join(project_path, folder_name)
        os.makedirs(folder_path, exist_ok=True)

        file_path = os.path.join(folder_path, filename)
        file.save(file_path)

        return jsonify({
            'status': 'อัปโหลดสำเร็จ',
            'saved_to': f'{project}/{folder_name}/{filename}'
        }), 200

    return jsonify({'error': 'ประเภทไฟล์ไม่รองรับ'}), 400

@app.route('/upload-form')
def upload_form():
    return send_from_directory('.', 'upload_form.html')


#=================================

@app.route('/api/data/<project>/<table>')
def get_table_data(project, table):
    db_path = db_files.get(project)
    if not db_path:
        return jsonify({"error": "ไม่พบโปรเจกต์"}), 404

    try:
        conn = sqlite3.connect(db_path)
        cur = conn.cursor()
        cur.execute(f"SELECT * FROM {table}")
        cols = [desc[0] for desc in cur.description]
        rows = cur.fetchall()
        conn.close()
        return jsonify({"columns": cols, "rows": rows})
    except Exception as e:
        return jsonify({"error": str(e)}), 500
    
@app.route('/api/projects')
def list_projects():
    return jsonify(list(db_files.keys()))

@app.route('/api/data/<project>')
def get_project_data(project):
    db_path = db_files.get(project)
    if not db_path:
        return jsonify({"error": "ไม่พบโปรเจกต์"}), 404

    conn = sqlite3.connect(db_path)
    cur = conn.cursor()
    cur.execute("SELECT name FROM sqlite_master WHERE type='table'")
    tables = cur.fetchall()

    result = {}
    for tname, in tables:
        try:
            cur.execute(f"SELECT * FROM {tname}")
            cols = [desc[0] for desc in cur.description]
            rows = cur.fetchall()

            # ✅ Map ชื่อ table → ชื่อที่ frontend ใช้
            mapped_keys = {
                "proxy_table": "change_proxy_table",
                "change_bio_table": "change_bio_table",
                "change_city_table": "change_city_table",
                "change_name_table": "change_name_table",
                "switch_for_bio_profile_table": "switch_for_bio_profile_table",
                "switch_for_lock_profile_table": "switch_for_lock_profile_table",
                "switch_for_unlock_profile_table": "switch_for_unlock_profile_table"
            }

            # ใส่ทั้ง key ดั้งเดิม และ key ที่ frontend ใช้
            result[tname] = {"columns": cols, "rows": rows}

            if tname in mapped_keys:
                result[mapped_keys[tname]] = {"columns": cols, "rows": rows}

        except Exception as e:
            print(f"[X] อ่านตาราง {tname} ล้มเหลว: {e}")
            continue

    conn.close()
    return jsonify(result)

@app.route('/api/update/<project>/<table>', methods=['POST'])
def update_table(project, table):
    db_path = db_files.get(project)
    if not db_path: return jsonify({"error": "ไม่พบโปรเจกต์"}), 404

    data = request.json
    threading.Thread(target=save_data_app_profile, args=(db_path, table, data)).start()
    return jsonify({"status": "updating in background"})

#test for api
@app.route('/api/update/<project>/<table>?media_id=<meta_id>', methods=['POST'])
def update_api_media(project, table, meta_id):
    db_path = db_files.get(project)
    if not db_path: return jsonify({"error": "ไม่พบโปรเจกต์"}), 404

    data = request.json
    conn = sqlite3.connect(db_path)
    cur = conn.cursor()

    cur.execute(f"DELETE FROM {table} WHERE media_id = ?", (meta_id,))
    cols = data["columns"]
    for row in data["rows"]:
        placeholders = ",".join(["?"] * len(cols))
        cur.execute(f"INSERT INTO {table} ({','.join(cols)}) VALUES ({placeholders})", row)

    conn.commit()
    conn.close()
    return jsonify({"status": "updating in background"})

def save_data_app_profile(db_path, table, data):
    try:
        conn = sqlite3.connect(db_path)
        cur = conn.cursor()

        cur.execute(f"DELETE FROM {table}")
        cols = data["columns"]
        for row in data["rows"]:
            placeholders = ",".join(["?"] * len(cols))
            cur.execute(f"INSERT INTO {table} ({','.join(cols)}) VALUES ({placeholders})", row)

        conn.commit()
        conn.close()
        print(f"[✓] อัปเดต {table} ใน {db_path}")
    except Exception as e:
        print(f"[X] Error:", e)

@app.route('/api/insert/<project>/change-bio', methods=['POST'])
def post_change_bio(project):
    db_path = db_files.get(project)
    bio_intro = request.args.get("bio_intro") or (request.json.get("bio_intro") if request.is_json else None)

    if not db_path:
        return jsonify({"error": "ไม่พบโปรเจกต์"}), 404
    if not bio_intro:
        return jsonify({"error": "กรุณาระบุ bio_intro"}), 400

    try:
        with sqlite3.connect(db_path) as conn:
            cur = conn.cursor()
            cur.execute("DELETE FROM change_bio_table")
            cur.execute("INSERT INTO change_bio_table (bio_intro) VALUES (?)", (bio_intro,))
            conn.commit()

        print(f"[✓] POST {project}: {bio_intro}")
        return jsonify({"status": "ok", "project": project, "bio_intro": bio_intro})

    except Exception as e:
        print("❌ Error:", e)
        return jsonify({"error": str(e)}), 500

@app.route('/api/get/<project>/change-bio', methods=['GET'])
def get_change_bio(project):
    db_path = db_files.get(project)
    if not db_path:
        return jsonify({"error": "ไม่พบโปรเจกต์"}), 404

    try:
        conn = sqlite3.connect(db_path)
        cur = conn.cursor()
        cur.execute("SELECT bio_intro FROM change_bio_table LIMIT 1")
        row = cur.fetchone()
        conn.close()

        return jsonify({"bio_intro": row[0] if row else None})

    except Exception as e:
        print("❌ Error:", e)
        return jsonify({"error": str(e)}), 500

@app.route('/api/insert/<project>/change-city', methods=['POST'])
def post_change_city(project):
    db_path = db_files.get(project)
    city_id = request.args.get("city_id") or (request.json.get("city_id") if request.is_json else None)

    if not db_path:
        return jsonify({"error": "ไม่พบโปรเจกต์"}), 404
    if not city_id:
        return jsonify({"error": "กรุณาระบุ city_id"}), 400

    try:
        with sqlite3.connect(db_path) as conn:
            cur = conn.cursor()

            # ✅ เช็กว่ามีตาราง change_city_table ก่อน
            cur.execute("SELECT name FROM sqlite_master WHERE type='table' AND name='change_city_table'")
            exists = cur.fetchone()
            if not exists:
                print(f"[X] ข้าม {project} เพราะไม่มี change_city_table")
                return jsonify({"error": f"{project} ไม่มี change_city_table"}), 400

            cur.execute("DELETE FROM change_city_table")
            cur.execute("INSERT INTO change_city_table (city_id) VALUES (?)", (city_id,))
            conn.commit()

        print(f"[✓] POST {project}: {city_id}")
        return jsonify({"status": "ok", "project": project, "city_id": city_id})

    except Exception as e:
        print("❌ Error:", e)
        return jsonify({"error": str(e)}), 500

@app.route('/api/get/<project>/change-city', methods=['GET'])
def get_change_city(project):
    db_path = db_files.get(project)
    if not db_path:
        return jsonify({"error": "ไม่พบโปรเจกต์"}), 404

    try:
        conn = sqlite3.connect(db_path)
        cur = conn.cursor()

        # ✅ เช็กตารางก่อน SELECT
        cur.execute("SELECT name FROM sqlite_master WHERE type='table' AND name='change_city_table'")
        exists = cur.fetchone()
        if not exists:
            return jsonify({"city_id": None})  # ไม่มีก็ส่ง null ไปเลย

        cur.execute("SELECT city_id FROM change_city_table LIMIT 1")
        row = cur.fetchone()
        conn.close()

        return jsonify({"city_id": row[0] if row else None})

    except Exception as e:
        print("❌ Error:", e)
        return jsonify({"error": str(e)}), 500

# ฟังก์ชันช่วย: อัปเดตค่าในตาราง switch ใดๆ
def update_switch_table(db_path, table_name, column_name, value):
    with sqlite3.connect(db_path) as conn:
        cur = conn.cursor()
        cur.execute(f"DELETE FROM {table_name}")
        cur.execute(f"INSERT INTO {table_name} ({column_name}) VALUES (?)", (value,))
        conn.commit()


# ✅ Update Switch
@app.route('/api/update/<project>/<switch_table>', methods=['POST'])
def update_switch(project, switch_table):
    db_path = db_files.get(project)
    if not db_path:
        return jsonify({"error": "ไม่พบโปรเจกต์"}), 404

    value = request.args.get("value") or (request.json.get("value") if request.is_json else None)
    if not value:
        return jsonify({"error": "กรุณาระบุ value"}), 400

    try:
        # ตรวจสอบ column ที่ใช้
        column_map = {
            "switch_for_bio_profile_table": "status",
            "switch_for_lock_profile_table": "status_id",
            "switch_for_unlock_profile_table": "status_id"
        }

        if switch_table not in column_map:
            return jsonify({"error": "ตารางไม่รองรับ"}), 400

        column_name = column_map[switch_table]
        update_switch_table(db_path, switch_table, column_name, value)
        print(f"[✓] Updated {switch_table} in {project} → {column_name} = {value}")
        return jsonify({"status": "ok", "table": switch_table, "value": value})

    except Exception as e:
        print("❌ Error:", e)
        return jsonify({"error": str(e)}), 500


# ✅ Get Switch Value
@app.route('/api/get/<project>/<switch_table>', methods=['GET'])
def get_switch(project, switch_table):
    db_path = db_files.get(project)
    if not db_path:
        return jsonify({"error": "ไม่พบโปรเจกต์"}), 404

    try:
        column_map = {
            "switch_for_bio_profile_table": "status",
            "switch_for_lock_profile_table": "status_id",
            "switch_for_unlock_profile_table": "status_id"
        }

        if switch_table not in column_map:
            return jsonify({"error": "ตารางไม่รองรับ"}), 400

        column_name = column_map[switch_table]

        with sqlite3.connect(db_path) as conn:
            cur = conn.cursor()
            cur.execute(f"SELECT {column_name} FROM {switch_table} LIMIT 1")
            row = cur.fetchone()

        return jsonify({column_name: row[0] if row else None})

    except Exception as e:
        print("❌ Error:", e)
        return jsonify({"error": str(e)}), 500

@app.route('/api/update/<project>/proxy-info', methods=['POST'])
def update_proxy_info(project):
    db_path = db_files.get(project)

    if not db_path:
        return jsonify({"error": "ไม่พบโปรเจกต์"}), 404

    if not request.is_json:
        return jsonify({"error": "ต้องเป็น JSON เท่านั้น"}), 400

    data = request.get_json()
    columns = data.get("columns")
    rows = data.get("rows")

    if not columns or not rows:
        return jsonify({"error": "ไม่มี columns หรือ rows"}), 400

    try:
        with sqlite3.connect(db_path) as conn:
            cur = conn.cursor()

            # ตรวจสอบตาราง
            cur.execute("SELECT name FROM sqlite_master WHERE type='table' AND name='proxy_table'")
            exists = cur.fetchone()
            if not exists:
                return jsonify({"error": f"{project} ไม่มี proxy_table"}), 400

            # เคลียร์ข้อมูลเก่า
            cur.execute("DELETE FROM proxy_table")

            # ใส่ข้อมูลใหม่ทั้งหมด
            for row in rows:
                values = tuple(row)
                placeholders = ",".join(["?"] * len(columns))
                cur.execute(f"INSERT INTO proxy_table ({','.join(columns)}) VALUES ({placeholders})", values)

            conn.commit()

        return jsonify({"status": "ok", "project": project, "rows": len(rows)})

    except Exception as e:
        print("❌ Error:", e)
        return jsonify({"error": str(e)}), 500

@app.route('/api/get/<project>/proxy-info', methods=['GET'])
def get_proxy_info(project):
    db_path = db_files.get(project)
    if not db_path:
        return jsonify({"error": "ไม่พบโปรเจกต์"}), 404

    try:
        conn = sqlite3.connect(db_path)
        cur = conn.cursor()

        cur.execute("SELECT * FROM proxy_table")
        row = cur.fetchone()
        conn.close()

        return jsonify(row)

    except Exception as e:
        print("❌ Error:", e)
        return jsonify({"error": str(e)}), 500
    
#====================================== respond log
@app.route('/api/get/<project>/respond-comment-comment', methods=['GET'])
def get_respond_comment_comment(project):
    db_path = db_files.get(project)

    if not db_path:
        return jsonify({"error": "ไม่พบโปรเจกต์"}), 404

    try:
        conn = sqlite3.connect(db_path)
        cur = conn.cursor()
        
        cur.execute("SELECT * FROM respond_for_comment_comment_table")
        rows = cur.fetchall()
        conn.close()
        
        return jsonify({"respond_for_comment_comment_table":rows[0] if rows else None})
    except Exception as e:
        print("❌ Error:", e)
        return jsonify({"error": str(e)}), 500

@app.route('/api/get/<project>/respond-comment-reel', methods=['GET'])
def get_respond_comment_reel(project):
    db_path = db_files.get(project)

    if not db_path:
        return jsonify({"error": "ไม่พบโปรเจกต์"}), 404

    try:
        conn = sqlite3.connect(db_path)
        cur = conn.cursor()
        
        cur.execute("SELECT * FROM respond_for_comment_reel_table")
        rows = cur.fetchall()
        conn.close()
        
        return jsonify({"respond_for_comment_reel_table":rows[0] if rows else None})    
    except Exception as e:
        print("❌ Error:", e)
        return jsonify({"error": str(e)}), 500

@app.route('/api/get/<project>/respond-comment', methods=['GET'])
def get_respond_comment(project):
    db_path = db_files.get(project)

    if not db_path:
        return jsonify({"error": "ไม่พบโปรเจกต์"}), 404

    try:
        conn = sqlite3.connect(db_path)
        cur = conn.cursor()
        
        cur.execute("SELECT * FROM respond_for_comment_table")
        rows = cur.fetchall()
        conn.close()
        
        return jsonify({"respond_for_comment_table":rows[0] if rows else None})
    except Exception as e:
        print("❌ Error:", e)
        return jsonify({"error": str(e)}), 500

@app.route('/api/get/<project>/respond-jewel', methods=['GET'])
def get_respond_jewel(project):
    db_path = db_files.get(project)

    if not db_path:
        return jsonify({"error": "ไม่พบโปรเจกต์"}), 404

    try:
        conn = sqlite3.connect(db_path)
        cur = conn.cursor()

        cur.execute("SELECT * FROM respond_for_jewel_table")
        rows = cur.fetchall()
        conn.close()

        return jsonify({"respond_for_jewel_table":rows[0] if rows else None})
    except Exception as e:
        print("❌ Error:", e)
        return jsonify({"error": str(e)}), 500

@app.route('/api/get/<project>/respond-like-before-comment', methods=['GET'])
def get_respond_like_before_comment(project):
    db_path = db_files.get(project)
    
    if not db_path:
        return jsonify({"error": "ไม่พบโปรเจกต์"}), 404
    
    try:
        conn = sqlite3.connect(db_path)
        cur = conn.cursor()

        cur.execute("SELECT * FROM respond_for_like_before_comment_table")
        rows = cur.fetchall()
        conn.close()
        
        return jsonify({"respond_for_like_before_comment_table":rows[0] if rows else None})
    except Exception as e:
        print("❌ Error:", e)
        return jsonify({"error": str(e)}), 500

@app.route('/api/get/<project>/respond-like-comment-only', methods=['GET'])
def get_respond_like_comment_only(project):
    db_path = db_files.get(project)

    if not db_path:
        return jsonify({"error": "ไม่พบโปรเจกต์"}), 404

    try:
        conn = sqlite3.connect(db_path)
        cur = conn.cursor()

        cur.execute("SELECT * FROM respond_for_like_comment_only_table")
        rows = cur.fetchall()
        conn.close()

        return jsonify({"respond_for_like_comment_only_table":rows[0] if rows else None})
    except Exception as e:
        print("❌ Error:", e)
        return jsonify({"error": str(e)}), 500

@app.route('/api/get/<project>/respond-like-only', methods=['GET'])
def get_respond_like_only(project):
    db_path = db_files.get(project)

    if not db_path:
        return jsonify({"error": "ไม่พบโปรเจกต์"}), 404
    
    try:
        conn = sqlite3.connect(db_path)
        cur = conn.cursor()

        cur.execute("SELECT * FROM respond_for_like_only_table")
        rows = cur.fetchall()
        conn.close()

        return jsonify({"respond_for_like_only_table":rows[0] if rows else None})
    except Exception as e:
        print("❌ Error:", e)
        return jsonify({"error": str(e)}), 500

@app.route('/api/get/<project>/respond-reel-only', methods=['GET'])
def get_respond_reel_only(project):
    db_path = db_files.get(project)

    if not db_path:
        return jsonify({"error": "ไม่พบโปรเจกต์"}), 404

    try:
        conn = sqlite3.connect(db_path)
        cur = conn.cursor()

        cur.execute("SELECT * FROM respond_for_like_reel_only_table")
        rows = cur.fetchall()
        conn.close()
        
        return jsonify({"respond_for_like_reel_only_table":rows[0] if rows else None})
    except Exception as e:
        print("❌ Error:", e)
        return jsonify({"error": str(e)}), 500

@app.route('/api/get/all-respond-log', methods=['GET'])
def get_all_respond_logs():
    result = {}
    for project, db_path in db_files.items():
        try:
            conn = sqlite3.connect(db_path)
            cur = conn.cursor()
            project_data = {}

            tables = {
                "respond_for_comment_comment_table": "respond_comment_comment",
                "respond_for_comment_reel_table": "respond_comment_reel",
                "respond_for_comment_table": "respond_comment",
                "respond_for_jewel_table": "respond_jewel",
                "respond_for_like_before_comment_table": "respond_like_before_comment",
                "respond_for_like_comment_only_table": "respond_like_comment_only",
                "respond_for_like_only_table": "respond_like_only",
                "respond_for_like_reel_only_table": "respond_reel_only"
            }

            for table_name, key_name in tables.items():
                try:
                    # ดึงทุกแถวเรียงจาก id ล่าสุดลงก่อน
                    cur.execute(f"SELECT id, respond_txt FROM {table_name} ORDER BY id DESC")
                    rows = cur.fetchall()
                    project_data[key_name] = [
                        {"id": row[0], "respond_txt": row[1]} for row in rows
                    ]
                except Exception as e:
                    project_data[key_name] = f"❌ {str(e)}"

            conn.close()
            result[project] = project_data

        except Exception as e:
            result[project] = {"error": str(e)}

    return jsonify(result)

@app.route('/api/get/check-acc', methods=['GET'])
def check_account():
    NEWS_BASE_DIR = os.path.dirname(os.path.abspath(__file__))
    NEWS_BASE_DIR = os.path.abspath(os.path.join(NEWS_BASE_DIR, ".."))
    # db_files = {}
    NEWS_BASE_DIR = "news.db"

    try:
        conn = sqlite3.connect(NEWS_BASE_DIR)
        conn.row_factory = sqlite3.Row
        cur = conn.cursor()
        cur.execute("SELECT * FROM account_dashboard")
        rows = cur.fetchall()
        conn.close()

        result = [dict(row) for row in rows]
        return jsonify(result)
    except Exception as e:
        print(f"❌ Error at /api/get/check-acc: {e}")
        return jsonify({"error": str(e)}), 500

@app.route('/api/insert/check-acc', methods=['POST'])
def insert_account():
    NEWS_BASE_DIR = os.path.dirname(os.path.abspath(__file__))
    NEWS_BASE_DIR = os.path.abspath(os.path.join(NEWS_BASE_DIR, ".."))
    # db_files = {}
    NEWS_BASE_DIR = "news.db"

    actor_id = request.args.get("actor_id")
    access_token = request.args.get("access_token")
    proxy = request.args.get("proxy")
    response = request.args.get("response")
    log = request.args.get("log")
    timestamp = request.args.get("timestamp")

    try:
        conn = sqlite3.connect(NEWS_BASE_DIR)
        cur = conn.cursor()

        cur.execute("""
            INSERT INTO account_dashboard (actor_id, access_token, proxy, response, log, timestamp)
            VALUES (?, ?, ?, ?, ?, ?)
        """, (actor_id, access_token, proxy, response, log, timestamp))
        conn.commit()
        return jsonify({"status": "ok"}), 200
    except Exception as e:
        print(f"❌ Error at /api/insert/check-acc: {e}")
        return jsonify({"error": str(e)}), 500
    finally:
        conn.close()

@app.route('/api/delete/check-acc', methods=['DELETE'])
def delete_account():
    NEWS_BASE_DIR = os.path.dirname(os.path.abspath(__file__))
    NEWS_BASE_DIR = os.path.abspath(os.path.join(NEWS_BASE_DIR, ".."))
    # db_files = {}
    NEWS_BASE_DIR = "news.db"

    try:
        conn = sqlite3.connect(NEWS_BASE_DIR)
        cur = conn.cursor()
        cur.execute("DELETE FROM account_dashboard")
        conn.commit()
        return jsonify({"DELETE": "ok"}), 200
    except Exception as e:
        print(f"❌ Error at /api/delete/check-acc: {e}")
        return jsonify({"error": str(e)}), 500
    finally:
        conn.close()
        
if __name__ == '__main__':
    scan_dbs()
    app.run(debug=True, host='0.0.0.0', port=5050)

