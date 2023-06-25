module github.com/lazybark/go-cloud-sync

go 1.19

replace github.com/lazybark/go-tls-server => ../go-tls-server

require (
	github.com/alexflint/go-arg v1.4.3
	github.com/fsnotify/fsnotify v1.6.0
	github.com/google/go-cmp v0.5.9
	github.com/lazybark/go-helpers v1.3.0
	github.com/lazybark/go-tls-server v1.0.3
	github.com/stretchr/testify v1.7.1
	gorm.io/driver/sqlite v1.5.2
	gorm.io/gorm v1.25.2-0.20230530020048-26663ab9bf55
)

require (
	github.com/alexflint/go-scalar v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gofrs/uuid v4.2.0+incompatible // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20220908164124-27713097b956 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)
