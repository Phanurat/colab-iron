import requests
import time
import os
import sqlite3
import math
from urllib.parse import quote_plus
import json
from collections import defaultdict

url = "http://127.0.0.1:5000"

PROMPT_TOTAL = 40
PROJECT_PER_PROMPT = 10

# โหลดตำแหน่งล่าสุด
def load_last_index():
    try:
        with open("last_project_index.txt", "r") as f:
            return int(f.read().strip())
    except:
        return 0

def save_last_index(index):
    with open("last_project_index.txt", "w") as f:
        f.write(str(index))

def get_project_range(project_list, start, count):
    return [project_list[(start + i) % len(project_list)] for i in range(count)]

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
        requests.post(f"{url}/api/update/news", json=payload)
        print(f"🔄 Updated row ID {rows_id} as unused")
    except Exception as e:
        print(f"❌ Failed to update row {rows_id}:", e)
    time.sleep(0.2)

def get_charactor_by_fixed_index(index):
    try:
        conn = sqlite3.connect("./promt.db")
        cur = conn.cursor()
        cur.execute("SELECT * FROM charactors LIMIT 1")
        row = cur.fetchone()
        columns = [desc[0] for desc in cur.description]
        conn.close()
        print(f"🧠 Using prompt index {index}: {columns[index]}")
        return [index, row[index]] if row and index < len(columns) else None
    except Exception as e:
        print(f"❌ Error reading prompt index {index}:", e)
        return None

def generate_comment(promt_text, topic_news):
    gen_url = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:streamGenerateContent?alt=sse"
    headers = {
        "Accept": "*/*",
        "Content-Type": "application/json",
        "x-goog-api-key": "AIzaSyAO_lmHQrYvFIeME5k0l8if66jTTtC_3HU"
    }
    body = {
        "contents": [{
            "parts": [
                {"text": promt_text},
                {"text": f"จากเนื้อหาข่าว: {topic_news} สร้างคอมเมนต์ 10 คอมเมนต์  Output Format: ให้ตอบเฉพาะ 10 คำตอบแยกบรรทัด ห้ามมีเลขลำดับ (1. 2. 3.) หรือคำว่า คอมเมนต์ / ข้อความ / ข้อใด ๆ ไม่ต้องเว้นบรรทัดระหว่างกัน ให้ตอบเป็นข้อความแต่ละบรรทัด 10 บรรทัด เท่านั้น "}
            ],
            "role": "user"
        }],
        "generationConfig": {
            "temperature": 1,
            "topP": 1,
            "topK": 1500,
            "maxOutputTokens": 8192
        }
    }
    try:
        res = requests.post(gen_url, headers=headers, json=body, stream=True)
        full_text = ""
        for line in res.iter_lines():
            if line and line.decode("utf-8").startswith("data: "):
                parts = json.loads(line.decode("utf-8")[6:]).get("candidates", [])[0].get("content", {}).get("parts", [])
                for part in parts:
                    full_text += part.get("text", "")
        comments = [line.strip() for line in full_text.split("\n") if line.strip()][:10]
        print(f"💬 Generated {len(comments)} comments from prompt")
        return comments
    except Exception as e:
        print("❌ Error generating comments:", e)
        return []

def get_project_folders():
    parent = os.path.abspath(os.path.join(os.path.dirname(os.path.abspath(__file__)), ".."))
    return sorted([f for f in os.listdir(parent) if f.startswith("acc")])

def group_by_key(rows):
    grouped = defaultdict(list)
    for row in rows:
        k = f"{row.get('link')}|{row.get('topic')}|{row.get('reaction')}|{row.get('status')}|{row.get('like_value')}|{row.get('comment_value')}"
        grouped[k].append(row)
    return grouped

