# Browser Support for Prototype

**User Story:** 155599293

There are dozen of browsers supported on multiple operating systems.  We need to manage the level of dev and testing effort required to support additional browsers and the overall impact to the user

## Considered Alternatives

* Support only latest browser versions supported by operating systems (Internet Explorer 11, Edge, Chrome, Safari, and Firefox for desktop operating systems, Safari for iOS, and Chrome for Android)

* Extend support to older versions of Internet Explorer (IE 9 & 10)

## Decision Outcome

* Chosen Alternative: **Support only latest browser versions supported by operating systems**
* These browsers support a vast majority of the potential users of the product
* Consequences. There is a small risk that a service member could access the system with an unsupported browser.  The user may not experience any issues, but possibly could incur display issues or prevent them from finishing the a process.  The user would have to start over using a supported browser.

## Pros and Cons of the Alternatives

### Support only latest browser versions supported by operating systems

* `+` Covers a majority of browsers current used for internet access
* `+` Limits testing scope and combinations by limiting browser versions
* `+` The mobile browsers selected represent a dominant market share on their respective mobile OS
* `-` Slight risk of a service member unable to complete their process by using an unsupported browser

### Extend support to older versions of Internet Explorer (IE 9 & 10)

* `+` Reduces the likelihood of a service member not completing the process buy using an unsupported browser
* `+` The mobile browsers selected represent a dominant market share on their respective mobile OS
* `-` Extra coding effort to support older IE versions
* `-` Extra testing effort to support older browsers

## Specific Minimum Version Requirements

The following browsers will be considered _minimum requirements_, meaning that issues that arise in browsers that don't meet these requirements will not be prioritized above other work. Problems that crop up in browsers that do meet these requirements will be prioritized based on their severity: cosmetic issues are generally less important, while functionality breakage is considered a critical bug.

### Minimum Browser Requirements

#### Windows 10

* IE 11, Edge 12+

#### Windows 10, macOS 10.12+

* Chrome 64+
* Firefox 58+

#### macOS 10.12+

* Safari 10+

#### iOS 10+

* Mobile Safari (Locked to OS version)

#### Android 7+

* Chrome for Android 56+

## Browsers for Development and Testing

The development team will be developing using the latest version of Google Chrome or Firefox and will be regularly testing the application in:

* Latest Google Chrome (macOS 10.13)
* Latest Firefox (macOS 10.13)
* Latest iOS Safari (iOS 11)
* Latest Chrome (Android 8)
* Internet Explorer 11 (Windows 10)
