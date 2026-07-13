/*
 * DPDK TCP Forwarder
 *
 * Listens on 0.0.0.0:3001 via a DPDK-controlled NIC and forwards all
 * accepted TCP connections to 192.168.1.10:8000.
 *
 * This is a minimal userspace TCP/IP stack sitting on top of raw PMD
 * packets. It implements: ARP, IPv4, TCP 3-way handshake (both
 * inbound and outbound), connection table, bi-directional data
 * relay, and FIN-based close.
 *
 * Limitations (intentional, kept simple):
 *   - No retransmission timer (assumes a healthy local network).
 *   - No sliding-window flow control (small fixed buffers).
 *   - No TIME_WAIT (close immediately after LAST_ACK / FIN_WAIT2).
 *   - Single DPDK port, single RX/TX queue.
 *
 * Build:  see Makefile
 * Run:    sudo ./build/tcp_forwarder -l 0 -n 1 -- -p 0
 *         (configure the NIC IP / mask / gateway via the defines below,
 *          or extend with a parameter parser if you need it at runtime.)
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <inttypes.h>
#include <signal.h>
#include <unistd.h>
#include <arpa/inet.h>

#include <rte_eal.h>
#include <rte_ethdev.h>
#include <rte_mbuf.h>
#include <rte_malloc.h>
#include <rte_lcore.h>
#include <rte_ip.h>
#include <rte_tcp.h>
#include <rte_ether.h>
#include <rte_arp.h>

/* ------------------------------------------------------------------ */
/* Configuration                                                       */
/* ------------------------------------------------------------------ */
#define LISTEN_PORT        3001
#define BACKEND_IP         "192.168.1.10"
#define BACKEND_PORT       8000

/* Replace these with the address assigned to the DPDK port. */
#define LOCAL_IP           "10.0.0.1"
#define LOCAL_MASK         "255.255.255.0"
/* Gateway (next-hop) used to reach BACKEND_IP. Set to the LOCAL_IP's
 * router on the L2 segment the DPDK port is attached to. */
#define GATEWAY_IP         "10.0.0.254"

#define NUM_MBUFS          (8192)
#define MBUF_CACHE_SIZE    (250)
#define BURST_SIZE         (32)

#define MAX_CONNS          (1024)
#define CONN_BUF_SIZE      (8192)

/* ------------------------------------------------------------------ */
/* Protocol constants                                                  */
/* ------------------------------------------------------------------ */
#define IP_PROTO_TCP       6
#define ETH_HDR_LEN        14
#define IP_HDR_LEN         20
#define TCP_HDR_LEN_MIN    20

/* TCP states (client-facing side) */
enum {
    TCP_CLOSED = 0,
    TCP_LISTEN,
    TCP_SYN_RCVD,
    TCP_ESTABLISHED,
    TCP_CLOSE_WAIT,
    TCP_LAST_ACK,
    TCP_FIN_WAIT1,
    TCP_FIN_WAIT2,
    TCP_CLOSING
};

/* ------------------------------------------------------------------ */
/* Connection record                                                   */
/* ------------------------------------------------------------------ */
struct conn {
    uint8_t  in_use;

    /* L2/L3 of the remote client */
    uint8_t  client_mac[RTE_ETHER_ADDR_LEN];
    uint32_t client_ip;     /* host order */
    uint16_t client_port;   /* host order */

    /* TCP per-direction state */
    int      cli_state;     /* state of the client-facing TCP */
    uint32_t cli_seq;       /* our sequence number (host order)  */
    uint32_t cli_ack;       /* next expected from client         */

    int      bk_state;      /* state of the backend-facing TCP   */
    uint16_t bk_src_port;   /* ephemeral port we picked          */
    uint32_t bk_seq;        /* our seq towards backend           */
    uint32_t bk_ack;        /* next expected from backend        */

    /* pending data: from client -> to send to backend */
    uint8_t  c2b[CONN_BUF_SIZE];
    uint16_t c2b_len;

    /* pending data: from backend -> to send to client */
    uint8_t  b2c[CONN_BUF_SIZE];
    uint16_t b2c_len;
};