def main():
    project_list = get_project_folders()
    project_len = len(project_list)
    start_index = load_last_index()
    news_data = check_dashboards()
    unused_rows = [r for r in news_data if r.get('log') == 'unused']
    grouped_data = group_by_key(unused_rows)

    for key, rows in grouped_data.items():
        link, topic, reaction, status, like_value, comment_value = key.split('|')
        like_value = int(like_value)
        comment_value = int(comment_value)
        rows_id = rows[0].get("id")

        projects_to_use = get_project_range(project_list, start_index, like_value)
        save_last_index((start_index + like_value) % project_len)

        ### สำหรับ status ที่มี comment
        if status in ["like_and_comment", "like_comment_only", "like_and_reply_comment", "like_reel_comment_reel"]:
            # สร้าง comment ล่วงหน้า
            all_comments = set()
            offset = 0
            while len(all_comments) < comment_value:
                prompt_index = ((start_index // PROJECT_PER_PROMPT) + offset) % PROMPT_TOTAL
                prompt_data = get_charactor_by_fixed_index(prompt_index)
                if not prompt_data:
                    break
                comments = generate_comment(prompt_data[1], topic)
                for c in comments:
                    if len(all_comments) < comment_value:
                        all_comments.add(c.strip())
                offset += 1
                time.sleep(1)
            all_comments = list(all_comments)

        # ตรวจสอบว่า status ไหนใช้ comment แล้วส่งไปที่ API ที่ถูกต้อง
        if status == "like_and_comment":
            endpoint = "like-and-comment"
            for i, project in enumerate(projects_to_use):
                comment_text = all_comments[i % len(all_comments)] if all_comments else ""
                api = f"{url}/api/update/{project}/{endpoint}?reaction_type={quote_plus(reaction)}&link={quote_plus(link)}&comment_text={quote_plus(comment_text)}"
                try:
                    requests.post(api)
                    check_unused(rows_id)
                    print(f"✅ {project} | {reaction} | {comment_text}")
                except Exception as e:
                    print(f"❌ Error posting to {project}:", e)
                time.sleep(0.3)

        elif status == "like_comment_only":
            endpoint = "like-comment-only"
            for i, project in enumerate(projects_to_use):
                comment_text = all_comments[i % len(all_comments)] if all_comments else ""
                api = f"{url}/api/update/{project}/{endpoint}?reaction_type={quote_plus(reaction)}&link={quote_plus(link)}&comment_text={quote_plus(comment_text)}"
                try:
                    requests.post(api)
                    check_unused(rows_id)
                    print(f"✅ {project} | {reaction} | {comment_text}")
                except Exception as e:
                    print(f"❌ Error posting to {project}:", e)
                time.sleep(0.3)

        elif status == "like_and_reply_comment":
            endpoint = "like-and-comment-reply-comment"
            for i, project in enumerate(projects_to_use):
                comment_text = all_comments[i % len(all_comments)] if all_comments else ""
                api = f"{url}/api/update/{project}/{endpoint}?reaction_type={quote_plus(reaction)}&link={quote_plus(link)}&comment_text={quote_plus(comment_text)}"
                try:
                    requests.post(api)
                    check_unused(rows_id)
                    print(f"✅ {project} | {reaction} | {comment_text}")
                except Exception as e:
                    print(f"❌ Error posting to {project}:", e)
                time.sleep(0.3)

        elif status == "like_reel_comment_reel":
            endpoint = "like-reel-and-comment-reel"
            for i, project in enumerate(projects_to_use):
                comment_text = all_comments[i % len(all_comments)] if all_comments else ""
                api = f"{url}/api/update/{project}/{endpoint}?reaction_type={quote_plus(reaction)}&link={quote_plus(link)}&comment_text={quote_plus(comment_text)}"
                try:
                    requests.post(api)
                    check_unused(rows_id)
                    print(f"✅ {project} | {reaction} | {comment_text}")
                except Exception as e:
                    print(f"❌ Error posting to {project}:", e)
                time.sleep(0.3)

        ### สำหรับ status ที่ไม่มี comment (เช่น like_only, shared_link)
        elif status == "like_only":
            for project in projects_to_use:
                api = f"{url}/api/update/{project}/like-only?reaction_type={quote_plus(reaction)}&link={quote_plus(link)}"
                try:
                    requests.post(api)
                    check_unused(rows_id)
                    print(f"👍 {project} | like_only")
                except Exception as e:
                    print(f"❌ Error posting like_only:", e)
                time.sleep(0.3)

        elif status == "like_page":
            for project in projects_to_use:
                api = f"{url}/api/update/{project}/like-page?reaction_type={quote_plus(reaction)}&link={quote_plus(link)}"
                try:
                    requests.post(api)
                    check_unused(rows_id)
                    print(f"📄 {project} | like_page")
                except Exception as e:
                    print(f"❌ Error posting like_page:", e)
                time.sleep(0.3)

        elif status == "like_reel_only":
            for project in projects_to_use:
                api = f"{url}/api/update/{project}/like-reel-only?reaction_type={quote_plus(reaction)}&link={quote_plus(link)}"
                try:
                    requests.post(api)
                    check_unused(rows_id)
                    print(f"🎥 {project} | like_reel_only")
                except Exception as e:
                    print(f"❌ Error posting like_reel_only:", e)
                time.sleep(0.3)

        elif status == "shared_link_text":
            for project in projects_to_use:
                # ใช้ generate_comment เพื่อสร้างข้อความ
                all_comments = set()
                offset = 0
                while len(all_comments) < 1:  # เราต้องการข้อความเดียวที่ไม่ซ้ำ
                    prompt_index = ((start_index // PROJECT_PER_PROMPT) + offset) % PROMPT_TOTAL
                    prompt_data = get_charactor_by_fixed_index(prompt_index)
                    if not prompt_data:
                        break
                    comments = generate_comment(prompt_data[1], topic)
                    for c in comments:
                        if len(all_comments) < 1:  # สร้างแค่ 1 ข้อความ
                            all_comments.add(c.strip())
                    offset += 1
                    time.sleep(1)

                comment_text = list(all_comments)[0] if all_comments else "ซัพพอร์ตนะ"  # ถ้าไม่มี comment ก็ใช้ข้อความเดิม

                api = f"{url}/api/update/{project}/shared-link-text?status_link={quote_plus(link)}&status_text={quote_plus(comment_text)}"
                try:
                    requests.post(api)
                    check_unused(rows_id)
                    print(f"🔗 {project} | shared_link_text | {comment_text}")
                except Exception as e:
                    print(f"❌ Error shared_link_text:", e)
                time.sleep(0.3)

        elif status == "shared_link":
            for project in projects_to_use:
                api = f"{url}/api/update/{project}/shared-link?link_link={quote_plus(link)}"
                try:
                    requests.post(api)
                    check_unused(rows_id)
                    print(f"🔗 {project} | shared_link")
                except Exception as e:
                    print(f"❌ Error shared_link:", e)
                time.sleep(0.3)

        elif status == "join_group":
            for project in projects_to_use:
                api = f"{url}/api/update/{project}/group-id?group_id={quote_plus(link)}"
                try:
                    requests.post(api)
                    check_unused(rows_id)
                    print(f"👥 {project} | join_group")
                except Exception as e:
                    print(f"❌ Error join_group:", e)
                time.sleep(0.3)

        else:
            print(f"⚠️ ไม่รู้จักสถานะ: {status}")

def main_loop():
    while True:
        try:
            main()
            time.sleep(10)
        except Exception as e:
            print("❌ ERROR main_loop:", e)
            time.sleep(5)

if __name__ == "__main__":
    main_loop()
