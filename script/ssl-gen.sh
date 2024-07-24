rm ./cert/*.pem

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ./cert/ca-key.pem -out ./cert/ca-cert.pem -subj "/C=FR/ST=Occitanie/L=Toulouse/O=Tech School/OU=Education/CN=*.techschool.guru/emailAddress=techschool.guru@gmail.com"

echo "CA's self-signed certificate"
openssl x509 -in ./cert/ca-cert.pem -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout ./cert/server-key.pem -out ./cert/server-req.pem -subj "/C=FR/ST=Ile de France/L=Paris/O=PC Book/OU=Computer/CN=*.pcbook.com/emailAddress=pcbook@gmail.com"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in ./cert/server-req.pem -days 60 -CA ./cert/ca-cert.pem -CAkey ./cert/ca-key.pem -CAcreateserial -out ./cert/server-cert.pem -extfile ./cert/server-ext.cnf

echo "Server's signed certificate"
openssl x509 -in ./cert/server-cert.pem -noout -text

# 4. Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout ./cert/client-key.pem -out ./cert/client-req.pem -subj "/C=FR/ST=Alsace/L=Strasbourg/O=PC Client/OU=Computer/CN=*.pcclient.com/emailAddress=pcclient@gmail.com"

# 5. Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -in ./cert/client-req.pem -days 60 -CA ./cert/ca-cert.pem -CAkey ./cert/ca-key.pem -CAcreateserial -out ./cert/client-cert.pem -extfile ./cert/client-ext.cnf

echo "Client's signed certificate"
openssl x509 -in ./cert/client-cert.pem -noout -text