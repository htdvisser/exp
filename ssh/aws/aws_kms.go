// Package aws provides an ssh.Signer on top of AWS KMS.
package aws

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/asn1"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	kmstypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"golang.org/x/crypto/ssh"
)

// Config is the configuration for the AWS client.
type Config struct {
	Region          string `json:"region,omitempty" yaml:"region,omitempty"`
	AccessKeyID     string `json:"access_key_id,omitempty" yaml:"access_key_id,omitempty"`
	SecretAccessKey string `json:"secret_access_key,omitempty" yaml:"secret_access_key,omitempty"`
	SessionToken    string `json:"session_token,omitempty" yaml:"session_token,omitempty"`
}

// BuildKMSClient builds an AWS KMS client from the configuration.
func (c Config) BuildKMSClient(ctx context.Context) (*kms.Client, error) {
	var opts []func(*config.LoadOptions) error
	if c.Region != "" {
		opts = append(opts, config.WithRegion(c.Region))
	}
	if c.AccessKeyID != "" {
		opts = append(opts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				c.AccessKeyID,
				c.SecretAccessKey,
				c.SessionToken,
			),
		))
	}
	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	client := kms.NewFromConfig(cfg)
	return client, nil
}

// KMSConfig is the configuration for authentication with a keypair stored in AWS KMS.
type KMSConfig struct {
	Config
	KeyID string `json:"key_id" yaml:"key_id"`
}

// Validate validates the configuration and returns an error if it is not valid.
func (c KMSConfig) Validate() error {
	if c.KeyID == "" {
		return fmt.Errorf("missing key id in AWS KMSConfig")
	}
	return nil
}

func buildSigner(ctx context.Context, client *kms.Client, keyID string) (*kmsSigner, error) {
	pkRes, err := client.GetPublicKey(ctx, &kms.GetPublicKeyInput{
		KeyId: &keyID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get public key from AWS KMS: %w", err)
	}
	if pkRes.KeyId != nil {
		keyID = *pkRes.KeyId
	}
	switch pkRes.CustomerMasterKeySpec {
	case kmstypes.CustomerMasterKeySpecRsa2048, kmstypes.CustomerMasterKeySpecRsa3072, kmstypes.CustomerMasterKeySpecRsa4096:
	case kmstypes.CustomerMasterKeySpecEccNistP256, kmstypes.CustomerMasterKeySpecEccNistP384, kmstypes.CustomerMasterKeySpecEccNistP521:
	default:
		return nil, fmt.Errorf("unsupported key type %q in public key returned from AWS KMS", pkRes.CustomerMasterKeySpec)
	}
	if pkRes.KeyUsage != kmstypes.KeyUsageTypeSignVerify {
		return nil, fmt.Errorf("key usage of public key returned from AWS KMS is %q, not \"SIGN_VERIFY\"", pkRes.KeyUsage)
	}
	pk, err := x509.ParsePKIXPublicKey(pkRes.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	sshPK, err := ssh.NewPublicKey(pk)
	if err != nil {
		return nil, fmt.Errorf("failed to convert public key to SSH public key: %w", err)
	}
	s := &kmsSigner{pubKey: sshPK, client: client, keyID: keyID}
	return s, nil
}

// Build builds an ssh.Signer from the configuration.
func (c KMSConfig) Build(ctx context.Context) (ssh.Signer, error) {
	client, err := c.BuildKMSClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to set up AWS KMS client: %w", err)
	}
	signer, err := buildSigner(ctx, client, c.KeyID)
	if err != nil {
		return nil, err
	}
	return signer, nil
}

type kmsSigner struct {
	pubKey ssh.PublicKey
	client *kms.Client
	keyID  string
}

func (s *kmsSigner) PublicKey() ssh.PublicKey {
	return s.pubKey
}

func (s *kmsSigner) Sign(rand io.Reader, data []byte) (*ssh.Signature, error) {
	return s.SignWithAlgorithm(rand, data, "")
}

func (s *kmsSigner) SignWithAlgorithm(_ io.Reader, data []byte, algorithm string) (*ssh.Signature, error) {
	pubKeyType := s.pubKey.Type()
	if algorithm == "" {
		algorithm = pubKeyType
		if algorithm == ssh.SigAlgoRSA {
			algorithm = ssh.SigAlgoRSASHA2512
		}
	}
	switch algorithm {
	case ssh.SigAlgoRSA, ssh.SigAlgoRSASHA2256, ssh.SigAlgoRSASHA2512:
		if pubKeyType != ssh.KeyAlgoRSA {
			return nil, fmt.Errorf("unsupported signature algorithm %q for key of type %q", algorithm, pubKeyType)
		}
	default:
		if algorithm != pubKeyType {
			return nil, fmt.Errorf("unsupported signature algorithm %q for key of type %q", algorithm, pubKeyType)
		}
	}
	var (
		hashFunc         crypto.Hash
		signingAlgorithm kmstypes.SigningAlgorithmSpec
	)
	switch algorithm {
	case ssh.SigAlgoRSA:
		return nil, fmt.Errorf("signature algorithm %q is not supported by AWS KMS", algorithm)
	case ssh.SigAlgoRSASHA2256:
		hashFunc, signingAlgorithm = crypto.SHA256, kmstypes.SigningAlgorithmSpecRsassaPkcs1V15Sha256
	case ssh.SigAlgoRSASHA2512:
		hashFunc, signingAlgorithm = crypto.SHA512, kmstypes.SigningAlgorithmSpecRsassaPkcs1V15Sha512
	case ssh.KeyAlgoECDSA256:
		hashFunc, signingAlgorithm = crypto.SHA256, kmstypes.SigningAlgorithmSpecEcdsaSha256
	case ssh.KeyAlgoECDSA384:
		hashFunc, signingAlgorithm = crypto.SHA384, kmstypes.SigningAlgorithmSpecEcdsaSha384
	case ssh.KeyAlgoECDSA521:
		hashFunc, signingAlgorithm = crypto.SHA512, kmstypes.SigningAlgorithmSpecEcdsaSha512
	default:
		return nil, fmt.Errorf("unsupported signature algorithm %q", algorithm)
	}

	h := hashFunc.New()
	h.Write(data)
	digest := h.Sum(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sig, err := s.client.Sign(ctx, &kms.SignInput{
		KeyId:            &s.keyID,
		Message:          digest,
		MessageType:      kmstypes.MessageTypeDigest,
		SigningAlgorithm: signingAlgorithm,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to sign data: %w", err)
	}

	signature := &ssh.Signature{
		Format: algorithm,
		Blob:   sig.Signature,
	}
	switch pubKeyType {
	case ssh.KeyAlgoECDSA256, ssh.KeyAlgoECDSA384, ssh.KeyAlgoECDSA521:
		type asn1Signature struct {
			R, S *big.Int
		}
		asn1Sig := new(asn1Signature)
		_, err := asn1.Unmarshal(signature.Blob, asn1Sig)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal ASN.1 signature: %w", err)
		}
		signature.Blob = ssh.Marshal(asn1Sig)
	}

	return signature, nil
}
