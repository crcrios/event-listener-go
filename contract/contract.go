package contract

import (
	"io/ioutil"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

// GetContractWithConfig creates a file system wallet.
//
//   Parameters:
//   conectionConfigFilePath: Path of the config file.
//   walletPath: File system path of the wallet.
//   orgMSPId: Org MSP ID.
//   walletIdentityLabel: Specifies the name to be associated with the identity.
//   userCertPath: Path of the user MSP cert.
//   privateKeyPath: Path of the user MSP private key.
//   channelName: Name of the channel.
//   contractName: Name of the contract.
//
//   Returns:
//   contract: Reference of the contract.
//   err: Error if exists.
func GetContractWithConfig(conectionConfigFilePath string, walletPath string, orgMSPId string, walletIdentityLabel string, userCertPath string, privateKeyPath string, channelName string, contractName string) (contract *gateway.Contract, err error) {
	wallet, err := CreateFileSystemWallet(walletPath, orgMSPId, walletIdentityLabel, userCertPath, privateKeyPath)
	if err != nil {
		return
	}

	gw, err := ConnectToGateway(conectionConfigFilePath, wallet, walletIdentityLabel)
	if err != nil {
		return
	}

	contract, err = GetContract(gw, channelName, contractName)
	if err != nil {
		return
	}
	return
}

// CreateFileSystemWallet creates a file system wallet.
//
//   Parameters:
//   walletPath: File system path of the wallet.
//   orgMSPId: Org MSP ID.
//   walletIdentityLabel: Specifies the name to be associated with the identity.
//   userCertPath: Path of the user MSP cert.
//   privateKeyPath: Path of the user MSP private key.
//
//   Returns:
//   wallet: Reference of the new wallet.
//   err: Error if exists.
func CreateFileSystemWallet(walletPath string, orgMSPId string, walletIdentityLabel string, userCertPath string, privateKeyPath string) (wallet *gateway.Wallet, err error) {
	wallet, err = gateway.NewFileSystemWallet(walletPath)
	if err != nil {
		return
	}

	if !wallet.Exists(walletIdentityLabel) {
		err = PopulateWallet(wallet, orgMSPId, walletIdentityLabel, userCertPath, privateKeyPath)
		if err != nil {
			return
		}
	}
	return
}

// PopulateWallet populates the identity and certs of the wallet.
//   Parameters:
//   wallet: Wallet reference.
//   orgMSPId: Org MSP ID.
//   walletIdentityLabel: Specifies the name to be associated with the identity.
//   userCertPath: Path of the user MSP cert.
//   privateKeyPath: Path of the user MSP private key.
//
//   Returns:
//   err: Error if exists.
func PopulateWallet(wallet *gateway.Wallet, orgMSPId string, walletIdentityLabel string, userCertPath string, privateKeyPath string) error {
	
	// _, err := ProvisionAwsIdentity(userCertPath string, privateKeyPath string)
	// if err != nil {
	// 	return err
	// }
	
	// Read the certificate
	cert, err := ioutil.ReadFile(filepath.Clean(userCertPath))
	if err != nil {
		return err
	}

	// Read the private key
	key, err := ioutil.ReadFile(filepath.Clean(privateKeyPath))
	if err != nil {
		return err
	}

	// Generate identity
	identity := gateway.NewX509Identity(orgMSPId, string(cert), string(key))

	return wallet.Put(walletIdentityLabel, identity)
}

// ConnectToGateway connect to a gateway defined by a network config file.
//   Parameters:
//   conectionConfigFilePath: Path of the config file.
//   wallet: Wallet reference.
//   walletIdentityLabel: Specifies the name to be associated with the identity.
//
//   Returns:
//   gw: Reference of the gateway created.
//   err: Error if exists.
func ConnectToGateway(conectionConfigFilePath string, wallet *gateway.Wallet, walletIdentityLabel string) (gw *gateway.Gateway, err error) {
	gw, err = gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(conectionConfigFilePath))),
		gateway.WithIdentity(wallet, walletIdentityLabel),
	)
	defer gw.Close()
	return
}

// GetContract get the contract of a channel.
//   Parameters:
//   gw: Reference of the gateway
//   channelName: Name of the channel.
//   contractName: Name of the contract.
//
//   Returns:
//   contract: Reference of the contract.
//   err: Error if exists.
func GetContract(gw *gateway.Gateway, channelName string, contractName string) (contract *gateway.Contract, err error) {
	network, err := gw.GetNetwork(channelName)
	if err != nil {
		return
	}

	contract = network.GetContract(contractName)
	return
}
