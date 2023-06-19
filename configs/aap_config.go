package configs

func GetAuthenticateAndPostConfig(cfgPath string) (*AuthenticateAndPostConfig, error) {
	config, err := parseConfig(cfgPath)
	if err != nil {
		return &AuthenticateAndPostConfig{}, err
	}

	return &config.AuthenticateAndPost, nil
}