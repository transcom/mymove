# Route Planner Guide

## Table of Contents

<!-- Table of Contents auto-generated with `bin/generate-md-toc.sh` -->

<!-- toc -->

* [Route planner interface](#route-planner-interface)
  * [Interface methods](#interface-methods)
* [Available route planners](#available-route-planners)
  * [HERE API](#here-api)
    * [TransitDistance](#transitdistance)
    * [Zip5TransitDistance](#zip5transitdistance)
    * [Zip3TransitDistance](#zip3transitdistance)
    * [LatLongTransitDistance](#latlongtransitdistance)
  * [Bing API](#bing-api)
  * [Test](#test)
* [Zip5ToLatLong](#zip5tolatlong)
* [Chamber](#chamber)
* [Testing](#testing)

Regenerate with "pre-commit run -a markdown-toc"

<!-- tocstop -->

## Route planner interface

In the MyMove project we have an interface for communicating with various Route planning APIs. Currently the project uses the HERE api for distance route calculations. The hope of this doc is to explain some of how it works. The interface and implementations can be found in `pkg/route`.

### Interface methods

The Route Planner interface requires the implementation of the following methods.

```go
// This method takes a source and destination `models.Address` and returns the distance as an `int`.
TransitDistance(source *models.Address, destination *models.Address) (int, error)

//This method takes a source and destination `LatLong` and returns the distance as an `int`
LatLongTransitDistance(source LatLong, destination LatLong) (int, error)

// This method takes a source and destination `Zip5` string and returns the distance as an `int`
Zip5TransitDistance(source string, destination string) (int, error)

// This method takes a source and destination `Zip3` string and returns the distance as an `int`
Zip3TransitDistance(source string, destination string) (int, error)
```

## Available route planners

Below is a list of the current route planner APIs that are implemented, though it appears as of this writing that only HERE API is used.

### HERE API

The HERE API is the currently used 3rd party api for doing route planning in the my move system.

#### TransitDistance

Turns the input addresses into `LatLong` data via HERE geocoder endpoint, and then it runs the `LatLongTransitDistance` method to determine the distance.

#### Zip5TransitDistance

Uses the `Zip5ToLatLong` method to turn the input Zip5's into `LatLong` data and then uses the `LatLongTransitDistance` method to determine the distance.

#### Zip3TransitDistance

`Zip3TransitDistance` is unimplemented and always returns a `NewUnsupportedPostalCodeError`

#### LatLongTransitDistance

This method takes a source and destination `LatLong` and uses the HERE API to calculate the distance between them.

### Bing API

This is an implementation of the Planner interface that uses Bing as a backend. However, as of this writing there are no uses of this class other than in tests.

### Test

This is an implementation of the Planner interface to be used in various tests.

## Zip5ToLatLong

The HERE API relies on zip to LatLong data from the [free zip code data project](https://github.com/midwire/free_zipcode_data). This data was used to create a static map of Zip5's to LatLong tuples in
`pkg/route/zip_locale.go`. This file also contains the method which can be used to lookup the data.

## Chamber

The following environment variables need to be set to make successful calls to the API. These values are stored in chamber so they are not be checked in.

```sh
HERE_MAPS_APP_ID
HERE_MAPS_ROUTING_ENDPOINT
HERE_MAPS_APP_CODE
HERE_MAPS_GEOCODE_ENDPOINT
```

## Testing

There are unit tests under `<planner_name>_test.go` which do some basic validation. There is an integration test which calls out to the real api in `planner_test.go`
