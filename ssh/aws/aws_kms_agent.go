package aws

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// KMSAgentConfig is the configuration for authentication with keypairs stored in AWS KMS.
type KMSAgentConfig struct {
	Config
	KeyIDs []string `json:"key_ids" yaml:"key_ids"`
}

// Validate validates the configuration and returns an error if it is not valid.
func (c KMSAgentConfig) Validate() error {
	if len(c.KeyIDs) == 0 {
		return fmt.Errorf("missing key ids in AWS KMSAgentConfig")
	}
	return nil
}

// Build builds an agent.Agent from the configuration.
func (c KMSAgentConfig) Build(ctx context.Context) (agent.Agent, error) {
	client, err := c.BuildKMSClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to set up AWS KMS client: %w", err)
	}
	signers := make([]*kmsSigner, len(c.KeyIDs))
	for i, keyID := range c.KeyIDs {
		signer, err := buildSigner(ctx, client, keyID)
		if err != nil {
			return nil, fmt.Errorf("failed to set up key %q: %w", keyID, err)
		}
		signers[i] = signer
	}
	return &kmsAgent{
		signers: signers,
	}, nil
}

type kmsAgent struct {
	signers []*kmsSigner
}

func (a *kmsAgent) List() ([]*agent.Key, error) {
	keys := make([]*agent.Key, len(a.signers))
	for i, signer := range a.signers {
		keys[i] = &agent.Key{
			Format:  signer.pubKey.Type(),
			Blob:    signer.pubKey.Marshal(),
			Comment: signer.keyID,
		}
	}
	return keys, nil
}

func (a *kmsAgent) Signers() ([]ssh.Signer, error) {
	signers := make([]ssh.Signer, len(a.signers))
	for i, signer := range a.signers {
		signers[i] = signer
	}
	return signers, nil
}

func (a *kmsAgent) selectSigner(key ssh.PublicKey) *kmsSigner {
	for _, signer := range a.signers {
		if bytes.Equal(signer.pubKey.Marshal(), key.Marshal()) {
			return signer
		}
	}
	return nil
}

func (a *kmsAgent) Sign(key ssh.PublicKey, data []byte) (*ssh.Signature, error) {
	return a.SignWithFlags(key, data, 0)
}

func (a *kmsAgent) SignWithFlags(key ssh.PublicKey, data []byte, flags agent.SignatureFlags) (*ssh.Signature, error) {
	signer := a.selectSigner(key)
	if signer == nil {
		return nil, fmt.Errorf("no signer matches public key")
	}
	algorithm := signer.pubKey.Type()
	switch {
	case algorithm == ssh.KeyAlgoRSA && flags&agent.SignatureFlagRsaSha256 != 0:
		algorithm = ssh.SigAlgoRSASHA2256
	case algorithm == ssh.KeyAlgoRSA && flags&agent.SignatureFlagRsaSha512 != 0:
		algorithm = ssh.SigAlgoRSASHA2512
	}
	return signer.SignWithAlgorithm(nil, data, algorithm)
}

var errNotImplemented = errors.New("not implemented")

func (*kmsAgent) Add(agent.AddedKey) error   { return errNotImplemented }
func (*kmsAgent) Remove(ssh.PublicKey) error { return errNotImplemented }
func (*kmsAgent) RemoveAll() error           { return errNotImplemented }
func (*kmsAgent) Lock([]byte) error          { return errNotImplemented }
func (*kmsAgent) Unlock([]byte) error        { return errNotImplemented }

func (*kmsAgent) Extension(string, []byte) ([]byte, error) { return nil, agent.ErrExtensionUnsupported }
