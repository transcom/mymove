# Architectural Decision Log

This log lists the architectural decisions for DP3 Infrastructure.

<!--Regenerate the content by using "npx adr-log -i". You can install it via "yarn". -->

<!-- adrlog -->

- [ADR-0000](0000-server-framework.md) - Use Truss' [Golang](https://golang.org/) web server skeleton to build API for DP3
- [ADR-0001](0001-go-orm.md) - Use [Pop](https://github.com/markbates/pop) as the ORM for 3M
- [ADR-0002](0002-go-package-management.md) - Use dep to manage go dependencies
- [ADR-0003](0003-go-path-and-project-layout.md) - Put mymove into the standard GOPATH, eliminate server and client directories
- [ADR-0004](0004-path-imports.md) - Use Both Absolute and Relative Paths for Imports
- [ADR-0005](0005-create-react-app.md) - Use [Create React App](https://github.com/facebook/create-react-app)
- [ADR-0006](0006-redux.md) - Use [Redux](https://redux.js.org) to manage state and [Redux Thunk](https://github.com/gaearon/redux-thunk) middleware to write action creators that return functions
- [ADR-0007](0007-swagger-client.md) - Use swagger-client to make calls to API from client
- [ADR-0008](0008-go-swagger.md) - Use go-swagger To Route, Parse, And Validate API Endpoints
- [ADR-0009](0009-form-creation-from-swagger.md) - Generate forms from swagger definitions of payload
- [ADR-0010](0010-isolate-test-access-to-database.md) - Isolate Test Access to Database
- [ADR-0011](0011-test-suites.md) - Test Suites
- [ADR-0012](0012-tsp-data-models.md) - The TSP Data Models
- [ADR-0013](0013-rest-api-updates.md) - REST API Updates
- [ADR-0015](0015-session-storage.md) - Session storage/handling
- [ADR-0016](0016-Browser-Support.md) - Browser Support for Prototype
- [ADR-0017](0017-react-router-redux-authentication.md) - Client side route restriction based on authentication
- [ADR-0018](0018-optional-field-interop.md) - Optional Field Interop

<!-- adrlogstop -->

For new ADRs, please use [template.md](template.md).

More information on MADR is available at <https://adr.github.io/madr/>.
General information about architectural decision records is available at <https://adr.github.io/>.
