package dns

import (
	"context"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/joomcode/diagnostic/tools/cli/logger"
	"github.com/miekg/dns"
)

func SystemLookupHost(ctx context.Context, log logger.Logger, host string) ([]net.IP, error) {
	systemResolver := net.Resolver{
		PreferGo: false,
	}

	cctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	fqdn := dns.Fqdn(host)
	log.Infof("Lookup host by system resolver: %q", fqdn)
	begin := time.Now()
	value, err := systemResolver.LookupHost(cctx, fqdn)
	rtt := time.Since(begin)
	if err != nil {
		log.Errorf("Lookup failed with error: %v", err)
		return nil, err
	}
	log.Infof("Lookup result [%v]:\n%s", rtt, strings.Join(value, "\n"))

	ips := make([]net.IP, 0, len(value))
	for _, host := range value {
		ip := net.ParseIP(host)
		if ip == nil {
			log.Errorf("Can't parse IP address: %q", host)
			continue
		}
		ips = append(ips, ip)
	}
	if len(ips) == 0 {
		log.Errorf("Can't found any hosts by name: %q", fqdn)
	}
	return ips, nil
}

func ServerLookupHost(ctx context.Context, log logger.Logger, host string, server net.IP, port int) ([]net.IP, error) {
	fqdn := dns.Fqdn(host)
	dnsServer := net.JoinHostPort(server.String(), strconv.Itoa(port))

	cctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	log.Infof("Lookup host by DNS server %q: %q", dnsServer, fqdn)

	c := new(dns.Client)

	var ips []net.IP
	for _, dnsType := range []uint16{dns.TypeA, dns.TypeAAAA} {
		var msg dns.Msg
		msg.SetQuestion(fqdn, dnsType)

		in, rtt, err := c.ExchangeContext(cctx, &msg, dnsServer)
		if err != nil {
			log.Errorf("Lookup %s failed with error: %v", dns.TypeToString[dnsType], err)
			continue
		}
		log.Infof("Lookup %s result [%v]:\n%s", dns.TypeToString[dnsType], dnsType, rtt, in.String())
		for _, rec := range in.Answer {
			switch r := rec.(type) {
			case *dns.A:
				ips = append(ips, r.A)
			case *dns.AAAA:
				ips = append(ips, r.AAAA)
			}
		}
	}
	if len(ips) == 0 {
		log.Errorf("Can't found any hosts by name: %q", fqdn)
	}
	return ips, nil
}

func SystemLookupTXT(ctx context.Context, log logger.Logger, record string) ([]string, error) {
	systemResolver := net.Resolver{
		PreferGo: false,
	}

	cctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	fqdn := dns.Fqdn(record)
	log.Infof("Lookup TXT DNS record by system resolver: %q", fqdn)
	begin := time.Now()
	value, err := systemResolver.LookupTXT(cctx, fqdn)
	rtt := time.Since(begin)
	if err != nil {
		log.Errorf("Lookup failed with error: %v", err)
		return nil, err
	}
	log.Infof("Lookup result [%v]:\n%s", rtt, strings.Join(value, "\n"))
	return value, nil
}

func ServerLookupTXT(ctx context.Context, log logger.Logger, record string, server net.IP, port int) ([]string, error) {
	fqdn := dns.Fqdn(record)
	dnsServer := net.JoinHostPort(server.String(), strconv.Itoa(port))

	cctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	log.Infof("Lookup TXT DNS record by DNS server %q: %q", dnsServer, fqdn)

	var msg dns.Msg
	msg.SetQuestion(fqdn, dns.TypeTXT)

	c := new(dns.Client)
	in, rtt, err := c.ExchangeContext(cctx, &msg, dnsServer)
	if err != nil {
		log.Errorf("Lookup failed with error: %v", err)
		return nil, err
	}
	log.Infof("Lookup result [%v]:\n%s", rtt, in.String())
	return nil, nil
}
