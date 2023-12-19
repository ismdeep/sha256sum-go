package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ismdeep/sha256sum-go/internal/core"
)

func main() {

	directory, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var checkFlag bool
	var name string
	c := cobra.Command{
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch checkFlag {
			case true:
				// verify
				if err := core.NewSHA256Sum(directory, name).Verify(); err != nil {
					return err
				}
				fmt.Println("Verification successful.")
			default:
				// generate
				if err := core.NewSHA256Sum(directory, name).Generate(); err != nil {
					return err
				}
				fmt.Println("Checksums generated successfully.")
			}
			return nil
		},
	}
	c.PersistentFlags().BoolVarP(&checkFlag, "check", "c", false, "check sha256sum")
	c.PersistentFlags().StringVarP(&name, "name", "n", "sha256sum.txt", "sha256sum file name")
	if err := c.Execute(); err != nil {
		panic(err)
	}
}
