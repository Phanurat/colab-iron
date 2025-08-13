comment_want = 10
project = 400
prompt_value = 40

with open("set_index.txt", "r") as file:
    value = int(file.readline().strip())
    # print(value)

    if (value_m%value_n) == 10:
        with open("set_index.txt", "w") as out_file:
            out_file.write("0")
