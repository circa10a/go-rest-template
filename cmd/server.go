package cmd

import (
	"fmt"

	"github.com/circa10a/go-rest-template/internal/server"
	"github.com/spf13/cobra"
)

var (
	serverAutoTLS   bool
	serverDomains   []string
	serverLogFormat string
	serverLogLevel  string
	serverMetrics   bool
	serverPort      int
	serverTLSCert   string
	serverTLSKey    string
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: fmt.Sprintf("Start the %s server", project),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := server.New(
			server.WithAutoTLS(serverAutoTLS),
			server.WithDomains(serverDomains),
			server.WithLogFormat(serverLogFormat),
			server.WithLogLevel(serverLogLevel),
			server.WithMetrics(serverMetrics),
			server.WithPort(serverPort),
			server.WithTLSCert(serverTLSCert),
			server.WithTLSKey(serverTLSKey),
			server.WithValidation(true),
		)
		if err != nil {
			return err
		}

		err = s.Start()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().BoolVarP(&serverAutoTLS, "auto-tls", "a", false, "Enable automatic TLS via Let's Encrypt. Requires port 80/443 open to the internet for domain validation.")
	serverCmd.Flags().StringVarP(&serverLogFormat, "log-format", "f", "text", "Server logging format. Supported values are 'text' and 'json'.")
	serverCmd.Flags().StringVarP(&serverLogLevel, "log-level", "l", "info", "Server logging level.")
	serverCmd.Flags().StringArrayVarP(&serverDomains, "domains", "d", []string{}, "Domains to issue certificate for. Must be used with --auto-tls.")
	serverCmd.Flags().BoolVarP(&serverMetrics, "metrics", "m", false, "Enable Prometheus metrics intrumentation.")
	serverCmd.Flags().IntVarP(&serverPort, "port", "p", 8080, "Port to listen on. Cannot be used in conjunction with --auto-tls since that will require listening on 80 and 443.")
	serverCmd.Flags().StringVarP(&serverTLSCert, "tls-certificate", "", "", "Path to custom TLS certificate. Cannot be used with --auto-tls.")
	serverCmd.Flags().StringVarP(&serverTLSKey, "tls-key", "", "", "Path to custom TLS key. Cannot be used with --auto-tls.")
}
