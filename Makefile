.PHONY: certs
certs: server_certs client_certs

.PHONY: server_certs
server_certs:
	mkdir -p tls && \
	openssl req -subj /C=/ST=/O=/L=/CN=localhost/OU=/ -x509 -nodes -days 3650  -newkey rsa:4096 -keyout tls/server_key.pem -out tls/server_cert.pem

.PHONY: client_certs
client_certs:
	mkdir -p tls && \
	openssl req -subj /C=/ST=/O=/L=/CN=localhost/OU=/ -x509 -nodes -days 3650  -newkey rsa:4096 -keyout tls/client_key.pem -out tls/client_cert.pem

.PHONY: clean
clean:
	rm -rf tls
