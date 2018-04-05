# Representing Dollar Values in Go and the Database

Care must be taken when representing dollars in code. Using floating point values can lead to unexpected rounding errors as floats don't usually equal the exact value you're trying to represent; they are often a tiny bit off from the value you think they are. After a few additions and multiplications on these numbers, these can add up and become visible errors.

## Considered Alternatives

* Represent dollar values as cents (integers)
* Use a some kind of `decimal` type

## Decision Outcome

* Chosen Alternative: *Represent dollar values as cents (integers)*
* This is a commonly used pattern and should thus be familiar to developers.
* We avoid the overhead of evaluating 3rd party libraries and learning how to integrate them with Pop, our ORM.

## Pros and Cons of the Alternatives

### Represent dollar values as cents (integers)

* `+` This is a commonly used pattern and should be familiar to developers, and easy to maintain.
* `+` No 3rd party libraries required.
* `-` Some dollar values in the product go down to the thousandth of a cent (ex. `7.66185`). This means that there will be some special values in "millicents."
* `-` Adds a burden on the developer to understand the semantics of what is in an `int` value.

### Use a some kind of `decimal` type

* `+` Adds a layer of type safety when using money values.
* `+` Postgres has this type built in: `Numeric(x,y)`
* `-` Go has no decimal type. We would have to find one.
* `-` Adds overhead in determining how to integrate the library with our ORM (how does Pop cast a decimal value from Postgres into whatever decimal class we adopt?)
* `-` Prevents us from using standard arithmetic operators as there is no overloading in Go.
