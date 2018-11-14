import React from 'react';
import { shallow } from 'enzyme';
import InvoiceTable from './InvoiceTable';

describe('InvoiceTable tests', () => {
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
      wrapper = shallow(<InvoiceTable shipmentLineItems={shipmentLineItems} totalAmount={10} />);
      expect(wrapper.find('table').length).toEqual(1);
    });
  });
});
