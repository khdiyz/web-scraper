package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/khdiyz/web-scraper/scraper"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	start := time.Now()

	// Open SQLite database
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	// Create devices table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS devices (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT,
        star_amount TEXT,
        comment TEXT,
        old_price TEXT,
        price TEXT,
        price_per_month TEXT,
        image_url TEXT
    )`)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	// Channel for receiving scraped data
	ch := make(chan []scraper.Device)

	// Scrape data in a goroutine
	go func() {
		devices, err := scraper.ScrapeDevices()
		fmt.Println("lenth of array:", len(devices))
		if err != nil {
			log.Println("Error scraping devices:", err)
			return
		}
		ch <- devices
		close(ch)
	}()

	// Write scraped data to database
	stmt, err := db.Prepare(`INSERT INTO devices (title, star_amount, comment, old_price, price, price_per_month, image_url) VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Fatal("Error preparing statement:", err)
	}
	defer stmt.Close()

	for devices := range ch {
		for _, d := range devices {
			_, err := stmt.Exec(d.Title, d.StarAmount, d.Comment, d.OldPrice, d.Price, d.PricePerMonth, d.ImageUrl)
			if err != nil {
				log.Println("Error inserting data into database:", err)
			}
		}
	}

	// Print scraping duration
	color.Blue("Scraping duration: %s\n", time.Since(start))
}
