import requests
import time
import os
import sqlite3
from urllib.parse import quote_plus
import json
from collections import defaultdict

url = "http://127.0.0.1:5000"


def get_project_folders():
    base_dir = os.path.dirname(os.path.abspath(__file__))
    parent_dir = os.path.abspath(os.path.join(base_dir, ".."))
    return sorted([f for f in os.listdir(parent_dir) if f.startswith("acc") and os.path.isdir(os.path.join(parent_dir, f))])


def check_dashboards():
    try:
        response = requests.get(f"{url}/api/get/news")
        if response.status_code == 200:
            return response.json()
        else:
            print("❌ Error fetching news:", response.status_code, response.text)
            return []
    except requests.exceptions.RequestException as e:
        print("❌ Error fetching news:", e)
        return []

def check_unused(rows_id):
    payload = {"log": "unused", "id": rows_id}
    try:
        response = requests.post(f"{url}/api/update/news", json=payload)
        if response.status_code == 200:
            print("🔄 อัปเดตสถานะ log สำเร็จ!")
    except Exception as e:
        print("❌ Error updating log:", e)
    time.sleep(1)

def get_charactor(id_prompt):
    db_path = os.path.join(os.path.dirname(os.path.abspath(__file__)), "./promt.db")
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    try:
        cursor.execute("SELECT * FROM charactor ORDER BY RANDOM() LIMIT 1")
        row = cursor.fetchone()
        if row and id_prompt < len(row):
            return [id_prompt, row[id_prompt]]
    except Exception as e:
        print(f"❌ DB error: {e}")
    finally:
        conn.close()
    return None

def generate_comment(prompt_text, topic_news):
    api_url = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:streamGenerateContent?alt=sse"
    headers = {
        "Content-Type": "application/json",
        "x-goog-api-key": "AIzaSyCIvoMfv-v54yLrgXaWu52t-L7eymSXFnA"
    }
    body = {
        "contents": [
            {
                "parts": [
                    {"text": prompt_text},
                    {"text": f"จากเนื้อหาข่าว: {topic_news}  สร้างคอมเมนต์ 10 คอมเมนต์..."}
                ],
                "role": "user"
            }
        ],
        "generationConfig": {
            "temperature": 1,
            "topP": 1,
            "topK": 1500,
            "maxOutputTokens": 8192
        }
    }
    try:
        response = requests.post(api_url, headers=headers, json=body, stream=True)
        if response.status_code != 200:
            print("❌ Gemini API error:", response.text)
            return []
        full_text = ""
        for line in response.iter_lines():
            if line and line.decode().startswith("data: "):
                json_data = json.loads(line.decode()[6:])
                parts = json_data.get("candidates", [])[0].get("content", {}).get("parts", [])
                for part in parts:
                    full_text += part.get("text", "")
        return [line.strip() for line in full_text.split("\n") if line.strip()][:10]
    except Exception as e:
        print("❌ Gemini Exception:", e)
        return []

def group_by_key(rows):
    grouped = defaultdict(list)
    for row in rows:
        key = f"{row.get('link')}|{row.get('topic')}|{row.get('reaction')}|{row.get('status')}"
        grouped[key].append(row)
    return grouped

def main():
    project_list = get_project_folders()
    print("📁 Loaded Projects:", project_list)

    news_data = check_dashboards()
    unused_rows = [row for row in news_data if row.get('log') == 'unused']
    if not unused_rows:
        print("✅ No unused rows found.")
        return

    grouped_data = group_by_key(unused_rows)

    for key, rows in grouped_data.items():
        link, topic, reaction, status = key.split('|')

        # ✅ เฉพาะกลุ่มที่ใช้คอมเมนต์
        if status in ["like_and_comment", "like_and_reply_comment", "like_reel_comment_reel"]:
            for id_prompt in range(1, 2):
                result = get_charactor(id_prompt)
                if not result:
                    continue
                _, prompt_text = result
                comments = generate_comment(prompt_text, topic)
                if not comments:
                    print(f"⚠️ ข้าม id_prompt={id_prompt} เพราะไม่ได้คอมเมนต์")
                    continue

                for i, comment_text in enumerate(comments):
                    for project in project_list:
                        # ➕ Insert ลง comment-dashboard ก่อน
                        insert_api = (
                            f"{url}/api/insert/comment-dashboard?"
                            f"comment={quote_plus(comment_text)}&"
                            f"log=unused&link={quote_plus(link)}&"
                            f"id_prompt={id_prompt}&reaction={quote_plus(reaction)}"
                        )
                        try:
                            requests.post(insert_api)

                        except Exception as e:
                            print(f"❌ Insert dashboard fail [{project}]:", e)

                        # ➕ Insert จริงลงโปรเจกต์
                        endpoint = {
                            "like_and_comment": "like-and-comment",
                            "like_and_reply_comment": "like-and-comment-reply-comment",
                            "like_reel_comment_reel": "like-reel-and-comment-reel"
                        }[status]

                        update_api = (
                            f"{url}/api/update/{project}/{endpoint}?"
                            f"reaction_type={quote_plus(reaction)}&"
                            f"link={quote_plus(link)}&"
                            f"comment_text={quote_plus(comment_text)}"
                        )
                        try:
                            requests.post(update_api)
                            print(f"✅ INSERT [{project}] → {comment_text}")
                        except Exception as e:
                            print(f"❌ Error update [{project}]:", e)
                        time.sleep(1)
                print(f"🎯 จบ id_prompt {id_prompt}")
                time.sleep(2)

        # ✅ ถ้าเป็น like-only อย่างเดียว
        elif status in ["like_only", "like_comment_only", "like_reel_only", "like_page"]:
            endpoint_map = {
                "like_only": "like-only",
                "like_comment_only": "like-comment-only",
                "like_reel_only": "like-reel-only",
                "like_page": "like-page"
            }
            endpoint = endpoint_map[status]
            for row in rows:
                for project in project_list:
                    update_api = (
                        f"{url}/api/update/{project}/{endpoint}?"
                        f"reaction_type={quote_plus(reaction)}&"
                        f"link={quote_plus(link)}"
                    )
                    try:
                        requests.post(update_api)
                        print(f"👍 LIKE [{project}] → {reaction} | {link}")
                    except Exception as e:
                        print(f"❌ Error LIKE [{project}]:", e)
                    time.sleep(1)

def main_loop():
    while True:
        try:
            main()
            time.sleep(10)
        except Exception as e:
            print("❌ เกิดข้อผิดพลาดใน main_loop:", e)
            time.sleep(5)

if __name__ == "__main__":
    main_loop()
