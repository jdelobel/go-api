package config

import (
	"testing"
)

func TestConfigLoad(t *testing.T) {

	configError := &Config{}
	configSuccess := &Config{}
	configSuccess.Database.Name = "my-db"
	configSuccess.Database.Host = "my-db"
	configSuccess.Database.User = "my-db"
	configSuccess.Database.Password = "my-db"
	tests := []struct {
		name string
		c    *Config
	}{
		{name: "Load error Configuration", c: configError},
		{name: "Load success Configuration", c: configSuccess},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.Load("config.json")
			if err != nil && tt.name != "Load error Configuration" {
				t.Fail()
			}
		})
	}

}

func BenchmarkConfigLoad(b *testing.B) {
	config := Config{}
	config.Database.Name = "my-db"
	config.Database.Name = "my-db"
	config.Database.Host = "my-db"
	config.Database.User = "my-db"
	config.Database.Password = "my-db"
	for n := 0; n < b.N; n++ {
		err := config.Load("config.json")
		if err != nil {
			b.Fail()
		}
	}
}
