# Swagger

This folder contains compiled Swagger (OpenAPI v2.0) specifications built from the API definitions in the
`./swagger-def/` directory. These files will be used by [Redoc](https://github.com/Redocly/redoc) to generate
documentation and [`go-swagger`](https://goswagger.io/) to generate the Go type files in `./pkg/gen/`. All YAML files in
this folder are **generated code** and should **NOT** be updated directly.

The HTML files are used by Redoc when previewing the documentation. The sub-folder `redoc/` contains an individual HTML
file for each documented API, in addition to the basic `_theme` JS and CSS files used for styling. The `index.html` file
on the top-level is integral for Redocly to know where to start building the docs and cannot be moved. Its only content
is a list of links to the individual API documentation files. This list should be updated whenever a new API HTML file
is added.
