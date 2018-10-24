package document

import (
	"github.com/transcom/mymove/pkg/di"
)

// AddProviders adds all the DI providers from the services package
func AddProviders(c *di.Container) {
	c.MustProvide(NewFetchDocumentService)
	c.MustProvide(NewFetchUploadService)
}
