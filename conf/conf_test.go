package conf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	conf, err := New("config.example.yaml")

	c := &Conf{
		App: App{
			Name:    "alteisen",
			Version: "v0.1.0",
		},
		HttpServ: HttpServ{
			Addr:         "6666",
			Mode:         "debug",
			ReadTimeout:  "10s",
			WriteTimeout: "10s",
		},
		Bot: Bot{
			Token: "1111111111",
		},
	}

	assert.Nil(t, err)
	assert.Equal(t, c, conf)

}
