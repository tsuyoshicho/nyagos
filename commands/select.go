package commands

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/zetamatta/nyagos/nodos"
)

func cmdShOpenWithDialog(ctx context.Context, cmd Param) (int, error) {
	for _, s := range cmd.Args()[1:] {
		fullpath, err := filepath.Abs(s)
		if err != nil {
			fmt.Fprintf(cmd.Err(), "%s: %s\n", s, err)
			continue
		}
		err = nodos.ShOpenWithDialog(fullpath, "")
		if err != nil {
			fmt.Fprintf(cmd.Err(), "%s: %s\n", s, err)
		}
	}
	return 0, nil
}
