package main

import (
	"context"
	"sap_segmentation/internal/client"
	"sap_segmentation/internal/config"
	"sap_segmentation/internal/db"
	"sap_segmentation/internal/importer"
	"sap_segmentation/internal/logger"
)

func main() {
	cfg := config.Parse()
	log := logger.New(cfg)
	client := client.New(cfg, log)
	db := db.New(cfg)
	imp := importer.New(cfg, log, client, db)

	ctx := context.Background()
	err := imp.Run(ctx)
	if err != nil {
		log.Error("importer finished with error: %s", err)
	}
	log.Info("success")
}
