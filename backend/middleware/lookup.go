package middleware

import (
	"backend/config"
	"log"

	"github.com/gofiber/fiber/v2"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/ip2location/ip2location-go/v9"
)

const (
	cacheSize = 1024
)

var (
	db    *ip2location.DB
	cache *lru.Cache[string, *ip2location.IP2Locationrecord]
)

// init setup db file containing GEO-locations per IP.
func init() {
	var err error

	ips, ok := config.Get("IP_FILE")
	if !ok {
		log.Fatalln("IP_FILE env not found.")
	}

	db, err = ip2location.OpenDB(ips)
	if err != nil {
		log.Fatalln("Failed to open ip2location DB:", err)
	}

	if cache, err = lru.New[string, *ip2location.IP2Locationrecord](cacheSize); err != nil {
		log.Fatalln("Failed to create lru-cache for ip2location DB.", err)
	}
}

// Lookup handler fetches all geolocation fields based on the queried IP address.
func Lookup() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ip := ctx.IP()
		record, ok := cache.Get(ip)

		// Not found in cache
		if !ok {
			dbRecord, _ := db.Get_all(ip)
			record = &dbRecord
			cache.Add(ip, record)
		}

		ctx.Locals("lookup", record)
		return ctx.Next()
	}
}
