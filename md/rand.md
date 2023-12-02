echo $RANDOM |md5sum |cut -c 1-8
openssl rand -base64 8
date +%s%N
head /dev/urandom | cksum
cat /proc/sys/kernel/random/uuid

yum -y install expect
mkpasswd -l 8

