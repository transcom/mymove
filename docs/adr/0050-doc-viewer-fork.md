# Fork & maintain react-file-viewer under @trussworks

**User Story:** [MB-2346 Orders Document viewer (PDF)](https://dp3.atlassian.net/browse/MB-2346)

The goal of this story is to build out a new document viewer component that implements the custom UI that has been designed, and improves user experience for viewing uploaded PDFs and images over the existing solution (which relies on the native browser PDF viewer).

Fortunately we've found an existing [file viewer library for React](https://github.com/plangrid/react-file-viewer) that is open source, supports many different file types, and has a well-structured and easy-to-understand source code.

Unfortunately, using this library out of the box is not ideal for a few reasons:

- There is no way to customize the HTML markup used for PDF controls (zoom) which is needed, not just to customize the styling but also to make sure we're meeting a11y standards
- There are no existing rotation controls and no way to add them without editing the source code
- The library does not appear to be actively maintained, the last release being September 27, 2017

It's my opinion that it would be beneficial for Truss in general to have a go-to library for displaying embedded binary files in React apps. Thankfully, using this library as a starting point gives us a solid foundation to extend and build off, and I don't believe will require a significant amount of overhead for completing the related stories. My suggestion is for us to fork the library under the @trussworks Github org, publish it to npm as @trussworks/react-file-viewer, and maintain it as open source, while expanding on it with features needed for MilMove’s implementation but keeping it abstract enough for other applications as well.

## Considered Alternatives

- Fork & maintain react-file-viewer under @trussworks
- Fork react-file-viewer and open PRs for improvements back to the original repo
- Copy and paste the react-file-viewer source code directly into MilMove
- Build our own document viewer from scratch

## Decision Outcome

- Chosen Alternative: Fork & maintain react-file-viewer under @trussworks
- With this option, we can immediately start editing & using the library to meet the requirements of MilMove, but also continue to maintain it for other future Truss projects as well.
- This means that Truss will have another JavaScript open source library to maintain, which does mean some overhead internally. However I think this is a beneficial area for us to develop more experience in, and establishes further practices around sharing common frontend code.

## Pros and Cons of the Alternatives

### Fork react-file-viewer and open PRs for improvements back to the original repo

- `+` Truss takes no ownership of the library and doesn't have to maintain it or take responsibility for publishing future releases
- `+` We might benefit from other contributions from the existing userbase if others also take this approach
- `-` The Github repo shows little to no activity by maintainers since fall of 2019, and this could indicate they are no longer interested in maintaining it.
- `-` We have less autonomy to take the library in a different direction (such as if we wanted to convert it to TypeScript)

### Copy and paste the react-file-viewer source code directly into MilMove

- `+` This is the most direct and (maybe) fastest way to get document viewer code into MilMove's codebase
- `+` We can edit the source code specifically for MilMove's needs without worrying about keeping it extensible
- `-` The existing code was built under a different environment (i.e., with different webpack config, different lint rules, different dependencies) and some number of adjustments would be required to consolidate that
- `-` We would not immediately get the benefits of a shared library for other Truss projects that might need similar functionality

### Build our own document viewer from scratch

- `+` We could use react-file-viewer as inspiration but write our own code directly in MilMove
- `-` This would probably be the most time consuming option, to set up everything from scratch in the MilMove code environment
- `-` We would not get the benefits of a shared library for other Truss projects that might need similar functionality
