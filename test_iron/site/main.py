import os
import sqlite3
import requests
import json
from collections import defaultdict
from urllib.parse import quote_plus
import time
import tkinter as tk
from tkinter import messagebox
from tkinter import ttk
import threading
import queue
import webview
import subprocess
import sys

def run_app():
    subprocess.Popen(["python", "app.py"])

# รัน server.py
def run_server():
    subprocess.Popen(["python", "server.py"])


threading.Thread(target=run_app, daemon=True).start()
threading.Thread(target=run_server, daemon=True).start()
time.sleep(2)
print("เซิร์ฟเวอร์พร้อมใช้งานแล้ว เริ่มโปรแกรมต่อไป...")

#==============================================================

url = "http://127.0.0.1:5000"
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
DB_PATH = os.path.join(BASE_DIR, "db/promt.db")

endpoint_map = {
    "like_and_comment": "like-and-comment",
    "like_and_reply_comment": "like-and-comment-reply-comment",
    "like_reel_comment_reel": "like-reel-and-comment-reel",
    "like_comment_only": "like-comment-only",
    "like_only": "like-only",
    "like_page": "like-page",
    "like_reel_only": "like-reel-only",
    "shared_link_text": "shared-link-text",
    "shared_link": "shared-link",
    "join_group": "group-id"
}
def all_prompt():
    with sqlite3.connect(DB_PATH) as conn:
        cursor = conn.cursor()
        cursor.execute("SELECT * FROM charactors")
        rows = cursor.fetchall()
        return [list(col) for col in zip(*rows)] if rows else ([], [], [], [])

def show_mapp_project():
    with sqlite3.connect(DB_PATH) as conn:
        cursor = conn.cursor()
        cursor.execute("SELECT * FROM mapp_prompt")
        rows = cursor.fetchall()
        return [list(col) for col in zip(*rows)] if rows else ([], [], [], [], [])

