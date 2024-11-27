package utils

import "github.com/spf13/viper"

const CONFIG_NAME = "config"

type ConfigKey string;

const (
  PostgresHost ConfigKey = "postgres.host";
  PostgresPort = "postgres.port"
  PostgresUser = "postgres.user";
  PostgresPassword = "postgres.password";
  PostgresDb = "postgres.db";
  ServerPort = "server.port";
  MetadataServer = "nft-server.address";
  NodeAddress = "node.address"
  AuthKey = "auth.admin_public_key";
  AuthEnabled = "auth.enabled"
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
  setViperDefaultWithKey(MetadataServer, "http://localhost:5138");
  setViperDefaultWithKey(NodeAddress, "http://localhost:7213");
  setViperDefaultWithKey(AuthKey, "");
  setViperDefaultWithKey(AuthEnabled, "true")
}

func setViperDefaultWithKey(key ConfigKey, value string) {
  viper.SetDefault(string(key), value);
}

func GetConfigString(key ConfigKey) string {
  return viper.GetString(string(key))
}

func GetConfigBool(key ConfigKey) bool {
  return viper.GetBool(string(key))
}
