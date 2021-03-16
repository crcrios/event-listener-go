package aws

import (
	
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/acmpca"
	"encoding/json"
	"io/ioutil"

)

type TLS_STUCT struct {
	Key        string `json:"key"`
	Cer        string `json:"cer"`
	Chain      string `json:"chain"`
}

var awsRegion = "us-east-1"
var awsProfile = "acm-pca-blockchain"
var awsArnCaCert = "arn:aws:acm-pca:us-east-1:872308410481:certificate-authority/ee2eadae-1a4e-4034-9f22-cc2626854c20"

func GetSecretValue(secretId string) (string, error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{Region: aws.String(awsRegion)},
	}))

	smc := secretsmanager.New(sess)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretId),
	}

	result, err := smc.GetSecretValue(input)
	if err != nil {
		return "", err
	}

	return *result.SecretString, nil
}

func GetCertificate(CertificateArn string, CertificateAuthorityArn string)  (string, error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{Region: aws.String(awsRegion)},
		Profile: awsProfile,
	}))

	apc := acmpca.New(sess)
	input := &acmpca.GetCertificateInput{
		CertificateArn: aws.String(CertificateArn),
		CertificateAuthorityArn : aws.String(CertificateAuthorityArn),
	}
	
	result, err := apc.GetCertificate(input)
	if err != nil {
		return "",err
	}

	return *result.Certificate,nil
}

func ProvisionTlsCertificates() (err error) {

	secretString,err := GetSecretValue("nu0094001-blockchain-dev-cli-tls")
	if err != nil {
		return
	}

	var tlsResult TLS_STUCT
	err = json.Unmarshal([]byte(secretString), &tlsResult)

	err = ioutil.WriteFile("certs/tls/ca.crt", []byte(tlsResult.Chain), 0644)
	if err != nil {
		return
	}
	err = ioutil.WriteFile("certs/tls/client.crt", []byte(tlsResult.Cer), 0644)
	if err != nil {
		return
	}
	err = ioutil.WriteFile("certs/tls/client.key", []byte(tlsResult.Key), 0644)
	if err != nil {
		return
	}

	return
}

func ProvisionMspCertificates() (err error) {

	secretString,err := GetSecretValue("nu0094001-blockchain-dev-ECDSA-peer-admin")
	if err != nil {
		return
	}

	certificate,err := GetCertificate(secretString, awsArnCaCert)
	if err != nil {
		return
	}

	err = ioutil.WriteFile("certs/msp/Admin@Bancolombia-cert.pem", []byte(certificate), 0644)
	if err != nil {
		return
	}

	secretString,err = GetSecretValue("nu0094001-blockchain-dev-ECDSA-Key-peer-admin")
	if err != nil {
		return
	}

	err = ioutil.WriteFile("certs/msp/priv_sk", []byte(secretString), 0644)
	if err != nil {
		return
	}

	return
}

//git clone https://github.com/crcrios/event-listener-go.git