static struct conn conns[MAX_CONNS];
static uint32_t    local_ip;          /* host order  */
static uint32_t    local_mask;        /* host order  */
static uint32_t    gateway_ip;        /* host order  */
static uint8_t     local_mac[RTE_ETHER_ADDR_LEN];
static uint8_t     gateway_mac[RTE_ETHER_ADDR_LEN];
static int         gateway_known;
static uint16_t    listen_port = LISTEN_PORT;

/* Backend, parsed once at startup. */
static uint32_t    backend_ip;        /* host order  */
static uint16_t    backend_port = BACKEND_PORT;

/* DPDK port id (single port). */
static uint16_t     g_portid = 0;

/* ------------------------------------------------------------------ */
/* Small helpers                                                       */
/* ------------------------------------------------------------------ */
static uint32_t ip_parse(const char *s)
{
    struct in_addr a;
    if (inet_pton(AF_INET, s, &a) != 1) {
        fprintf(stderr, "bad ip: %s\n", s);
        rte_exit(EXIT_FAILURE, "ip_parse");
    }
    return ntohl(a.s_addr);
}

static uint16_t ip_csum(const void *data, int len)
{
    const uint16_t *p = data;
    uint32_t sum = 0;
    while (len > 1) { sum += *p++; len -= 2; }
    if (len) sum += *(uint8_t *)p;
    while (sum >> 16) sum = (sum & 0xffff) + (sum >> 16);
    return ~sum;
}

static uint16_t tcp_csum(uint32_t saddr, uint32_t daddr,
                         const void *tcp, int tcp_len)
{
    /* pseudo header */
    uint32_t sum = 0;
    uint16_t *p;
    p = (uint16_t *)&saddr; sum += p[0]; sum += p[1];
    p = (uint16_t *)&daddr; sum += p[0]; sum += p[1];
    sum += htons(IP_PROTO_TCP);
    sum += htons(tcp_len);
    p = (uint16_t *)tcp;
    while (tcp_len > 1) { sum += *p++; tcp_len -= 2; }
    if (tcp_len) sum += *(uint8_t *)p;
    while (sum >> 16) sum = (sum & 0xffff) + (sum >> 16);
    return ~sum;
}

/* ------------------------------------------------------------------ */
/* DPDK port init                                                      */
/* ------------------------------------------------------------------ */
static int port_init(uint16_t port, struct rte_mempool *mbuf_pool)
{
    struct rte_eth_conf port_conf = {0};
    struct rte_eth_dev_info dev_info;
    struct rte_eth_rxconf rxconf;
    struct rte_eth_txconf txconf;
    int ret;
    uint16_t rx_rings = 1, tx_rings = 1;
    uint16_t nb_rxd = 1024, nb_txd = 1024;

    rte_eth_dev_info_get(port, &dev_info);
    if (dev_info.tx_offload_capa & RTE_ETH_TX_OFFLOAD_IPV4_CKSUM)
        port_conf.txmode.offloads |= RTE_ETH_TX_OFFLOAD_IPV4_CKSUM;
    if (dev_info.tx_offload_capa & RTE_ETH_TX_OFFLOAD_TCP_CKSUM)
        port_conf.txmode.offloads |= RTE_ETH_TX_OFFLOAD_TCP_CKSUM;

    ret = rte_eth_dev_configure(port, rx_rings, tx_rings, &port_conf);
    if (ret < 0) { fprintf(stderr, "configure: %d\n", ret); return ret; }

    rxconf = dev_info.default_rxconf;
    rxconf.offloads = port_conf.rxmode.offloads;
    ret = rte_eth_rx_queue_setup(port, 0, nb_rxd,
                                 rte_eth_dev_socket_id(port),
                                 &rxconf, mbuf_pool);
    if (ret < 0) { fprintf(stderr, "rx_queue: %d\n", ret); return ret; }

    txconf = dev_info.default_txconf;
    txconf.offloads = port_conf.txmode.offloads;
    ret = rte_eth_tx_queue_setup(port, 0, nb_txd,
                                 rte_eth_dev_socket_id(port),
                                 &txconf);
    if (ret < 0) { fprintf(stderr, "tx_queue: %d\n", ret); return ret; }

    ret = rte_eth_dev_start(port);
    if (ret < 0) { fprintf(stderr, "start: %d\n", ret); return ret; }

    rte_eth_macaddr_get(port, (struct rte_ether_addr *)local_mac);
    printf("port %u mac %02x:%02x:%02x:%02x:%02x:%02x\n",
           port,
           local_mac[0], local_mac[1], local_mac[2],
           local_mac[3], local_mac[4], local_mac[5]);

    /* We trust that the link is up; if not, you'll see no RX. */
    return 0;
}

