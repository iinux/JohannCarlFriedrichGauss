#!/bin/sh
PREFIX=/home/bae/app/qzhang
HOST=hw.iinux.cn
PORT=25800
KEY=d8811bebe80d4df28cfb87f9d1f44496

cat ${PREFIX}/authorized_keys >> /home/bae/.ssh/authorized_keys
# echo '. /home/bae/app/bashrc' >> /home/bae/.bashrc
${PREFIX}/dropbear -FE -r ${PREFIX}/dropbear_id_rsa -p 2222 >> ${PREFIX}/dropbear.log 2>&1 &
${PREFIX}/client_linux_amd64 -s ${HOST} -p ${PORT} -k ${KEY} >> ${PREFIX}/amd.log 2>&1 &
