# https://www.vultr.com/resources/faq/#downloadspeedtests

times=5
for i in "tx-us" "il-us" "ga-us" "tor-ca" "lax-ca-us" "mex-mx" "nj-us" "sjo-ca-us" "fl-us" "wa-us" "man-uk" "lon-gb" "ams-nl" "par-fr" "mad-es" "sto-se" "fra-de" "waw-pl" "scl-cl" "sao-br" "hnd-jp" "hon-hi-us" "osk-jp" "sel-kor" "tlv-il" "del-in" "bom-in" "syd-au" "blr-in" "mel-au" "jnb-za" "sgp"; do
    ping -q -c $times $i-ping.vultr.com
    #ping6 -q -c $times $i-ipv6.vultr.com
done
