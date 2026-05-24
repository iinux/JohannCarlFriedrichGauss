import socket

KEY = b"0123456789abcdef"
IV = b"abcdefghijklmnop"
PLAINTEXT = b"Hello AF_ALG!!!"

def pkcs7_pad(data, block_size=16):
    pad_len = block_size - (len(data) % block_size)
    return data + bytes([pad_len]) * pad_len

def pkcs7_unpad(data):
    return data[:-data[-1]]

# 创建 AF_ALG socket
tfm = socket.socket(socket.AF_ALG, socket.SOCK_SEQPACKET, 0)
tfm.bind(("skcipher", "cbc(aes)"))
tfm.setsockopt(socket.SOL_ALG, socket.ALG_SET_KEY, KEY)

plaintext = pkcs7_pad(PLAINTEXT)

# ===== 加密 =====
op, _ = tfm.accept()
op.sendmsg_afalg([plaintext], op=socket.ALG_OP_ENCRYPT, iv=IV)
ciphertext = op.recv(len(plaintext))
op.close()

print("plaintext :", PLAINTEXT)
print("ciphertext:", ciphertext.hex())

# ===== 解密 =====
dec_op, _ = tfm.accept()
dec_op.sendmsg_afalg([ciphertext], op=socket.ALG_OP_DECRYPT, iv=IV)
decrypted = dec_op.recv(len(ciphertext))
dec_op.close()

decrypted = pkcs7_unpad(decrypted)
print("decrypted :", decrypted)

tfm.close()
