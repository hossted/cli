package hossted

func Dev() error {
	err := ChangeMOTD("example.com")
	if err != nil {
		return err
	}

	return nil
}
