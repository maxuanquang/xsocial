package configs

func GetWebConfig(cfgPath string) (*WebConfig, error) {
	allConfig, err := parseConfig(cfgPath)
	if err != nil {
		return &WebConfig{}, err
	}

	return &allConfig.WebConfig, nil
}