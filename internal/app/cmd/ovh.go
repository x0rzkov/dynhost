package cmd

import (
	"errors"
	"github.com/Wojteek/dynhost/internal/app"
	"github.com/Wojteek/dynhost/internal/app/service"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// OVHCommand the definition of the command
func OVHCommand() *cli.Command {
	return &cli.Command{
		Name: "ovh",

		Usage: "OVH provider",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "auth-username",
				Usage:    "The authentication username of the DynHost option",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "auth-password",
				Usage:    "The authentication password of the DynHost option",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "hostname",
				Usage:    "The hostname of the DynHost option",
				Required: true,
			},
		},

		Before: func(ctx *cli.Context) error {
			if authUsername := ctx.String("auth-username"); len(authUsername) == 0 {
				return errors.New("auth-username is required")
			}

			if authPassword := ctx.String("auth-password"); len(authPassword) == 0 {
				return errors.New("auth-password is required")
			}

			if hostname := ctx.String("hostname"); len(hostname) == 0 {
				return errors.New("hostname is required")
			}

			if dataPath := ctx.String("data"); len(dataPath) == 0 {
				return errors.New("data is required")
			}

			return nil
		},

		Action: func(ctx *cli.Context) error {
			authUsername := ctx.String("auth-username")
			authPassword := ctx.String("auth-password")
			hostname := ctx.String("hostname")
			dataPath := ctx.String("data")
			timer := ctx.Duration("timer")

			processCommand := &app.ProcessCommand{
				DataPath: dataPath,
				Timer:    timer,
			}

			var changedIPCallback app.ChangedIPCallback = func(currentIP string) error {
				r := service.OVH{
					IP:       currentIP,
					Hostname: hostname,
					Credentials: service.CredentialsOVH{
						Username: authUsername,
						Password: authPassword,
					},
				}

				if _, err := r.UpdateRecordRequest(); err != nil {
					return err
				}

				return nil
			}
			var updateIP = app.UpdateIP(processCommand, changedIPCallback)

			log.WithFields(log.Fields{
				"service": "ovh",
				"timer":   timer,
			}).Info("The DynHost is running")

			if timer == 0 {
				_ = updateIP()
			} else {
				app.Timer(timer, updateIP)
			}

			return nil
		},
	}
}
