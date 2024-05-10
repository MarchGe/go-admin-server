package gen

import (
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/common/utils"
	"github.com/spf13/cobra"
)

var length int

const chars = "123456789123456789abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ"

var Gen = &cobra.Command{
	Use:   "gen",
	Short: "generate key etc.",
	Long:  "This command is used to generate some information, etc random string.",
	Args: func(cmd *cobra.Command, args []string) error {
		if length <= 0 {
			return errors.New("use --rand(-r) to specify length of random string\n")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(utils.RandomStringFrom(length, chars))
	},
}

func init() {
	Gen.PersistentFlags().IntVarP(&length, "rand", "r", 0, "generate a random string of specified length.")
}
