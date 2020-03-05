package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	mto "github.com/transcom/mymove/pkg/gen/primeclient/move_task_order"
)

func fetchMTOs(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//Create the logger
	//Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	err := CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	primeGateway, cacStore, err := CreateClient(cmd, v, args)
	if err != nil {
		return err
	}

	if cacStore != nil {
		defer cacStore.Close()
	}

	var params mto.FetchMTOUpdatesParams
	params.SetTimeout(time.Second * 30)
	resp, errFetchMTOUpdates := primeGateway.MoveTaskOrder.FetchMTOUpdates(&params)
	if errFetchMTOUpdates != nil {
		// If the response cannot be parsed as JSON you may see an error like
		// is not supported by the TextConsumer, can be resolved by supporting TextUnmarshaler interface
		// Likely this is because the API doesn't return JSON response for BadRequest OR
		// The response type is not being set to text
		log.Fatal(errFetchMTOUpdates.Error())
	}

	payload := resp.GetPayload()
	if payload != nil {
		payload, errJSONMarshall := json.Marshal(payload)
		if errJSONMarshall != nil {
			log.Fatal(errJSONMarshall)
		}
		fmt.Println(string(payload))
	} else {
		log.Fatal(resp.Error())
	}

	return nil
}
