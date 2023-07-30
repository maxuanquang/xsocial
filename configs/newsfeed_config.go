package configs

func GetNewsfeedConfig(cfgPath string) (*NewsfeedConfig, error) {
	config, err := parseConfig(cfgPath)
	if err != nil {
		return &NewsfeedConfig{}, err
	}

	return &config.Newsfeed, nil
}