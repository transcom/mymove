# Use React Lazy for code splitting

As things stand we use a standard React import pattern when we need to utilize
exported components in other source files. One such example of this scenario occurs when
we use routes to dynamically determine what content to present our users with, such as in
`scenes/office/index.jsx` or `scenes/MyMove/index.jsx`.

## Example using current method

```javascript
import EditContactInfo from 'scenes/Review/EditContactInfo';
...
<Switch>
...
  <ValidatedPrivateRoute exact path="/moves/review/edit-contact-info" component={EditContactInfo} />
...
</Switch>
```

## Considered Alternatives

* Leave things as they are
* Use [React Loadable](https://github.com/jamiebuilds/react-loadable)
* Use [Loadable Components](https://github.com/smooth-code/loadable-components)
* Use [React Lazy](https://reactjs.org/docs/code-splitting.html#reactlazy)

## Decision Outcome

* Chosen Alternative: **Use React Lazy**

This option allows us to ensure we are only loading what's required into the users browser for components as they are
first rendered. The technique can be effectively implemented in places where routes as used to dynamically decide
what component to offer the customer. Introducing this into those places allows us to gain performance benefit without
the introduction of significant complexity.

Using this specific tool allows us to make use of the pattern in a way that should be compatible for upcoming
performance focused React features that we may find desirous to use in the future.

### Example using React Lazy

```javascript
...
const DocumentViewer = lazy(() => import('./DocumentViewer'));
...
<Switch>
...
    <Suspense fallback={<LoadingPlaceholder />}>
      <RenderWithOrWithoutHeader
        component={DocumentViewer}
        withHeader={false}
        tag={DivOrMainTag}
        {...props}
      />
    </Suspense>
...
</Switch>
```

In this example we use `Lazy()` combined with `import()` to bring in our needed component. This combination is what
allows us to ensure that it isn't presented to the browser until first render, which in this case should be when the
customer navigates to this route. The `Suspense` tag provides an alternative thing to render while our content loads.

## Pros and Cons of the Alternatives

### Leave things as they are

* `+` No changes needed to be done.
* `-` No increase in front end performance

### Use React Loadable

* `+` Performance increase for frontend app
* `-` Risk introduced due to library no longer being maintained
* `-` Increase in verbosity for frontend routes

### Use Loadable Components

* `+` Performance increase for frontend app
* `+` React team recommends this library for some use cases (server side rendering)
* `-` Increase in verbosity for frontend routes
* `-` Abstraction layer on top of a relatively new React features. Unknown how this might interact with
development by React team on [Concurrent Mode](https://reactjs.org/docs/concurrent-mode-intro.html), which we
may want to use in the future.

### Use React Lazy

* `+` Performance increase for frontend app
* `+` Developed by React Team
* `+` Designed to eventually interact with [Concurrent Mode](https://reactjs.org/docs/concurrent-mode-intro.html).
* `-` Increase in verbosity for frontend routes
