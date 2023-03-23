package hossted

func Schedule(env string) error {

	config, _ := GetConfig()

	if config.Update == true {
		Ping(env) //call hossted ping-send dockers info
		Scan(env) //call hossted scan-send vulnerabilities info
	}

	return nil
}
