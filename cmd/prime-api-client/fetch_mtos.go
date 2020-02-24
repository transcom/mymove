package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
	mto "github.com/transcom/mymove/pkg/gen/primeclient/move_task_order"
)

func checkFetchMTOsConfig(v *viper.Viper, logger *log.Logger) error {
	err := cli.CheckCAC(v)
	if err != nil {
		return err
	}

	err = cli.CheckPrimeAPI(v)
	if err != nil {
		return err
	}

	err = cli.CheckVerbose(v)
	if err != nil {
		return err
	}

	return nil
}

func initFetchMTOsFlags(flag *pflag.FlagSet) {
	cli.InitCACFlags(flag)
	cli.InitPrimeAPIFlags(flag)
	cli.InitVerboseFlags(flag)

	flag.SortFlags = false
}

func fetchMTOs(cmd *cobra.Command, args []string) error {
	err := cmd.ParseFlags(args)
	if err != nil {
		return errors.Wrap(err, "Could not parse args")
	}
	flags := cmd.Flags()
	v := viper.New()
	err = v.BindPFlags(flags)
	if err != nil {
		return errors.Wrap(err, "Could not bind flags")
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	//Create the logger
	//Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	err = checkFetchMTOsConfig(v, logger)
	if err != nil {
		logger.Fatal(err)
	}

	primeGateway, err := CreateClient(v)
	if err != nil {
		logger.Fatal(err)
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