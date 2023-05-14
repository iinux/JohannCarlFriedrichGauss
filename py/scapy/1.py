from scapy.all import *


# https://scapy.readthedocs.io/en/latest/usage.html#simple-one-liners
send(IP(dst="198.46.152.123")/ICMP())
