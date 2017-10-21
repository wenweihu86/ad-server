package adserver

import (
	"github.com/spf13/viper"
	"fmt"
)

type GlobalConf struct {
	GeoBlockFileName string
	GeoLocationFileName string
	AdFileName string
}

var GlobalConfObject *GlobalConf

func init()  {
	GlobalConfObject = new(GlobalConf)
}

func LoadGlobalConf(configPath, configFileName string)  {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configFileName)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
		panic(-1)
	}
	err := viper.Unmarshal(GlobalConfObject)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
		panic(-1)
	}
	fmt.Printf("GeoBlockFileName=%s GeoLocationFileName=%s AdFileName=%s\n",
		GlobalConfObject.GeoBlockFileName,
	    GlobalConfObject.GeoLocationFileName,
	    GlobalConfObject.AdFileName)
}
