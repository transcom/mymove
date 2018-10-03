import React from 'react';
import { Provider } from 'react-redux';

import configureStore from 'redux-mock-store';
import { mount } from 'enzyme';

import PreApprovalRequestForm from '.';
// import SwaggerField from 'shared/JsonSchemaForm/JsonSchemaField';
// import Form from 'redux-form';

const simpleSchema = {
  properties: {
    id: {
      type: 'string',
      format: 'uuid',
      example: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
    },
    shipment_id: {
      type: 'string',
      format: 'uuid',
      example: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
    },
    accessorial_id: {
      type: 'string',
      format: 'uuid',
      example: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
    },
    quantity_1: {
      type: 'number',
      format: 'float',
      title: 'Base Quantity',
      description: 'Accessorial base quantity',
      minimum: 0,
      example: 16.7,
    },
    quantity_2: {
      type: 'integer',
      title: '2nd Quantity',
      description: 'Accessorial base quantity',
      minimum: 0,
      example: 10000,
    },
    location: {
      $ref: '#/definitions/AccessorialLocation',
    },
    notes: {
      type: 'string',
      title: 'Notes',
      example:
        'Mounted deer head measures 23" x 34" x 27"; crate will be 16.7 cu ft',
    },
    status: {
      $ref: '#/definitions/AccessorialStatus',
    },
    submitted_date: {
      type: 'string',
      title: 'Submitted Date',
      format: 'date-time',
      example: '2018-10-21T00:00:00.000Z',
    },
    approved_date: {
      type: 'string',
      title: 'Approved Date',
      format: 'date-time',
      example: '2018-10-21T00:00:00.000Z',
    },
    created_at: {
      type: 'string',
      format: 'date-time',
    },
    updated_at: {
      type: 'string',
      format: 'date-time',
    },
  },
};
const accessorials = [
  {
    id: 'sdlfkj',
    code: 'F9D',
    item: 'Long Haul',
  },
];
const submit = jest.fn();
const mockStore = configureStore();
let store;
let wrapper;
beforeEach(() => {
  store = mockStore({});
  //mount appears to be necessary to get inner components to load (i.e. tests fail with shallow)
  wrapper = mount(
    <Provider store={store}>
      <PreApprovalRequestForm
        ship_accessorial_schema={simpleSchema}
        accessorials={accessorials}
        onSubmit={submit}
      />
    </Provider>,
  );
});

it('renders without crashing', () => {
  // eslint-disable-next-line
  expect(wrapper.exists('div.usa-grid')).toBe(true);
  // Check that it renders swagger field content
  expect(wrapper.find('.half-width').length).toBe(1);
});
