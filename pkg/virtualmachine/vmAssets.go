package virtualmachine

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
)

type User struct {
	Name              string   `yaml:"name"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys"`
	Sudo              string   `yaml:"sudo"`
	Groups            string   `yaml:"groups"`
	Shell             string   `yaml:"shell"`
}

type CloudConfig struct {
	Users      []User `yaml:"users"`
	WriteFiles []struct {
		Encoding   string `yaml:"encoding"`
		Content    string `yaml:"content"`
		Path       string `yaml:"path"`
		Permission string `yaml:"permissions"`
	} `yaml:"write_files"`
	RunCmd []string `yaml:"runcmd"`
}

// Create ssh key for a virtual machine
func CreateVirtualmachineSSHKeyPair(outputKeyPath string) (err error) {
	// Generate RSA key.
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// Extract public component.
	// pub := key.Public()

	// Encode private key to PKCS#1 ASN.1 PEM.
	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	sshPubKey, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		return err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(sshPubKey)

	err = writeKeyToFile(keyPEM, outputKeyPath+"/private_key.pem")
	if err != nil {
		return err
	}

	err = writeKeyToFile(pubKeyBytes, outputKeyPath+"/public_key.pem")
	if err != nil {
		return err
	}

	return nil
}

// writePemToFile writes keys to a file
func writeKeyToFile(keyBytes []byte, saveFileTo string) error {
	err := os.WriteFile(saveFileTo, keyBytes, 0600)
	if err != nil {
		return err
	}

	log.Printf("Key saved to: %s", saveFileTo)
	return nil
}

func InjectSSHKeyIntoUserDataYamlFile(userDataYamlFile []byte, sshPubKey string) (newUserDataYamlFile []byte) {
	// Unmarshal YAML data into a struct
	var config CloudConfig
	err := yaml.Unmarshal(userDataYamlFile, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
	}

	// Add a new SSH key to the first user's "ssh_authorized_keys" array
	config.Users[0].SSHAuthorizedKeys = append(config.Users[0].SSHAuthorizedKeys, strings.Replace(sshPubKey, "\n", "", -1))

	// Marshal the modified struct back to YAML
	modifiedYAML, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("Error marshalling YAML: %v", err)
	}

	// Add #cloud-config at the start
	cloudConfigLine := []byte("#cloud-config\n\n")
	modifiedYAML = append(cloudConfigLine, modifiedYAML...)

	return modifiedYAML

}
