import stun

nat_type, external_ip, external_port = stun.get_ip_info()

if __name__ == '__main__':
    print(nat_type, external_ip, external_port)