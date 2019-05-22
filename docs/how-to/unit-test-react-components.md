# How To Unit Test React Components

Unit testing React components requires familiarity with two libraries:

* [Jest](https://jestjs.io/docs/en/expect), our test framework and runner
* [Enzyme](https://airbnb.io/enzyme/docs/api/), a React component testing library

*Note that the version of Jest we're using is dictated by `create-react-app` and tends to lag behind the latest releases.*

The basic pattern for unit testing React components is:

1. "Render" the component using either `shallow` or `mount`. Pass in whatever `props` are needed to by the component.
2. Make assertions against the component state using Jest and Enzyme.

It's best to avoid testing network requests or Redux when unit testing a component. If your component does either of these things, it is probably worth restructuring your code using the container pattern so that the inner component can be easily tested.

Remember that the value returned by either `shallow`, `mount`, or `render` will be a different object in each case and have different behavior as a result.

## Shallow Rendering

Shallow rendering renders a component but not any of its children. You'll want to keep this in mind when writing your assertions.

```javascript
import React from 'react';
import { shallow } from 'enzyme';
import HHGWeightWarning from './HHGWeightWarning';

describe('HHG with too high a weight estimate', function() {
  const shipment = { weight_estimate: 12000 };
  const entitlements = { weight: 10000 };
  const wrapper = shallow(<HHGWeightWarning shipment={shipment} entitlements={entitlements} />);

  it('shows a warning if the estimated weight is too high', function() {
    expect(wrapper.text()).toContain(
      'Your weight estimate of 12,000 is 2,000 lbs over your maximum entitlement of 10,000 lbs.',
    );
  });
});
```

If you need to see if a React component exists in the render tree, use the Enzyme API to check for it:

```javascript
describe('with valid weights', function() {
  const shipment = { weight_estimate: 1000, progear_weight_estimate: 200, spouse_progear_weight_estimate: 200 };
  const entitlements = { weight: 2000, pro_gear: 300, pro_gear_spouse: 300 };
  const wrapper = shallow(<HHGWeightWarning shipment={shipment} entitlements={entitlements} />);

  it('shows no alerts', function() {
    expect(wrapper.containsMatchingElement(Alert)).toEqual(false);
  });
});
```

Alternatively, you can use `html()` to render a node to HTML and make assertions against the markup.

## Full Rendering

While providing a more realistic perspective on the application's behavior, using `mount` can be problematic. One cause of this is that certain components (such as `Link` from `react-router`) can't be rendered unless they appear beneath another component or provider. In such a case, you'll need to setup any providers or parent components and render the component under test inside them.

For example, if a component needs to be rendered within a `Provider` that provides the `store` prop, you'll need to render it like so in a test:

```javascript
import React from 'react';
import { Provider } from 'react-redux';
import ReactDOM from 'react-dom';
import configureStore from 'redux-mock-store';
import { mount } from 'enzyme';

describe('HomePage tests', () => {
  const mockStore = configureStore();
  let store;
  describe('When the user has never logged in before', function() {
    store = mockStore({});
    const wrapper = mount(
      <Provider store={store}>
        <ComponentThatNeedsAccessToTheStore />
      </Provider>,
    );

    // assertions go here
  });
});
```

Generally, however, it is recommended to use the container pattern to separate the data access and rendering concerns of a component and to focus unit tests on the inner component.

## Static Rendering

Static rendering renders a React component to HTML and provides a [nice API](https://github.com/cheeriojs/cheerio/tree/aa90399c9c02f12432bfff97b8f1c7d8ece7c307#api) for traversing the resulting markup.

```javascript
import React from 'react';
import { render } from 'enzyme';
import HHGWeightWarning from './HHGWeightWarning';

describe('with valid weights', function() {
  const shipment = { weight_estimate: 1000, progear_weight_estimate: 200, spouse_progear_weight_estimate: 200 };
  const entitlements = { weight: 2000, pro_gear: 300, pro_gear_spouse: 300 };
  const wrapper = render(<HHGWeightWarning shipment={shipment} entitlements={entitlements} />);

  it('shows no alerts', function() {
    expect(wrapper.find('.usa-alert').length).toEqual(0);
  });
});
```
