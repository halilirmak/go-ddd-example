package seed

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
	postgres_dto "github.com/cryptoPickle/go-ddd-example/internal/infrastructure/db/postress/dto"
)

func ItemSeeds(numberOfSeeds int) Seeds {
	var (
		refs       []string
		seeds      Seeds
		currencies = []string{"EUR", "GBP", "USD"}
	)

	for range 10 {
		refs = append(refs, faker.FirstName())
	}

	for i := range numberOfSeeds {
		productName := fmt.Sprintf("%s-%s", faker.FirstName(), faker.FirstName())
		s := Seed{
			Name: fmt.Sprintf("Create item seed for %s", productName),
			Run: func(db *gorm.DB) error {
				return CreateItem(db, productName, currencies[i%3], refs[i%10], randomMoney())
			},
		}
		seeds.Add(s)
	}
	return seeds
}

func CreateItem(db *gorm.DB, name, currency, sellerRef string, price float32) error {
	item := postgres_dto.Item{
		ID:              uuid.New(),
		ItemName:        name,
		SellerReference: sellerRef,
		Price:           price,
		Currency:        currency,
	}
	fmt.Printf("%+v\n", item)
	return db.Create(&item).Error
}

func randomMoney() float32 {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	price := float32(r.Intn(10000))/100 + r.Float32()
	return float32(int(price*100)) / 100
}
