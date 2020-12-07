package impl

import (
	"log"

	"github.com/urfave/cli/v2"

	"github.com/lexycore/gitlab-go-import/version"
)

const (
	envPrefix = "GO_IMPORT_"

	bindAddrDefault = "127.0.0.1:8008" // Default binding address
)

// CreateCLI creates a new cli Application with Name, Usage, Version and Actions
func CreateCLI() *cli.App {
	app := cli.NewApp()
	app.Name = version.Description
	app.Usage = version.Usage
	app.Version = version.Version()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "bind",
			Aliases: []string{"b"},
			Value:   bindAddrDefault,
			EnvVars: []string{envPrefix + "BIND_ADDR"},
			Usage:   "bind address",
		},
		&cli.StringFlag{
			Name:    "gitlab-url",
			Aliases: []string{"g"},
			EnvVars: []string{envPrefix + "GITLAB_URL"},
			Usage:   "GitLab server URL",
		},
		&cli.StringFlag{
			Name:    "gitlab-token",
			Aliases: []string{"gt"},
			EnvVars: []string{envPrefix + "GITLAB_TOKEN"},
			Usage:   "GitLab private token with API access",
		},
	}
	app.Action = func(ctx *cli.Context) error {
		config := &appConfig{
			BindAddr:    ctx.String("bind"),
			GitLabURL:   ctx.String("gitlab-url"),
			GitLabToken: ctx.String("gitlab-token"),
		}

		server := NewServer(config)
		log.Printf("service version: %s", version.Version())
		if err := server.Serve(); err != nil {
			return err
		}
		return nil
	}
	return app
}
