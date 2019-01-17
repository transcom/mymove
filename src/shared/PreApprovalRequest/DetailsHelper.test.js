import React from 'react';
import { shallow } from 'enzyme';
import { getDetailComponent } from './DetailsHelper';

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
describe('testing getDetailComponent()', () => {
  describe('returns default details component', () => {
    const DetailComponent = getDetailComponent();
    wrapper = shallow(<DetailComponent swagger={simpleSchema} />);

    it('renders without crashing', () => {
      // eslint-disable-next-line
      expect(wrapper.exists('div')).toBe(true);
    });
  });

  describe('returns 105B/E details component', () => {
    let DetailComponent = getDetailComponent('105B', true);
    wrapper = shallow(<DetailComponent />);
    it('renders 105B details without crashing', () => {
      // eslint-disable-next-line
      expect(wrapper.exists('div')).toBe(true);
    });

    DetailComponent = getDetailComponent('105E', true);
    wrapper = shallow(<DetailComponent />);
    it('renders 105E details without crashing', () => {
      // eslint-disable-next-line
      expect(wrapper.exists('div')).toBe(true);
    });
  });
});
