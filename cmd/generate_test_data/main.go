package main

import (
	"log"

	"github.com/gobuffalo/pop"
	"github.com/namsral/flag"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"

	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/testdatagen"
	tdgs "github.com/transcom/mymove/pkg/testdatagen/scenario"
	"github.com/transcom/mymove/pkg/uploader"
)

// Hey, refactoring self: you can pull the UUIDs from the objects rather than
// querying the db for them again.
func main() {
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, which configures the database.")
	scenario := flag.Int("scenario", 0, "Specify which scenario you'd like to run. Current options: 1, 2, 3, 4, 5, 6, 7.")
	namedScenario := flag.String("named-scenario", "", "It's like a scenario, but more descriptive.")
	flag.Parse()

	//DB connection
	err := pop.AddLookupPaths(*config)
	if err != nil {
		log.Panic(err)
	}
	db, err := pop.Connect(*env)
	if err != nil {
		log.Panic(err)
	}

	if *scenario == 1 {
		tdgs.RunAwardQueueScenario1(db)
	} else if *scenario == 2 {
		tdgs.RunAwardQueueScenario2(db)
	} else if *scenario == 4 {
		err = tdgs.RunPPMSITEstimateScenario1(db)
	} else if *scenario == 5 {
		err = tdgs.RunRateEngineScenario1(db)
	} else if *scenario == 6 {
		query := `DELETE FROM transportation_service_provider_performances;
				  DELETE FROM transportation_service_providers;
				  DELETE FROM traffic_distribution_lists;
				  DELETE FROM tariff400ng_zip3s;
				  DELETE FROM tariff400ng_zip5_rate_areas;
				  DELETE FROM tariff400ng_service_areas;
				  DELETE FROM tariff400ng_linehaul_rates;
				  DELETE FROM tariff400ng_shorthaul_rates;
				  DELETE FROM tariff400ng_full_pack_rates;
				  DELETE FROM tariff400ng_full_unpack_rates;`

		err = db.RawQuery(query).Exec()
		if err != nil {
			log.Panic(err)
		}
		err = tdgs.RunRateEngineScenario2(db)
	} else if *scenario == 7 {
		// Create TSPs with shipments divided among them
		numTspUsers := 2
		numShipments := 25
		numShipmentOfferSplit := []int{15, 10}
		// TSPs should never be able to see DRAFT or SUBMITTED or AWARDING shipments.
		status := []models.ShipmentStatus{"AWARDED", "ACCEPTED", "APPROVED", "IN_TRANSIT", "DELIVERED", "COMPLETED"}
		_, _, _, err := testdatagen.CreateShipmentOfferData(db, numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
		if err != nil {
			log.Panic(err)
		}
		// Create an office user
		testdatagen.MakeDefaultOfficeUser(db)
		log.Print("Success! Created TSP test data.")
	} else if *namedScenario == tdgs.E2eBasicScenario.Name {
		// Initialize logger
		logger, err := zap.NewDevelopment()
		if err != nil {
			log.Panic(err)
		}

		// Initialize storage and uploader
		zap.L().Info("Using memory storage backend")
		fsParams := storage.NewMemoryParams("tmp", "testdata", logger)
		storer := storage.NewMemory(fsParams)
		loader := uploader.NewUploader(db, logger, storer)

		tdgs.E2eBasicScenario.Run(db, loader, logger, storer)
		log.Print("Success! Created e2e test data.")
	} else {
		flag.PrintDefaults()
	}
	if err != nil {
		log.Panic(err)
	}
}
