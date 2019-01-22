import React from 'react';
import { shallow } from 'enzyme';
import { UnbilledTable } from './UnbilledTable';
import * as CONSTANTS from 'shared/constants.js';
import { no_op } from 'shared/utils.js';

describe('UnbilledTable tests', () => {
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
    CONSTANTS.isOfficeSite = true;
  });

  describe('When shipmentLineItems exist', () => {
    it('renders without crashing', () => {
      wrapper = shallow(
        <UnbilledTable
          lineItems={shipmentLineItems}
          lineItemsTotal={10}
          approvePayment={no_op}
          cancelPayment={no_op}
          allowPayments={true}
          createInvoiceStatus={null}
        />,
      );
      expect(wrapper.find('div.invoice-panel-header-cont').length).toEqual(1);
      expect(wrapper.find('button').text()).toMatch('Approve Payment');
      expect(wrapper.find('LineItemTable').length).toEqual(1);
    });

    it('displays payment confirmation', () => {
      wrapper = shallow(
        <UnbilledTable
          lineItems={shipmentLineItems}
          lineItemsTotal={10}
          approvePayment={no_op}
          cancelPayment={no_op}
          allowPayments={true}
          createInvoiceStatus={null}
        />,
      );
      expect(wrapper.find('button').text()).toMatch('Approve Payment');
      wrapper.find('button').simulate('click');
      wrapper.update();
      expect(wrapper.find('span.warning--header').text()).toMatch("Please make sure you've double-checked everything.");
    });
  });
});
