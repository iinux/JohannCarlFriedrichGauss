#!/usr/bin/env python3
"""打印指定的 Unicode 字符。"""

chars = [
    ("U+0027", "\u0027"),  # APOSTROPHE '
    ("U+2019", "\u2019"),  # RIGHT SINGLE QUOTATION MARK '
    ("U+02BC", "\u02BC"),  # MODIFIER LETTER APOSTROPHE ʼ
    ("U+0289", "\u0289"),  # LATIN LETTER SMALL CAPITAL BARRED E ʉ
]

for name, ch in chars:
    print(f"{name} -> {ch}  (codepoint: {ord(ch)})")