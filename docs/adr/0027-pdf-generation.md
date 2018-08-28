# PDF Generation

**User Story:** _[ticket/issue-number]_ <!-- optional -->

MyMove will be required to generate several PDFs to provide all parties with the paperwork
required by existing processes and systems. These PDFs will need to be generated programmatically,
using data in the application's database.

* We'd like to minimize the amount of additional infrastructure work needed to support this functionality.
* Accordingly, we'd like to keep the solution within the existing ecosystems already in use by the project:
  * Go
  * JavaScript
  * (as a last resort) Python 3
* The generated PDFs should not be editable using standard software. Specifically, any form elements
  in the PDF should not be editable. _(It emerged later on that we might want to be able to generate PDFs containing fillable forms.)_
* PDFs will also need to be created as a compilation of existing PDFs and images that have been uploaded as documentation.
* Any libraries we adopt need to be actively maintained.

In the service of identifying solutions that would satisfy these requirements (and recognizing that the combination of two
specialized tools might be required), we divided the problem into two parts:

* Generate non-editable PDFs forms with the correct data filled out
* Merge multiple PDFs and images into a single PDF

## Considered Alternatives

### Generating non-editable PDFs

* Draw the PDF form manually using a library such as [gofpdf](https://github.com/jung-kurt/gofpdf)
* Create the PDF using HTML and convert it to PDF using a tool such as [gotenberg](https://github.com/thecodingmachine/gotenberg)
* Use [ReportLab](https://www.reportlab.com/) to fill out a pre-existing PDF containing editable fields
* Use [PDFKit](https://www.pdflabs.com/tools/pdftk-server/) to fill out a pre-existing PDF containing editable fields
* Use [FDFMerge](https://appligent.com/server-software/fdfmerge/) to fill out a pre-existing PDF containing editable fields
* Use [Unidoc](https://github.com/unidoc/unidoc) to fill out a preexisting PDF containing editable fields

### Merging PDFs and Images

* [gotenberg](https://github.com/thecodingmachine/gotenberg)
* [Unidoc](https://github.com/unidoc/unidoc)
* [poppler](https://poppler.freedesktop.org/)
* [pdfcpu](https://github.com/hhrutter/pdfcpu)

## Decision Outcome

### Generating non-editable PDFs

Chosen alternative: Draw the PDF form manually using [gofpdf](https://github.com/jung-kurt/gofpdf)

* `+` Drawing PDFs in this way is a relatively straightforward process
* `+` An investment of a few days appeared to be enough to have an initial form fully generated
* `+` Written in Go and requires only other Go libraries as dependencies
* `+` Does not require a deep understanding of the PDF format
* `+` Does not require immediate infrastructure changes to run in production
* `+` Open source
* `-` Requires learning a new drawing API unique to this library
* `-` Can require a lot of code for a large PDF
* `-` Changes to the form design will require code changes

### Merging PDFs and Images

Chosen alternative: [pdfcpu](https://github.com/hhrutter/pdfcpu)

* `+` Written in Go and requires only other Go libraries as dependencies
* `+` Was able to wrap images in PDFs using our test files
* `+` Was able to combine multiple PDFs into one using our test files
* `+` Exposes its command line functionality through a simple Go API
* `+` Open source
* `-` Requires images to be converted to PDFs first before merging the resulting files together

## Pros and Cons of the Alternatives <!-- optional -->

### Generating non-editable PDFs

Create the PDF using HTML and convert it to PDF using a tool such as [gotenberg](https://github.com/thecodingmachine/gotenberg)

* `+` Leverages the team's existing knowledge of HTML and CSS
* `+` Does not require a deep understanding of the PDF format
* `+` Open source
* `-` Runs as a separate service, which will require infrastructure work to support
* `-` Changes to the form design will require code changes
* `-` Aligning page breaks and other formatting may be difficult

Use [ReportLab](https://www.reportlab.com/) to fill out a pre-existing PDF containing editable fields

* `+` Installs via Pip
* `+` Does not require a deep understanding of the PDF format
* `+` Has commercial support
* `+` Open source
* `-` Requires a commercial license
* `-` Reportlab Account required for installation
* `-` Does not support "active form elements" and can not fill out PDF forms
* `-` Would require running Python in production (it is currently only used during development and deployment)

Use [PDFKit](https://www.pdflabs.com/tools/pdftk-server/) to fill out a pre-existing PDF containing editable fields

* `+` Popular solution for filling out PDF forms
* `-` Requires generating an FDF file to be merged with a PDF
* `-` Built using [GCJ](https://en.wikipedia.org/wiki/GNU_Compiler_for_Java), support for which is being removed from many Linux distributions
* `-` Has not been updated recently, and the development plans are not publicly available
* `-` Source is only available as an archive download
* `-` Would require infrastructure work to deploy

Use [FDFMerge](https://appligent.com/server-software/fdfmerge/) to fill out a pre-existing PDF containing editable fields

* `+` Fully supports filling in PDF forms
* `+` Has commercial support
* `-` Requires a commercial license
* `-` Closed source
* `-` OSX version is not up to date, which would complicate development
* `-` Available only as a binary download
* `-` Would require infrastructure work to deploy

Use [Unidoc](https://github.com/unidoc/unidoc) to fill out a preexisting PDF containing editable fields

* `+` Written in Go and requires only other Go libraries as dependencies
* `+` Open source
* `-` Requires a commercial license
* `-` Form functionality is very low-level, accordingly:
* `-` Requires a deep understanding of the PDF file format

### Merging PDFs and Images

[gotenberg](https://github.com/thecodingmachine/gotenberg)

* `+` Open source
* `-` Runs as a separate service, which will require infrastructure work to support

[Unidoc](https://github.com/unidoc/unidoc)

* `+` Written in Go and requires only other Go libraries as dependencies
* `+` Open source
* `-` Requires a commercial license
* `-` Form functionality is very low-level, accordingly:
* `-` Requires a deep understanding of the PDF file format

[poppler](https://poppler.freedesktop.org/)

* `+` Open source
* `-` Command line program, which will require infrastructure work to support
* `-` Requires images to be converted to PDFs first before merging the resulting files together
