package viper

import "github.com/spf13/viper"

type Loader struct {
	Data       *viper.Viper
	configPath string
	configName string
	configType string
}

func NewLoader(path, name, cfgType string) *Loader {
	subViper := viper.New()
	return &Loader{
		Data:       subViper,
		configPath: path,
		configName: name,
		configType: cfgType,
	}
}

func (c *Loader) Get(key string) any {
	return c.Data.Get(key)
}

func (c *Loader) SetDefault(key string, val any) {
	c.Data.SetDefault(key, val)
}

func (c *Loader) Read() error {
	c.Data.AddConfigPath(c.configPath)
	c.Data.SetConfigName(c.configName)
	c.Data.SetConfigType(c.configType)
	if err := c.Data.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
