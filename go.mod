module github.com/idelchi/envsync

go 1.24.2

require (
	github.com/BurntSushi/toml v1.5.0
	github.com/goccy/go-yaml v1.17.1
	github.com/idelchi/godyl v0.0.12
	github.com/spf13/cobra v1.9.1
)

replace github.com/idelchi/godyl => ../../godyl/dev

require (
	github.com/bmatcuk/doublestar/v4 v4.8.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
)
