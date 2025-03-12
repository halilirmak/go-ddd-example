package seed

import (
	"log"

	"gorm.io/gorm"
)

type Seed struct {
	Name string
	Run  func(*gorm.DB) error
}

type Seeds []Seed

func (s *Seeds) Add(seed Seed) {
	*s = append(*s, seed)
}

func (s *Seeds) RunAll(db *gorm.DB) {
	for _, seed := range *s {
		if err := seed.Run(db); err != nil {
			log.Fatalf("running seed '%s', failed with error: %s", seed.Name, err)
		}
	}
}
