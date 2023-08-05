package configs

func GetNewsfeedPublishingConfig(cfgPath string) (*NewsfeedPublishingConfig, error) {
	config, err := parseConfig(cfgPath)
	if err != nil {
		return &NewsfeedPublishingConfig{}, err
	}

	return &config.NewsfeedPublishing, nil
}