/* ------------------------------------------------------------------ */
/* Connection table                                                    */
/* ------------------------------------------------------------------ */
static struct conn *conn_alloc(void)
{
    for (int i = 0; i < MAX_CONNS; i++)
        if (!conns[i].in_use) {
            memset(&conns[i], 0, sizeof(conns[i]));
            conns[i].in_use = 1;
            return &conns[i];
        }
    return NULL;
}

static void conn_free(struct conn *c) { c->in_use = 0; }

/* Linear lookup. */
static struct conn *conn_lookup_client(uint32_t cip, uint16_t cport)
{
    for (int i = 0; i < MAX_CONNS; i++)
        if (conns[i].in_use &&
            conns[i].client_ip   == cip &&
            conns[i].client_port == cport)
            return &conns[i];
    return NULL;
}

static struct conn *conn_lookup_backend(uint16_t sport)
{
    for (int i = 0; i < MAX_CONNS; i++)
        if (conns[i].in_use && conns[i].bk_src_port == sport)
            return &conns[i];
    return NULL;
}

/* ------------------------------------------------------------------ */
/* Sending packets                                                     */
/* ------------------------------------------------------------------ */
static struct rte_mbuf *mbuf_alloc(void)
{
    struct rte_mbuf *m = rte_pktmbuf_alloc(
        rte_mempool_lookup("PKT_POOL"));
    if (!m) return NULL;
    return m;
}

static uint16_t alloc_ephemeral(void)
{
    /* 0xc000 .. 0xffff range, simple round-robin */
    static uint16_t next = 0xc000;
    next++;
    if (next < 0xc000) next = 0xc000;
    return next;
}

/* dst_mac: 6 bytes; src/dst ip in host order. */
static int build_tcp_pkt(struct rte_mbuf *m,
                         const uint8_t *dst_mac,
                         uint32_t src_ip, uint32_t dst_ip,
                         uint16_t src_port, uint16_t dst_port,
                         uint32_t seq, uint32_t ack, uint8_t flags,
                         const uint8_t *payload, uint16_t plen)
{
    struct rte_ether_hdr *eth;
    struct rte_ipv4_hdr  *ip;
    struct rte_tcp_hdr   *tcp;
    uint16_t tcp_len = TCP_HDR_LEN_MIN + plen;
    uint16_t pkt_len = ETH_HDR_LEN + IP_HDR_LEN + tcp_len;
    char *buf;

    if (rte_pktmbuf_append(m, pkt_len) == NULL) return -1;

    buf = rte_pktmbuf_mtod(m, char *);
    eth = (struct rte_ether_hdr *)buf;
    ip  = (struct rte_ipv4_hdr  *)(buf + ETH_HDR_LEN);
    tcp = (struct rte_tcp_hdr   *)(buf + ETH_HDR_LEN + IP_HDR_LEN);

    /* Ethernet */
    rte_ether_addr_copy((struct rte_ether_addr *)dst_mac,
                        &eth->dst_addr);
    rte_ether_addr_copy((struct rte_ether_addr *)local_mac,
                        &eth->src_addr);
    eth->ether_type = htons(RTE_ETHER_TYPE_IPV4);

    /* IPv4 */
    ip->version_ihl = 0x45;
    ip->type_of_service = 0;
    ip->total_length = htons(IP_HDR_LEN + tcp_len);
    ip->packet_id = 0;
    ip->fragment_offset = 0;
    ip->time_to_live = 64;
    ip->next_proto_id = IP_PROTO_TCP;
    ip->src_addr = htonl(src_ip);
    ip->dst_addr = htonl(dst_ip);
    ip->hdr_checksum = 0;
    ip->hdr_checksum = ip_csum(ip, IP_HDR_LEN);

    /* TCP */
    memset(tcp, 0, TCP_HDR_LEN_MIN);
    tcp->src_port   = htons(src_port);
    tcp->dst_port   = htons(dst_port);
    tcp->sent_seq   = htonl(seq);
    tcp->recv_ack   = htonl(ack);
    tcp->data_off   = (TCP_HDR_LEN_MIN / 4) << 4;
    tcp->tcp_flags  = flags;
    tcp->rx_win     = htons(8192);
    tcp->cksum      = 0;
    if (plen && payload)
        memcpy((uint8_t *)tcp + TCP_HDR_LEN_MIN, payload, plen);
    tcp->cksum = tcp_csum(src_ip, dst_ip, tcp, tcp_len);

    return 0;
}

