package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Nnachevvv/passmag/storage"
	"github.com/spf13/cobra"
)

// NewRemoveCmd creates a new removeCmd
func NewRemoveCmd() *cobra.Command {
	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove password from your password manager",
		Long:  `Remove password from your password manager from given host`,
		RunE: func(cmd *cobra.Command, args []string) error {
			u, err := EnterSession()
			if err != nil {
				return err
			}

			var removeName string
			prompt := &survey.Password{Message: "Enter name of password you want to remove:"}

			err = survey.AskOne(prompt, &removeName, survey.WithStdio(Stdio.In, Stdio.Out, Stdio.Err))
			if err != nil {
				return err
			}

			s, err := storage.Load(u.VaultData)
			if err != nil {
				return err
			}

			err = s.Remove(removeName)
			if err != nil {
				return err
			}

			byteData, err := json.Marshal(s)
			if err != nil {
				return fmt.Errorf("failed to marshal map : %w", err)
			}

			err = Crypt.EncryptFile(u.VaultPath, byteData, u.VaultPwd)

			if err != nil {
				return fmt.Errorf("failed to encrypt sessionData : %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "successfully removed password")
			s.SyncStorage(u.Password, MongoDB, Client)

			return nil
		},
	}
	return removeCmd
}
