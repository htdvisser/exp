FROM alpine

RUN apk add --update --no-cache bash openssh

RUN addgroup -S testgroup && adduser -S testuser -G testgroup -s /bin/bash && echo 'testuser:testpassword' | chpasswd

ADD ssh_host_* /etc/ssh/
ADD authorized_keys /home/testuser/.ssh/authorized_keys

CMD ["/usr/sbin/sshd", "-D", "-e", "-f", "/root/test/sshd_config"]
