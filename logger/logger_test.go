package logger

import (
	"testing"
)

var conf = Conf{Level: "DEBUG", App: "boxid-segments-importer-test", Version: "1.0.0"}

func BenchmarkStdOut(b *testing.B) {
	var err error
	for n := 0; n < b.N; n++ {
		err = Init(conf)
		if err != nil {
			b.Fatalf("Failed to init config %s", err)
		}
	}
}

func Test(t *testing.T) {
	tests := []struct {
		name   string
		logger func(t *testing.T)
	}{
		{name: "test conf stdout logger", logger: testConfStdOut},
		{name: "test conf graylog logger", logger: testConfGraylog},
		{name: "test conf graylog logger error", logger: testGraylogConfError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logger(t)
		})
	}
}

func testConfStdOut(t *testing.T) {
	err := Init(conf)
	if err != nil {
		t.Error(err)
	}
	if Log == nil {
		t.Error(err)
	}
	Log = nil
}

func testConfGraylog(t *testing.T) {
	conf.Host = "127.0.0.1"
	conf.Port = "25215"
	err := Init(conf)
	if err != nil {
		t.Errorf("Failed to init logger %s", err)
	}
	if Log == nil {
		t.Fail()
	}
	conf.Host = ""
	conf.Port = ""
	Log = nil
}

func testGraylogConfError(t *testing.T) {
	conf.Host = "fakehost"
	conf.Port = "25215"

	err := Init(conf)
	if err == nil {
		t.Errorf("Init logger should not be succeded %s", err)

	}
	if Log != nil {
		t.Fail()
	}
	conf.Host = ""
	conf.Port = ""
	Log = nil
}
