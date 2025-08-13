import requests

def delete_data():
    start_project = 1
    end_project = 101
    for i in range(start_project, end_project):
        project = f"acc{i:03d}"
        table = "shared_link_table"
        # table = "subscribee_id_table"
        url = f"http://127.0.0.1:5000/api/clear/{project}/{table}"

        try:
            response = requests.post(url)
            print(f"[{project}] ✅ Status {response.status_code}: {response.json()}")
        except Exception as e:
            print(f"[{project}] ❌ Error: {e}")

def input_add_uid():
    start_project = 1
    end_project = 51
    uid_add = 10000000000
    for i in range(start_project, end_project):
        project = f"acc{i:03d}"
        url = f"http://127.0.0.1:5000/api/insert/{project}/uid-add?uid_add={uid_add}"

        try:
            response = requests.post(url)
            print(f"[{project}] ✅ Status {response.status_code}: {response.json()}]")
        except Exception as e:
            print(f"[{project}] ❌ Error: {e}")

def delete_add_uid():
    start_project = 1
    end_project = 51
    for i in range(start_project, end_project):
        project = f"acc{i:03d}"
        url = f"http://127.0.0.1:5000/api/delete/{project}/uid-add"

        try:
            response = requests.delete(url)
            print(f"[{project}] ✅ Status {response.status_code}: {response.json()}]")
        except Exception as e:
            print(f"[{project}] ❌ Error: {e}")

delete_data()

# delete_add_uid()
# input_add_uid()