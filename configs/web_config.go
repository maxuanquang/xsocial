package configs

func GetWebConfig(cfgPath string) (*WebConfig, error) {
	config, err := parseConfig(cfgPath)
	if err != nil {
		return &WebConfig{}, err
	}

	return &config.Web, nil
}
