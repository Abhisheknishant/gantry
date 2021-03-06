package cmd // import "github.com/ad-freiburg/gantry/cmd"

import (
	"fmt"

	"github.com/ad-freiburg/gantry/preprocessor"
	"github.com/spf13/cobra"
)

func init() {
	preprocessorCmd.AddCommand(preprocessorStatementsCmd)
}

var preprocessorStatementsCmd = &cobra.Command{
	Use:   "statements",
	Short: "",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Available preprocessor statements:")
		p, err := preprocessor.NewPreprocessor()
		if err != nil {
			return err
		}
		for _, f := range p.Functions() {
			fmt.Printf("\n%s\n", f.Usage())
		}
		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {},
}
