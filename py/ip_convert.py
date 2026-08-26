def int_to_ip(num: int) -> str:
    if not 0 <= num <= 0xFFFFFFFF:
        raise ValueError("num must be in range [0, 4294967295]")
    parts = []
    for shift in (24, 16, 8, 0):
        parts.append(str((num >> shift) & 0xFF))
    return ".".join(parts)


def ip_to_int(ip: str) -> int:
    segments = ip.split(".")
    if len(segments) != 4:
        raise ValueError("ip must have 4 dot-separated segments")
    num = 0
    for seg in segments:
        if not seg.isdigit():
            raise ValueError(f"invalid segment: {seg}")
        value = int(seg)
        if not 0 <= value <= 255:
            raise ValueError(f"segment out of range: {seg}")
        num = (num << 8) | value
    return num


if __name__ == "__main__":
    import sys

    if len(sys.argv) != 2:
        print("用法: python ip_convert.py <ip地址 或 整数>")
        sys.exit(1)

    arg = sys.argv[1]
    if arg.isdigit():
        print(int_to_ip(int(arg)))
    else:
        print(ip_to_int(arg))
