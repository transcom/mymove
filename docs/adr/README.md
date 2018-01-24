# Architectural Decision Log

This log lists the architectural decisions for DP3 server- and client-side code.

<!-- adrlog -- Regenerate the content by using "adr-log -i". You can install it via "npm install -g adr-log" -->

- [ADR-0000](0000-server-framework.md) - Use Truss' [golang](https://golang.org/) web server skeleton to build API for dp3
- [ADR-0001](0001-go-orm.md) - Use Pop as the ORM for 3M
- [ADR-0002](0002-go-package-management.md) - Use dep to manage go dependencies
- [ADR-0003](0003-go-path-and-project-layout.md) - Put mymove into the standard gopath, eliminte server and client directories
- [ADR-0004](0004-path-imports.md) - Use Both Absolute and Relative Paths for Imports
- [ADR-0005](0005-create-react-app.md) - Use Create React App
- [ADR-0006](0006-redux.md) - Use Redux to manage state and Redux Thunk middleware to write action creators that return functions

<!-- adrlogstop -->

For new ADRs, please use [template.md](template.md).

More information on MADR is available at <https://adr.github.io/madr/>.
General information about architectural decision records is available at <https://adr.github.io/>.
