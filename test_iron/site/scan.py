import csv

# กำหนดชื่อไฟล์ input และ output
input_file = "test400a1.txt"
output_file = "output.csv"

# อ่านไฟล์และแยกข้อมูล
rows = []
with open(input_file, "r", encoding="utf-8") as f:
    for line in f:
        line = line.strip()
        if "|" in line:
            id_part, token_part = line.split("|", 1)
            rows.append([id_part, token_part])

# เขียนเป็น CSV
with open(output_file, "w", newline="", encoding="utf-8") as f:
    writer = csv.writer(f)
    writer.writerow(["ID", "Token"])  # หัวตาราง
    writer.writerows(rows)

print(f"✅ แยกข้อมูลสำเร็จ! บันทึกไว้ใน: {output_file}")
