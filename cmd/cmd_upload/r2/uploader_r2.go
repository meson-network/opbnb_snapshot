package r2

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/meson-network/opbnb_snapshot/basic"
	"github.com/meson-network/opbnb_snapshot/basic/color"
	"github.com/meson-network/opbnb_snapshot/src/uploader/uploader_r2"
)

func Uploader_r2(clictx *cli.Context) error {

	fmt.Println(color.Green(basic.Logo))

	originDir, thread, bucketName, additional_path,
		accountId, accessKeyId, accessKeySecret, retry_times := ReadParam(clictx)

	return uploader_r2.Upload_r2(originDir, thread, bucketName, additional_path,
		accountId, accessKeyId, accessKeySecret, retry_times)
}
