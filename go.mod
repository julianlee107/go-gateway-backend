module github.com/julianlee107/gatewayScaffold

go 1.13

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/gin-gonic/gin v1.6.3
	github.com/go-ini/ini v1.57.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.2.0
	github.com/jinzhu/gorm v1.9.16
	github.com/julianlee107/go-common v0.0.0-20200812130750-93fe7e0701cc
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.7
	gopkg.in/ini.v1 v1.57.0 // indirect
)

replace github.com/julianlee107/go-common => ../go-common
