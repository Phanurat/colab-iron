import os
import shutil
import tkinter as tk
from tkinter import ttk, messagebox
import sqlite3

# ------------------- CONFIG -------------------
BASE_DIR = os.path.dirname(os.path.abspath(__file__))  # ที่อยู่ของ main.py
TEMPLATE_FOLDER = os.path.join(BASE_DIR, "Template_structure")  # โฟลเดอร์ template
TARGET_ROOT = os.path.abspath(os.path.join(BASE_DIR, ".."))  # D:\Irondom\irondome_exe_path_base
# ----------------------------------------------

# 🔢 ฟังก์ชันหาหมายเลข accXXX ที่ยังไม่ถูกใช้
def get_next_unused_acc_folder():
    used_numbers = set()
    for name in os.listdir(TARGET_ROOT):
        if name.startswith("acc") and name[3:].isdigit():
            used_numbers.add(int(name[3:]))

    i = 1
    while True:
        if i not in used_numbers:
            return f"acc{str(i).zfill(3)}"
        i += 1

# 🗂️ คัดลอก Template_structure ไปยัง accXXX ใหม่
def copy_template_structure(new_folder_name):
    src = TEMPLATE_FOLDER
    dst = os.path.join(TARGET_ROOT, new_folder_name)
    shutil.copytree(src, dst)
    return dst

# 🗃️ เพิ่ม UID + ACCESS_TOKEN ลงใน fb_comment_system.db
def insert_into_db(db_path, actor_id, access_token):
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    cursor.execute("""
        CREATE TABLE IF NOT EXISTS app_profiles (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            actor_id TEXT,
            access_token TEXT
        );
    """)
    cursor.execute("""
        INSERT INTO app_profiles (actor_id, access_token)
        VALUES (?, ?)
    """, (actor_id, access_token))
    conn.commit()
    conn.close()

# 🚀 เมื่อกดปุ่มเพิ่มโปรเจกต์
def handle_submit():
    input_text = input_box.get("1.0", tk.END).strip()
    lines = [line.strip() for line in input_text.splitlines() if line.strip()]

    if not lines:
        messagebox.showerror("❌ ไม่มีข้อมูล", "กรุณากรอก UID|ACCESS_TOKEN อย่างน้อย 1 บรรทัด")
        return

    success_count = 0
    fail_lines = []
    log_lines = []

    for i, line in enumerate(lines, start=1):
        if "|" not in line:
            fail_lines.append((i, line))
            continue

        parts = line.split("|", 1)
        if len(parts) != 2 or not parts[0].strip() or not parts[1].strip():
            fail_lines.append((i, line))
            continue

        uid, token = parts[0].strip(), parts[1].strip()
        folder_name = get_next_unused_acc_folder()

        try:
            new_path = copy_template_structure(folder_name)
            db_path = os.path.join(new_path, "fb_comment_system.db")
            insert_into_db(db_path, uid, token)
            success_count += 1
            log_lines.append(f"{folder_name} → UID: {uid}")
        except Exception as e:
            fail_lines.append((i, f"{line} → ERROR: {str(e)}"))

    summary = f"✅ เพิ่มสำเร็จ {success_count} รายการ\n"
    if log_lines:
        summary += "\n📁 โฟลเดอร์ที่สร้าง:\n" + "\n".join(log_lines)
    if fail_lines:
        summary += f"\n\n❌ ข้อมูลผิดพลาด {len(fail_lines)} บรรทัด:\n"
        for idx, content in fail_lines:
            summary += f"  [บรรทัด {idx}] {content}\n"

    messagebox.showinfo("ผลลัพธ์", summary)
    input_box.delete("1.0", tk.END)
    refresh_project_list()


# 🗑️ ลบหลายโปรเจกต์ที่เลือกใน listbox
def delete_selected_projects():
    selected_indices = listbox.curselection()
    if not selected_indices:
        messagebox.showwarning("ยังไม่ได้เลือก", "กรุณาเลือกโปรเจกต์ที่จะลบ")
        return

    selected_names = [listbox.get(i) for i in selected_indices]
    confirm = messagebox.askyesno("ยืนยันการลบ", f"คุณแน่ใจหรือไม่ว่าต้องการลบ {len(selected_names)} โปรเจกต์?\n\n" + "\n".join(selected_names))
    if not confirm:
        return

    success = []
    failed = []
    for name in selected_names:
        path = os.path.join(TARGET_ROOT, name)
        try:
            shutil.rmtree(path)
            success.append(name)
        except Exception as e:
            failed.append((name, str(e)))

    msg = f"✅ ลบสำเร็จ {len(success)} โปรเจกต์:\n" + "\n".join(success)
    if failed:
        msg += f"\n\n❌ ลบไม่สำเร็จ {len(failed)} โปรเจกต์:\n"
        for name, err in failed:
            msg += f"  - {name}: {err}\n"

    messagebox.showinfo("ผลลัพธ์", msg)
    refresh_project_list()

# 🔁 โหลดรายชื่อโปรเจกต์ใน listbox
def refresh_project_list():
    listbox.delete(0, tk.END)
    acc_folders = sorted([f for f in os.listdir(TARGET_ROOT) if f.startswith("acc") and os.path.isdir(os.path.join(TARGET_ROOT, f))])
    for name in acc_folders:
        listbox.insert(tk.END, name)

# -------------------- GUI --------------------
root = tk.Tk()
root.title("🧠 UID Project Manager")

frame = ttk.Frame(root, padding=20)
frame.pack()

# กรอกข้อมูลเพิ่ม UID
ttk.Label(frame, text="📝 ใส่ UID|ACCESS_TOKEN แยกบรรทัดละ 1 ชุด:").pack()
input_box = tk.Text(frame, height=10, width=70)
input_box.pack(padx=10, pady=10)

submit_btn = ttk.Button(frame, text="📦 สร้าง accXXX ทั้งหมด", command=handle_submit)
submit_btn.pack(pady=10)

# ลบโปรเจกต์
ttk.Label(frame, text="🗑️ เลือกโปรเจกต์ที่ต้องการลบ (เลือกได้หลายอัน):").pack(pady=(20, 5))
listbox = tk.Listbox(frame, width=40, height=10, selectmode="multiple")
listbox.pack()

delete_btn = ttk.Button(frame, text="🧹 ลบโปรเจกต์ที่เลือก", command=delete_selected_projects)
delete_btn.pack(pady=5)

refresh_project_list()
root.mainloop()
