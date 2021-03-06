package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Nnachevvv/passmag/storage"
	"github.com/spf13/viper"
	"golang.org/x/crypto/argon2"
)

// User contains current logged user
type User struct {
	Password  []byte
	VaultPwd  []byte
	VaultPath string
	VaultData []byte
}

//EnterSession prompts to enter Session Key and ask for master password
func EnterSession() (User, error) {
	var path string
	var err error
	if !viper.IsSet("password.path") {
		path, err = storage.FilePath()
		if err != nil {
			return User{}, err
		}
	} else {
		path = viper.GetString("password.path")
	}

	if err := storage.VaultExist(path); err != nil {
		return User{}, err
	}

	var sessionKey, masterPassword string
	if !viper.IsSet("PASS_SESSION") || viper.Get("PASS_SESSION") == "" {
		prompt := &survey.Input{Message: "Please enter your session key:"}
		survey.AskOne(prompt, &sessionKey, survey.WithValidator(survey.Required), survey.WithStdio(Stdio.In, Stdio.Out, Stdio.Err))
	} else {
		sessionKey = viper.GetString("PASS_SESSION")
	}

	prompt := &survey.Password{Message: "Enter your master password:"}
	survey.AskOne(prompt, &masterPassword, survey.WithValidator(survey.Required), survey.WithStdio(Stdio.In, Stdio.Out, Stdio.Err))
	u := User{Password: []byte(masterPassword),
		VaultPwd:  argon2.IDKey([]byte(masterPassword), []byte(sessionKey), 1, 64*1024, 4, 32),
		VaultPath: path}

	vaultData, err := Crypt.DecryptFile(u.VaultPath, u.VaultPwd)
	if err != nil {
		return User{}, fmt.Errorf("failed to load your vault try again : %w", err)
	}

	u.VaultData = vaultData
	return u, err
}
