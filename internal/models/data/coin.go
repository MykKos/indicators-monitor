package data

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	TokenPrice struct {
		Hash       string  `gorm:"primaryKey;<-"`
		Token      string  `gorm:"<-"`
		OpenPrice  float64 `gorm:"<-"`
		ClosePrice float64 `gorm:"<-"`
		MinPrice   float64 `gorm:"<-"`
		MaxPrice   float64 `gorm:"<-"`

		TF        string  `gorm:"<-"`
		CEX       string  `gorm:"<-"`
		Timestamp float64 `gorm:"<-"`
	}
)

func FindPrices(tp TokenPrice, db *gorm.DB) []TokenPrice {
	var prices []TokenPrice
	db.Where(&tp).Find(&prices)
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Timestamp > prices[j].Timestamp
	})

	return prices
}

func (tp *TokenPrice) Save(db *gorm.DB) {
	tp.GenHash()
	db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(tp)
}

func (tp *TokenPrice) GenHash() {
	hashLine := fmt.Sprintf("%s.%s.%s.%f", tp.CEX, tp.TF, tp.Token, tp.Timestamp)
	hash := md5.Sum([]byte(hashLine))
	tp.Hash = hex.EncodeToString(hash[:])
}
