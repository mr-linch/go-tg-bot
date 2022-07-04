package config

type HTTP struct {
	Listen string `default:":8080" usage:"http listen address"`
}
