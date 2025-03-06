/*
Copyright © 2025 Seednode <seednode@seedno.de>
*/

package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	ReleaseVersion string = "1.0.0"
)

var (
	bind        string
	exitOnError bool
	port        uint16
	profile     bool
	tlsCert     string
	tlsKey      string
	verbose     bool
	version     bool
)

func main() {
	cmd := &cobra.Command{
		Use:   "subnetter",
		Short: "Serves a tool for learning IP subnetting.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			initializeConfig(cmd)

			if tlsCert == "" && tlsKey != "" || tlsCert != "" && tlsKey == "" {
				return errors.New("TLS certificate and keyfile must both be specified to enable HTTPS")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := servePage()

			return err
		},
	}

	cmd.Flags().StringVarP(&bind, "bind", "b", "0.0.0.0", "address to bind to")
	cmd.Flags().BoolVar(&exitOnError, "exit-on-error", false, "shut down webserver on error, instead of just printing the error")
	cmd.Flags().Uint16VarP(&port, "port", "p", 8080, "port to listen on")
	cmd.Flags().BoolVar(&profile, "profile", false, "register net/http/pprof handlers")
	cmd.Flags().StringVar(&tlsCert, "tls-cert", "", "path to TLS certificate")
	cmd.Flags().StringVar(&tlsKey, "tls-key", "", "path to TLS keyfile")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "log requests to stdout")
	cmd.Flags().BoolVarP(&version, "version", "V", false, "display version and exit")

	cmd.Flags().SetInterspersed(true)

	cmd.CompletionOptions.HiddenDefaultCmd = true

	cmd.SilenceErrors = true
	cmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	cmd.SetVersionTemplate("subnetter v{{.Version}}\n")
	cmd.Version = ReleaseVersion

	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func initializeConfig(cmd *cobra.Command) {
	v := viper.New()

	v.SetEnvPrefix("subnetter")

	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	v.AutomaticEnv()

	bindFlags(cmd, v)
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := strings.ReplaceAll(f.Name, "-", "_")

		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
