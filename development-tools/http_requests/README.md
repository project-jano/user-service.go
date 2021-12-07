# Generate test CSR

````sh
openssl ecparam -name prime256v1 -genkey -noout -out test_private_key.pem
openssl req -out test_csr.csr -new -key test_private_key.pem -days 365 -subj "/C=ES/ST=Madrid/L=Madrid/O=ApplicationName/CN=UserId"
````