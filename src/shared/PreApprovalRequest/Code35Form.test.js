import React from 'react';
import { shallow } from 'enzyme';
import { Code35Form } from './Code35Form';

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
    },
    description: {
      type: 'string',
      format: 'textarea',
      title: 'Notes',
      example: 'Mounted deer head',
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

let wrapper;
describe('code 35A details component', () => {
  describe('renders', () => {
    wrapper = shallow(<Code35Form ship_line_item_schema={simpleSchema} />);

    it('without crashing', () => {
      expect(wrapper.exists('SwaggerField')).toBe(true);
    });
  });
});
