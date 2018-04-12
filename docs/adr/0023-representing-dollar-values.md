# Representing Dollar Values in Go and the Database

Care must be taken when representing dollars in code. Using floating point values can lead to unexpected rounding errors as floats don't usually equal the exact value you're trying to represent; they are often a tiny bit off from the value you think they are. After a few additions and multiplications on these numbers, these can add up and become visible errors.

## Considered Alternatives

* Represent USD values as cents (integers)
* Use a some kind of `decimal` type

## Decision Outcome

* Chosen Alternative: *Represent dollar values as cents (integers)*
* This is a commonly used pattern and should thus be familiar to developers.
* We avoid the overhead of evaluating 3rd party libraries and learning how to integrate them with Pop, our ORM.
* We will create a `cent` type (just an alias of an int) in order to help the compiler remind us of what our variables represent.

## Pros and Cons of the Alternatives

### Represent USD values as cents (integers)

What this means is that in the code base we will be representing "$1.00" as an integer: `100`. Additionally, we will create a type in Go such that it's clearer that we're working with cents. Doing so doesn't add much burden; it is possible to use integer arithmetic on these values.

For example:

```go
type cents int

// Dollars returns the dollar value of c in whole numbers.
func (c cents) Dollars() int {
    return int(c / 100)
}

func main() {
    var x cents
    var y cents
    x = 100
    y = 200

    fmt.Printf("%d + %d = %d\n", x, y, x+y)
    fmt.Printf("%d cents in dollars is $%d\n", x, x.Dollars())
}
```

* `+` This is a commonly used pattern and should be familiar to developers, and easy to maintain.
* `+` No 3rd party libraries required.
* `-` We might end up recreating the functionality of a decimal type.

### Use a some kind of `decimal` type

* `+` Adds a layer of type safety when using money values.
* `+` Postgres has this type built in: `Numeric(x,y)`
* `-` Go has no decimal type. We would have to find one.
* `-` Adds overhead in determining how to integrate the library with our ORM (how does Pop cast a decimal value from Postgres into whatever decimal class we adopt?)
* `-` Prevents us from using standard arithmetic operators as there is no overloading in Go.
