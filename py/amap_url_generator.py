def generate_amap_uri(coords):
    # 提取起点、途经点和终点
    start_coord = coords[0]
    end_coord = coords[-1]
    via_coords = coords[1:-1]

    # 解析起终点坐标
    slon, slat = start_coord
    dlon, dlat = end_coord

    # 如果有途经点则解析
    if via_coords:
        vialons = "|".join(str(lon) for lon, _ in via_coords)
        vialats = "|".join(str(lat) for _, lat in via_coords)
        vian = len(via_coords)
        # 生成途经点编号 (1,2,3,...)
        vianames = "|".join(str(i+1) for i in range(len(via_coords)))
    else:
        vialons = vialats = ""
        vian = 0
        vianames = ""

    # 构造高德URL
    uri = f"amapuri://route/plan/?dlon={dlon}&dlat={dlat}&dev=0&t=0"

    if vialons:
        uri += f"&vialons={vialons}&vialats={vialats}&vian={vian}&vianames={vianames}"

    if slon != 'no_start':
        uri += f"&slon={slon}&slat={slat}"

    return uri


# 示例输入处理
input_data = """
"""

# 处理输入数据
coords = []
for line in input_data.strip().split("\n"):
    if line.strip():  # 跳过空行
        if line == 'no_start':
            coords.append((line, line))
            continue
        lon, lat = line.split(",")
        coords.append((lon.strip(), lat.strip()))

# 生成URL
amap_url = generate_amap_uri(coords)
print("生成的高德URL:")
print(amap_url)
