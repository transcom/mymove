# Swagger Definitions

This folder contains the files that together define MilMove's APIs. All API specifications that are in active
development or production are being built from these files.

The entrypoints for each of these APIs are the top-level YAML files in this folder. Each distinct YAML file on this
level corresponds to a complete API definition.

## Contents

The sub-folders in `swagger-def` represent the top-level sections of a Swagger (OpenAPI v2.0) specification. These
sections are:

* [`info`](info/) [ [docs](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/2.0.md#info-object) ] -
  Contains metadata about the API. In particular, this is where the top-level descriptions for our APIs will live.

* [`tags`](tags/) [ [docs](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/2.0.md#tag-object) ] -
  Used to organize operations/endpoints. Each tag can have a description that will be visible under the section header
  in the reference docs. Tag component names are `camelCase`

* [`definitions`](definitions/) [ [docs](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/2.0.md#schema-object) ] -
  Reusable schema objects that define input/output data types. Some examples of objects that should be placed here are
  Shipments, Service Items, and Payment Requests. Definition component names are `PascalCase`

* [`responses`](responses/) [ [docs](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/2.0.md#response-object) ] -
  Each file describes a single response from an API Operation. Response component names are `PascalCase`

* [`parameters`](parameters/) [ [docs](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/2.0.md#parameter-object) ] -
  Each file describes a single operation parameter. Operation parameters include:
  * headers,
  * cookies,
  * request bodies,
  * query strings,
  * path variables. For example, in `/shipments/{shipmentId}`, the path parameter is `shipmentId`.

  Parameter component names are `camelCase`

* [`paths`](paths/) [ [docs](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/2.0.md#paths-object) ] -
  Defines each endpoint/operation. A path can have one operation per HTTP method.

Refer to the full [OpenAPI v2.0 Specification](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/2.0.md#openapi-specification)
for more detail about each section.

Each folder can contain either `.yaml` or `.md` files. When adding a YAML file, the filename and path will dictate where
this excerpt appears in the compiled specification. This filename must be a valid component name, and it is possible to
create filenames that are incompatible with the Swagger format. Please keep this in mind when adding new definitions and
sub-folders. Files located in the [`paths`](paths/README.md) directory have additional rules for their filenames.

Markdown files are not limited by filename, but they should be organized logically. `.md` files may only be referenced
by certain keys, such as `description`, in the YAML specification. To learn more about embedding markdown, read
[Embedded Markdown](https://redoc.ly/docs/api-reference-docs/embedded-markdown/) in the official Redocly documentation.

Note that some Redocly features are for premium (paid) users only - to confirm if a feature is supported in the
community edition, look for the "**Supported in Redoc CE**" tag (although not all free features will be marked this way).
