package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../../..")
)

var Config config

type config struct {
	SDKVersion string `yaml:"sdk_version"`
	SDKDataDir string `yaml:"sdk_data_dir"`

	Etcd struct {
		EtcdSchema string   `yaml:"etcdSchema"`
		EtcdAddr   []string `yaml:"etcdAddr"`
	}

	IsSkipDatabase bool `yaml:"is_skip_database"`

	Mysql struct {
		DBAddress      []string `yaml:"dbMysqlAddress"`
		DBUserName     string   `yaml:"dbMysqlUserName"`
		DBPassword     string   `yaml:"dbMysqlPassword"`
		DBDatabaseName string   `yaml:"dbMysqlDatabaseName"`
		DBTableName    string   `yaml:"DBTableName"`
		DBMsgTableNum  int      `yaml:"dbMsgTableNum"`
		DBMaxOpenConns int      `yaml:"dbMaxOpenConns"`
		DBMaxIdleConns int      `yaml:"dbMaxIdleConns"`
		DBMaxLifeTime  int      `yaml:"dbMaxLifeTime"`
	}

	Mongo struct {
		DBUri         string `yaml:"dbUri"`
		DBAddress     string `yaml:"dbAddress"`
		DBDirect      bool   `yaml:"dbDirect"`
		DBTimeout     int    `yaml:"dbTimeout"`
		DBDatabase    string `yaml:"dbDatabase"`
		DBSource      string `yaml:"dbSource"`
		DBUserName    string `yaml:"dbUserName"`
		DBPassword    string `yaml:"dbPassword"`
		DBMaxPoolSize int    `yaml:"dbMaxPoolSize"`
	}

	Redis struct {
		DBAddress     []string `yaml:"dbAddress"`
		DBMaxIdle     int      `yaml:"dbMaxIdle"`
		DBMaxActive   int      `yaml:"dbMaxActive"`
		DBIdleTimeout int      `yaml:"dbIdleTimeout"`
		DBPassWord    string   `yaml:"dbPassWord"`
		EnableCluster bool     `yaml:"enableCluster"`
	}

	RpcRegisterIP string `yaml:"rpcRegisterIP"`
	ListenIP      string `yaml:"listenIP"`

	WalletApi struct {
		GinPort  []int  `yaml:"walletApiPort"`
		ListenIP string `yaml:"listenIP"`
	} `yaml:"wallet_api"`

	AdminApi struct {
		GinPort  []int  `yaml:"adminApiPort"`
		ListenIP string `yaml:"listenIP"`
	} `yaml:"admin_api"`

	RpcPort struct {
		WalletPort  []int `yaml:"walletRPCPort"`
		AdminPort   []int `yaml:"adminRPCPort"`
		BitcoinPort []int `yaml:"btcRPCPort"`
		EthPort     []int `yaml:"ethRPCPort"`
		TronPort    []int `yaml:"tronRPCPort"`
		PushPort    []int `yaml:"pushRPCPort"`
	} `yaml:"rpc_port"`

	RpcRegisterName struct {
		WalletRPC  string `yaml:"walletRPCName"`
		AdminRPC   string `yaml:"adminRPCName"`
		BitcoinRPC string `yaml:"btcRPCName"`
		EthRPC     string `yaml:"ethRPCName"`
		TronRPC    string `yaml:"tronRPCName"`
		PushRPC    string `yaml:"pushRPCName"`
	} `yaml:"rpc_register_name"`

	Log struct {
		StorageLocation       string   `yaml:"storageLocation"`
		RotationTime          int      `yaml:"rotationTime"`
		RemainRotationCount   uint     `yaml:"remainRotationCount"`
		RemainLogLevel        uint     `yaml:"remainLogLevel"`
		GormLogLevel          uint     `yaml:"gormLogLevel"`
		ElasticSearchSwitch   bool     `yaml:"elasticSearchSwitch"`
		ElasticSearchAddr     []string `yaml:"elasticSearchAddr"`
		ElasticSearchUser     string   `yaml:"elasticSearchUser"`
		ElasticSearchPassword string   `yaml:"elasticSearchPassword"`
	}

	Push struct {
		Jpns struct {
			AppKey       string `yaml:"appKey"`
			MasterSecret string `yaml:"masterSecret"`
			PushUrl      string `yaml:"pushUrl"`
			PushIntent   string `yaml:"pushIntent"`
			IsProduct    bool   `yaml:"isProduct"`
			Enable       bool   `yaml:"enable"`
		}
		Getui struct {
			PushUrl      string `yaml:"pushUrl"`
			AppKey       string `yaml:"appKey"`
			Enable       bool   `yaml:"enable"`
			Intent       string `yaml:"intent"`
			MasterSecret string `yaml:"masterSecret"`
		}
	}

	TokenPolicy struct {
		AccessSecret      string `yaml:"accessSecret"`
		AccessSecretGAuth string `yaml:"accessSecretGAuth"`
		AccessExpire      int64  `yaml:"accessExpire"`
	}

	AdminUser2FAuthEnable bool   `yaml:"adminUser2FAuthEnable"`
	TotpIssuerName        string `yaml:"totpIssuerName"`
	SwaggerEnable         bool   `yaml:"swaggerEable"`
	MinimumBlockNumber    int64  `yaml:"minimumBlockNumber"`

	Manager struct {
		AppManagerUid []string `yaml:"appManagerUid"`
		Secrets       []string `yaml:"secrets"`
		Actions       []Action `yaml:"actions"`
		Currencies    []string `yaml:"currencies"`
	}
}
type Action struct {
	Pid  int    `yaml:"pid"`
	Name string `yaml:"name"`
}

func init() {
	cfgName := os.Getenv("CONFIG_NAME")
	if len(cfgName) == 0 {
		cfgName = Root + "/config/config.yaml"
	}

	bytes, err := ioutil.ReadFile(cfgName)
	if err == nil {
		if err = yaml.Unmarshal(bytes, &Config); err != nil {
			panic(err.Error())
		}
	}
	// } else {
	// 	panic(err.Error())
	// }
}
