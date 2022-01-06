package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"

	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/terminal"
	"htdvisser.dev/exp/ssh/aws"
)

var config aws.KMSAgentConfig

func init() {
	flag.StringVar(&config.Region, "aws.region", "", "AWS region")
	flag.StringVar(&config.AccessKeyID, "aws.access-key-id", "", "AWS Access Key ID")
	flag.StringVar(&config.SecretAccessKey, "aws.secret-access-key", "", "AWS Secret Access Key")
	flag.StringVar(&config.SessionToken, "aws.session-token", "", "AWS Session Token")
	flag.StringVar(&config.AssumeRoleARN, "aws.assume-role-arn", "", "AWS Role ARN to assume")
	flag.Func("aws.kms.key-id", "AWS KMS Key ID (can be specified more than once)", func(keyID string) error {
		config.KeyIDs = append(config.KeyIDs, keyID)
		return nil
	})
}

var (
	socketPath = flag.String("socket", "aws-kms-agent.sock", "socket path")
	rm         = flag.Bool("rm", false, "remove old socket")
)

func main() {
	flag.Parse()

	if err := config.Validate(); err != nil {
		log.Println(err)
		flag.Usage()
		os.Exit(2)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := run(ctx); err != nil {
		log.Println(err)
		stop() // os.Exit does not run defers.
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	kmsAgent, err := config.Build()
	if err != nil {
		return fmt.Errorf("failed to build agent: %w", err)
	}

	socketDir := filepath.Dir(*socketPath)
	if err = os.MkdirAll(socketDir, 0777); err != nil {
		return fmt.Errorf("failed to create folder %q for socket: %w", socketDir, err)
	}
	if stat, err := os.Stat(*socketPath); err == nil {
		if stat.Mode() != fs.ModeSocket {
			return fmt.Errorf("socket path %q unavailable", *socketPath)
		}
		if !*rm {
			return fmt.Errorf("socket path %q still has old socket", *socketPath)
		}
		if err = os.Remove(*socketPath); err != nil {
			return fmt.Errorf("failed to remove old socket at %q: %w", *socketPath, err)
		}
	}
	lis, err := net.Listen("unix", *socketPath)
	if err != nil {
		return fmt.Errorf("failed to listen on socket path %q: %w", *socketPath, err)
	}

	go func() {
		<-ctx.Done()
		lis.Close()
	}()

	if terminal.IsTerminal(int(os.Stdin.Fd())) {
		absPath, err := filepath.Abs(*socketPath)
		if err != nil {
			absPath = *socketPath
		}

		fmt.Fprintln(os.Stdout, "# You can now use the agent by running:")
		fmt.Fprintf(os.Stdout, "export SSH_AUTH_SOCK=%q", absPath)
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "# Press ctrl-c to exit")
		fmt.Fprintln(os.Stdout)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept connection: %w", err)
		}
		go func() {
			if err := agent.ServeAgent(kmsAgent, conn); err != nil && !errors.Is(err, io.EOF) {
				log.Printf("Error in ServeAgent for conn %v: %v", conn.RemoteAddr(), err)
			}
		}()
	}
}
