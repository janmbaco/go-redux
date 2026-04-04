module github.com/janmbaco/go-redux/v2/examples/event-sourced-wallet

go 1.24.13

require (
	github.com/google/uuid v1.6.0
	github.com/janmbaco/go-infrastructure/v2 v2.1.4
	github.com/janmbaco/go-redux/v2 v2.0.0
	gorm.io/gorm v1.31.1
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.33.0 // indirect
)

replace github.com/janmbaco/go-redux/v2 => ../..
