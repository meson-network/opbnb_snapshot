package greenfield

import "github.com/urfave/cli/v2"

func GetFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "dir", Required: true},
		&cli.StringFlag{Name: "bucket_name", Required: true},
		&cli.StringFlag{Name: "additional_path", Required: false},
		&cli.IntFlag{Name: "thread", Required: false},
		&cli.IntFlag{Name: "retry_times", Required: false},
		&cli.StringFlag{Name: "private_key", Required: true},
		&cli.StringFlag{Name: "rpc_addr", Required: true},
		&cli.StringFlag{Name: "chain_id", Required: true},
	}
}

func ReadParam(clictx *cli.Context) (string, int, string, string,
	string, string, string, int) {

	originDir := clictx.String("dir")
	thread := clictx.Int("thread")
	bucketName := clictx.String("bucket_name")
	additional_path := clictx.String("additional_path")
	privateKey := clictx.String("private_key")
	rpcAddr := clictx.String("rpc_addr")
	chainId := clictx.String("chain_id")
	retry_times := clictx.Int("retry_times")

	return originDir, thread, bucketName, additional_path,
		privateKey, rpcAddr, chainId, retry_times
}
