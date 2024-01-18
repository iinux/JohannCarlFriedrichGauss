file_path = 'rb.data'


def start():
    tag = 'start'
    previous_line = ''
    with open(file_path, 'r') as file:
        for line in file:
            columns = line.strip().split(' ')
            last_column = columns[-1]
            if tag != last_column and previous_line != '':
                check_time_and_print(previous_line)

            tag = last_column
            previous_line = line


# 判断时间范围并打印符合条件的行
def check_time_and_print(line):
    line = line.strip()
    columns = line.split(' ')
    arrive_time_str = columns[-4]

    if '18:00' <= arrive_time_str <= '18:39':
        print(columns[0][11:] + ' ' + columns[1][0:8])


if __name__ == '__main__':
    start()
