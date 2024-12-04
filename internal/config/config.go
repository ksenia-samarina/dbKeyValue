package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	GRPC        GRPCConfig        `yaml:"grpc"`
	MemTable    MemTableConfig    `yaml:"mem_table_config"`
	BloomFilter BloomFilterConfig `yaml:"bloom_filter_config"`
	StoragePath string            `yaml:"storagePath"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type MemTableConfig struct {
	BtreeSize int  `yaml:"btreeSize"`
	MaxSize   uint `yaml:"maxMemTableSize"`
}

type BloomFilterConfig struct {
	BloomFilterN  uint    `yaml:"bloomFilterN"`  // n items will be stored in bloom filter
	BloomFilterFp float64 `yaml:"bloomFilterFp"` // false positive estimate
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty: " + err.Error())
	}
	return &cfg
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
