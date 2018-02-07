import React from 'react';
import { mount, shallow } from 'enzyme';
import { DD1299 } from '.';
import configureStore from 'redux-mock-store';
import { Provider } from 'react-redux';

const loadSchema = () => {};
const schema = {};
const uiSchema = {};
const mockStore = configureStore();

describe('When there is a form creation error', () => {
  let wrapper;

  beforeEach(() => {
    const hasSchemaError = true;
    const hasCreateError = false;
    const hasCreateSuccess = false;
    wrapper = shallow(
      <DD1299
        hasSchemaError={hasSchemaError}
        hasCreateSuccess={hasCreateSuccess}
        hasCreateError={hasCreateError}
        schema={schema}
        uiSchema={uiSchema}
        loadSchema={loadSchema}
      />,
    );
  });

  it('renders an error alert', () => {
    const alerts = wrapper.find('Alert');
    expect(alerts.length).toBe(1);
    expect(alerts.first().prop('type')).toBe('error');
  });
});

describe('When a form is successfully created', () => {
  let wrapper;

  beforeEach(() => {
    const hasSchemaError = false;
    const hasCreateError = false;
    const hasCreateSuccess = true;
    wrapper = shallow(
      <DD1299
        hasSchemaError={hasSchemaError}
        hasCreateSuccess={hasCreateSuccess}
        hasCreateError={hasCreateError}
        schema={schema}
        uiSchema={uiSchema}
        loadSchema={loadSchema}
      />,
    );
  });

  it('renders a success alert', () => {
    const alerts = wrapper.find('Alert');
    expect(alerts.length).toBe(1);
    expect(alerts.first().prop('type')).toBe('success');
  });
});
describe('When a form fails to be created', () => {
  let wrapper;

  beforeEach(() => {
    const hasSchemaError = false;
    const hasCreateError = true;
    const hasCreateSuccess = false;
    //provider and store are necessary here since this renders the redux form
    const store = mockStore({});
    wrapper = mount(
      <Provider store={store}>
        <DD1299
          hasSchemaError={hasSchemaError}
          hasCreateSuccess={hasCreateSuccess}
          hasCreateError={hasCreateError}
          schema={schema}
          uiSchema={uiSchema}
          loadSchema={loadSchema}
        />
      </Provider>,
    );
  });

  it('renders a failure alert', () => {
    const alerts = wrapper.find('Alert');
    expect(alerts.length).toBe(1);
    expect(alerts.first().prop('type')).toBe('error');
  });
});
