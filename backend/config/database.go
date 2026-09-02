package config

import (
	"log"
	"server/prisma/db"
)

func ConnectDB() (*db.PrismaClient, error) {
	client := db.NewClient()

	if err := client.Prisma.Connect(); err != nil {
		return nil, err
	}

	log.Println("Connected to database")
	return client, nil
}
