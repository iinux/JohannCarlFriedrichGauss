from PIL import Image, ImageDraw, ImageFont
import os

rows = [
    ("biz_code", "biz_message", "业务状态参考"),
    ("00000000", "成功", "成功"),
    ("00000001", "未知系统异常，请稍后重试", "不确定，建议重试"),
]

font_candidates = [
    "/System/Library/Fonts/PingFang.ttc",
    "/System/Library/Fonts/STHeiti Light.ttc",
    "/System/Library/Fonts/Hiragino Sans GB.ttc",
    "/Library/Fonts/Arial Unicode.ttf",
]

font_path = None
for p in font_candidates:
    if os.path.exists(p):
        font_path = p
        break

font_size = 18
header_font_size = 20
font = ImageFont.truetype(font_path, font_size)
header_font = ImageFont.truetype(font_path, header_font_size)

col_widths = [140, 560, 200]
row_height = 38
header_height = 46
padding_x = 12

total_width = sum(col_widths) + 2
total_height = header_height + row_height * (len(rows) - 1) + 2

img = Image.new("RGB", (total_width, total_height), "white")
draw = ImageDraw.Draw(img)

header_bg = (66, 103, 178)
header_fg = "white"
alt_bg = (245, 247, 250)
border = (200, 200, 200)
status_color = {
    "成功": (40, 167, 69),
    "失败": (220, 53, 69),
    "不确定，建议重试": (255, 153, 0),
}

y = 0
for i, row in enumerate(rows):
    x = 0
    is_header = (i == 0)
    if is_header:
        h = header_height
        draw.rectangle([0, y, total_width, y + h], fill=header_bg)
    else:
        h = row_height
        if i % 2 == 0:
            draw.rectangle([0, y, total_width, y + h], fill=alt_bg)

    for j, cell in enumerate(row):
        cw = col_widths[j]
        cur_font = header_font if is_header else font
        fill = header_fg if is_header else "black"
        if not is_header and j == 2 and cell in status_color:
            fill = status_color[cell]

        bbox = draw.textbbox((0, 0), cell, font=cur_font)
        text_h = bbox[3] - bbox[1]
        text_y = y + (h - text_h) // 2 - bbox[1]
        draw.text((x + padding_x, text_y), cell, font=cur_font, fill=fill)

        draw.line([(x + cw, y), (x + cw, y + h)], fill=border, width=1)
        x += cw

    draw.line([(0, y + h), (total_width, y + h)], fill=border, width=1)
    y += h

draw.rectangle([0, 0, total_width - 1, total_height - 1], outline=border, width=1)

out = "biz_code_table.png"
os.makedirs(os.path.dirname(out), exist_ok=True)
img.save(out)
print(out)
