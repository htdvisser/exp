package sshclient

import (
	"context"
	"net"
	"os"
	"testing"
)

func buildConnectConfig() ConnectConfig {
	return ConnectConfig{
		Network: "tcp",
		Address: "localhost:2222",
		HostKey: HostKeyConfig{
			Source: "insecure_ignore",
		},
		Username: "testuser",
		AuthMethods: []AuthMethodConfig{
			{
				Method:   "password",
				Password: "testpassword",
			},
		},
	}
}

func TestDial(t *testing.T) {
	c := buildConnectConfig()
	if err := c.Validate(); err != nil {
		t.Errorf("Config failed to validate: %v", err)
	}

	conn, err := net.Dial(c.Network, c.Address)
	if err != nil {
		t.Skip("SSH Server not running")
	}
	defer conn.Close()

	t.Run("KnownHosts", func(t *testing.T) {
		c := buildConnectConfig()
		c.HostKey = HostKeyConfig{Source: "known_hosts"}
		c.HostKey.KnownHosts.File = "testdata/known_hosts"
		if err := c.Validate(); err != nil {
			t.Errorf("Config failed to validate: %v", err)
		}
		cli, err := c.Dial(context.Background())
		if err != nil {
			t.Fatalf("runssh_test: dial failed: %v", err)
		}
		defer cli.Close()
	})

	t.Run("Password", func(t *testing.T) {
		c := buildConnectConfig()
		if err := c.Validate(); err != nil {
			t.Errorf("Config failed to validate: %v", err)
		}
		cli, err := c.Dial(context.Background())
		if err != nil {
			t.Fatalf("runssh_test: dial failed: %v", err)
		}
		defer cli.Close()
	})

	t.Run("PrivateKey", func(t *testing.T) {
		for _, keyType := range []string{"ecdsa", "ed25519", "rsa"} {
			t.Run(keyType, func(t *testing.T) {
				c := buildConnectConfig()
				c.AuthMethods = []AuthMethodConfig{
					{
						Method: "private_keys",
						PrivateKeys: []PrivateKeyConfig{
							{File: "testdata/id_" + keyType},
						},
					},
				}
				if err := c.Validate(); err != nil {
					t.Errorf("Config failed to validate: %v", err)
				}
				cli, err := c.Dial(context.Background())
				if err != nil {
					t.Fatalf("runssh_test: dial failed: %v", err)
				}
				defer cli.Close()
			})
		}
	})

	t.Run("AWSKMS", func(t *testing.T) {
		var (
			awsRegion          = os.Getenv("AWS_REGION")
			awsAccessKeyID     = os.Getenv("AWS_ACCESS_KEY_ID")
			awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
			awsKMSKeyID        = os.Getenv("AWS_KMS_KEY_ID")
		)
		if awsKMSKeyID == "" {
			t.Skip("AWS KMS Key ID not set")
		}
		c := buildConnectConfig()
		c.AuthMethods = []AuthMethodConfig{
			{
				Method: "aws_kms",
				AWSKMS: AWSKMSConfig{
					Region:          awsRegion,
					AccessKeyID:     awsAccessKeyID,
					SecretAccessKey: awsSecretAccessKey,
					KeyID:           awsKMSKeyID,
				},
			},
		}
		if err := c.Validate(); err != nil {
			t.Errorf("Config failed to validate: %v", err)
		}
		cli, err := c.Dial(context.Background())
		if err != nil {
			t.Fatalf("runssh_test: dial failed: %v", err)
		}
		defer cli.Close()
	})

}
