package common

import (
	"github.com/spf13/viper"
)

// SubViper : From the singleton viper, creates another config
// object which holds only a part of of the config
func SubViper(args ...string) *viper.Viper {
	v := viper.New()
	for _, arg := range args {
		v.Set(arg, viper.Get(arg))
	}
	return v
}

// MockViper : Creates a dummy viper object to use
// in tests instead of the real one, the expected args
// are key1, val1, key2, val2, etc...
func MockViper(args ...interface{}) *viper.Viper {
	v := viper.New()
	key := ""
	isValue := false
	for _, arg := range args {
		if isValue {
			v.Set(key, arg)
		} else {
			key = arg.(string)
		}
		isValue = !isValue
	}
	return v
}
