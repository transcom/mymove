import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';

import { reduxifyForm } from '.';
import configureStore from 'redux-mock-store';
import { shallow, mount, render } from 'enzyme';
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
    },
  },
};
const uiSchema = {
  order: ['firstName', 'lastName', 'demographics'],
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
  expect(wrapper.exists(<form className="default" />)).toBe(true);
});

it('renders 4 Field components', () => {
  expect(wrapper.find('Field').length).toBe(4);
});

it('renders select when there is an enum', () => {
  expect(wrapper.find('select').length).toBe(1);
});

it('renders date when format is date', () => {
  expect(wrapper.find('input[type="date"]').length).toBe(1);
});
