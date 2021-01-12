package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nnachevv/passmag/crypt"
	"github.com/nnachevv/passmag/storage"
	"github.com/nnachevv/passmag/user"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// the questions to ask
var addQs = []*survey.Question{
	{
		Name:   "name",
		Prompt: &survey.Input{Message: "Enter name adress:"},
	},
	{
		Name:   "password",
		Prompt: &survey.Password{Message: "Enter your password:"},
		Validate: func(val interface{}) error {
			if str, ok := val.(string); !ok || len(str) < 8 {
				return errors.New("password should be longer than 8 characters")
			}
			return nil
		},
	},
}

var add = &cobra.Command{

	Use:   "add",
	Short: "Initialize email, password and master password for your password manager",
	Long:  `Set master password`,
	RunE: func(cmd *cobra.Command, args []string) error {
		u, err := user.EnterSession()
		if err != nil {
			return err
		}

		err = addPasswords(u)
		if err != nil {
			return err
		}
		fmt.Println("succesfully added")

		return nil
	},
}

func addPasswords(u user.User) error {
	answers := struct {
		Name     string
		Password string
	}{}

	err := survey.Ask(addQs, &answers)
	if err != nil {
		return err
	}

	s, err := storage.Load(u.VaultData)
	if err != nil {
		return err
	}

	err = s.Add(answers.Name, answers.Password)
	if err != nil {
		var confirm bool
		editConfirm := &survey.Confirm{Message: "Do you want to edit name with newly password"}
		survey.AskOne(editConfirm, &confirm, survey.WithValidator(survey.Required))
		if confirm {
			s.Edit(answers.Name, answers.Password)
		}
		return nil
	}

	byteData, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal map : %w", err)
	}

	err = crypt.EncryptFile(u.VaultPath, byteData, u.VaultPwd)

	if err != nil {
		return fmt.Errorf("failed to encrypt sessionData : %w", err)
	}

	err = SyncVault(s, u.Password)
	if err != nil && err != ErrCreateUser {
		return err
	}

	return nil
}
