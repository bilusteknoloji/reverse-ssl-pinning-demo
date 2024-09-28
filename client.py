import requests

client_cert = (
    'certs/client/client-cert-signed.pem',
    'certs/client/client-key.pem',
)
server_ca_cert = 'certs/client/client-ca.pem'
url = 'https://localhost:8443'

try:
    response = requests.get(url, cert=client_cert, verify=server_ca_cert)
    print(f'Server response:\n{response.text}')

except requests.exceptions.SSLError as ssl_error:
    print(f'SSL Error: {ssl_error}')

except Exception as e:
    print(f'Error: {e}')
