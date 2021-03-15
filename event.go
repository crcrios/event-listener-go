package main

import (
	"log"

	// "event/contract"

	"event/aws"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

const (
	orgMSP              = "Org0"
	walletPath          = "wallet"
	walletIdentityLabel = "appUser"
	userCertPath        = "certs\\msp\\Admin@Org0-cert.pem"
	privateKeyPath      = "certs\\msp\\priv_sk"
	channelName         = "mychannel"
	contractName        = "basic"
	conectionConfigPath = "network\\connection-dev.yaml"
)

func main() {
	//api.HandleRequests()

	// contract, err := contract.GetContractWithConfig(conectionConfigPath, walletPath, orgMSP, walletIdentityLabel, userCertPath, privateKeyPath, channelName, contractName)
	// if err != nil {
	// 	log.Fatalf("Failed to get contract: %v", err)
	// }

	log.Println("Inicio!")

	aws.ProvisionTlsCertificates()
	aws.ProvisionMspCertificates()
	//eventListener(contract)
	//callFunction(contract)

	log.Println("Funciona!")
}

func callFunction(contract *gateway.Contract) {
	//result, err := contract.SubmitTransaction("Clear")
	//result, err := contract.EvaluateTransaction("QueryCouchDB", "{\"selector\":{}}")
	result, err := contract.EvaluateTransaction("GetAllAssets", "{\"selector\":{}}")

	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}

	log.Println(string(result))
}

func eventListener(contract *gateway.Contract) {
	log.Println("Event listener started...")
	eventID := "Clear"

	reg, notifier, err := contract.RegisterEvent(eventID)
	if err != nil {
		log.Printf("Failed to register contract event: %s", err)
		return
	}

	for i := range notifier {
		log.Println(i.EventName + " - " + string(i.Payload))
	}

	defer contract.Unregister(reg)
}
