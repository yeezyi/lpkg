package config

type Cfg interface {
	Get(key string) interface{}
	SetDefault(key string, val interface{})
	Read() error
}

type config struct {
	sourceMap map[string]Cfg
}

var cfg = &config{sourceMap: make(map[string]Cfg)}

func Get(source, key string) interface{} {
	s, ok := cfg.sourceMap[source]
	if !ok {
		return nil
	}
	return s.Get(key)
}

func SetDefault(source, key string, val interface{}) {
	s, ok := cfg.sourceMap[source]
	if !ok {
		return
	}
	s.SetDefault(key, val)
}

func AddSource(kind string, obj Cfg) error {
	if err := obj.Read(); err != nil {
		return err
	}
	cfg.sourceMap[kind] = obj
	return nil
}
