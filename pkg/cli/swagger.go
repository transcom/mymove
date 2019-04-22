package cli

import "github.com/spf13/pflag"

const (
	// SwaggerFlag is the Public Swagger Flag
	SwaggerFlag string = "swagger"
	// InternalSwaggerFlag is the Internal Swagger Flag
	InternalSwaggerFlag string = "internal-swagger"
	// OrdersSwaggerFlag is the Orders Swagger Flag
	OrdersSwaggerFlag string = "orders-swagger"
	// DPSSwaggerFlag is the DPS Swagger Flag
	DPSSwaggerFlag string = "dps-swagger"
	// ServeSwaggerUIFlag is the Serve Swagger UI Flag
	ServeSwaggerUIFlag string = "serve-swagger-ui"
)

// InitSwaggerFlags initializes the Swagger command line flags
func InitSwaggerFlags(flag *pflag.FlagSet) {
	flag.String(SwaggerFlag, "swagger/api.yaml", "The location of the public API swagger definition")
	flag.String(InternalSwaggerFlag, "swagger/internal.yaml", "The location of the internal API swagger definition")
	flag.String(OrdersSwaggerFlag, "swagger/orders.yaml", "The location of the Orders API swagger definition")
	flag.String(DPSSwaggerFlag, "swagger/dps.yaml", "The location of the DPS API swagger definition")
	flag.Bool(ServeSwaggerUIFlag, false, "Whether to serve swagger UI for the APIs")
}
