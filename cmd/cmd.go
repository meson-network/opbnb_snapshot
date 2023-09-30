package cmd

import (
	"fmt"

	"github.com/meson-network/opbnb_snapshot/cmd/cmd_download"
	"github.com/meson-network/opbnb_snapshot/cmd/cmd_endpoint"
	"github.com/meson-network/opbnb_snapshot/cmd/cmd_split"
	"github.com/meson-network/opbnb_snapshot/cmd/cmd_upload/greenfield"
	"github.com/meson-network/opbnb_snapshot/cmd/cmd_upload/r2"
	"github.com/urfave/cli/v2"
)

func ConfigCmd() *cli.App {

	return &cli.App{
		CommandNotFound: func(context *cli.Context, s string) {
			fmt.Println("command not find, use -h or --help show help")
		},

		Commands: []*cli.Command{
			{
				Name:  "download",
				Usage: "multithread download and merge files",
				Flags: cmd_download.GetFlags(),
				Action: func(clictx *cli.Context) error {
					cmd_download.Download(clictx)
					return nil
				},
			},
			{
				Name:  "split",
				Usage: "split data file to small files",
				Flags: cmd_split.GetFlags(),
				Action: func(clictx *cli.Context) error {
					cmd_split.Split(clictx)
					return nil
				},
			},
			{
				Name:  "upload",
				Usage: "upload files",
				Subcommands: []*cli.Command{
					{
						Name:  "r2",
						Flags: r2.GetFlags(),
						Usage: "upload to cloudflare R2 storage",
						Action: func(clictx *cli.Context) error {
							r2.Uploader_r2(clictx)
							return nil
						},
					},
					{
						Name:  "greenfield",
						Flags: greenfield.GetFlags(),
						Usage: "upload to Greenfield storage",
						Action: func(clictx *cli.Context) error {
							greenfield.Uploader_greenfield(clictx)
							return nil
						},
					},
				},
			},
			{
				Name:  "endpoint",
				Usage: "set endpoint",
				Subcommands: []*cli.Command{
					{
						Name:  "add",
						Flags: cmd_endpoint.GetFlags(),
						Usage: "add new endpoints",
						Action: func(clictx *cli.Context) error {
							cmd_endpoint.AddEndpoint(clictx)
							return nil
						},
					},
					{
						Name:  "remove",
						Flags: cmd_endpoint.GetFlags(),
						Usage: "remove endpoints",
						Action: func(clictx *cli.Context) error {
							cmd_endpoint.RemoveEndpoint(clictx)
							return nil
						},
					},
					{
						Name:  "set",
						Flags: cmd_endpoint.GetFlags(),
						Usage: "reset endpoints",
						Action: func(clictx *cli.Context) error {
							cmd_endpoint.SetEndpoint(clictx)
							return nil
						},
					},
					{
						Name:  "clear",
						Flags: cmd_endpoint.GetFlags(),
						Usage: "remove all exist endpoints",
						Action: func(clictx *cli.Context) error {
							cmd_endpoint.ClearEndpoint(clictx)
							return nil
						},
					},
					{
						Name:  "print",
						Flags: cmd_endpoint.GetFlags(),
						Usage: "remove all exist endpoints",
						Action: func(clictx *cli.Context) error {
							cmd_endpoint.PrintEndpoint(clictx)
							return nil
						},
					},
				},
			},
		},
	}
}
