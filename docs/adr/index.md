# Architectural Decision Log

This log lists the architectural decisions for DP3 Infrastructure.

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
- [ADR-0031](0031-css-tooling.md) - *CSS Tooling*
- [ADR-0032](0032-csrf-protection.md) - CSRF Protection for the Application
- [ADR-0033](0033-service-object-layer.md) - *Service Object Layer*
- [ADR-0034](0034-working-with-mocks-generation-and-assertion.md) - *Working With Mocks: Generation and Assertion*
- [ADR-0035](0035-use-query-builder.md) - Use Query Builder for for Admin Interface
- [ADR-0036](0036-go-integration.md) - Use Separate Integration Package for Go Integration Tests
- [ADR-0037](0037-go-path-and-project-layout-revisited.md) - Put mymove outside of standard GOPATH
- [ADR-0038](0038-soft-delete.md) - Use Soft Delete Instead of Hard Delete
- [ADR-0039](0039-react-lazy-and-code-splitting.md) - Use React Lazy for code splitting
- [ADR-0040](0040-role-base-authorization.md) - Add Role-Based Authorization
- [ADR-0041](0041-front-end-form-library.md) - Front End Form Library
- [ADR-0042](0042-optimistic-locking.md) - Use If-Match / E-tags for optimistic locking
- [ADR-0043](0043-prime-time.md) - *Handling time in the Prime API*
- [ADR-0044](0044-params-styling.md) - Use camelCase for API params
- [ADR-0045](0045-nesting-swagger-paths.md) - Nesting Swagger paths in the Prime API with multiple IDs
- [ADR-0046](0046-use-nodenv.md) - Use [nodenv](https://github.com/nodenv/nodenv) to manage Node versions in development
- [ADR-0047](0047-build-only-pull-requests-in-circleci.md) - Use CircleCI to build only Pull Requests and master
- [ADR-0048](0048-frontend-file-org.md) - Use a consistent file structure for front-end code
- [ADR-0049](0049-etag-for-child-updates.md) - Do not update child records using parent's E-tag
- [ADR-0050](0050-doc-viewer-fork.md) - Fork & maintain react-file-viewer under @trussworks
- [ADR-0051](0051-swagger-date-formats.md) - Use only Swagger supported formats for dates
- [ADR-0052](0052-use-data-testid.md) - Use `data-testid` as an attribute for finding components in tests
- [ADR-0053](0053-use-react-query-office-app.md) - Use React Query for Office App API interactions
- [ADR-0054](0054-use-CSS-to-highlight-unfinished-features.md) - Use CSS to highlight unfinished features
- [ADR-0055](0055-consolidate-moves-and-mtos.md) - Consolidate moves and move task orders into one database table
- [ADR-0056](0056-use-asdf-to-manage-golang-versions-in-development.md) - Use ASDF To Manage Golang Versions In Development
- [ADR-0057](0057-lodash.md) - Deprecate use of lodash over time
- [ADR-0058](0058-replace-loki-with-happo.md) - Use Happo for visual regression testing
- [ADR-0059](0059-use-snapshot-to-cleanup-loadtesting.md) - Use snapshot to cleanup load testing
- [ADR-0060](0060-move-state-for-service-counseling.md) - Move state for service counseling

<!-- adrlogstop -->

For new ADRs, please use [template.md](template.md).

More information on MADR is available at <https://adr.github.io/madr/>.
General information about architectural decision records is available at <https://adr.github.io/>.
