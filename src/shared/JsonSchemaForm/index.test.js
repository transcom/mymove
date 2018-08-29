import React from 'react';
import { Provider } from 'react-redux';

import { reduxifyForm } from '.';
import configureStore from 'redux-mock-store';
import { mount } from 'enzyme';

const simpleSchema = {
  title: 'A registration form',
  description: 'A simple form example.',
  type: 'object',
  required: ['firstName', 'lastName'],
  properties: {
    firstName: {
      type: 'string',
      title: 'First name',
    },
    lastName: {
      type: 'string',
      title: 'Last name',
    },
    birthday: {
      type: 'string',
      format: 'date',
      title: 'Birthday',
    },
    sex: {
      type: 'string',
      title: 'sex',
      enum: ['Male', 'Female', 'Non-binary', 'Other'],
      'x-display-value': {
        Male: 'male',
        Female: 'female',
        'Non-binary': 'non-binary',
        Other: 'other',
      },
    },
    address: {
      $$ref: '#/definitions/Address',
      properties: {
        address1: {
          type: 'string',
          title: 'Address 1',
        },
        city: {
          type: 'string',
          title: 'City',
        },
      },
    },
  },
};
const uiSchema = {
  order: ['firstName', 'lastName', 'demographics', 'address'],
  definitions: {
    Address: {
      order: ['address1', 'city'],
    },
  },
  groups: {
    demographics: {
      title: 'demographics',
      fields: ['birthday', 'sex'],
    },
  },
};

// since JsonSchemaForm is using redux-form <Field /> components, so reduxifyForm must be called and a store must be provided
const TestForm = reduxifyForm('test');
const mockStore = configureStore();
let store;
let wrapper;
beforeEach(() => {
  store = mockStore({});
  //mount appears to be necessary to get inner components to load (i.e. tests fail with shallow)
  wrapper = mount(
    <Provider store={store}>
      <TestForm schema={simpleSchema} uiSchema={uiSchema} />
    </Provider>,
  );
});

it('renders without crashing', () => {
  // eslint-disable-next-line
  expect(wrapper.exists('form.default')).toBe(true);
});

it('renders 6 Field components', () => {
  expect(wrapper.find('Field').length).toBe(6);
});

it('renders select when there is an enum', () => {
  expect(wrapper.find('select').length).toBe(1);
});

it('renders date when format is date', () => {
  expect(wrapper.find('div.DayPickerInput').length).toBe(1);
});

it('renders a referenced field group', () => {
  expect(wrapper.text()).toContain('Address 1');
});
