package postgres

import (
	"log"

	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgres struct {
	dsn string
	db  *gorm.DB
}

func New(connectionString string) *postgres {
	return &postgres{
		dsn: connectionString,
	}
}

func (p *postgres) Start() {
	db, err := gorm.Open(gormpostgres.Open(p.dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("error on opening connection to postgres. err=%v", err)
	}

	p.db = db
}

func (p *postgres) Stop() {
	db, err := p.db.DB()
	if err != nil {
		log.Printf("error on getting gorm's db instance. err=%v", err)
	}

	err = db.Close()
	if err != nil {
		log.Printf("error on closing gorm's db instance. err=%v", err)
	}
}

// GetDb access db instance of postgres
func (p *postgres) GetDb() *gorm.DB {
	return p.db
}