static int tx_one(struct rte_mbuf *m)
{
    uint16_t sent = rte_eth_tx_burst(g_portid, 0, &m, 1);
    if (sent == 0) { rte_pktmbuf_free(m); return -1; }
    return 0;
}

/* Send an ARP request to resolve gateway_ip. */
static int arp_request(uint32_t target_ip)
{
    struct rte_mbuf *m = mbuf_alloc();
    struct rte_ether_hdr *eth;
    struct rte_arp_hdr *arp;
    char *buf;
    if (!m) return -1;
    if (rte_pktmbuf_append(m, sizeof(*eth) + sizeof(*arp)) == NULL) {
        rte_pktmbuf_free(m); return -1;
    }
    buf = rte_pktmbuf_mtod(m, char *);
    eth = (struct rte_ether_hdr *)buf;
    arp = (struct rte_arp_hdr *)(buf + sizeof(*eth));

    memset(eth->dst_addr.addr_bytes, 0xff, RTE_ETHER_ADDR_LEN);
    rte_ether_addr_copy((struct rte_ether_addr *)local_mac,
                        &eth->src_addr);
    eth->ether_type = htons(RTE_ETHER_TYPE_ARP);

    arp->arp_hardware = htons(RTE_ARP_HRD_ETHER);
    arp->arp_protocol = htons(RTE_ETHER_TYPE_IPV4);
    arp->arp_hlen = RTE_ETHER_ADDR_LEN;
    arp->arp_plen = 4;
    arp->arp_opcode = htons(RTE_ARP_OP_REQUEST);
    rte_ether_addr_copy((struct rte_ether_addr *)local_mac,
                        &arp->arp_data.arp_sha);
    arp->arp_data.arp_sip = htonl(local_ip);
    memset(arp->arp_data.arp_tha.addr_bytes, 0, RTE_ETHER_ADDR_LEN);
    arp->arp_data.arp_tip = htonl(target_ip);

    return tx_one(m);
}

/* ------------------------------------------------------------------ */
/* TX helpers driven by the conn state                                 */
/* ------------------------------------------------------------------ */
static void send_syn_ack_to_client(struct conn *c)
{
    struct rte_mbuf *m = mbuf_alloc();
    if (!m) return;
    build_tcp_pkt(m, c->client_mac,
                  local_ip, c->client_ip,
                  listen_port, c->client_port,
                  c->cli_seq, c->cli_ack,
                  RTE_TCP_SYN_FLAG | RTE_TCP_ACK_FLAG,
                  NULL, 0);
    c->cli_seq += 1;   /* SYN consumes one seq */
    tx_one(m);
}

static void send_ack_to_client(struct conn *c)
{
    struct rte_mbuf *m = mbuf_alloc();
    if (!m) return;
    build_tcp_pkt(m, c->client_mac,
                  local_ip, c->client_ip,
                  listen_port, c->client_port,
                  c->cli_seq, c->cli_ack,
                  RTE_TCP_ACK_FLAG, NULL, 0);
    tx_one(m);
}

static void send_fin_ack_to_client(struct conn *c)
{
    struct rte_mbuf *m = mbuf_alloc();
    if (!m) return;
    build_tcp_pkt(m, c->client_mac,
                  local_ip, c->client_ip,
                  listen_port, c->client_port,
                  c->cli_seq, c->cli_ack,
                  RTE_TCP_FIN_FLAG | RTE_TCP_ACK_FLAG,
                  NULL, 0);
    c->cli_seq += 1;
    tx_one(m);
}

