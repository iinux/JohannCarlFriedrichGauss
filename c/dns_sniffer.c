#include <stdio.h>
#include <pcap.h>
#include <netinet/ip6.h>
#include <netinet/udp.h>
#include <arpa/inet.h>

#define DNS_PORT 53

void parse_dns(const u_char *dns_payload) {
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

    putchar('\n');
}

void packet_handler(u_char *user_data, const struct pcap_pkthdr *pkthdr, const u_char *packet) {
    struct ip6_hdr *ip6h;
    struct udphdr *udph;
    u_char *dns_payload;

    ip6h = (struct ip6_hdr *)(packet + 14); // Skip Ethernet header
    udph = (struct udphdr *)(packet + 14 + sizeof(struct ip6_hdr)); // Skip IPv6 header
    dns_payload = (u_char *)(packet + 14 + sizeof(struct ip6_hdr) + sizeof(struct udphdr));

    if (ntohs(udph->uh_dport) == DNS_PORT && dns_payload[2] & 0x80) {
        parse_dns(dns_payload);
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

    // Set a filter to capture only DNS response packets
    struct bpf_program fp;
    char filter_exp[] = "udp port 53 and dst host <your_ip_address>";
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
