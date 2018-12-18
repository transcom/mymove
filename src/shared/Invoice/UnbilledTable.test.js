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
    CONSTANTS.isOfficeSite = false;
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
          createInvoiceStatus={{
            error: null,
            isLoading: false,
            isSuccess: false,
          }}
        />,
      );
      expect(wrapper.find('InvoicePayment').length).toEqual(1);
      expect(wrapper.find('LineItemTable').length).toEqual(1);
    });
  });
});
