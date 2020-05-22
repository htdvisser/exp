all: authorized_keys known_hosts docker

clean:
	docker-compose down
	rm -f id_* ssh_host_* known_hosts authorized_keys

authorized_keys: clientkeys
	cat id_*.pub > authorized_keys

clientkeys: id_ecdsa.pub id_ed25519.pub id_rsa.pub

id_ecdsa.pub:
	ssh-keygen -t ecdsa -f id_ecdsa -N "" -C "Test Client"

id_ed25519.pub:
	ssh-keygen -t ed25519 -f id_ed25519 -N "" -C "Test Client"

id_rsa.pub:
	ssh-keygen -t rsa -b 4096 -f id_rsa -N "" -C "Test Client"

id_awskms.pub: id_awskms.pub.pem
	ssh-keygen -f $< -i -mPKCS8 > $@

id_awskms.pub.pem:
	aws kms get-public-key --key-id $$AWS_KMS_KEY_ID --query PublicKey | jq -r '"-----BEGIN PUBLIC KEY-----\n" + . + "\n-----END PUBLIC KEY-----"' > $@

known_hosts: hostkeys
	cat ssh_host_*_key.pub | sed -e 's/^/localhost:2222 /' > known_hosts

hostkeys: ssh_host_ecdsa_key.pub ssh_host_ed25519_key.pub ssh_host_rsa_key.pub

ssh_host_ed25519_key.pub:
	ssh-keygen -t ed25519 -f ssh_host_ed25519_key -N "" -C ""

ssh_host_ecdsa_key.pub:
	ssh-keygen -t ecdsa -f ssh_host_ecdsa_key -N "" -C ""

ssh_host_rsa_key.pub:
	ssh-keygen -t rsa -b 4096 -f ssh_host_rsa_key -N "" -C ""

docker:
	docker-compose build
