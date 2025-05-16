package cli

import (
	"github.com/idelchi/envsync/internal/profile"
	"github.com/spf13/cobra"
)

// Add returns the cobra command for adding or updating a variable.
func Add() *cobra.Command {
	return &cobra.Command{
		Use:   "add <profile> <key> <val>",
		Short: "Add or update variable",
		Args:  cobra.ExactArgs(3),
		RunE: func(_ *cobra.Command, a []string) error {
			s, err := profile.Load(fileFlag)
			if err != nil {
				return err
			}
			return s.SetVar(a[0], a[1], a[2])
		},
	}
}
