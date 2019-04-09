import React from 'react';
import { mount } from 'enzyme';
import LineItemTable from './LineItemTable';

describe('LineItemTable tests', () => {
  let wrapper;
  const shipmentLineItems = [
    {
      id: 'sldkjf',
      tariff400ng_item: { code: '105D', item: 'Reg Shipping' },
      amount: 1,
      quantity_1: 1,
      location: 'Destination',
    },
    {
      id: 'sldsdff',
      tariff400ng_item: { code: '105D', item: 'Reg Shipping' },
      location: 'Destination',
      amount: 1,
      quantity_1: 1,
    },
  ];

  describe('When shipmentLineItems exist', () => {
    it('renders without crashing', () => {
      wrapper = mount(
        <LineItemTable shipmentLineItems={shipmentLineItems} totalAmount={10} shipmentStatus="delivered" />,
      );
      expect(wrapper.find('table').length).toEqual(1);
      expect(wrapper.find('tr').length).toEqual(shipmentLineItems.length + 2);
    });
  });
});
