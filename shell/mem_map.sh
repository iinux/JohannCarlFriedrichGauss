# Set pid of nginx master process here
pid=8192

# generate gdb commands from the process's memory mappings using awk
cat /proc/$pid/maps | awk '$6 !~ "^/" {split ($1,addrs,"-"); print "dump memory mem_" addrs[1] " 0x" addrs[1] " 0x" addrs[2] ;}END{print "quit"}' > gdb-commands

# use gdb with the -x option to dump these memory regions to mem_* files
gdb -p $pid -x gdb-commands

# look for some (any) nginx.conf text
grep worker_connections mem_*
grep server_name mem_*

# https://serverfault.com/questions/361421/dump-nginx-config-from-running-process
