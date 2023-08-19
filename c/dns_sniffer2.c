#include <stdio.h>
#include <pcap.h>
#include <netinet/ip.h>
#include <netinet/ip6.h>
#include <netinet/udp.h>
#include <netinet/in.h>
#include <arpa/inet.h>

#define DNS_PORT 53

void parse_dns(const u_char *dns_payload, char* dst_str) {
    int index = 12; // Skip DNS header

    while (dns_payload[index] != 0) {
        if (dns_payload[index] >= 0xC0) {
            index += 2; // Compressed name, skip two bytes
            break;
        }

        int label_length = dns_payload[index++];
        for (int i = 0; i < label_length; ++i) {
            putchar(dns_payload[index++]);
        }

        if (dns_payload[index] != 0) {
            putchar('.');
        }
    }

    int dns_type = dns_payload[index+2];
    if (dns_type == 1) {
        printf(" 4");
    } else if (dns_type == 28) {
        printf(" 6");
    } else if (dns_type == 2) {
        printf(" NS");
    } else if (dns_type == 5) {
        printf(" CNAME");
    } else if (dns_type == 6) {
        printf(" SOA");
    } else if (dns_type == 11) {
        printf(" WKS");
    } else if (dns_type == 12) {
        printf(" PTR");
    } else if (dns_type == 13) {
        printf(" HINFO");
    } else if (dns_type == 15) {
        printf(" MX");
    } else if (dns_type == 16) {
        printf(" TXT");
    } else if (dns_type == 17) {
        printf(" RP");
    } else if (dns_type == 18) {
        printf(" AFSDB");
    } else if (dns_type == 24) {
        printf(" SIG");
    } else if (dns_type == 25) {
        printf(" KEY");
    } else if (dns_type == 29) {
        printf(" LOC");
    } else if (dns_type == 33) {
        printf(" SRV");
    } else if (dns_type == 35) {
        printf(" NAPTR");
    } else if (dns_type == 37) {
        printf(" CERT");
    } else if (dns_type == 39) {
        printf(" DNAME");
    } else if (dns_type == 42) {
        printf(" APL");
    } else if (dns_type == 43) {
        printf(" DS");
    } else if (dns_type == 44) {
        printf(" SSHFP");
    } else if (dns_type == 45) {
        printf(" IPSECKEY");
    } else if (dns_type == 46) {
        printf(" RRSIG");
    } else if (dns_type == 47) {
        printf(" NSEC");
    } else if (dns_type == 48) {
        printf(" DNSKEY");
    } else if (dns_type == 49) {
        printf(" DHCID");
    } else if (dns_type == 50) {
        printf(" NSEC3");
    } else if (dns_type == 51) {
        printf(" NSEC3PARAM");
    } else if (dns_type == 55) {
        printf(" HIP");
    } else if (dns_type == 59) {
        printf(" CDS");
    } else if (dns_type == 60) {
        printf(" CDNSKEY");
    } else if (dns_type == 61) {
        printf(" OPENPGPKEY");
    } else if (dns_type == 65) {
        printf(" HTTPS");
    } else if (dns_type == 99) {
        printf(" SPF");
    } else if (dns_type == 249) {
        printf(" TKEY");
    } else if (dns_type == 250) {
        printf(" TSIG");
    } else if (dns_type == 252) {
        printf(" AXFR");
    } else if (dns_type == 256) {
        printf(" URI");
    } else if (dns_type == 257) {
        printf(" CAA");
    } else if (dns_type == 32768) {
        printf(" TA");
    } else if (dns_type == 32769) {
        printf(" DLV");
    }

    printf(" %s ", dst_str);

    putchar('\n');
}

void packet_handler(u_char *user_data, const struct pcap_pkthdr *pkthdr, const u_char *packet) {
    struct ip *iph;
    struct ip6_hdr *ip6h;
    struct udphdr *udph;
    u_char *dns;
    char ipv6_str[INET6_ADDRSTRLEN];
    char* dst_str = ipv6_str;

    iph = (struct ip *)(packet + 14); // Skip Ethernet header

    if (iph->ip_v == 4) {
        udph = (struct udphdr *) (packet + 14 + iph->ip_hl * 4); // Skip IP header
        dns = (u_char * )(packet + 14 + iph->ip_hl * 4 + sizeof(struct udphdr));
        dst_str = inet_ntoa(iph->ip_dst);
    } else {
        ip6h = (struct ip6_hdr *)(packet + 14); // Skip Ethernet header
        udph = (struct udphdr *)(packet + 14 + sizeof(struct ip6_hdr)); // Skip IPv6 header
        dns = (u_char *)(packet + 14 + sizeof(struct ip6_hdr) + sizeof(struct udphdr));
        inet_ntop(AF_INET6, &(ip6h->ip6_dst), ipv6_str, INET6_ADDRSTRLEN);
    }

    if (ntohs(udph->uh_dport) == DNS_PORT) {
        parse_dns(dns, dst_str);
    }
}

int main() {
    char errbuf[PCAP_ERRBUF_SIZE];
    pcap_t *handle;
    char *dev;

    // Get a network device to capture packets
    dev = pcap_lookupdev(errbuf);
    if (dev == NULL) {
        printf("Error finding network device: %s\n", errbuf);
        return 1;
    }

    // Open the network device for packet capture
    handle = pcap_open_live(dev, BUFSIZ, 1, 1000, errbuf);
    if (handle == NULL) {
        printf("Error opening network device: %s\n", errbuf);
        return 1;
    }

    // Set a filter to capture only DNS packets
    struct bpf_program fp;
    char filter_exp[] = "udp port 53";
    if (pcap_compile(handle, &fp, filter_exp, 0, PCAP_NETMASK_UNKNOWN) == -1) {
        printf("Error compiling filter: %s\n", pcap_geterr(handle));
        return 1;
    }
    if (pcap_setfilter(handle, &fp) == -1) {
        printf("Error setting filter: %s\n", pcap_geterr(handle));
        return 1;
    }

    // Start capturing packets and call packet_handler for each captured packet
    pcap_loop(handle, 0, packet_handler, NULL);

    // Close the packet capture handle
    pcap_close(handle);

    return 0;
}
