# Optional Field Interop

Some of the fields in our API are optional. go-swagger represents optional fields using pointers, if the field is not present, the pointer is set to null. Pop (and the go database package in general) supports using pointers to represent null but are more geared toward using a special nullable field struct (NullString, NullInt, etc) where the struct has `value` and `valid` fields. The valid field indicates whether the optional field is present.

We have decided for now to adopt go-swagger's use of pointers throughout, defining our models as having pointers for all their optional fields.

## Considered Alternatives

* Use pointers to represent optional fields
* Use null structs to represent optional fields

## Decision Outcome

* Chosen Alternative: **Use pointers to represent optional fields**

Pointers are what go-swagger uses for optional fields so it makes conversions between the two representations of our data simpler. Additionally, pointers enforce the behavior of an Optional type where attempting to use the value of it when it is set to null results in a runtime error rather than silently succeeding in the case of the NullField.

## Pros and Cons of the Alternatives

### Use pointers to represent optional fields

* `+` This makes converting from the go-swagger representation of data to our model's representation of data simpler as for the basic types one can simply be set to the other.
* `+` Pointers correctly represent a Maybe/Optional type where they are either null or hold a value.
* `-` The purpose of pointers in go is to enable pass by reference or pass by value semantics, this use of them is not in any way related to that, so it is possible we will end up with some unexpected pointer related bugs by using them this way.

### Use null structs to represent optional fields

* `+` This is how the database package and pop are documented to handle null fields.
* `+` We don't mix our desire for optional fields with the semantics of reference and value, our models would not have pointers in them needlessly, cutting down on pointer bugs.
* `-` Converting from the go-swagger representation to our models would require more helper functions to make easy.
* `-` A NullString is only an OK maybe type. There is nothing stopping you from using the value in an instance where the valid bit is false. At least with a pointer your program will crash if you try the same thing.
