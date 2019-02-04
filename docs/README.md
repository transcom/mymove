# Programming Guide

The intention for this document is to share our collective knowledge on best practices and allow everyone working on the DOD MyMove project to write code in compatible styles.

If you are looking to understand choices made in this project, see the list of [ADRs](https://github.com/transcom/mymove/tree/master/docs/adr).

## Programming Guides

* [Front-end / React](frontend.md) guide
* [Back-end / Go](backend.md) guide
* [Database / Postgres](database.md) guide

<!--index-->

## HOWTOs

* [Call Swagger Endpoints from React](how-to/access-swagger-endpoints-from-react.md#how-to-call-swagger-endpoints-from-react)
* [Backup and Restore the Development Database](how-to/backup-and-restore-dev-database.md#how-to-backup-and-restore-the-development-database)
* [display dates and times](how-to/display-dates-and-times.md#how-to-display-dates-and-times)
* [Generate Mocks with Mockery](how-to/generate-mocks-with-mockery.md#how-to-generate-mocks-with-mockery)
* [Instrument Data in Honeycomb](how-to/instrument-data-in-honeycomb.md#how-to-instrument-data-in-honeycomb)
* [revert a change](how-to/revert-a-change.md#how-to-revert-a-change)
* [Run Against S3 Locally](how-to/run-against-s3-locally.md#how-to-run-against-s3-locally)
* [Run End to End (Cypress) Tests](how-to/run-e2e-tests.md#how-to-run-end-to-end-cypress-tests)
* [Run Go Tests](how-to/run-go-tests.md#how-to-run-go-tests)
* [Run JavaScript (Jest) Tests](how-to/run-js-tests.md#how-to-run-javascript-jest-tests)
* [Store Data in Redux](how-to/store-data-in-redux.md#how-to-store-data-in-redux)
* [Store UI State in Redux](how-to/store-ui-state-in-redux.md#how-to-store-ui-state-in-redux)
* [Troubleshoot GEX Connection](how-to/troubleshoot-gex-connection.md#how-to-troubleshoot-gex-connection)
* [Unit Test React Components](how-to/unit-test-react-components.md#how-to-unit-test-react-components)
* [Upgrade Go Version](how-to/upgrade-go-version.md#how-to-upgrade-go-version)

## ADRs

* 0000 [Use Truss' Golang web server skeleton to build API for DP3](adr/0000-server-framework.md#use-truss-golang-web-server-skeleton-to-build-api-for-dp3)
* 0001 [Use Pop as the ORM for 3M](adr/0001-go-orm.md#use-pop-as-the-orm-for-3m)
* 0002 [Use dep to manage go dependencies](adr/0002-go-package-management.md#use-dep-to-manage-go-dependencies)
* 0003 [Put mymove into the standard GOPATH, eliminate server and client directories](adr/0003-go-path-and-project-layout.md#put-mymove-into-the-standard-gopath-eliminate-server-and-client-directories)
* 0004 [Use Both Absolute and Relative Paths for Imports](adr/0004-path-imports.md#use-both-absolute-and-relative-paths-for-imports)
* 0005 [Use Create React App](adr/0005-create-react-app.md#use-create-react-app)
* 0006 [Use Redux to manage state and Redux Thunk middleware to write action creators that return functions](adr/0006-redux.md#use-redux-to-manage-state-and-redux-thunk-middleware-to-write-action-creators-that-return-functions)
* 0007 [Use swagger-client to make calls to API from client](adr/0007-swagger-client.md#use-swagger-client-to-make-calls-to-api-from-client)
* 0008 [Use go-swagger To Route, Parse, And Validate API Endpoints](adr/0008-go-swagger.md#use-go-swagger-to-route-parse-and-validate-api-endpoints)
* 0009 [Generate forms from swagger definitions of payload](adr/0009-form-creation-from-swagger.md#generate-forms-from-swagger-definitions-of-payload)
* 0010 [Isolate Test Access to Database](adr/0010-isolate-test-access-to-database.md#isolate-test-access-to-database)
* 0011 [Test Suites](adr/0011-test-suites.md#test-suites)
* 0012 [The TSP Data Models](adr/0012-tsp-data-models.md#the-tsp-data-models)
* 0013 [REST API Updates](adr/0013-rest-api-updates.md#rest-api-updates)
* 0014 [Go Dependency Management](adr/0014-go-dependency-management.md#go-dependency-management)
* 0015 [Session storage/handling](adr/0015-session-storage.md#session-storage-handling)
* 0016 [Browser Support for Prototype](adr/0016-Browser-Support.md#browser-support-for-prototype)
* 0017 [Client side route restriction based on authentication](adr/0017-react-router-redux-authentication.md#client-side-route-restriction-based-on-authentication)
* 0018 [Optional Field Interop](adr/0018-optional-field-interop.md#optional-field-interop)
* 0019 [_Range Slider React Component_](adr/0019-client-rangeslider.md#range-slider-react-component)
* 0020 [Using Swagger to manage server route authentication](adr/0020-swagger-auth.md#using-swagger-to-manage-server-route-authentication)
* 0021 [Temporary use and plan for expunging Social Security Numbers in the prototype](adr/0021-ssn-use.md#temporary-use-and-plan-for-expunging-social-security-numbers-in-the-prototype)
* 0022 [Chose Excelize package to parse XLSX files](adr/0022-xlsx-lib.md#chose-excelize-package-to-parse-xlsx-files)
* 0023 [Representing Dollar Values in Go and the Database](adr/0023-representing-dollar-values.md#representing-dollar-values-in-go-and-the-database)
* 0024 [Model Authorization and Handler Design](adr/0024-model-authorization-and-handler-design.md#model-authorization-and-handler-design)
* 0025 [Client Side Feature Flags using Custom JavaScript](adr/0025-client-side-feature-flags.md#client-side-feature-flags-using-custom-javascript)
* 0026 [Use Snyk Vulnerability Scanning](adr/0026-use-snyk-vulnerability-scanning.md#use-snyk-vulnerability-scanning)
* 0027 [PDF Generation](adr/0027-pdf-generation.md#pdf-generation)
* 0028 [Config Management](adr/0028-config-management.md#config-management)
* 0029 [Honeycomb Integration](adr/0029-honeycomb-integration.md#honeycomb-integration)
* 0030 [IAM Authentication for Database](adr/0030-rds-iam.md#iam-authentication-for-database)
* 0031 [*CSS Tooling*](adr/0031-css-tooling.md#css-tooling)
* 0032 [CSRF Protection for the Application](adr/0032-csrf-protection.md#csrf-protection-for-the-application)

<!--endindex-->