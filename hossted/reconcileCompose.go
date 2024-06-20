package hossted

import (
	"fmt"

	"github.com/hossted/cli/hossted/service/common"
	"github.com/hossted/cli/hossted/service/compose"
)

func ReconcileCompose() error {
	emailID, err := common.GetEmail()
	if err != nil {
		return err
	}

	resp, err := common.GetLoginResponse()
	if err != nil {
		return err
	}

	// get OrgID
	orgID, err := compose.GetOrgID()
	if err != nil {
		fmt.Println(err)
		return err
	}

	/*
		================================ TBD: To be discussed ================================
		sendComposeInfo inside compose.ComposeReconciler uses HOSSTED_API_URL env variable which will not
		be present during cron ReconcileCompose. So for now, the below function is not sending compose info.
	*/
	err = compose.ComposeReconciler(orgID, emailID, resp.Token)
	if err != nil {
		return err
	}

	return nil
}
