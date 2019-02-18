import React from 'react';
import { Provider } from 'react-redux';

import configureStore from 'redux-mock-store';
import { mount } from 'enzyme';

import { PreApprovalForm } from './PreApprovalForm';

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
    tariff400ng_item_id: {
      type: 'string',
      format: 'uuid',
      example: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
    },
    item_dimensions: {
      type: 'object',
      format: 'dimensions',
      properties: {
        length: {
          type: 'integer',
          format: 'dimension',
        },
        width: {
          type: 'integer',
          format: 'dimension',
        },
        height: {
          type: 'integer',
          format: 'dimension',
        },
      },
    },
    crate_dimensions: {
      type: 'object',
      format: 'dimensions',
      properties: {
        length: {
          type: 'integer',
          format: 'dimension',
        },
        width: {
          type: 'integer',
          format: 'dimension',
        },
        height: {
          type: 'integer',
          format: 'dimension',
        },
      },
    },
    quantity_1: {
      type: 'integer',
      format: 'basequantity',
      title: 'Base Quantity',
      description: 'Line Item base quantity',
      minimum: 0,
      example: 167000,
    },
    quantity_2: {
      type: 'integer',
      format: 'basequantity',
      title: '2nd Quantity',
      description: 'Line Item base quantity',
      minimum: 0,
      example: 10000,
    },
    location: {
      type: 'string',
      title: 'Location',
      enum: ['ORIGIN', 'DESTINATION'],
      'x-display-value': { ORIGIN: 'Origin', DESTINATION: 'Destination' },
    },
    notes: {
      type: 'string',
      title: 'Notes',
      example: 'Mounted deer head measures 23" x 34" x 27"; crate will be 16.7 cu ft',
    },
    status: {
      $ref: '#/definitions/ShipmentLineItemStatus',
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
const tariff400ng_items = [
  {
    id: 'sdlfkj',
    code: 'F9D',
    item: 'Long Haul',
  },
];

const simple105TariffItem = {
  code: '105B',
};

const simpleNon105TariffItem = {
  code: '28A',
};

const filteredLocations = ['ORIGIN', 'DESTINATION'];
const submit = jest.fn();
const mockStore = configureStore();
let store;
let wrapper;

describe('PreApprovalForm tests', () => {
  beforeEach(() => {
    store = mockStore({});
  });
  describe('When a PreApprovalForm is loaded', () => {
    beforeEach(() => {
      //mount appears to be necessary to get inner components to load (i.e. tests fail with shallow)
      wrapper = mount(
        <Provider store={store}>
          <PreApprovalForm
            ship_line_item_schema={simpleSchema}
            tariff400ngItems={tariff400ng_items}
            onSubmit={submit}
            tariff400ngItem={simpleNon105TariffItem}
          />
        </Provider>,
      );
    });
    it('renders without crashing', () => {
      // eslint-disable-next-line
      expect(wrapper.exists('div.usa-grid-full')).toBe(true);
      // Check that it renders swagger field content
      expect(wrapper.find('.half-width').length).toBe(6);
    });
  });
  describe('When a LocationSearch box is loaded with multiple possible locations', () => {
    beforeEach(() => {
      wrapper = mount(
        <Provider store={store}>
          <PreApprovalForm
            ship_line_item_schema={simpleSchema}
            tariff400ngItems={tariff400ng_items}
            onSubmit={submit}
            filteredLocations={filteredLocations}
            tariff400ngItem={simpleNon105TariffItem}
          />
        </Provider>,
      );
    });
    it('shows a dropdown', () => {
      expect(wrapper.exists('select')).toBe(true);
    });
  });
  describe('When a LocationSearch box is loaded with only one possible location', () => {
    beforeEach(() => {
      wrapper = mount(
        <Provider store={store}>
          <PreApprovalForm
            ship_line_item_schema={simpleSchema}
            tariff400ngItems={tariff400ng_items}
            onSubmit={submit}
            filteredLocations={['ORIGIN']}
            tariff400ngItem={simpleNon105TariffItem}
          />
        </Provider>,
      );
    });
    it('shows text', () => {
      expect(wrapper.find('.location-select div').text()).toEqual('Origin');
    });
  });
  describe('When code 105B/105E is chosen', () => {
    beforeEach(() => {
      wrapper = mount(
        <Provider store={store}>
          <PreApprovalForm
            ship_line_item_schema={simpleSchema}
            tariff400ngItems={tariff400ng_items}
            onSubmit={submit}
            filteredLocations={['ORIGIN']}
            tariff400ng_item_code={'105B'}
            tariff400ngItem={simple105TariffItem}
            context={{ flags: { robustAccessorial: true } }}
          />
        </Provider>,
      );
    });
    it('renders dimensions forms without crashing', () => {
      expect(wrapper.find('.dimensions-form').length).toBe(2);
    });
  });
});
