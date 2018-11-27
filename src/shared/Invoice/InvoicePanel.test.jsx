import React from 'react';
import { mount } from 'enzyme';

import { InvoicePanel } from './InvoicePanel';
import { isOfficeSite } from 'shared/constants.js';
import * as CONSTANTS from 'shared/constants.js';

describe('InvoicePanel tests', () => {
  describe('When no items exist', () => {
    let wrapper;
    const shipmentLineItems = [];
    wrapper = mount(
      <InvoicePanel unbilledShipmentLineItems={shipmentLineItems} lineItemsTotal={0} shipmentStatus={'DELIVERED'} />,
    );

    it('renders without crashing', () => {
      expect(wrapper.find('.empty-content').length).toEqual(1);
    });
  });

  describe('Approve Payment button shows on delivered state and office app', () => {
    CONSTANTS.isOfficeSite = true;
    const shipmentLineItems = [
      {
        id: 'sldkjf',
        tariff400ng_item: { code: '105D', item: 'Reg Shipping' },
        amount: 1,
        quantity_1: 1,
        location: 'DESTINATION',
      },
      {
        id: 'sldsdff',
        tariff400ng_item: { code: '105D', item: 'Reg Shipping' },
        location: 'DESTINATION',
        amount: 1,
        quantity_1: 1,
      },
    ];
    let wrapper = mount(
      <InvoicePanel unbilledShipmentLineItems={shipmentLineItems} lineItemsTotal={0} shipmentStatus={'DELIVERED'} />,
    );

    it('renders enabled "Approve Payment" button', () => {
      expect(isOfficeSite).toBe(true);
      expect(wrapper.props().shipmentStatus).toBe('DELIVERED');

      wrapper.update();
      expect(
        wrapper
          .children()
          .containsMatchingElement(<button className="button button-secondary">Approve Payment</button>),
      ).toBeTruthy();
    });

    it('renders disabled "Approve Payment" button', () => {
      expect(isOfficeSite).toBe(true);
      expect(wrapper.props().shipmentStatus).toBe('DELIVERED');

      wrapper.update();
      expect(wrapper.find('button').prop('disabled')).toBeTruthy();
    });
  });

  describe('When line items exist', () => {
    let wrapper;
    const shipmentLineItems = [
      {
        id: 'sldkjf',
        tariff400ng_item: { code: '105D', item: 'Reg Shipping' },
        amount: 1,
        quantity_1: 1,
        location: 'DESTINATION',
      },
      {
        id: 'sldsdff',
        tariff400ng_item: { code: '105D', item: 'Reg Shipping' },
        location: 'DESTINATION',
        amount: 1,
        quantity_1: 1,
      },
    ];
    wrapper = mount(
      <InvoicePanel unbilledShipmentLineItems={shipmentLineItems} lineItemsTotal={0} shipmentStatus={'DELIVERED'} />,
    );
    it('renders the table', () => {
      expect(wrapper.find('table').length).toEqual(1);
    });
  });
});
