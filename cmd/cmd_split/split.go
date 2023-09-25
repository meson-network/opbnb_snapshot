package cmd_split

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/meson-network/opbnb_snapshot/basic"
	"github.com/meson-network/opbnb_snapshot/basic/color"
	"github.com/meson-network/opbnb_snapshot/src/split"
)

func Split(clictx *cli.Context) error {
	fmt.Println(color.Green(basic.Logo))

	originFilePath, destDir, sizeStr, thread := ReadParam(clictx)

	return split.Split(originFilePath, destDir, sizeStr, thread)
}
