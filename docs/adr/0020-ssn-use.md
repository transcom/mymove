# *Temporary use and plan for expunging Social Security Numbers in the prototype*

**User Stories:**
* Service members can access moving services without having their identity at risk of being stolen.
* Couselor can locate a specific move record quickly.
* Finance can link an account, a move, and a service member to each other.

## Considered Alternatives

* Generating a unique move ID
* Relying on DODID
* Using SSNs temporarily
* Using SSNs permanently

## Decision Outcome

Chosen Alternative: *Using SSNs temporarily*
* Will use DODIDs wherever possible
* Will use SSNs where not possible so Finance can still have access to the information they need to do their work
* Once DODIDs are present in Orders across all branches, we would like to switch to this method and delete all SSNs from our infrastructure.
* SSN is not to be used as a key field for any data object so it can be masked and ultimately purged without breaking anything

## Pros and Cons of the Alternatives

### *Generating a unique move ID*

* `+` *Diminishes attack surface for identity theft (pro)* : not using SSNs prevents movers from seeing the nubmer, as well as lessening the impact if the database ever experiences a breach.
* `+` *Easily referenced by all stakeholders (pro)* : just like airline codes!
* `-` *Creates a number which needs to be referenced by additional government databases (con)* : we would need to link to government databases (specifically, to Orders and Accounting), and integration is hard + time consuming.

### *Relying on DODID*
Linking DODID and SSN to each other through a different gov database

* `+` *Diminishes attack surface for identity theft (pro)* : not using SSNs prevents movers from seeing the nubmer, as well as lessening the impact if the database ever experiences a breach.
* `+` *DODIDs are already present in most DOD databases (pro)* : easy to reference and cross-reference, including with finance, so long as the service member has a DODID.
* `-` *DODIDs are not printed on every branch's orders (con)* : we cannot consistently use them across the branches of service as of yet, as they do not appear on all orders.

### *Using SSNs temporarily*
Use DODIDs where possible and as rolled out, but use SSNs where they don't appear.
* `+` *Diminishes attack surface for identity theft (pro)* : not using SSNs when possible prevents all movers from seeing the nubmer, as well as lessening the impact if the database ever experiences a breach.
* `+` *Still allows some attack surface for identity theft (con)* : some service members will still be at risk
* `+` *DODIDs are already present in most DOD databases (pro)* : easy to reference and cross-reference, including with finance, so long as the service member has a DODID.
* `+` *Finance can access accounts of members who do not have a DODID associated with their orders (pro)* : and not everyone has a DODID on their Orders.

### *Using SSNs permanently*

* `+` *fails to diminish attack surface for identity theft of the member (con)*
* `+` *finance is able to look up each service member's account (pro)*
