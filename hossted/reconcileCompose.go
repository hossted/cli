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

	// get OrgID and HosstedApiUrl
	orgID, hosstedAPIUrl, projectName, osUUID, err := compose.GetOrgIDHosstedApiUrl()
	if err != nil {
		fmt.Println(err)
		return err
	}

	/*
		================================ TBD: To be discussed ================================
		sendComposeInfo inside compose.ComposeReconciler uses HOSSTED_API_URL env variable which will not
		be present during cron ReconcileCompose. So for now, the below function is not sending compose info.
	*/
	err = compose.ReconcileCompose(osUUID, orgID, emailID, resp.Token, projectName, hosstedAPIUrl, "true")
	if err != nil {
		return err
	}

	return nil
}
