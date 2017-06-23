package main

import (
	"fmt"
	"net/http"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func writeOK(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := gorm.Open("mysql", "root:98526@/echeng?charset=utf8")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	// db.AutoMigrate(&Shift{})

	// Create
	// db.Create(&Shift{Code: "L1212", Price: 1000})

	// Read
	var shift EcBusLineShift
	// db.First(&shift, 1) // find shift with id 1
	db.First(&shift, "bus_line_id = ?", "2") // find shift with code l1212

	// Update - update shift's price to 2000
	// db.Model(&shift).Update("default_price", 2000)

	// Delete - delete shift
	// db.Delete(&shift)

	price := shift.GetPrice()
	priceStr := fmt.Sprintf("%.2f", price)
	w.Write([]byte(priceStr))
	w.WriteHeader(200)
}

type EcBusLineShift struct {
	gorm.Model
	DepartureAt string
	ArrivedAt string
	DefaultPrice int
}

func (shift EcBusLineShift) GetPrice() float32 {
	return float32(shift.DefaultPrice) / 100
}

