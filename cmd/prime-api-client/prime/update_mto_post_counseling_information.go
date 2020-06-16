package prime

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/cmd/prime-api-client/utils"
	mto "github.com/transcom/mymove/pkg/gen/primeclient/move_task_order"
)

// InitUpdatePostCounselingInfoFlags declares which flags are enabled
func InitUpdatePostCounselingInfoFlags(flag *pflag.FlagSet) {
	flag.String(utils.FilenameFlag, "", "Name of the file being passed in")

	flag.SortFlags = false
}

func checkUpdatePostCounselingInfoConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := utils.CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	if v.GetString(utils.FilenameFlag) == "" && (len(args) < 1 || len(args) > 0 && !utils.ContainsDash(args)) {
		logger.Fatal(errors.New("update-post-counseling-info expects a file to be passed in"))
	}

	return nil
}

// UpdatePostCounselingInfo creates a gateway and sends the request to the endpoint
func UpdatePostCounselingInfo(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//  Create the logger
	//  Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := utils.ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkUpdatePostCounselingInfoConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err)
	}

	primeGateway, cacStore, errCreateClient := utils.CreatePrimeClient(v)
	if errCreateClient != nil {
		return errCreateClient
	}

	// Defer closing the store until after the API call has completed
	if cacStore != nil {
		defer cacStore.Close()
	}

	// Decode json from file that was passed into MTOShipment
	filename := v.GetString(utils.FilenameFlag)
	var reader *bufio.Reader
	if filename != "" {
		file, fileErr := os.Open(filepath.Clean(filename))
		if fileErr != nil {
			logger.Fatal(fileErr)
		}
		reader = bufio.NewReader(file)
	}

	if len(args) > 0 && utils.ContainsDash(args) {
		reader = bufio.NewReader(os.Stdin)
	}

	jsonDecoder := json.NewDecoder(reader)
	var postCounselingInfo mto.UpdateMTOPostCounselingInformationBody
	err = jsonDecoder.Decode(&postCounselingInfo)
	if err != nil {
		return fmt.Errorf("decoding data failed: %w", err)
	}

	params := mto.UpdateMTOPostCounselingInformationParams{
		MoveTaskOrderID: postCounselingInfo.MoveTaskOrderID,
		Body:            postCounselingInfo,
		IfMatch:         v.GetString(utils.ETagFlag),
	}
	params.SetTimeout(time.Second * 30)

	resp, errUpdatePostCounselingInfo := primeGateway.MoveTaskOrder.UpdateMTOPostCounselingInformation(&params)
	if errUpdatePostCounselingInfo != nil {
		// If the response cannot be parsed as JSON you may see an error like
		// is not supported by the TextConsumer, can be resolved by supporting TextUnmarshaler interface
		// Likely this is because the API doesn't return JSON response for BadRequest OR
		// The response type is not being set to text
		logger.Fatal(errUpdatePostCounselingInfo.Error())
	}

	payload := resp.GetPayload()
	if payload != nil {
		payload, errJSONMarshall := json.Marshal(payload)
		if errJSONMarshall != nil {
			logger.Fatal(errJSONMarshall)
		}
		fmt.Println(string(payload))
	} else {
		logger.Fatal(resp.Error())
	}

	return nil
}