def generate_comment(prompt_text, topic_news):
    gen_url = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:streamGenerateContent?alt=sse"
    headers = {
        "Accept": "*/*",
        "Content-Type": "application/json",
        "x-goog-api-key": "AIzaSyAO_lmHQrYvFIeME5k0l8if66jTTtC_3HU"
    }
    body = {
        "contents": [{
            "parts": [
                {"text": prompt_text},
                {"text": f"""From the news content: {topic_news}  Generate 10 comments  Output Format: Only return 10 separate lines of responses. Do not include numbering (1. 2. 3.) or words like comment / message / any label. Do not leave blank lines between them. Each line should be a comment. Return exactly 10 lines."""}
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

def save_comment_to_db(project_id, topic, prompt_id, prompt_text, comment_text, reaction_type, status, link, log="success"):
    with sqlite3.connect(DB_PATH) as conn:
        cursor = conn.cursor()
        cursor.execute("""
            INSERT INTO comment_dashboard (project_id, topic, prompt_id, prompt_text, comment, log, timestamp, reaction_type, status, link)
            VALUES (?, ?, ?, ?, ?, ?, datetime('now'), ?, ?, ?)
        """, (project_id, topic, prompt_id, prompt_text, comment_text, log, reaction_type, status, link))
        conn.commit()
        return cursor.lastrowid

def all_comment_data():
    with sqlite3.connect(DB_PATH) as conn:
        cursor = conn.cursor()
        cursor.execute("SELECT * FROM comment_dashboard")
        rows = cursor.fetchall()
        return [list(col) for col in zip(*rows)] if rows else ([], [], [], [], [], [], [], [], [], [], [])

def insert_comment_db(id):
    with sqlite3.connect(DB_PATH) as conn:
        try:
            cursor = conn.cursor()
            cursor.execute(f"UPDATE comment_dashboard SET log = 'used' WHERE id = {id} AND log = 'success';")
            if cursor.rowcount > 0:
                conn.commit()
                print("✅ DB อัปเดตเป็น 'used' เรียบร้อย")
                return True
            else:
                print("⚠️ ไม่ได้อัปเดตอะไรเลย (อาจเคยอัปเดตไปแล้ว)")
                return False
        except Exception as e:
            print(f"❌ Error updating comment_dashboard: {e}")
            return False
        finally:
            cursor.close()

def read_index(max_prompt):
    try:
        with open("set_index.txt", "r") as f:
            idx = int(f.read().strip())
            return idx if idx < max_prompt else 0
    except:
        return 0

def write_index(idx):
    with open("set_index.txt", "w") as f:
        f.write(str(idx))

def main(topic_news, status, reaction_type, link, total_comment, update_progress=None):
    rows_id, project_id_list, prompt_id_list, prompt_role_list, prompt_content_list = show_mapp_project()
    total_projects = len(set(project_id_list))
    projects_per_prompt = 10
    max_prompt_allowed = min(len(set(prompt_id_list)), total_projects // projects_per_prompt)
    if max_prompt_allowed == 0:
        max_prompt_allowed = 1

    all_maps = [{
        "project_id": pid,
        "prompt_id": prid,
        "prompt_role": prole,
        "prompt_content": pcontent
    } for pid, prid, prole, pcontent in zip(project_id_list, prompt_id_list, prompt_role_list, prompt_content_list)]

    grouped = defaultdict(list)
    prompt_meta = {}
    for item in all_maps:
        grouped[item["prompt_id"]].append(item["project_id"])
        if item["prompt_id"] not in prompt_meta:
            prompt_meta[item["prompt_id"]] = {
                "prompt_role": item["prompt_role"],
                "prompt_content": item["prompt_content"]
            }

    generated_count = 0
    new_comment_ids = []
    start_index = read_index(max_prompt_allowed)
    prompt_ids = list(grouped.keys())

    idx = start_index
    while generated_count < total_comment:
        prompt_idx = idx % max_prompt_allowed
        prompt_id = prompt_ids[prompt_idx]
        projects = grouped[prompt_id]
        prompt_content = prompt_meta[prompt_id]["prompt_content"]

        comments = [prompt_content] * len(projects) if status in [
            "like_comment_only", "like_only", "like_page", "like_reel_only", "shared_link_text", "shared_link", "join_group"
        ] else generate_comment(prompt_content, topic_news)

        available = min(len(comments), len(projects), total_comment - generated_count)
        if available == 0:
            idx += 1
            continue

        for i in range(available):
            project = projects[i]
            comment_text = comments[i]
            comment_id = save_comment_to_db(project, topic_news, prompt_id, prompt_content, comment_text, reaction_type, status, link)
            new_comment_ids.append(comment_id)
            generated_count += 1
            if update_progress:
                update_progress(generated_count / total_comment * 100)
            if generated_count >= total_comment:
                break

        idx += 1
        write_index(idx % max_prompt_allowed)

    endpoint = endpoint_map.get(status)
    if not endpoint:
        print(f"⚠️ Status '{status}' not recognized.")
        return

    for id_val in new_comment_ids:
        id_list, project_id, topic, link_list, prompt_id, prompt_text, comment, log, reaction_type_list, status_list, timestamp = all_comment_data()
        idx = id_list.index(id_val)
        proj = project_id[idx]
        endpoint = endpoint_map.get(status)
        if not endpoint:
            print(f"❌ Unknown endpoint for status '{status}'")
            continue

        base_api = f"{url}/api/update/{proj}/{endpoint}"
        params = {}

        # Mapping แบบละเอียดตาม endpoint
        if endpoint in ["like-and-comment", "like-and-comment-reply-comment", "like-reel-and-comment-reel", "like-comment-only"]:
            params["reaction_type"] = reaction_type
            params["link"] = link_list[idx]
            params["comment_text"] = comment[idx]

        elif endpoint in ["like-only", "like-reel-only", "like-page"]:
            params["reaction_type"] = reaction_type
            params["link"] = link_list[idx]

        elif endpoint == "shared-link":
            params["link_link"] = link_list[idx]

        elif endpoint == "shared-link-text":
            params["status_link"] = link_list[idx]
            params["status_text"] = "ทดสอบระบบ"
        
        elif endpoint == "like-page":
            params["link_page"] = link_list[idx]

        elif endpoint == "group-id":
            params["group_id"] = link_list[idx]

        else:
            print(f"⚠️ ยังไม่รองรับ endpoint: {endpoint}")
            continue

        # Build final URL
        final_url = base_api + "?" + "&".join(f"{k}={quote_plus(str(v))}" for k, v in params.items())

        try:
            response = requests.post(final_url)
            response.raise_for_status()
            insert_comment_db(id_val)
        except Exception as e:
            print(f"❌ Error posting to {proj}: {e}")

        time.sleep(0.3)


# ---------------- GUI Queue ----------------
task_queue = queue.Queue()
processing_task = False

def process_tasks():
    global processing_task
    if processing_task:
        return
    processing_task = True

    while not task_queue.empty():
        task = task_queue.get()
        def update_progress(percent):
            task["progress_var"].set(percent)
            task["label"].config(text=f"🔋 {int(percent)}%")
        try:
            main(task["topic"], task["status"], task["reaction"], task["link"], task["total"], update_progress)
            task["label"].config(text="✅ เสร็จแล้ว")
            task["progress_var"].set(100)
        except Exception as e:
            task["label"].config(text=f"❌ Error: {e}")
    processing_task = False

def enqueue_task():
    topic = topic_entry.get()
    status = status_var.get()
    reaction = reaction_entry.get()
    link = link_entry.get()
    try:
        total = int(total_entry.get())
    except ValueError:
        messagebox.showerror("Error", "จำนวนคอมเมนต์ต้องเป็นตัวเลข")
        return

    if not topic or not status or not reaction or not link:
        messagebox.showerror("Error", "กรุณากรอกข้อมูลให้ครบถ้วน")
        return

    progress_var = tk.DoubleVar()
    task_frame = ttk.Frame(task_list_frame, relief="solid", borderwidth=1)
    task_frame.pack(fill="x", pady=5, padx=5)

    title = f"📰 Topic: {topic} | 💬 {total} คอมเมนต์ | 🔁 {status}"
    ttk.Label(task_frame, text=title, font=("Arial", 11, "bold")).pack(anchor="w", padx=5, pady=2)
    progress_bar = ttk.Progressbar(task_frame, variable=progress_var, maximum=100)
    progress_bar.pack(fill="x", padx=5)
    label = ttk.Label(task_frame, text="🌟 0%")
    label.pack(anchor="e", padx=5)

    task = {
        "topic": topic,
        "status": status,
        "reaction": reaction,
        "link": link,
        "total": total,
        "progress_var": progress_var,
        "label": label
    }
    task_queue.put(task)
    threading.Thread(target=process_tasks).start()

    topic_entry.delete(0, tk.END)
    status_var.set("like_and_comment")
    reaction_entry.delete(0, tk.END)
    link_entry.delete(0, tk.END)
    total_entry.delete(0, tk.END)

def clear_tasks():
    for widget in task_list_frame.winfo_children():
        widget.destroy()
    while not task_queue.empty():
        try:
            task_queue.get(False)
        except queue.Empty:
            continue
        task_queue.task_done()

# ---------------- GUI ----------------
root = tk.Tk()
root.title("💬 Comment Generator Queue System")
root.geometry("650x750")

form_frame = ttk.Frame(root)
form_frame.pack(pady=10, padx=10, fill="x")

topic_entry = ttk.Entry(form_frame, width=50)
total_entry = ttk.Entry(form_frame, width=10)
reaction_entry = ttk.Entry(form_frame, width=50)
link_entry = ttk.Entry(form_frame, width=50)
status_var = tk.StringVar(value="like_and_comment")

widgets = [
    ("หัวข้อข่าว (Topic):", topic_entry),
    ("Status:", ttk.Combobox(form_frame, textvariable=status_var, values=list(endpoint_map.keys()), state="readonly")),
    ("Reaction Type:", reaction_entry),
    ("ลิงก์โพสต์ (Link):", link_entry),
    ("จำนวนคอมเมนต์:", total_entry)
]

for i, (label_text, widget) in enumerate(widgets):
    ttk.Label(form_frame, text=label_text).grid(row=i, column=0, sticky="w", pady=3)
    widget.grid(row=i, column=1, sticky="ew", pady=3)
form_frame.columnconfigure(1, weight=1)

#=========================================================================================
def open_app():
    webview.create_window("My Local App", "http://127.0.0.1:5050", width=1200, height=800)
    webview.start()

open_app_btn = ttk.Button(root, text="App Profile", command=open_app)
open_app_btn.pack(pady=10)

def add_project():
    subprocess.Popen(["python", "add_project.py"])

open_app_btn = ttk.Button(root, text="Manage Project", command=add_project)
open_app_btn.pack(pady=10)

def restart_app():
    os.execl(sys.executable, sys.executable, "main.py")

restart_app_btn = ttk.Button(root, text="🔄 Restart App", command=restart_app)
restart_app_btn.pack(pady=10)
#=========================================================================================

submit_btn = ttk.Button(root, text="🚀 เพิ่มเข้าคิว Generate", command=enqueue_task)
submit_btn.pack(pady=10)

ttk.Label(root, text="📋 รายการ Task ทั้งหมด", font=("Arial", 13, "bold")).pack(anchor="w", padx=12)

task_list_canvas = tk.Canvas(root, height=400)
task_list_canvas.pack(fill="both", expand=True, padx=12)

scrollbar = ttk.Scrollbar(root, orient="vertical", command=task_list_canvas.yview)
scrollbar.place(in_=task_list_canvas, relx=1.0, rely=0, relheight=1.0, anchor='ne')
task_list_canvas.configure(yscrollcommand=scrollbar.set)

container = ttk.Frame(task_list_canvas)
task_list_canvas.create_window((0, 0), window=container, anchor="nw")
container.bind("<Configure>", lambda e: task_list_canvas.configure(scrollregion=task_list_canvas.bbox("all")))
task_list_frame = container

clear_btn = ttk.Button(root, text="🗑️ ล้างรายการทั้งหมด", command=clear_tasks)
clear_btn.pack(pady=5)

root.mainloop()