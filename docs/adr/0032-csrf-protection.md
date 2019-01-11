# CSRF Protection for the Application

**User Story:** Story [#162096596](https://www.pivotaltracker.com/story/show/162096596)

We want to be able to protect our application against [Cross-Site Request Forgery (CSRF)](https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)) attacks. While CSRF is no longer in the [OWASP top ten list](https://www.owasp.org/index.php/Category:OWASP_Top_Ten_Project) of security risks, our application does not use a full fledge Golang framework that includes CSRF protection. Therefore, we added protection to address the vulnerability.

## Decision Drivers

* Ease of implementation
* Good library documentation
* Works with our current Go web framework: Goji

## Considered Alternatives

* Double submit cookie method with justinas/nosurf
* Double submit cookie method with gorilla/csrf

## Decision Outcome

* Chosen Alternative: **Double-submit cookie method with gorilla/csrf**

* **Justification:** Gorilla/CSRF library has a fairly simple implementation that works with our framework
* Good documentation for JavaScript applications
* Works with Goji
* Generates unique-per-request (masked) tokens as a mitigation against the [BREACH](http://breachattack.com/) attack
* Uses [secure cookies](https://en.wikipedia.org/wiki/Secure_cookie) to store the unmasked csrf token session
* CSRF token session is stateless via [Double Submit Cookie method](https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)_Prevention_Cheat_Sheet#Double_Submit_Cookie) meaning that multiple browser tabs won't cause a user problems as their per-request token is compared with the base (unmasked) token.

* **Consequences:** HTTP requests will now need to include `x-csrf-token` header
* Every HTTP requests that modifies data will need the header

## Pros and Cons of the Alternatives

### *Double-submit cookie method with justinas/nosurf*

* `+` Works with Goji
* `+` Double submit cookie method
* `-` Documentation lacking for JavaScript applications
* `-` Doesn't use secure cookies
