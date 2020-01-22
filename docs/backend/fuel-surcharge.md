# Route Planner Guide

## Table of Contents

<!-- Table of Contents auto-generated with `bin/generate-md-toc.sh` -->

<!-- toc -->

* [Fuel surcharge interface](#fuel-surcharge-interface)
  * [Interface methods](#interface-methods)
* [Available fuel surcharge APIs](#available-fuel-surcharge-apis)
  * [PriceDomesticFuelSurcharge](#pricedomesticfuelsurcharge)
* [Testing](#testing)

Regenerate with "pre-commit run -a markdown-toc"

<!-- tocstop -->

## Fuel surcharge interface

In the MyMove project we have an interface for communicating with fuel surcharge APIs. The interface and implementations can be found in `pkg/services/fuelprice`.

### Interface methods

The fuel surcharge interface requires the implementation of the following methods.

```go
//This method takes a route planner, weight of type `unit.Pound`, source and destination `Zip3` strings and returns the fuel surcharge as type `unit.Cents`
PriceDomesticFuelSurcharge(planner route.Planner, weight unit.Pound, source string, destination string) (unit.Cents, error)
```

## Available fuel surcharge APIs

Below is a list of the current fuel surcharge APIs. Note that none are currently implemented.

### PriceDomesticFuelSurcharge

`PriceDomesticFuelSurcharge` is unimplemented and always returns an error.

This function was added as a placeholder for calculating domestic fuel surcharge prices. It will be extended when we start the fuel surcharge epic.

## Testing

There are unit tests in `diesel_fuel_price_storer_test.go`.