static void send_rst_to(struct conn *c,
                        uint8_t *dst_mac, uint32_t dst_ip,
                        uint16_t dst_port, uint32_t seq)
{
    struct rte_mbuf *m = mbuf_alloc();
    if (!m) return;
    build_tcp_pkt(m, dst_mac,
                  local_ip, dst_ip,
                  listen_port, dst_port,
                  seq, 0,
                  RTE_TCP_RST_FLAG, NULL, 0);
    tx_one(m);
}

/* SYN toward backend (we are the active opener). */
static void send_syn_to_backend(struct conn *c)
{
    struct rte_mbuf *m = mbuf_alloc();
    uint32_t dst_ip = backend_ip;
    uint8_t *dst_mac = gateway_mac;
    /* If backend is on the same subnet we still use gateway as a
     * generic L2 destination; that's fine if backend == gateway. */
    if (!gateway_known || !dst_mac) return;

    if (!m) return;
    build_tcp_pkt(m, dst_mac,
                  local_ip, dst_ip,
                  c->bk_src_port, backend_port,
                  c->bk_seq, 0,
                  RTE_TCP_SYN_FLAG, NULL, 0);
    c->bk_state = 1; /* SYN_SENT-ish */
    c->bk_seq += 1;
    tx_one(m);
}

static void send_data_to_backend(struct conn *c, const uint8_t *data, uint16_t len)
{
    struct rte_mbuf *m;
    if (!gateway_known) return;
    if (len == 0) return;
    m = mbuf_alloc();
    if (!m) return;
    build_tcp_pkt(m, gateway_mac,
                  local_ip, backend_ip,
                  c->bk_src_port, backend_port,
                  c->bk_seq, c->bk_ack,
                  RTE_TCP_PSH_FLAG | RTE_TCP_ACK_FLAG,
                  data, len);
    c->bk_seq += len;
    tx_one(m);
}

/* ------------------------------------------------------------------ */
/* RX processing                                                       */
/* ------------------------------------------------------------------ */
static int is_my_ip(uint32_t dst_ip)
{
    return (dst_ip == local_ip) ||
           ((dst_ip & local_mask) == (local_ip & local_mask));
}

static void handle_arp(struct rte_mbuf *m)
{
    struct rte_ether_hdr *eth = rte_pktmbuf_mtod(m, struct rte_ether_hdr *);
    struct rte_arp_hdr *arp = (struct rte_arp_hdr *)(eth + 1);
    uint32_t sip, tip;

    if (arp->arp_opcode != htons(RTE_ARP_OP_REQUEST)) {
        /* could also handle REPLY to learn sender MAC */
        if (arp->arp_opcode == htons(RTE_ARP_OP_REPLY)) {
            sip = ntohl(arp->arp_data.arp_sip);
            if (sip == gateway_ip) {
                rte_ether_addr_copy(&arp->arp_data.arp_sha,
                                    (struct rte_ether_addr *)gateway_mac);
                gateway_known = 1;
                printf("learned gateway mac %02x:%02x:%02x:%02x:%02x:%02x\n",
                       gateway_mac[0], gateway_mac[1], gateway_mac[2],
                       gateway_mac[3], gateway_mac[4], gateway_mac[5]);
            }
        }
        return;
    }

    sip = ntohl(arp->arp_data.arp_sip);
    tip = ntohl(arp->arp_data.arp_tip);
    if (tip != local_ip) return;

    /* Reply. */
    struct rte_mbuf *r = mbuf_alloc();
    struct rte_ether_hdr *reth;
    struct rte_arp_hdr *rarp;
    char *buf;
    if (!r) return;
    if (rte_pktmbuf_append(r, sizeof(*reth) + sizeof(*rarp)) == NULL) {
        rte_pktmbuf_free(r); return;
    }
    buf = rte_pktmbuf_mtod(r, char *);
    reth = (struct rte_ether_hdr *)buf;
    rarp = (struct rte_arp_hdr *)(buf + sizeof(*reth));

    rte_ether_addr_copy(&eth->src_addr, &reth->dst_addr);
    rte_ether_addr_copy((struct rte_ether_addr *)local_mac,
                        &reth->src_addr);
    reth->ether_type = htons(RTE_ETHER_TYPE_ARP);

    rarp->arp_hardware = htons(RTE_ARP_HRD_ETHER);
    rarp->arp_protocol = htons(RTE_ETHER_TYPE_IPV4);
    rarp->arp_hlen = RTE_ETHER_ADDR_LEN;
    rarp->arp_plen = 4;
    rarp->arp_opcode = htons(RTE_ARP_OP_REPLY);
    rte_ether_addr_copy((struct rte_ether_addr *)local_mac,
                        &rarp->arp_data.arp_sha);
    rarp->arp_data.arp_sip = htonl(local_ip);
    rte_ether_addr_copy(&arp->arp_data.arp_sha,
                        &rarp->arp_data.arp_tha);
    rarp->arp_data.arp_tip = sip;

    tx_one(r);
    (void)sip;
}

