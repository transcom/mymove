package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/db/sequence"
	"github.com/transcom/mymove/pkg/edi"
	"github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/server"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/invoice"
)

// Call this from command line with go run cmd/generate_shipment_edi/main.go --shipmentID <UUID> --approver <email>
// Must use a shipment that is delivered, but not yet approved for payment (that does not already have a submitted invoice)
func main() {
	flag := pflag.CommandLine
	flag.String("shipmentID", "", "The ID of the shipment to invoice")
	flag.String("approver", "", "The office approver e-mail")
	flag.Bool("gex", false, "Choose to send the file to gex")

	// EDI Invoice Config
	flag.String("gex-basic-auth-username", "", "GEX api auth username")
	flag.String("gex-basic-auth-password", "", "GEX api auth password")
	flag.String("gex-url", "", "URL for sending an HTTP POST request to GEX")

	flag.String("dod-ca-package", "", "Path to PKCS#7 package containing certificates of all DoD root and intermediate CAs")
	flag.String("move-mil-dod-ca-cert", "", "The DoD CA certificate used to sign the move.mil TLS certificate.")
	flag.String("move-mil-dod-tls-cert", "", "The DoD-signed TLS certificate for various move.mil services.")
	flag.String("move-mil-dod-tls-key", "", "The private key for the DoD-signed TLS certificate for various move.mil services.")

	flag.String("edi", "", "The filepath to an edi file to send to GEX")
	flag.String("transaction-name", "test", "The required name sent in the url of the gex api request")
	flag.Parse(os.Args[1:])

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	logger, err := logging.Config("development", true)
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	shipmentIDString := v.GetString("shipmentID")
	approverEmail := v.GetString("approver")
	sendToGex := v.GetBool("gex")
	transactionName := v.GetString("transaction-name")
	if shipmentIDString == "" || approverEmail == "" {
		log.Fatal("Usage: go run cmd/generate_shipment_edi/main.go --shipmentID <29cb984e-c70d-46f0-926d-cd89e07a6ec3> --approver <officeuser1@example.com> --gex false")
	}

	db, err := pop.Connect("development")
	if err != nil {
		log.Fatal(err)
	}

	shipmentID := uuid.Must(uuid.FromString(shipmentIDString))
	shipment, err := invoice.FetchShipmentForInvoice{DB: db}.Call(shipmentID)
	if err != nil {
		log.Fatal(err)
	}

	approver, err := models.FetchOfficeUserByEmail(db, approverEmail)
	if err != nil {
		log.Fatalf("Could not fetch office user with e-mail %s: %v", approverEmail, err)
	}

	var invoiceModel models.Invoice
	verrs, err := invoice.CreateInvoice{DB: db, Clock: clock.New()}.Call(*approver, &invoiceModel, shipment)
	if err != nil {
		log.Fatal(err)
	}
	if verrs.HasAny() {
		log.Fatal(verrs)
	}

	var sendToGexHTTP services.GexSender
	if sendToGex {
		certificates, rootCAs, err := initDODCertificates(v, logger)
		if certificates == nil || rootCAs == nil || err != nil {
			log.Fatal("Error in getting tls certs", err)
		}
		tlsConfig := &tls.Config{Certificates: certificates, RootCAs: rootCAs}
		url := v.GetString("gex-url")
		if len(url) == 0 {
			log.Fatal("Not sending to GEX because no URL set. Set GEX_URL in your envrc.local.")
		}
		sendToGexHTTP = invoice.NewGexSenderHTTP(
			url,
			true,
			tlsConfig,
			v.GetString("gex-basic-auth-username"),
			v.GetString("gex-basic-auth-password"),
		)
	}

	resp, err := processInvoice(db, shipment, invoiceModel, sendToGex, &transactionName, sendToGexHTTP)
	if resp != nil {
		fmt.Printf("status code: %v\n", resp.StatusCode)
	}
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func processInvoice(db *pop.Connection, shipment models.Shipment, invoiceModel models.Invoice, sendToGex bool, transactionName *string, gexSender services.GexSender) (resp *http.Response, err error) {
	defer func() {
		if err != nil || (resp != nil && resp.StatusCode != 200) {
			// Update invoice record as failed
			invoiceModel.Status = models.InvoiceStatusSUBMISSIONFAILURE
			verrs, deferErr := db.ValidateAndSave(&invoiceModel)
			if deferErr != nil {
				log.Fatal(deferErr)
			}
			if verrs.HasAny() {
				log.Fatal(verrs)
			}
		} else {
			// Update invoice record as submitted
			shipmentLineItems := shipment.ShipmentLineItems
			verrs, deferErr := invoice.UpdateInvoiceSubmitted{DB: db}.Call(&invoiceModel, shipmentLineItems)
			if deferErr != nil {
				log.Fatal(deferErr)
			}
			if verrs.HasAny() {
				log.Fatal(verrs)
			}
		}
	}()

	var icnSequencer sequence.Sequencer
	if sendToGex {
		icnSequencer, err = sequence.NewRandomSequencer(ediinvoice.ICNRandomMin, ediinvoice.ICNRandomMax)
		if err != nil {
			log.Fatal("Could not create random sequencer for ICN", err)
		}
	} else {
		icnSequencer = sequence.NewDatabaseSequencer(db, ediinvoice.ICNSequenceName)
	}

	invoice858C, err := ediinvoice.Generate858C(shipment, invoiceModel, db, false, icnSequencer, clock.New())
	if err != nil {
		return nil, err
	}

	if sendToGex {
		fmt.Println("Sending to GEX. . .")
		invoice858CString, err := invoice858C.EDIString()
		if err != nil {
			return nil, err
		}
		resp, err := gexSender.SendToGex(invoice858CString, *transactionName)
		if resp == nil || err != nil {
			log.Fatal("Gex Sender had no response", err)
		}

		fmt.Printf("status code: %v, error: %v\n", resp.StatusCode, err)
	}
	ediWriter := edi.NewWriter(os.Stdout)
	err = ediWriter.WriteAll(invoice858C.Segments())
	return nil, err
}

//TODO: Infra will work to refactor and reduce duplication (also found in webserver/main.go)
func initDODCertificates(v *viper.Viper, logger *zap.Logger) ([]tls.Certificate, *x509.CertPool, error) {

	tlsCert := v.GetString("move-mil-dod-tls-cert")
	if len(tlsCert) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s is missing", "move-mil-dod-tls-cert")
	}

	caCert := v.GetString("move-mil-dod-ca-cert")
	if len(caCert) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s is missing", "move-mil-dod-ca-cert")
	}

	//Append move.mil cert with CA certificate chain
	cert := bytes.Join(
		[][]byte{
			[]byte(tlsCert),
			[]byte(caCert),
		},
		[]byte("\n"),
	)

	key := []byte(v.GetString("move-mil-dod-tls-key"))
	if len(key) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s is missing", "move-mil-dod-tls-key")
	}

	keyPair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return make([]tls.Certificate, 0), nil, errors.Wrap(err, "failed to parse DOD keypair for server")
	}

	pathToPackage := v.GetString("dod-ca-package")
	if len(pathToPackage) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Wrap(&errInvalidPKCS7{Path: pathToPackage}, fmt.Sprintf("%s is missing", "dod-ca-package"))
	}

	pkcs7Package, err := ioutil.ReadFile(pathToPackage) // #nosec
	if err != nil {
		return make([]tls.Certificate, 0), nil, errors.Wrap(err, fmt.Sprintf("%s is invalid", "dod-ca-package"))
	}

	if len(pkcs7Package) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Wrap(&errInvalidPKCS7{Path: pathToPackage}, fmt.Sprintf("%s is an empty file", "dod-ca-package"))
	}

	dodCACertPool, err := server.LoadCertPoolFromPkcs7Package(pkcs7Package)
	if err != nil {
		return make([]tls.Certificate, 0), dodCACertPool, errors.Wrap(err, "Failed to parse DoD CA certificate package")
	}

	return []tls.Certificate{keyPair}, dodCACertPool, nil

}

//TODO: Infra will refactor to reduce duplication
type errInvalidPKCS7 struct {
	Path string
}

//TODO: Infra will refactor to reduce duplication
func (e *errInvalidPKCS7) Error() string {
	return fmt.Sprintf("invalid DER encoded PKCS7 package: %s", e.Path)
}
