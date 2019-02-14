# Architectural Decision Log

This log lists the architectural decisions for DP3 Infrastructure.

<!--Regenerate the content by using "npx adr-log -i". You can install it via "yarn". -->

<!-- adrlog -->

- [ADR-0000](0000-server-framework.md) - Use Truss' [Golang](https://golang.org/) web server skeleton to build API for DP3
- [ADR-0001](0001-go-orm.md) - Use [Pop](https://github.com/gobuffalo/pop) as the ORM for 3M
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
- [ADR-0014](0014-go-dependency-management.md) - Go Dependency Management
- [ADR-0015](0015-session-storage.md) - Session storage/handling
- [ADR-0016](0016-Browser-Support.md) - Browser Support for Prototype
- [ADR-0017](0017-react-router-redux-authentication.md) - Client side route restriction based on authentication
- [ADR-0018](0018-optional-field-interop.md) - Optional Field Interop
- [ADR-0019](0019-client-rangeslider.md) - _Range Slider React Component_
- [ADR-0020](0020-swagger-auth.md) - Using Swagger to manage server route authentication
- [ADR-0021](0021-ssn-use.md) - Temporary use and plan for expunging Social Security Numbers in the prototype
- [ADR-0022](0022-xlsx-lib.md) - Chose Excelize package to parse XLSX files
- [ADR-0023](0023-representing-dollar-values.md) - Representing Dollar Values in Go and the Database
- [ADR-0024](0024-model-authorization-and-handler-design.md) - Model Authorization and Handler Design
- [ADR-0025](0025-client-side-feature-flags.md) - Client Side Feature Flags using Custom JavaScript
- [ADR-0026](0026-use-snyk-vulnerability-scanning.md) - Use Snyk Vulnerability Scanning
- [ADR-0027](0027-pdf-generation.md) - PDF Generation
- [ADR-0028](0028-config-management.md) - Config Management
- [ADR-0029](0029-honeycomb-integration.md) - Honeycomb Integration
- [ADR-0030](0030-rds-iam.md) - IAM Authentication for Database
- [ADR-0031](0031-css-tooling.md) - CSS Tooling
- [ADR-0032](0032-csrf-protection.md) - CSRF Protection for the Application
- [ADR-0033](0033-service-object-layer.md) - Mocking Approach and Tooling

<!-- adrlogstop -->

For new ADRs, please use [template.md](template.md).

More information on MADR is available at <https://adr.github.io/madr/>.
General information about architectural decision records is available at <https://adr.github.io/>.
