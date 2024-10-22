package utils

import "github.com/spf13/viper"

const CONFIG_NAME = "config"

type ConfigKey string;

const (
	PostgresHost ConfigKey = "postgres.host";
	PostgresPort = "postgres.port"
	PostgresUser = "postgres.user";
	PostgresPassword = "postgres.password";
	PostgresDb = "postgres.db"

	ServerPort = "server.port"

	MetadataServer = "nft-server.address"
)

func ReadConfig() error {
	setDefaults();
	viper.SetConfigName("config");
	viper.SetConfigType("toml");
  viper.AddConfigPath(".");
  viper.AddConfigPath("./.local");
  return viper.ReadInConfig();
}

func setDefaults() {
	setViperDefaultWithKey(PostgresHost, "localhost");
	setViperDefaultWithKey(PostgresPort, "5432");
	setViperDefaultWithKey(PostgresUser, "admin");
	setViperDefaultWithKey(PostgresPassword, "admin");
	setViperDefaultWithKey(PostgresDb, "wallets");
	setViperDefaultWithKey(ServerPort, "80");
	setViperDefaultWithKey(MetadataServer, "http://localhost:5138")
}

func setViperDefaultWithKey(key ConfigKey, value string) {
	viper.SetDefault(string(key), value);
}

func GetConfigString(key ConfigKey) string {
	return viper.GetString(string(key))
}
