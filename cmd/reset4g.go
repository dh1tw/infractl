package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dh1tw/infractl/microtik"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// reset4gCmd represents the reset4g command
var reset4gCmd = &cobra.Command{
	Use:   "4g-reset",
	Short: "Hard power reset of a 4G USB modem connected to the microtik router",
	Long: `This command performs a hard power reset of a 4G stick connected
to the internal USB port of a microtik routerboard. The power will be cut for
5 seconds.

You can save the details of your microtik router in the config file under the
the key [microtik].
`,
	Run: reset4g,
}

func init() {
	rootCmd.AddCommand(reset4gCmd)
	reset4gCmd.Flags().StringP("address", "a", "192.168.0.1", "address of your microtik router")
	reset4gCmd.Flags().IntP("port", "p", 8728, "API port of your microtik router")
	reset4gCmd.Flags().StringP("username", "U", "admin", "username for your microtik router")
	reset4gCmd.Flags().StringP("password", "P", "admin", "password for your microtik router")

}

func reset4g(cmd *cobra.Command, args []string) {
	// Try to read config file
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		if strings.Contains(err.Error(), "Not Found in") {
			fmt.Println("no config file found")
		} else {
			fmt.Println("Error parsing config file", viper.ConfigFileUsed())
			fmt.Println(err)
			os.Exit(1)
		}
	}

	viper.BindPFlag("microtik.address", cmd.Flags().Lookup("address"))
	viper.BindPFlag("microtik.port", cmd.Flags().Lookup("port"))
	viper.BindPFlag("microtik.username", cmd.Flags().Lookup("username"))
	viper.BindPFlag("microtik.password", cmd.Flags().Lookup("password"))

	mConfig := microtik.Config{
		Address:  viper.GetString("microtik.address"),
		Port:     viper.GetInt("microtik.port"),
		Username: viper.GetString("microtik.username"),
		Password: viper.GetString("microtik.password"),
	}

	mt := microtik.New(mConfig)

	// before we can reset the 4G modem, we must make sure that the ADSL route
	// is active. Otherwise, when the 4G route would become unavailable after the reset
	// and no other route is available, microtik generates a new dynamical route which
	// messes up the configuration.
	if err := mt.SetRoute("adsl", "disabled=false"); err != nil {
		log.Fatal(err)
	}

	if err := mt.Reset4G(); err != nil {
		log.Fatal(err)
	}
	log.Println("4G reset successfully initiated")
}
