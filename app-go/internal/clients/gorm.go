package clients

import "gorm.io/gorm"

// IDb specify expectation for db provider. Right now i'm
// using gorm directly for the sake of simplicity
type IDb interface {
	GetDb() *gorm.DB
}
