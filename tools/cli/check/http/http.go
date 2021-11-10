package http

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

func CheckHTTP(ctx context.Context, uri string, remoteHost net.IP) error {
	fmt.Println(uri, remoteHost)
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	fmt.Println(u)

	cctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "curl/7.74.0")

	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		fmt.Println("BEGIN DIAL")
		defer fmt.Println("END DIAL")
		return DialWithHostOverride(ctx, network, addr, remoteHost)
	}
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: dialContext,
			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				serverName, _, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}

				conn, err := dialContext(ctx, network, addr)
				if err != nil {
					return nil, err
				}
				rootCAs, err := x509.SystemCertPool()
				fmt.Println(err)
				tlsConfig := &tls.Config{
					InsecureSkipVerify: true,
					ServerName:         serverName,
					RootCAs:            rootCAs,
				}
				tlsConfig.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
					return VerifyPeerCertificate(ctx, rawCerts, tlsConfig)
				}
				c := tls.Client(conn, tlsConfig)
				if err := c.HandshakeContext(ctx); err != nil {
					_ = conn.Close()
					return nil, err
				}
				return c, nil
				//return tls.DialWithDialer(&zeroDialer, network, addr, nil)
			},
		},
	}
	res, err := client.Do(req.WithContext(cctx))

	fmt.Println(err)
	fmt.Println(res)

	return nil
}

func DialWithHostOverride(ctx context.Context, network, addr string, remoteHost net.IP) (net.Conn, error) {
	if remoteHost != nil {
		_, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		addr = net.JoinHostPort(remoteHost.String(), port)
	}
	var zeroDialer net.Dialer
	return zeroDialer.DialContext(ctx, network, addr)
}

func VerifyPeerCertificate(ctx context.Context, rawCerts [][]byte, config *tls.Config) error {
	rootCAs, err := x509.SystemCertPool()
	fmt.Println(err)

	certs := make([]*x509.Certificate, len(rawCerts))
	for i, asn1Data := range rawCerts {
		cert, err := x509.ParseCertificate(asn1Data)
		if err != nil {
			return errors.New("tls: failed to parse certificate from server: " + err.Error())
		}
		certs[i] = cert
	}

	opts := x509.VerifyOptions{
		Roots:         rootCAs,
		CurrentTime:   time.Now().UTC(),
		DNSName:       config.ServerName,
		Intermediates: x509.NewCertPool(),
	}
	for _, cert := range certs[1:] {
		opts.Intermediates.AddCert(cert)
	}
	_, err = certs[0].Verify(opts)
	if err != nil {
		//	c.sendAlert(alertBadCertificate)
		//return err
		fmt.Println(err)
	}
	for _, cert := range certs {

		//* Server certificate:
		//*  subject: CN=ip.bozaro.ru
		//*  start date: Sep  7 15:44:01 2021 GMT
		//*  expire date: Dec  6 15:44:00 2021 GMT
		//*  subjectAltName: host "ip.bozaro.ru" matched cert's "ip.bozaro.ru"
		//*  issuer: C=US; O=Let's Encrypt; CN=R3
		//*  SSL certificate verify ok.

		fmt.Println(cert.DNSNames, cert.Subject, cert.NotBefore, cert.NotAfter, cert.BasicConstraintsValid, cert.IsCA, cert.VerifyHostname(config.ServerName))
	}
	fmt.Println(certs)
	return nil
}
