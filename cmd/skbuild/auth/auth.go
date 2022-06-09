package auth

import (
	"github.com/skema-dev/skemabuild/internal/auth"
	"github.com/skema-dev/skemabuild/internal/pkg/console"
	"github.com/spf13/cobra"
)

const (
	authDescription     = "Authentication from git provider"
	authLongDescription = "skbuild auth"
)

func NewCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "auth",
		Short: authDescription,
		Long:  authLongDescription,
		Run: func(c *cobra.Command, args []string) {
			providerType, _ := c.Flags().GetString("type")
			var provider auth.AuthProvider
			switch providerType {
			case "github":
				provider = auth.NewGithubAuthProvider()
			default:
				panic("incorrect provider type")
			}
			provider.StartAuthProcess()
			provider.SaveTokenToFile()

			savedToken := provider.GetLocalToken()
			console.Info(savedToken)
		},
	}

	cmd.Flags().StringP("type", "t", "github", "auth provider: github")

	return cmd
}
