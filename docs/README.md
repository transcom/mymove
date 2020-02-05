# Programming Guide

The intention for this document is to share our collective knowledge on best practices and allow everyone working on the DOD MyMove project to write code in compatible styles.

If you are looking to understand choices made in this project, see the list of [ADRs](https://github.com/transcom/mymove/tree/master/docs/adr).

## Programming Guides

* [Front-end / React](frontend.md) guide
* [Back-end / Go](backend.md) guide
  * [Service Objects / Go](service-objects.md) guide

## Pairing Sessions

* [Pairing](pairing.md). A list of past pairing recordings.

## Metrics

* [Metrics](metrics.md). Documentation for application metrics.

## Security

* [Anti-Virus](anti_virus.md). Documentation for the anti-virus solutions employed.

<!--index-->

## Database

* [Backup and Restore the Development Database](database/backup-and-restore-dev-database.md#how-to-backup-and-restore-the-development-database)
* [Database Guides](database/database.md#database-guides)
* [Migrate the Database](database/migrate-the-database.md#how-to-migrate-the-database)
* [Soft Delete](database/soft-delete.md#how-to-soft-delete)

## HOWTOs

* [Access a Global Application Variable](how-to/access-global-variables.md#how-to-access-a-global-application-variable)
* [Call Swagger Endpoints from React](how-to/access-swagger-endpoints-from-react.md#how-to-call-swagger-endpoints-from-react)
* [Add Application Logging](how-to/add-application-logging.md#how-to-add-application-logging)
* [Automatically add JIRA ID to Commit Message](how-to/automatically-add-jira-id-to-commit-message.md#how-to-automatically-add-jira-id-to-commit-message)
* [Create An ECS Scheduled Task](how-to/create-an-ecs-scheduled-task.md#how-to-create-an-ecs-scheduled-task)
* [Create or Deactivate Users](how-to/create-or-deactivate-users.md#how-to-create-or-deactivate-users)
* [Deploy to Experimental](how-to/deploy-to-experimental.md#how-to-deploy-to-experimental)
* [display dates and times](how-to/display-dates-and-times.md#how-to-display-dates-and-times)
* [Generate Mocks with Mockery](how-to/generate-mocks-with-mockery.md#how-to-generate-mocks-with-mockery)
* [handle back-end errors](how-to/handle-backend-errors.md#how-to-handle-back-end-errors)
* [Make a Sample Prime API Call](how-to/make-a-sample-prime-api-call.md#how-to-make-a-sample-prime-api-call)
* [Manage Dependabot](how-to/manage-dependabot.md#how-to-manage-dependabot)
* [Manage Dependencies With go mod](how-to/manage-dependencies-with-go-mod.md#how-to-manage-dependencies-with-go-mod)
* [Manage Docker Locally](how-to/manage-docker-locally.md#how-to-manage-docker-locally)
* [revert a change](how-to/revert-a-change.md#how-to-revert-a-change)
* [Run Acceptance Tests](how-to/run-acceptance-tests.md#how-to-run-acceptance-tests)
* [Run Against S3 & CDN Locally](how-to/run-against-s3-locally.md#how-to-run-against-s3-cdn-locally)
* [Run End to End (Cypress) Tests](how-to/run-e2e-tests.md#how-to-run-end-to-end-cypress-tests)
* [Run Go Tests](how-to/run-go-tests.md#how-to-run-go-tests)
* [Run JavaScript (Jest) Tests](how-to/run-js-tests.md#how-to-run-javascript-jest-tests)
* [Run and troubleshoot pre-commit hooks](how-to/run-pre-commit-hooks.md#run-and-troubleshoot-pre-commit-hooks)
* [Run server_test job in CircleCI container locally](how-to/run-server-test-circle-ci.md#run-server-test-job-in-circleci-container-locally)
* [Use and Run Storybook](how-to/run-storybook.md#how-to-use-and-run-storybook)
* [Searching for Application Errors](how-to/search-for-application-errors.md#how-to-searching-for-application-errors)
* [Setup Postman to make Mutual TLS API Calls](how-to/setup-postman-to-make-mutual-tls-api-calls.md#how-to-setup-postman-to-make-mutual-tls-api-calls)
* [Store Data in Redux](how-to/store-data-in-redux.md#how-to-store-data-in-redux)
* [Store UI State in Redux](how-to/store-ui-state-in-redux.md#how-to-store-ui-state-in-redux)
* [Test Virus Scanning](how-to/test-virus-scanning.md#how-to-test-virus-scanning)
* [Troubleshoot GEX Connection](how-to/troubleshoot-gex-connection.md#how-to-troubleshoot-gex-connection)
* [Unit Test React Components](how-to/unit-test-react-components.md#how-to-unit-test-react-components)
* [Upgrade Go Version](how-to/upgrade-go-version.md#how-to-upgrade-go-version)
* [Upload Electronic Orders Using your CAC](how-to/upload-electronic-orders.md#how-to-upload-electronic-orders-using-your-cac)
* [Upload Electronic Orders Using your CAC](how-to/use-mtls-with-cac.md#how-to-upload-electronic-orders-using-your-cac)
* [View ECS Service Logs](how-to/view-ecs-service-logs.md#how-to-view-ecs-service-logs)

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
* 0033 [*Service Object Layer*](adr/0033-service-object-layer.md#service-object-layer)
* 0034 [*Working With Mocks: Generation and Assertion*](adr/0034-working-with-mocks-generation-and-assertion.md#working-with-mocks-generation-and-assertion)
* 0035 [Use Query Builder for for Admin Interface](adr/0035-use-query-builder.md#use-query-builder-for-for-admin-interface)
* 0036 [Use Separate Integration Package for Go Integration Tests](adr/0036-go-integration.md#use-separate-integration-package-for-go-integration-tests)
* 0037 [Put mymove outside of standard GOPATH](adr/0037-go-path-and-project-layout-revisited.md#put-mymove-outside-of-standard-gopath)
* 0038 [Use Soft Delete Instead of Hard Delete](adr/0038-soft-delete.md#use-soft-delete-instead-of-hard-delete)
* 0039 [Use React Lazy for code splitting](adr/0039-react-lazy-and-code-splitting.md#use-react-lazy-for-code-splitting)
* 0040 [Add Role-Based Authorization](adr/0040-role-base-authorization.md#add-role-based-authorization)
* 0041 [Front End Form Library](adr/0041-front-end-form-library.md#front-end-form-library)
* 0042 [Use Last-Modified / If-Unmodified-Since for optimistic locking](adr/0042-optimistic-locking.md#use-last-modified-if-unmodified-since-for-optimistic-locking)
* 0043 [*Handling time in the Prime API*](adr/0043-prime-time.md#handling-time-in-the-prime-api)

<!--endindex-->
