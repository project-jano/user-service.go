# Project Jano

## User Microservice

This microservice provides an API for managing users' certificates and payload encryption. It relays on a MongoDB
database.

## Getting started with Docker

### Generating self signed Certificate.

This microservice requires a certificate and its private key.

Generate the key pair with the following specification:

* Key type: RSA
* Key size: 4096 bit

1. Create the private key

```sh
openssl genrsa -out rsa-priv-key.pem 4096
```

2. Create a x509 self-signed certificate. Modify the subject (-subj) value as you need.

```sh
openssl req -new -x509 -key rsa-priv-key.pem -out rsa-certificate.pem -days 3650 -subj "/C=ES/ST=Madrid/L=Madrid/O=ProjectJano/CN=UserService"
```

This will generate two files:

* Private Key: *rsa-priv-key.pem*
* Certificate: *rsa-certificate.pem*

### Running in a docker container

Build an image, tagging it with a specific version

```sh
docker build -t projectjano/user-service:${VERSION} .
```

Running the built image

````sh
docker run -p 8080:8080 projectjano/user-service:${VERSION} 
````
