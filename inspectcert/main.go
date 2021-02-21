package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/pflag"
	"htdvisser.dev/exp/pflagenv"
	"htdvisser.dev/exp/tlsconfig"
)

var config struct {
	server tlsconfig.ClientConfig
}

var flags = func() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("inspectcert", pflag.ContinueOnError)
	flagSet.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: inspectcert [FLAGS] [HOST[:PORT]]")
		fmt.Fprintln(os.Stderr, "Inspect the TLS certificate chain for some server")
		flagSet.PrintDefaults()
	}

	flagSet.AddFlagSet(config.server.Flags("", &tlsconfig.ClientConfig{
		ServerCA: tlsconfig.CAConfig{CACert: ""},
	}))

	return flagSet
}()

var (
	serial   = flags.Bool("cert.serial", false, "print certificate serial")
	subject  = flags.Bool("cert.subject", true, "print certificate subject")
	issuer   = flags.Bool("cert.issuer", true, "print certificate issuer")
	validity = flags.Bool("cert.validity", true, "print certificate validity")
	dns      = flags.Bool("cert.dns", true, "print certificate DNS names")
	raw      = flags.Bool("cert.raw", false, "print raw certificate")
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := pflagenv.NewParser().ParseEnv(flags); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flags.Usage()
		os.Exit(2)
	}

	if err := flags.Parse(os.Args[1:]); err != nil {
		if err != pflag.ErrHelp {
			fmt.Fprintln(os.Stderr, err)
			flags.Usage()
			os.Exit(2)
		}
		os.Exit(0)
	}

	if len(flags.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "missing [HOST[:PORT]] argument")
		flags.Usage()
		os.Exit(2)
	}

	host, port, err := net.SplitHostPort(flags.Arg(0))
	if err != nil {
		host = flags.Arg(0)
		port = "443"
	}
	target := net.JoinHostPort(host, port)

	dialer := net.Dialer{}

	conn, err := dialer.DialContext(ctx, "tcp", target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to \"%s\": %s\n", target, err)
		os.Exit(1)
	}
	defer conn.Close()

	tlsConfig, err := config.server.Load(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load TLS config: %s\n", err)
		os.Exit(2)
	}

	if tlsConfig.ServerName == "" {
		tlsConfig.ServerName = host
	}

	tlsConn := tls.Client(conn, tlsConfig)

	if err = tlsConn.Handshake(); err != nil {
		var unknownAuthorityError x509.UnknownAuthorityError
		if errors.As(err, &unknownAuthorityError) {
			fmt.Fprintln(os.Stderr,
				color.New(color.BgRed, color.FgWhite).Sprint("Error:"),
				color.New(color.FgRed).Sprint("Unknown Authority"),
			)
			printCert(unknownAuthorityError.Cert)
			os.Exit(1)
		}
		var certificateInvalidError x509.CertificateInvalidError
		if errors.As(err, &certificateInvalidError) && certificateInvalidError.Reason == x509.Expired {
			fmt.Fprintln(os.Stderr,
				color.New(color.BgRed, color.FgWhite).Sprint("Error:"),
				color.New(color.FgRed).Sprint("Expired Certificate"),
			)
			printCert(certificateInvalidError.Cert)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "Failed to complete TLS handshake with \"%s\": %s\n", target, err)
		os.Exit(1)
	}

	chains := tlsConn.ConnectionState().VerifiedChains

	if len(chains) > 1 {
		fmt.Fprintln(os.Stderr,
			color.New(color.BgYellow, color.FgBlack).Sprint("Warning:"),
			color.New(color.FgYellow).Sprint("Multiple Chains"),
		)
	}

	for i, chain := range chains {
		if i > 0 {
			fmt.Fprintln(os.Stdout, "---")
		}
		for _, cert := range chain {
			printCert(cert)
		}
	}
}

func printCert(cert *x509.Certificate) {
	fmt.Fprint(os.Stdout, "- ")

	var indent bool
	item := func() {
		if indent {
			fmt.Fprint(os.Stdout, "  ")
		} else {
			indent = true
		}
	}

	if *serial {
		item()
		fmt.Fprint(os.Stdout, "Serial:   ")
		for i, b := range cert.SerialNumber.Bytes() {
			if i > 0 {
				fmt.Fprint(os.Stdout, color.New(color.Faint).Sprint(":"))
			}
			fmt.Fprintf(os.Stdout, "%02x", b)
		}
		fmt.Fprintln(os.Stdout)
	}

	if *subject {
		item()
		fmt.Fprintf(os.Stdout, "Subject:  %s\n", cert.Subject)
	}

	if *issuer {
		item()
		fmt.Fprintf(os.Stdout, "Issuer:   %s\n", cert.Issuer)
	}

	if *validity {
		item()
		fmt.Fprintln(os.Stdout, "Validity: ")
		notBeforeColor := color.GreenString
		if time.Now().Before(cert.NotBefore) {
			notBeforeColor = color.RedString
		}
		fmt.Fprintf(os.Stdout, "    Not Before: %s\n", notBeforeColor(cert.NotBefore.UTC().Format(time.RFC3339)))
		notAfterColor := color.GreenString
		if time.Now().After(cert.NotAfter) {
			notAfterColor = color.RedString
		}
		fmt.Fprintf(os.Stdout, "    Not After:  %s\n", notAfterColor(cert.NotAfter.UTC().Format(time.RFC3339)))
	}

	if *dns {
		if len(cert.DNSNames) > 0 {
			item()
			fmt.Fprint(os.Stdout, "DNS:      ")
			if len(cert.DNSNames) > 3 {
				fmt.Fprintln(os.Stdout)
			} else {
				fmt.Fprint(os.Stdout, "[")
			}
			for i, name := range cert.DNSNames {
				if len(cert.DNSNames) > 3 {
					fmt.Fprint(os.Stdout, "    - ")
				} else if i > 0 {
					fmt.Fprint(os.Stdout, ", ")
				}
				fmt.Fprint(os.Stdout, name)
				if len(cert.DNSNames) > 3 {
					fmt.Fprintln(os.Stdout)
				}
			}
			if len(cert.DNSNames) <= 3 {
				fmt.Fprintln(os.Stdout, "]")
			}
		}
	}

	if *raw {
		item()
		fmt.Fprintln(os.Stdout, "Raw: |")
		pemBlock := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		})
		for _, line := range bytes.Split(bytes.TrimSpace(pemBlock), []byte("\n")) {
			fmt.Fprint(os.Stdout, "    ")
			os.Stdout.Write(line)
			fmt.Fprintln(os.Stdout)
		}
	}
}
