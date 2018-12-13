import React from 'react';
import { shallow } from 'enzyme';
import LineItemTable from './LineItemTable';
import * as CONSTANTS from 'shared/constants.js';

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

  beforeEach(() => {
    CONSTANTS.isOfficeSite = false;
  });

  describe('When shipmentLineItems exist', () => {
    it('renders without crashing', () => {
      wrapper = shallow(
        <LineItemTable shipmentLineItems={shipmentLineItems} totalAmount={10} shipmentStatus="delivered" />,
      );
      expect(wrapper.find('table').length).toEqual(1);
    });
  });
});