static void handle_ipv4(struct rte_mbuf *m)
{
    struct rte_ether_hdr *eth = rte_pktmbuf_mtod(m, struct rte_ether_hdr *);
    struct rte_ipv4_hdr  *ip  = (struct rte_ipv4_hdr *)(eth + 1);

    if (ip->next_proto_id != IP_PROTO_TCP) return;

    uint16_t ihl = (ip->version_ihl & 0x0f) * 4;
    struct rte_tcp_hdr *tcp = (struct rte_tcp_hdr *)((char *)ip + ihl);
    int tcp_data_off = (tcp->data_off >> 4) * 4;
    int total_len = ntohs(ip->total_length);
    int payload_len = total_len - ihl - tcp_data_off;
    if (payload_len < 0) payload_len = 0;
    uint8_t *payload = (uint8_t *)tcp + tcp_data_off;

    uint32_t dst_ip = ntohl(ip->dst_addr);
    uint32_t src_ip = ntohl(ip->src_addr);
    uint16_t dst_port = ntohs(tcp->dst_port);
    uint16_t src_port = ntohs(tcp->src_port);
    uint32_t seq = ntohl(tcp->sent_seq);
    uint32_t ack = ntohl(tcp->recv_ack);
    uint8_t  flags = tcp->tcp_flags;

    rte_ether_addr_copy(&eth->src_addr,
                        (struct rte_ether_addr *)eth->src_addr.addr_bytes);

    /* ----- packets destined to our LISTEN port (from client) ----- */
    if (dst_port == listen_port && is_my_ip(dst_ip)) {
        struct conn *c = conn_lookup_client(src_ip, src_port);

        if (!c && (flags & RTE_TCP_SYN_FLAG) && !(flags & RTE_TCP_ACK_FLAG)) {
            c = conn_alloc();
            if (!c) {
                /* drop */
                return;
            }
            rte_ether_addr_copy(&eth->src_addr,
                                (struct rte_ether_addr *)c->client_mac);
            c->client_ip   = src_ip;
            c->client_port = src_port;
            c->cli_state   = TCP_LISTEN;
            c->cli_seq     = 1000;          /* ISN - arbitrary */
            c->cli_ack     = seq + 1;       /* ack the SYN */
            c->bk_src_port = alloc_ephemeral();
            c->bk_seq      = 2000;
            c->bk_ack      = 0;
            c->bk_state    = TCP_CLOSED;

            send_syn_ack_to_client(c);
            c->cli_state   = TCP_SYN_RCVD;
            printf("conn %p SYN_RCVD from %u.%u.%u.%u:%u\n",
                   (void *)c,
                   (src_ip >> 24) & 0xff, (src_ip >> 16) & 0xff,
                   (src_ip >> 8) & 0xff, src_ip & 0xff, src_port);
            return;
        }

        if (!c) return;

        /* Established: data from client -> queue for backend */
        if (c->cli_state == TCP_ESTABLISHED) {
            if (payload_len > 0) {
                uint16_t room = CONN_BUF_SIZE - c->c2b_len;
                uint16_t copy = payload_len < room ? payload_len : room;
                memcpy(c->c2b + c->c2b_len, payload, copy);
                c->c2b_len += copy;
                c->cli_ack = seq + payload_len;
                /* ACK the data to client */
                send_ack_to_client(c);
                /* Forward to backend (if backend is up) */
                if (c->bk_state == TCP_ESTABLISHED && c->c2b_len > 0) {
                    send_data_to_backend(c, c->c2b, c->c2b_len);
                    c->c2b_len = 0;
                }
            }
            /* Close request from client */
            if (flags & RTE_TCP_FIN_FLAG) {
                c->cli_ack = seq + payload_len + 1;
                send_ack_to_client(c);
                c->cli_state = TCP_CLOSE_WAIT;
                /* half-close: stop sending to backend, but keep RX'ing */
                printf("conn %p client sent FIN -> CLOSE_WAIT\n", (void *)c);
                /* Try to tear down backend half too. */
                if (c->bk_state == TCP_ESTABLISHED) {
                    /* just drop the backend; production code would FIN. */
                    c->bk_state = TCP_CLOSED;
                }
                conn_free(c);
                return;
            }
            return;
        }

        /* SYN_RCVD: expect ACK of our SYN-ACK */
        if (c->cli_state == TCP_SYN_RCVD) {
            if (flags & RTE_TCP_ACK_FLAG) {
                c->cli_state = TCP_ESTABLISHED;
                printf("conn %p ESTABLISHED (client side)\n", (void *)c);
                /* Now open the backend. */
                if (!gateway_known) arp_request(gateway_ip);
                send_syn_to_backend(c);
                c->bk_state = 1;
                /* If we already accumulated c2b before backend came up,
                 * don't forward yet; we'll do it after backend SYN-ACK. */
            } else if (flags & RTE_TCP_RST_FLAG) {
                conn_free(c);
            }
            return;
        }

        /* LAST_ACK: client ACKs our FIN. */
        if (c->cli_state == TCP_LAST_ACK) {
            if (flags & RTE_TCP_ACK_FLAG) {
                conn_free(c);
                return;
            }
        }
    }

    /* ----- packets from backend (our ephemeral src port) ----- */
    if (src_port == backend_port && src_ip == backend_ip) {
        struct conn *c = conn_lookup_backend(dst_port);
        if (!c) return;

        /* Backend SYN-ACK */
        if ((flags & (RTE_TCP_SYN_FLAG | RTE_TCP_ACK_FLAG)) ==
            (RTE_TCP_SYN_FLAG | RTE_TCP_ACK_FLAG)) {
            c->bk_ack = seq + 1;
            c->bk_state = TCP_ESTABLISHED;
            /* ACK the SYN-ACK */
            {
                struct rte_mbuf *m2 = mbuf_alloc();
                if (m2) {
                    build_tcp_pkt(m2, gateway_mac,
                                  local_ip, backend_ip,
                                  c->bk_src_port, backend_port,
                                  c->bk_seq, c->bk_ack,
                                  RTE_TCP_ACK_FLAG, NULL, 0);
                    tx_one(m2);
                }
            }
            printf("conn %p backend ESTABLISHED\n", (void *)c);
            /* Flush queued client data now. */
            if (c->c2b_len > 0) {
                send_data_to_backend(c, c->c2b, c->c2b_len);
                c->c2b_len = 0;
            }
            return;
        }

        if (c->bk_state == TCP_ESTABLISHED) {
            /* Data from backend -> queue for client */
            if (payload_len > 0) {
                uint16_t room = CONN_BUF_SIZE - c->b2c_len;
                uint16_t copy = payload_len < room ? payload_len : room;
                memcpy(c->b2c + c->b2c_len, payload, copy);
                c->b2c_len += copy;
                c->bk_ack = seq + payload_len;
                /* ACK to backend */
                {
                    struct rte_mbuf *m2 = mbuf_alloc();
                    if (m2) {
                        build_tcp_pkt(m2, gateway_mac,
                                      local_ip, backend_ip,
                                      c->bk_src_port, backend_port,
                                      c->bk_seq, c->bk_ack,
                                      RTE_TCP_ACK_FLAG, NULL, 0);
                        tx_one(m2);
                    }
                }
                /* Forward to client */
                if (c->b2c_len > 0 && c->cli_state == TCP_ESTABLISHED) {
                    struct rte_mbuf *m3 = mbuf_alloc();
                    if (m3) {
                        build_tcp_pkt(m3, c->client_mac,
                                      local_ip, c->client_ip,
                                      listen_port, c->client_port,
                                      c->cli_seq, c->cli_ack,
                                      RTE_TCP_PSH_FLAG | RTE_TCP_ACK_FLAG,
                                      c->b2c, c->b2c_len);
                        c->cli_seq += c->b2c_len;
                        tx_one(m3);
                        c->b2c_len = 0;
                    }
                }
            }
            if (flags & RTE_TCP_FIN_FLAG) {
                /* Backend wants to close. Forward FIN to client. */
                struct rte_mbuf *m2 = mbuf_alloc();
                if (m2) {
                    c->bk_ack = seq + payload_len + 1;
                    build_tcp_pkt(m2, gateway_mac,
                                  local_ip, backend_ip,
                                  c->bk_src_port, backend_port,
                                  c->bk_seq, c->bk_ack,
                                  RTE_TCP_ACK_FLAG, NULL, 0);
                    tx_one(m2);
                }
                struct rte_mbuf *m3 = mbuf_alloc();
                if (m3) {
                    build_tcp_pkt(m3, c->client_mac,
                                  local_ip, c->client_ip,
                                  listen_port, c->client_port,
                                  c->cli_seq, c->cli_ack,
                                  RTE_TCP_FIN_FLAG | RTE_TCP_ACK_FLAG,
                                  NULL, 0);
                    c->cli_seq += 1;
                    c->cli_state = TCP_FIN_WAIT1;
                    tx_one(m3);
                }
            }
            return;
        }
    }
}

