# REST API Updates

A large part of the functionality of the system is exposed as a [RESTful](https://en.wikipedia.org/wiki/Representational_state_transfer)
 API. Such APIs map standard HTTP methods (GET, POST, PUT, etc) to methods of manipulating the stateful objects
 the API is manipulating, e.g. Ruby on Rails uses the conventions [here](http://guides.rubyonrails.org/routing.html#crud-verbs-and-actions).

Whilst there is no official definition of the mapping used for a RESTful API there are some very common patterns like using
 `POST /widgets` to create a new Widget, the value of the Widget being in the body of the request. Likewise `GET /widgets/15435`
 will retrieve the Widget with ID 15435.

`PUT` is generally used to update or replace an object, i.e. `PUT /widgets/15435` would update or replace Widget 15435
with the new value in the body of the request.

Where the consensus breaks down is when a request wants to update or change just part of the state of an object, e.g.
change the description of Widget 15435 to be, say, "Recently updated - now in blue". Here there are a range of patterns
 in play, some being easy to implement but technically ambiguous/confusing to those that are 'correct' as proscribed in
  RFCs but in practice ignored and putting a greater burden on the clients of the API.

In particular, we want to consider a couple of use cases, viz:

* Updating the Actual Delivery Date on a shipment, i.e. a TSP delivers the shipment and wants to record the delivery date
* Accepting an Awarded Shipment, i.e. a shipment is awarded to a TSP and an agent for the TSP wants to accept that
 shipment. In this case the client may not (or possibly cannot) exactly understand what state in the object should change
 but clearly knows the action (`accept`) that they wish to perform.

## Considered Alternatives

* **Allowing state change via POST/PUT passing partial objects**, e.g.

  ```http request
  POST /shipments/1243 HTTP/1.1
  Content-Type: application/json

  {
    "status": 'ACCEPTED'
  }
  ```

* **Using PUT along with field specifiers to update state**, e.g. `PUT /shipments/13425/status` with 'ACCEPTED' as the body.

  ```http request
  POST /shipments/1243/status HTTP/1.1
  Content-Type: application/json

  {
    "value": 'ACCEPTED'
  }
  ```

* **Using PATCH to update objects by passing partial JSON objects and adding 'action' URLs**, e.g. `POST /shipments/13425/accept`
 to surface more complex state changes

  ```http request
  PATCH /shipments/1243 HTTP/1.1
  Content-Type: application/json

  {
    "status": 'ACCEPTED'
  }
  ```

  and if acceptance is a more complex operation which involves updating the state of the shipment more subtly than above, e.g.

  ```http request
  POST /shipments/1243/accept HTTP/1.1
  Content-Type: application/json

  {
    "reason": "Can accommodate this move"
  }
  ```

  NOTE: This is not explicitly setting either an 'accept' nor a 'reason' property of a shipment, but can be thought of as an rpc, e.g.

  ```javascript
  shipment.accept("Can accommodate this move");
  ```

* **Using PATCH along with [JSON Patch](https://tools.ietf.org/html/rfc6902) or [JSON Merge Patch](https://tools.ietf.org/html/rfc7386)**
 to update objects

  The [canonical way to use PATCH](https://en.wikipedia.org/wiki/Patch_verb#Patching_resources) is with an atomic
  description of the change, ideally using either [JSON Patch](https://tools.ietf.org/html/rfc6902) or
  [JSON Merge Patch](https://tools.ietf.org/html/rfc7386), e.g.

  ```http request
  PATCH /shipment/1243 HTTP/1.1
  Content-Type: application/json-patch+json

   [
     { "op": "replace", "path": "/status", "value": "ACCEPTED" },
     { "op": "add", "path": "/accept_reason", "value": "Can accomodate this move" }
   ]
  ```

## Decision Outcome

### Chosen Alternative: *Use `PATCH` with partial JSON objects (falling back to `POST`) to allow updates and `action` URLS for more complex operations*

* Justification: While using `PATCH` with partial objects (application/json) is frowned upon by some
 [commentators](http://williamdurand.fr/2014/02/14/please-do-not-patch-like-an-idiot/), it is common practice,
 [see Github](https://developer.github.com/v3/pulls/#update-a-pull-request) and is simple enough.

  It does not preclude adding support for one of the PATCH standards later as these use a different, explicit, content type, e.g.
 `application/json-patch+json`.

* Consequences: There is no good way to remove a field from an object without PUTting a new version of the object. In addition, people may "Well Actually" the API.

## Pros and Cons of the Alternatives

### *Allowing state change via POST/PUT passing partial objects*

Using POST for this is not problematic (in fact it's a good fallback when PATCH is not supported) but using PUT is ambiguous

* `+` Conceptually easy to understand
* `-` Ambiguous if using PUT
* `-` Provokingly non-standard
* `-` Relies on all updates to be done by explicitly altering fields on the objects - has no support for actions like `accept`
* `-` No easy way to remove fields from an object.

### *Using PUT along with field specifiers to update state*

* `+` Conceptually easy to understand
* `+/-` Arguably standard
* `-` Relies on all updates to be done by explicitly altering fields on the objects - has no support for actions like `accept`
* `-` Not common/familiar

### *Using PATCH to update objects by passing partial JSON objects and adding 'action' URLs*

* `+` Conceptually easy to understand
* `+` Offers low and high granularity changes
* `+` In line with common practice
* `+` Compatible with future support of JSON-PATCH or JSON-MERGE-PATCH
* `+/-` Arguably standard/not the best way to used patch
* `-` No easy way to remove fields from an object.

### *Using PATCH along with [JSON Patch](https://tools.ietf.org/html/rfc6902) or [JSON Merge Patch](https://tools.ietf.org/html/rfc7386)*

* `+` Strict adherence to best practices
* `+` Offers low and high granularity changes
* `-` Not common practice so unlikely to be easy for clients to implement
* `-` No support from Swagger codegen tools
