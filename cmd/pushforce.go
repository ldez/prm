package cmd

import (
	"fmt"
	"os"

	"github.com/ldez/prm/config"
	"github.com/ldez/prm/types"
)

func PushForce(options *types.PushForceOptions) error {

	// get configuration
	confs, err := config.ReadFile()
	if err != nil {
		return err
	}

	repoDir, err := os.Getwd()
	if err != nil {
		return err
	}

	con, err := config.Find(confs, repoDir)
	if err != nil {
		return err
	}

	if pr, err := con.FindPullRequests(options.Number); err == nil {
		fmt.Println("push force", pr)

		err := pr.PushForce()
		if err != nil {
			return err
		}
	}

	return nil
}
