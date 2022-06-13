package hossted

func Dev() error {
	err := ChangeMOTD("dev.com")
	if err != nil {
		return err
	}

	return nil
}
