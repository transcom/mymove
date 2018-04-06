# Temporary use and plan for expunging Social Security Numbers in the prototype

**User Stories:**

* Service members can access moving services without having their identity at risk of being stolen.
* Counselor can locate a specific move record quickly.
* Finance can link an account, a move, and a service member to each other.

## Considered Alternatives

* Generating a unique move ID
* Relying on EDIPI
* Using SSNs temporarily
* Using SSNs permanently

## Decision Outcome

Chosen Alternative: _Using SSNs temporarily_

* Will use EDIPIs wherever possible
* Will use SSNs where not possible so Finance can still have access to the information they need to do their work
* Once EDIPIs are present in Orders across all branches, we would like to switch to this method and delete all SSNs from our infrastructure.
* SSN is not to be used as a key field for any data object so it can be masked and ultimately purged without breaking anything

## Pros and Cons of the Alternatives

### _Generating a unique move ID_

* `+` _Diminishes attack surface for identity theft (pro)_ : not using SSNs prevents movers from seeing the number, as well as lessening the impact if the database ever experiences a breach.
* `+` _Easily referenced by all stakeholders (pro)_ : just like airline codes!
* `-` _Creates a number which needs to be referenced by additional government databases (con)_ : we would need to link to government databases (specifically, to Orders and Accounting), and integration is hard + time consuming.

### _Relying on EDIPI_

Linking EDIPI and SSN to each other through a different gov database

* `+` _Diminishes attack surface for identity theft (pro)_ : not using SSNs prevents movers from seeing the number, as well as lessening the impact if the database ever experiences a breach.
* `+` _EDIPIs are already present in most DOD databases (pro)_ : easy to reference and cross-reference, including with finance, so long as the service member has a DODID.
* `-` _EDIPIs are not printed on every branch's orders (con)_ : we cannot consistently use them across the branches of service as of yet, as they do not appear on all orders.

### _Using SSNs temporarily_

Use EDIPIs where possible and as rolled out, but use SSNs where they don't appear.

* `+` _Diminishes attack surface for identity theft (pro)_ : not using SSNs when possible prevents all movers from seeing the number, as well as lessening the impact if the database ever experiences a breach.
* `+` _Still allows some attack surface for identity theft (con)_ : some service members will still be at risk
* `+` _EDIPIs are already present in most DOD databases (pro)_ : easy to reference and cross-reference, including with finance, so long as the service member has a DODID.
* `+` _Finance can access accounts of members who do not have a EDIPI associated with their orders (pro)_ : and not everyone has a DODID on their Orders.

### _Using SSNs permanently_

* `+` _fails to diminish attack surface for identity theft of the member (con)_
* `+` _finance is able to look up each service member's account (pro)_
