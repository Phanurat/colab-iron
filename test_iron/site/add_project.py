import os
import shutil
import sqlite3
from PyQt5.QtWidgets import (
    QApplication, QWidget, QVBoxLayout, QLabel, QTextEdit,
    QPushButton, QListWidget, QMessageBox, QListWidgetItem, QAbstractItemView
)
from PyQt5.QtCore import Qt
import sys

# ------------------- CONFIG -------------------
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
TEMPLATE_FOLDER = os.path.join(BASE_DIR, "Template_structure")
TARGET_ROOT = os.path.abspath(os.path.join(BASE_DIR, ".."))
# ----------------------------------------------

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

def copy_template_structure(new_folder_name):
    src = TEMPLATE_FOLDER
    dst = os.path.join(TARGET_ROOT, new_folder_name)
    shutil.copytree(src, dst)
    return dst

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

class UIDProjectManager(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("🧠 UID Project Manager (PyQt5)")
        self.setFixedSize(600, 700)

        self.layout = QVBoxLayout()
        self.setLayout(self.layout)

        self.setup_ui()
        self.refresh_project_list()

    def setup_ui(self):
        self.label_input = QLabel("📝 ใส่ UID|ACCESS_TOKEN (1 บรรทัด/ชุด):")
        self.layout.addWidget(self.label_input)

        self.text_input = QTextEdit()
        self.layout.addWidget(self.text_input)

        self.btn_submit = QPushButton("📦 สร้าง accXXX ทั้งหมด")
        self.btn_submit.clicked.connect(self.handle_submit)
        self.layout.addWidget(self.btn_submit)

        self.label_list = QLabel("🗂️ รายการโปรเจกต์ (เลือกเพื่อลบ):")
        self.layout.addWidget(self.label_list)

        self.list_projects = QListWidget()
        self.list_projects.setSelectionMode(QAbstractItemView.MultiSelection)
        self.layout.addWidget(self.list_projects)

        self.btn_delete = QPushButton("🧹 ลบโปรเจกต์ที่เลือก")
        self.btn_delete.clicked.connect(self.delete_selected_projects)
        self.layout.addWidget(self.btn_delete)

    def refresh_project_list(self):
        self.list_projects.clear()
        acc_folders = sorted(
            [f for f in os.listdir(TARGET_ROOT) if f.startswith("acc") and os.path.isdir(os.path.join(TARGET_ROOT, f))]
        )
        for name in acc_folders:
            self.list_projects.addItem(QListWidgetItem(name))

    def handle_submit(self):
        input_text = self.text_input.toPlainText().strip()
        lines = [line.strip() for line in input_text.splitlines() if line.strip()]

        if not lines:
            QMessageBox.critical(self, "❌ ไม่มีข้อมูล", "กรุณากรอก UID|ACCESS_TOKEN อย่างน้อย 1 บรรทัด")
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
        
        QMessageBox.information(self, "ผลลัพธ์", summary)
        self.text_input.clear()
        self.refresh_project_list()

    def delete_selected_projects(self):
        selected_items = self.list_projects.selectedItems()
        if not selected_items:
            QMessageBox.warning(self, "ยังไม่ได้เลือก", "กรุณาเลือกโปรเจกต์ที่จะลบ")
            return

        selected_names = [item.text() for item in selected_items]
        confirm = QMessageBox.question(
            self,
            "ยืนยันการลบ",
            f"คุณแน่ใจหรือไม่ว่าต้องการลบ {len(selected_names)} โปรเจกต์?\n\n" + "\n".join(selected_names)
        )

        if confirm != QMessageBox.Yes:
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

        # ✅ แสดงผลการลบ
        msg = f"✅ ลบสำเร็จ {len(success)} โปรเจกต์:\n" + "\n".join(success)
        if failed:
            msg += f"\n\n❌ ลบไม่สำเร็จ {len(failed)} โปรเจกต์:\n"
            for name, err in failed:
                msg += f"  - {name}: {err}"

        QMessageBox.information(self, "ผลลัพธ์", msg)  # <-- เพิ่มบรรทัดนี้

        self.refresh_project_list()  # <-- รีเฟรชรายการหลังลบ


if __name__ == "__main__":
    import sys
    app = QApplication(sys.argv)
    window = UIDProjectManager()
    window.show()
    sys.exit(app.exec_())
