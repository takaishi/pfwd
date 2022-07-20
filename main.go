package main

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"io/ioutil"
	"os"
	"os/exec"
)

type Config struct {
	Gateways map[string]GatewayConfig
}
type GatewayConfig struct {
	IdentityFile     string `yaml:"identity_file"`
	LoginName        string `yaml:"login_name"`
	LocalBindAddress string `yaml:"local_bind_address"`
	Port             string `yaml:"port"`
	Destination      string `yaml:"destination"`
}

func main() {
	if err := _main(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func _main() error {
	profile := os.Args[1]
	config, err := loadConfig()
	if err != nil {
		return err
	}
	gwConfig := config.Gateways[profile]

	if err := startForward(profile, gwConfig); err != nil {
		return err
	}

	if err := execCommand(os.Args[2], os.Args[3:]); err != nil {
		return err
	}

	if err := stopForward(profile, gwConfig); err != nil {
		return err
	}

	return nil
}

func loadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(fmt.Sprintf("%s/.sshwrap.yaml", home))
	if err != nil {
		return nil, err
	}
	config := Config{}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func startForward(profile string, gwConfig GatewayConfig) error {
	args := []string{"-M", "-f", "-N", "-S", fmt.Sprintf("/tmp/.%s", profile)}
	if gwConfig.IdentityFile != "" {
		args = append(args, []string{"-i", gwConfig.IdentityFile}...)
	}
	if gwConfig.LoginName != "" {
		args = append(args, []string{"-l", gwConfig.LoginName}...)
	}
	if gwConfig.Port != "" {
		args = append(args, []string{"-p", gwConfig.Port}...)
	} else {
		args = append(args, []string{"-p", "22"}...)
	}
	if gwConfig.LocalBindAddress != "" {
		args = append(args, []string{"-L", gwConfig.LocalBindAddress}...)
	}
	args = append(args, gwConfig.Destination)
	cmd := exec.Command("ssh", args...)
	fmt.Printf("%+v\n", cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func stopForward(profile string, gwConfig GatewayConfig) error {
	args := []string{"-S", fmt.Sprintf("/tmp/.%s", profile)}
	if gwConfig.IdentityFile != "" {
		args = append(args, []string{"-i", gwConfig.IdentityFile}...)
	}
	if gwConfig.LoginName != "" {
		args = append(args, []string{"-l", gwConfig.LoginName}...)
	}
	args = append(args, []string{"-O", "exit", gwConfig.Destination}...)
	cmd := exec.Command("ssh", args...)
	fmt.Printf("%+v\n", cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func execCommand(name string, args []string) error {
	cmd := exec.Command(name, args...)
	fmt.Printf("%+v\n", cmd)
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s\n", out)
	if err != nil {
		return err
	}
	return nil
}
