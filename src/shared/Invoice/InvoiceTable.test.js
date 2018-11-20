import React from 'react';
import { shallow } from 'enzyme';
import InvoiceTable from './InvoiceTable';
import { isOfficeSite } from 'shared/constants.js';
import * as CONSTANTS from 'shared/constants.js';

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

  beforeEach(() => {
    CONSTANTS.isOfficeSite = false;
  });

  describe('When shipmentLineItems exist', () => {
    it('renders without crashing', () => {
      wrapper = shallow(<InvoiceTable shipmentLineItems={shipmentLineItems} totalAmount={10} />);
      expect(wrapper.find('table').length).toEqual(1);
    });

    it('renders with Approve Payment button in Office app', () => {
      CONSTANTS.isOfficeSite = true;
      wrapper = shallow(
        <InvoiceTable shipmentLineItems={shipmentLineItems} totalAmount={10} shipmentStatus={'DELIVERED'} />,
      );

      expect(isOfficeSite).toBe(true);
      expect(wrapper.find('button').prop('disabled')).toBeTruthy();
    });
  });
});
