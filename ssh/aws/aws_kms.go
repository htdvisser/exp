// Package aws provides an ssh.Signer on top of AWS KMS.
package aws

import (
	"crypto"
	"crypto/x509"
	"encoding/asn1"
	"fmt"
	"io"
	"math/big"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"golang.org/x/crypto/ssh"
)

// Config is the configuration for the AWS client.
type Config struct {
	Region          string `json:"region,omitempty" yaml:"region,omitempty"`
	AccessKeyID     string `json:"access_key_id,omitempty" yaml:"access_key_id,omitempty"`
	SecretAccessKey string `json:"secret_access_key,omitempty" yaml:"secret_access_key,omitempty"`
	SessionToken    string `json:"session_token,omitempty" yaml:"session_token,omitempty"`
	AssumeRoleARN   string `json:"assume_role_arn,omitempty" yaml:"assume_role_arn,omitempty"`
}

// BuildKMSClient builds an AWS KMS client from the configuration.
func (c Config) BuildKMSClient() (*kms.KMS, error) {
	awsConfig := aws.NewConfig()
	if c.Region != "" {
		awsConfig = awsConfig.WithRegion(c.Region)
	}
	if c.AccessKeyID != "" {
		awsConfig = awsConfig.WithCredentials(credentials.NewStaticCredentials(
			c.AccessKeyID,
			c.SecretAccessKey,
			c.SessionToken,
		))
	}
	ses, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create new AWS session: %w", err)
	}
	if c.AssumeRoleARN != "" {
		awsConfig = awsConfig.WithCredentials(stscreds.NewCredentials(
			ses, c.AssumeRoleARN,
		))
		ses.Config.MergeIn(awsConfig)
	}
	client := kms.New(ses)
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

type typedPublicKey struct {
	ssh.PublicKey
	keyType string
}

func (pk typedPublicKey) Type() string {
	return pk.keyType
}

func buildSigner(client *kms.KMS, keyID string) (*kmsSigner, error) {
	pkRes, err := client.GetPublicKey(&kms.GetPublicKeyInput{
		KeyId: &keyID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get public key from AWS KMS: %w", err)
	}
	if pkRes.KeyId != nil {
		keyID = *pkRes.KeyId
	}
	if pkRes.CustomerMasterKeySpec == nil {
		return nil, fmt.Errorf("missing key type in public key returned from AWS KMS")
	}
	switch *pkRes.CustomerMasterKeySpec {
	case "RSA_2048", "RSA_3072", "RSA_4096":
	case "ECC_NIST_P256", "ECC_NIST_P384", "ECC_NIST_P521":
	default:
		return nil, fmt.Errorf("unsupported key type %q in public key returned from AWS KMS", *pkRes.CustomerMasterKeySpec)
	}
	if pkRes.KeyUsage == nil {
		return nil, fmt.Errorf("missing key usage in public key returned from AWS KMS")
	}
	if *pkRes.KeyUsage != "SIGN_VERIFY" {
		return nil, fmt.Errorf("key usage of public key returned from AWS KMS is %q, not \"SIGN_VERIFY\"", *pkRes.KeyUsage)
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
	if s.pubKey.Type() == ssh.KeyAlgoRSA {
		s.pubKey = typedPublicKey{PublicKey: sshPK, keyType: ssh.SigAlgoRSASHA2256}
	}
	return s, nil
}

// Build builds an ssh.Signer from the configuration.
func (c KMSConfig) Build() (ssh.Signer, error) {
	client, err := c.BuildKMSClient()
	if err != nil {
		return nil, fmt.Errorf("failed to set up AWS KMS client: %w", err)
	}
	signer, err := buildSigner(client, c.KeyID)
	if err != nil {
		return nil, err
	}
	return signer, nil
}

type kmsSigner struct {
	pubKey ssh.PublicKey
	client *kms.KMS
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
	} else if algorithm != s.pubKey.Type() {
		return nil, fmt.Errorf("unsupported signature algorithm %q for key of type %q", algorithm, pubKeyType)
	}
	var (
		hashFunc         crypto.Hash
		signingAlgorithm string
	)
	switch algorithm {
	case ssh.SigAlgoRSA:
		return nil, fmt.Errorf("signature algorithm %q is not supported by AWS KMS", algorithm)
	case ssh.SigAlgoRSASHA2256:
		hashFunc, signingAlgorithm = crypto.SHA256, kms.SigningAlgorithmSpecRsassaPkcs1V15Sha256
	case ssh.SigAlgoRSASHA2512:
		hashFunc, signingAlgorithm = crypto.SHA512, kms.SigningAlgorithmSpecRsassaPkcs1V15Sha512
	case ssh.KeyAlgoECDSA256:
		hashFunc, signingAlgorithm = crypto.SHA256, kms.SigningAlgorithmSpecEcdsaSha256
	case ssh.KeyAlgoECDSA384:
		hashFunc, signingAlgorithm = crypto.SHA384, kms.SigningAlgorithmSpecEcdsaSha384
	case ssh.KeyAlgoECDSA521:
		hashFunc, signingAlgorithm = crypto.SHA512, kms.SigningAlgorithmSpecEcdsaSha512
	default:
		return nil, fmt.Errorf("unsupported signature algorithm %q", algorithm)
	}

	h := hashFunc.New()
	h.Write(data)
	digest := h.Sum(nil)

	sig, err := s.client.Sign(&kms.SignInput{
		KeyId:            &s.keyID,
		Message:          digest,
		MessageType:      aws.String("DIGEST"),
		SigningAlgorithm: aws.String(signingAlgorithm),
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
