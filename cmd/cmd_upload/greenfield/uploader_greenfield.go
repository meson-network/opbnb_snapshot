package greenfield

import (
	"fmt"

	"github.com/meson-network/opbnb_snapshot/basic"
	"github.com/meson-network/opbnb_snapshot/basic/color"
	"github.com/meson-network/opbnb_snapshot/src/uploader/uploader_greenfield"
	"github.com/urfave/cli/v2"
)

func Uploader_greenfield(clictx *cli.Context) error {

	fmt.Println(color.Green(basic.Logo))

	originDir, thread, bucketName, additional_path,
		privateKey, rpcAddr, chainId, retry_times := ReadParam(clictx)

	return uploader_greenfield.Upload_greenfield(originDir, thread, bucketName, additional_path,
		privateKey, rpcAddr, chainId, retry_times)

}
