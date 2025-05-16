package cli

import (
	"strings"

	"github.com/idelchi/envsync/internal/profile"
	"github.com/spf13/cobra"
)

// Convert returns the cobra command for converting the config format.
func Convert() *cobra.Command {
	return &cobra.Command{
		Use:   "convert <yaml|toml|json>",
		Short: "Convert the current config to another format and write to file",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			ext := "." + strings.ToLower(args[0])

			store, err := profile.Load(fileFlag)
			if err != nil {
				return err
			}

			store.Path = "envsync" + ext
			store.Ext = ext
			return store.Save()
		},
	}
}
