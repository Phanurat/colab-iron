import sqlite3
import os

def get_charactor(id_prompt):
    base_dir = os.path.dirname(os.path.abspath(__file__))
    char_db_path = os.path.join(base_dir, "promt.db")
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

def save_prompts_to_txt(prompts, filename="selected_prompts.txt"):
    try:
        with open(filename, "w", encoding="utf-8") as f:
            for idx, text in prompts:
                f.write(f"🧠 Prompt {idx}: {text}\n")
        print(f"✅ บันทึก prompt ทั้งหมดลงในไฟล์ {filename} แล้ว")
    except Exception as e:
        print(f"❌ Error saving to file: {e}")

def main():
    targets = [1, 39]  # index 1 และ 39 (แทน id_prompt 1 และ 40)
    results = []

    for i in targets:
        result = get_charactor(i)
        if result:
            print(f"🧠 Prompt {result[0]}: {result[1]}")
            results.append(result)

    if results:
        save_prompts_to_txt(results)

if __name__ == "__main__":
    main()