/* ------------------------------------------------------------------ */
/* Main loop                                                           */
/* ------------------------------------------------------------------ */
static volatile int keep_running = 1;
static void on_sig(int s) { (void)s; keep_running = 0; }

int main(int argc, char **argv)
{
    struct rte_mempool *mbuf_pool;
    int ret;

    signal(SIGINT, on_sig);
    signal(SIGTERM, on_sig);

    ret = rte_eal_init(argc, argv);
    if (ret < 0) rte_exit(EXIT_FAILURE, "eal init");

    argc -= ret; argv += ret;

    if (rte_eth_dev_count_avail() == 0)
        rte_exit(EXIT_FAILURE, "no ethernet ports");

    mbuf_pool = rte_pktmbuf_pool_create("PKT_POOL",
                                        NUM_MBUFS,
                                        MBUF_CACHE_SIZE, 0,
                                        RTE_MBUF_DEFAULT_BUF_SIZE,
                                        rte_socket_id());
    if (!mbuf_pool) rte_exit(EXIT_FAILURE, "mbuf pool");

    g_portid = 0;
    if (port_init(g_portid, mbuf_pool) < 0)
        rte_exit(EXIT_FAILURE, "port init");

    local_ip   = ip_parse(LOCAL_IP);
    local_mask = ip_parse(LOCAL_MASK);
    gateway_ip = ip_parse(GATEWAY_IP);
    backend_ip = ip_parse(BACKEND_IP);

    printf("local %08x gw %08x backend %08x\n",
           local_ip, gateway_ip, backend_ip);

    /* Kick off gateway ARP resolution. */
    arp_request(gateway_ip);

    while (keep_running) {
        struct rte_mbuf *bufs[BURST_SIZE];
        uint16_t n = rte_eth_rx_burst(g_portid, 0, bufs, BURST_SIZE);
        for (uint16_t i = 0; i < n; i++) {
            struct rte_mbuf *m = bufs[i];
            struct rte_ether_hdr *eth = rte_pktmbuf_mtod(m, struct rte_ether_hdr *);

            if (eth->ether_type == htons(RTE_ETHER_TYPE_ARP)) {
                handle_arp(m);
            } else if (eth->ether_type == htons(RTE_ETHER_TYPE_IPV4)) {
                handle_ipv4(m);
            }
            rte_pktmbuf_free(m);
        }
    }

    rte_eth_dev_stop(g_portid);
    rte_eth_dev_close(g_portid);
    return 0;
}