package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "generate":
		generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)

		host := generateCmd.String("host", "", "The SSH server host (e.g., example.com)")
		alias := generateCmd.String("alias", "", "The SSH alias for the host")
		keyName := generateCmd.String("keyName", "", "Name of the key file (default is id_rsa)")
		passphrase := generateCmd.String("passphrase", "", "SSH key passphrase")

		generateCmd.Usage = func() {
			fmt.Println("Usage: crane-ssh generate [options]")
			fmt.Println("\nOptions:")
			generateCmd.PrintDefaults()
			fmt.Println("\nExample:")
			fmt.Println("crane-ssh generate --host=github.com --alias=github.com --keyName=id_rsa")
		}
		generateCmd.Parse(os.Args[2:])

		if *host == "" {
			fmt.Print("Enter SSH host (which will be saved in the ~/.ssh/config file): ")
			*host = readInput()
		}

		if *alias == "" {
			fmt.Print("Enter alias for the SSH host (which will be saved in the ~/.ssh/config file): ")
			*alias = readInput()
		}

		if *keyName == "" {
			fmt.Print("Enter file name in which to save the key (default: id_rsa saved in ~/.ssh/): ")
			*keyName = readInput()
		}

		if *passphrase == "" {
			fmt.Print("Enter passphrase (empty for no passphrase): ")
			*passphrase = readInput()
		}

		if *host == "" || *alias == "" {
			log.Fatal("Both host and alias are required for 'generate' command.")
		}
		if *keyName == "" {
			*keyName = "id_rsa"
		}

		runGenerate(*host, *alias, *keyName, *passphrase)
	default:
		printHelp()
	}
}

func runGenerate(host, alias, keyName, passphrase string) {
	sshDir := filepath.Join(os.Getenv("HOME"), ".ssh")
	pubKeyPath := filepath.Join(sshDir, keyName+".pub")

	if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
		fmt.Println("No SSH key found. Generating a new one...")

		err = generateSSHKey(sshDir, keyName, passphrase)
		if err != nil {
			log.Fatalf("Failed to generate SSH key: %v", err)
		}
	} else {
		fmt.Println("SSH key already exists.")
	}

	pubKey, err := os.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatalf("Failed to read public key: %v", err)
	}

	err = addToSSHConfig(sshDir, host, alias, keyName)
	if err != nil {
		log.Fatalf("Failed to update SSH config: %v", err)
	}

	err = clipboard.WriteAll(string(pubKey))
	if err != nil {
		fmt.Printf("Failed to copy public key to clipboard: %v\n", err)
		fmt.Printf("Falling back to showing public key created:\n\n")
		fmt.Println(string(pubKey))
	} else {
		fmt.Println("Public key copied to clipboard!")
	}
}

func generateSSHKey(sshDir, keyName, passphrase string) error {
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return fmt.Errorf("failed to create .ssh directory: %v", err)
	}

	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-f", filepath.Join(sshDir, keyName), "-N", passphrase)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func addToSSHConfig(sshDir, host, alias, keyName string) error {
	configFilePath := filepath.Join(sshDir, "config")

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		_, err := os.Create(configFilePath)
		if err != nil {
			return fmt.Errorf("failed to create SSH config file: %v", err)
		}
	}

	if hostExistsInConfig(configFilePath, alias) {
		fmt.Printf("Host %s already exists in the SSH config.\n", alias)
		return nil
	}

	configEntry := fmt.Sprintf("\nHost %s\n  HostName %s\n  IdentityFile %s\n  Preferredauthentications publickey\n  IdentitiesOnly yes\n", alias, host, filepath.Join(sshDir, keyName))

	file, err := os.OpenFile(configFilePath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open SSH config file for writing: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(configEntry)
	if err != nil {
		return fmt.Errorf("failed to write to SSH config file: %v", err)
	}

	fmt.Printf("Host %s (%s) added to SSH config.\n", alias, host)
	return nil
}

func hostExistsInConfig(configFilePath, alias string) bool {
	file, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("Failed to open SSH config file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "Host") && strings.Contains(line, alias) {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading SSH config file: %v", err)
	}

	return false
}

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func printHelp() {
	fmt.Println("Usage: crane-ssh <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  generate   Generate a new SSH key pair, configure SSH settings in ~/.ssh/config and copies public key to clipboard")
	fmt.Println("\nUse \"crane-ssh <command> --help\" for more information about a command.")
}
