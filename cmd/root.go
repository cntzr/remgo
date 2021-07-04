package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"

	"github.com/spf13/viper"
)

var (
	cfgFile string
	dataDir string
	dbgMode bool
	langMsg string
	prnt    = message.NewPrinter(language.English)

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "remgo",
		Short: "calendar and reminder with CLI",
		Long: `remgo simply manages dates. It uses a CLI
and offers overviews of appointments for single days
or periods. The default output goes to a terminal with
color support. for future releases it will generate plain 
output for mail exchange, too. It makes a nice figure as 
part of pipes supporting at least stdout.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.config/remgo/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&dataDir, "data-dir", "~/.reminders", "data directory for reminder files (default is ~/.reminders)")
	rootCmd.PersistentFlags().StringVar(&langMsg, "language", "", "language for messages (de|en), default is en")
	rootCmd.PersistentFlags().BoolVarP(&dbgMode, "debug", "", false, "debug mode, default is off")

	// i18n support >>
	for _, e := range entries {
		tag := language.MustParse(e.tag)
		switch msg := e.msg.(type) {
		case string:
			message.SetString(tag, e.key, msg)
		case catalog.Message:
			message.Set(tag, e.key, msg)
		case []catalog.Message:
			message.Set(tag, e.key, msg...)
		}
	}
	// << i18n support

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name "config" (without extension).
		viper.AddConfigPath(home + "/.config/remgo")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// Set config
		if viper.IsSet("debug") && !dbgMode {
			dbgMode = viper.GetBool("debug")
		}
		if viper.IsSet("data-dir") && dataDir == "" {
			dataDir = viper.GetString("data-dir")
		} else {
			dataDir = home + "/.reminders"
		}
		if viper.IsSet("language") && langMsg == "" {
			langMsg = viper.GetString("language")
			if langMsg == "de" {
				prnt = message.NewPrinter(language.German)
			}
		}
		if dbgMode {
			// Output in debug mode only
			prnt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
			prnt.Printf("Reading from data directory: %s\n", dataDir)
			prnt.Printf("Language for messages: %s\n", langMsg)
		}
	}
}
