// natcheck implements the NAT behavior tests from RFC 5780 section 4.
// The test flow is based on Pion STUN's MIT-licensed nat-behaviour example.
package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/pion/stun/v3"
)

type behavior string

const (
	endpointIndependent  behavior = "endpoint-independent"
	addressDependent     behavior = "address-dependent"
	addressPortDependent behavior = "address-and-port-dependent"
	inconclusive         behavior = "inconclusive"
)

type probeResult struct {
	mapped *stun.XORMappedAddress
	other  *stun.OtherAddress
	origin *stun.ResponseOrigin
}

type checker struct {
	conn    *net.UDPConn
	timeout time.Duration
}

func main() {
	server := flag.String("server", "stun.cloudflare.com:3478", "STUN server (prefer one supporting RFC 5780)")
	fallback := flag.String("fallback-server", "stun.l.google.com:19302", "second STUN server used when RFC 5780 is unavailable")
	timeout := flag.Duration("timeout", 3*time.Second, "timeout for each STUN request")
	flag.Parse()

	primary, err := net.ResolveUDPAddr("udp4", *server)
	if err != nil {
		fatal("resolve primary STUN server", err)
	}
	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		fatal("open UDP socket", err)
	}
	defer conn.Close()

	c := &checker{conn: conn, timeout: *timeout}
	fmt.Println("local socket:", conn.LocalAddr())
	fmt.Println("primary STUN server:", primary)

	baseline, err := c.probe(primary, nil)
	if err != nil {
		fatal("primary STUN probe", err)
	}
	fmt.Println("public endpoint:", baseline.mapped)
	if baseline.origin != nil {
		fmt.Println("response origin:", baseline.origin)
	}

	mapping := inconclusive
	filtering := inconclusive
	if baseline.other != nil {
		fmt.Println("RFC 5780 other address:", baseline.other)
		mapping = c.mappingBehavior(primary, baseline)
		filtering = c.filteringBehavior(primary)
	} else {
		fmt.Println("primary server does not provide OTHER-ADDRESS; using limited two-server mapping test")
		mapping = c.fallbackMapping(*fallback, baseline)
	}

	fmt.Println()
	fmt.Println("NAT mapping behavior:", mapping)
	fmt.Println("NAT filtering behavior:", filtering)
	fmt.Println("likely NAT type:", legacyNATType(mapping, filtering, isLocalEndpoint(conn, baseline.mapped)))
	if mapping == inconclusive || filtering == inconclusive {
		fmt.Println("note: an inconclusive result means the STUN server lacks RFC 5780 support or a probe timed out")
	}
}

func (c *checker) mappingBehavior(primary *net.UDPAddr, baseline probeResult) behavior {
	other := &net.UDPAddr{IP: baseline.other.IP, Port: primary.Port}
	second, err := c.probe(other, nil)
	if err != nil {
		fmt.Println("mapping test II:", err)
		return inconclusive
	}
	fmt.Printf("mapping via other IP, primary port: %s\n", second.mapped)
	if sameMapped(baseline.mapped, second.mapped) {
		return endpointIndependent
	}

	other.Port = baseline.other.Port
	third, err := c.probe(other, nil)
	if err != nil {
		fmt.Println("mapping test III:", err)
		return inconclusive
	}
	fmt.Printf("mapping via other IP and port: %s\n", third.mapped)
	if sameMapped(second.mapped, third.mapped) {
		return addressDependent
	}
	return addressPortDependent
}

func (c *checker) filteringBehavior(primary *net.UDPAddr) behavior {
	// RFC 5780 CHANGE-REQUEST: change IP and port.
	if _, err := c.probe(primary, []byte{0, 0, 0, 6}); err == nil {
		return endpointIndependent
	} else if !isTimeout(err) {
		fmt.Println("filtering test II:", err)
		return inconclusive
	}

	// CHANGE-REQUEST: change port only.
	if _, err := c.probe(primary, []byte{0, 0, 0, 2}); err == nil {
		return addressDependent
	} else if isTimeout(err) {
		return addressPortDependent
	} else {
		fmt.Println("filtering test III:", err)
		return inconclusive
	}
}

func (c *checker) fallbackMapping(server string, baseline probeResult) behavior {
	addr, err := net.ResolveUDPAddr("udp4", server)
	if err != nil {
		fmt.Println("resolve fallback STUN server:", err)
		return inconclusive
	}
	fmt.Println("fallback STUN server:", addr)
	result, err := c.probe(addr, nil)
	if err != nil {
		fmt.Println("fallback STUN probe:", err)
		return inconclusive
	}
	fmt.Println("public endpoint via fallback:", result.mapped)
	if sameMapped(baseline.mapped, result.mapped) {
		return endpointIndependent
	}
	return addressPortDependent
}

func (c *checker) probe(destination *net.UDPAddr, changeRequest []byte) (probeResult, error) {
	request := stun.MustBuild(stun.TransactionID, stun.BindingRequest)
	if changeRequest != nil {
		request.Add(stun.AttrChangeRequest, changeRequest)
	}
	if _, err := c.conn.WriteToUDP(request.Raw, destination); err != nil {
		return probeResult{}, err
	}

	deadline := time.Now().Add(c.timeout)
	buffer := make([]byte, 2048)
	for {
		if err := c.conn.SetReadDeadline(deadline); err != nil {
			return probeResult{}, err
		}
		n, _, err := c.conn.ReadFromUDP(buffer)
		if err != nil {
			return probeResult{}, err
		}
		message := &stun.Message{Raw: append([]byte(nil), buffer[:n]...)}
		if err := message.Decode(); err != nil || message.TransactionID != request.TransactionID {
			continue
		}

		var result probeResult
		mapped := &stun.XORMappedAddress{}
		if err := mapped.GetFrom(message); err != nil {
			return probeResult{}, fmt.Errorf("response has no XOR-MAPPED-ADDRESS: %w", err)
		}
		result.mapped = mapped
		other := &stun.OtherAddress{}
		if other.GetFrom(message) == nil {
			result.other = other
		}
		origin := &stun.ResponseOrigin{}
		if origin.GetFrom(message) == nil {
			result.origin = origin
		}
		return result, nil
	}
}

func sameMapped(a, b *stun.XORMappedAddress) bool {
	return a != nil && b != nil && a.Port == b.Port && a.IP.Equal(b.IP)
}

func isTimeout(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

func isLocalEndpoint(conn *net.UDPConn, mapped *stun.XORMappedAddress) bool {
	if mapped == nil || mapped.Port != conn.LocalAddr().(*net.UDPAddr).Port {
		return false
	}
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return false
	}
	for _, address := range addresses {
		ip, _, err := net.ParseCIDR(address.String())
		if err == nil && ip.Equal(mapped.IP) {
			return true
		}
	}
	return false
}

func legacyNATType(mapping, filtering behavior, noNAT bool) string {
	if noNAT {
		return "no NAT (public address is assigned locally)"
	}
	if mapping == addressDependent || mapping == addressPortDependent {
		return "symmetric NAT / endpoint-dependent mapping"
	}
	if mapping != endpointIndependent {
		return "unknown"
	}
	switch filtering {
	case endpointIndependent:
		return "full-cone-like"
	case addressDependent:
		return "address-restricted-cone-like"
	case addressPortDependent:
		return "port-restricted-cone-like"
	default:
		return "endpoint-independent mapping; filtering unknown"
	}
}

func fatal(operation string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", operation, err)
	os.Exit(1)
}
