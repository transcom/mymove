# Architectural Decision Log

This log lists the architectural decisions for DP3 Infrastructure.

<!-- adrlog -- Regenerate the content by using "adr-log -i". You can install it via "npm install -g adr-log" -->

- [ADR-0000](0000-server-framework.md) - Use Truss' [golang](https://golang.org/) web server skeleton to build API for dp3
- [ADR-0001](0001-go-orm.md) - Use [Pop](https://github.com/markbates/pop) as Go ORM for dp3
- [ADR-0002](0002-go-package-management.md) - Use dep to manage go dependencies
- [ADR-0003](0003-go-path-and-project-layout.md) - Put mymove into the standard gopath, eliminte server and client directories

<!-- adrlogstop -->

For new ADRs, please use [template.md](template.md).

More information on MADR is available at <https://adr.github.io/madr/>.
General information about architectural decision records is available at <https://adr.github.io/>.
