# Using React-App-Rewired

**User Story:** *[MB-9033](https://dp3.atlassian.net/browse/MB-9033)* :lock:

## Background

The MilMove client application uses **Create-React-App** and **React-Scripts**
for its build toolchain. This build toolchain has many benefits to developing
the client application, but also has limitations when it comes to updating
the configuration of various build tools that are used in the development of
the client application. These tools such as **webpack**, **ESLint**, and
**Babel** are configured with pre-determined configurations that are
in-accessible without *ejecting* from the **Create-React-App** toolchain.

Facebook's own documentation mentions that *If you aren't satisfied with the
build tool and configuration choices, you can `eject` at any time. This command
will remove the single build dependency from your project.* This is not an
option currently on MilMove as ejecting would add some maintainability overhead
that the development team may address at a later time or a later ADR.

## Decision drivers

Some of the issues we've seen is that in order to best control our build
toolchain there are certain configurations that need to be updated. For
instance, [**webpack** 5 removed the Node polyfills][wp5-migrate] that are used
by the MilMove client application. This causes our client application to break
in unique ways around `process` not being defined. Refactoring our client
application is a possible solution, but there are ways to have **webpack** 5 use
the Node polyfills that we need by defining it in a **webpack** configuration
file. Our issue here is that **Create-React-App** and **React-Scripts** prevent
us from modifying these scripts as they are not exposed to users of those
libraries.

[wp5-migrate]:
https://webpack.js.org/migrate/5/#run-a-single-build-and-follow-advice

This is true for other tools such as **ESLint** as well and has been an issue
previously on the client application for the project. Sometimes the
**React-Scripts** and **Create-React-App** tools update and support some level
of customization but it leaves the MilMove engineering team in a holding pattern
without a clear path forward besides following contributions upstream or
contributing upstream to the **Create-React-App** project.

It's also true that Facebook's **Create-React-App** development team is
notorious for denying any changes for specific configuration updates. As stated
above, their own documentation states that `ejecting` or maintaining a fork of
**Create-React-App** are the only viable solutions for customization using their
build toolchain.

## Considered Alternatives

> **bold denotes chosen**

* *Do nothing*
* *Eject the **Create-React-App** configuration*
* **Dynamically patch the configurations**

## Decision Outcome

* Chosen Alternative: *Dynamically patch the configurations*

We can leverage [**React-App-Rewired** in order to perform dynamic patching of
the **webpack** configuration][gh-rar], and other configurations as needed, in
order to bring back these Node polyfills for supporting the `process` object
that the client application currently relies on. This is also useful in case
there are other configuration options that would like to be modified or extended
in the future. The **React-App-Rewired** library supports modifying **webpack**
& **Jest** configurations.

[gh-rar]: https://github.com/timarney/react-app-rewired

## Pros and Cons of the Alternatives <!-- optional -->

### *Do nothing*

* `+` *[argument 1 pro]*
* `+` *[argument 2 pro]*
* `-` *[argument 1 con]*
* *[...]* <!-- numbers of pros and cons can vary -->

### *Eject the Create-React-App configuration*

* `+` *[argument 1 pro]*
* `+` *[argument 2 pro]*
* `-` *[argument 1 con]*
* *[...]* <!-- numbers of pros and cons can vary -->

### *Dynamically patch the configurations*

* `+` *[argument 1 pro]*
* `+` *[argument 2 pro]*
* `-` *[argument 1 con]*
* *[...]* <!-- numbers of pros and cons can vary -->
