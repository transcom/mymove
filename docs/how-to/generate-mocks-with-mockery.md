# How To Generate Mocks with Mockery

[Mockery](https://github.com/vektra/mockery) provides the ability to easily generate mocks for golang interfaces. It removes the boilerplate coding required to use mocks.

 *In Golang, mocks can only be created on interfaces - not structs. So, it is important that for whichever mock you are trying to generate, it should correspond to the appropriate interface.*

## Auto-generating mocks with `go generate`

 The `make mocks_generate` command will regenerate mocks for all interfaces tagged with the appropriate `go generate` command. To add an interface to the list of auto-generated mocks, just add a
 `go:generate` comment like below and update the name with your interface name.

```.go
// AccessCodeClaimer is the service object interface for ValidateAccessCode
//go:generate mockery -name AccessCodeClaimer
type AccessCodeClaimer interface {
    ClaimAccessCode(code string, serviceMemberID uuid.UUID) (*models.AccessCode, *validate.Errors, error)
}
```
