package config

type Config struct {
	StrategyToUse string
}

func GetConfig() Config {
	return Config{
		StrategyToUse: "plug",
	}
}
