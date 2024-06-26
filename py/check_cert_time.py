import ssl
import socket
from datetime import datetime
from cryptography import x509
from cryptography.hazmat.backends import default_backend


def check_ssl_certificate(hostname, port=443):
    context = ssl.create_default_context()
    with socket.create_connection((hostname, port)) as sock:
        with context.wrap_socket(sock, server_hostname=hostname) as secure_sock:
            der_cert = secure_sock.getpeercert(binary_form=True)

    certificate = x509.load_der_x509_certificate(der_cert, default_backend())

    not_before = certificate.not_valid_before
    not_after = certificate.not_valid_after
    current_time = datetime.utcnow()

    print(f"Certificate validity period:")
    print(f"Not Before: {not_before}")
    print(f"Not After: {not_after}")
    print(f"Current Time: {current_time}")

    if not_before <= current_time <= not_after:
        print("The certificate is currently valid.")
    else:
        print("The certificate is not currently valid.")


# Usage
hostname = "www.example.com"
check_ssl_certificate(hostname)
