
echo "creating server private key and node id"
openssl genpkey -algorithm x25519 -out x25519-priv.pem
openssl pkey -in x25519-priv.pem -pubout -out x25519-pub.pem

