from datetime import datetime
from zoneinfo import ZoneInfo

# เวลาท้องถิ่น
local_time = datetime.now()
print(f"Local Time: {local_time}")

# เวลาของ Bangkok (ใช้ zoneinfo)
bangkok_time = datetime.now(ZoneInfo("Asia/Bangkok"))
print(f"Bangkok Time: {bangkok_time}")

# เวลาของ UTC
utc_time = datetime.now(ZoneInfo("UTC"))
print(f"UTC Time: {utc_time}")

# แปลงเวลาจาก UTC เป็นเวลา Bangkok
utc_time = datetime.now(ZoneInfo("UTC"))
bangkok_time_converted = utc_time.astimezone(ZoneInfo("Asia/Bangkok"))
print(f"Converted UTC to Bangkok Time: {bangkok_time_converted}")
