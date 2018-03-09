# Browser Support for Prototype

**User Story:** 155599293

There are dozen of browsers supported on multiple operating systems.  We need to manage the level of dev and testing effort required to support additional browsers and the overall impact to the user

## Considered Alternatives

* *Windows 11, Edge, Chrome, Safari and Firefox for desktop operating systems and Safari for IOS and Chrome for Android*
* *Also support older Internet Explorer versions (IE 9 & 10)*

## Decision Outcome

* Windows 11, Edge, Chrome, Safari and Firefox for desktop operating systems that support the browser and Safari for IOS and Chrome for Android *[alternative 1]*
* *These browsers support a vast majority of the potential users of the product*
* Consequences. There is a small risk that a service member could access the system with an unsupported browser.  The user may not experience any issues, but possibly could incur display issues or prevent them from finishing the a process.  The user would have to start over using a supported browser. *

### *Latest browser versions supported by operating systems*

* `+` *Covers a majority of browsers current used for internet access*
* `+` *Limits testing scope and combinations by limiting browser versions*
* `+` *The mobile browsers selected represent a dominant market share on their respective mobile OS*
* `-` *Slight risk of a service member unable to complete their process by using an unsupported browser*

### *Increase browser versions to increase user reach*

* `+` *Reduces the likelihood of a service member not completing the process buy using an unsupported browser*
* `+` *The mobile browsers selected represent a dominant market share on their respective mobile OS*
* `-` *Extra coding effort to support older IE versions*
* `-` *Extra testing effort to support older browsers*
