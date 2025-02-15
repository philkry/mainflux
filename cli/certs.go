package cli

import (
	"github.com/spf13/cobra"
)

var cmdCerts = []cobra.Command{
	{
		Use:   "get <cert_serial> <user_auth_token>",
		Short: "Get certificate",
		Long:  `Gets a certificate for a given cert ID.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				logUsage(cmd.Use)
				return
			}
			cert, err := sdk.ViewCert(args[0], args[1])
			if err != nil {
				logError(err)
				return
			}
			logJSON(cert)
		},
	},
	{
		Use:   "revoke <thing_id> <user_auth_token>",
		Short: "Revoke certificate",
		Long:  `Revokes a certificate for a given thing ID.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				logUsage(cmd.Use)
				return
			}
			rtime, err := sdk.RevokeCert(args[0], args[1])
			if err != nil {
				logError(err)
				return
			}
			logRevokedTime(rtime)
		},
	},
}

// NewCertsCmd returns certificate command.
func NewCertsCmd() *cobra.Command {
	var ttl string

	issueCmd := cobra.Command{
		Use:   "issue <thing_id> <user_auth_token> [--ttl=8760h]",
		Short: "Issue certificate",
		Long:  `Issues new certificate for a thing`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				logUsage(cmd.Use)
				return
			}

			thingID := args[0]

			c, err := sdk.IssueCert(thingID, ttl, args[1])
			if err != nil {
				logError(err)
				return
			}
			logJSON(c)
		},
	}

	issueCmd.Flags().StringVar(&ttl, "ttl", "8760h", "certificate time to live in duration")

	cmd := cobra.Command{
		Use:   "certs [issue | get | revoke ]",
		Short: "Certificates management",
		Long:  `Certificates management: issue, get or revoke certificates for things"`,
	}

	cmdCerts = append(cmdCerts, issueCmd)

	for i := range cmdCerts {
		cmd.AddCommand(&cmdCerts[i])
	}

	return &cmd
}
