import random

data = {
    '东周': (-770, -256),
    '春秋': (-770, -476),
    '战国': (-475, -221),
    '秦': (-221, -207),
    '西楚': (-206, -202),
    '西汉': (-202, 9),
    '新': (9, 23),
    '东汉': (25, 220),
    '曹魏': (220, 266),
    '蜀汉': (221, 263),
    '孙吴': (229, 280),
    '西晋': (266, 316),
    '东晋': (317, 420),
    '五胡十六国': (304, 439, False),
    '南北朝': (420, 589),
    '南朝宋': (420, 479, False),
    '南朝齐': (479, 502, False),
    '南朝梁': (502, 557, False),
    '南朝陈': (557, 589, False),
    '北魏': (386, 534, False),
    '西魏': (535, 557, False),
    '东魏': (534, 550, False),
    '北周': (557, 581, False),
    '北齐': (550, 577, False),
    '隋': (581, 619),
    '唐': (618, 907),
    '武周': (690, 705),
    '五代十国': (907, 979),
    '后梁': (907, 923, False),
    '后唐': (923, 937, False),
    '后晋': (936, 947, False),
    '后汉': (947, 951, False),
    '后周': (951, 960, False),
    '十国': (907, 979, False),
    '辽': (916, 1125),
    '西辽': (1124, 1218, False),
    '北宋': (960, 1127),
    '南宋': (1127, 1279),
    '西夏': (1038, 1227),
    '金': (1115, 1234),
    '大蒙古国': (1206, 1635, False),
    '元': (1271, 1368),
    '北元': (1368, 1388, False),
    '明': (1368, 1644),
    '南明': (1644, 1662, False),
    '后金': (1616, 1636, False),
    '清': (1636, 1912),
}


def print_all():
    for k, v in data.items():
        print(k, v)


def aa():
    while True:
        k = random.choice(list(data.keys()))
        answer = data[k]
        if len(answer) > 2:
            continue
        else:
            break

    while True:
        i = input(k + '\n')
        if i == 'all':
            print_all()
            continue
        else:
            break

    iss = i.split(' ')
    if int(iss[0]) == answer[0] and int(iss[1]) == answer[1]:
        print('you are correct')
    else:
        print('your are wrong, answer is %d->%d' % (answer[0], answer[1]))


if __name__ == '__main__':
    while True:
        aa()
