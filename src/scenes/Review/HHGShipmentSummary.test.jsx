import React from 'react';
import { shallow } from 'enzyme';
import HHGShipmentSummary from './HHGShipmentSummary';
import Address from './Address';

describe('HHG without a delivery address', function() {
  const shipment = { move_dates_summary: { pack: [], pickup: [], transit: [], delivery: [] } };
  const entitlements = {};
  const wrapper = shallow(<HHGShipmentSummary shipment={shipment} entitlements={entitlements} movePath="" />);

  it('warns if there is no delivery address', function() {
    const warning = wrapper.first('.delivery-notice');
    expect(warning.text()).toContain("If you don't have a delivery address before the movers arrive");
  });

  it('shows that no delivery address was entered', function() {
    const tr = wrapper.first('.delivery-address');
    expect(tr.text()).toContain('Delivery Address:Not entered');
  });
});

describe('HHG with a delivery address', function() {
  const shipment = {
    has_delivery_address: true,
    delivery_address: {
      street_address_1: '123 some street',
    },
    move_dates_summary: { pack: [], pickup: [], transit: [], delivery: [] },
  };
  const entitlements = {};
  const wrapper = shallow(<HHGShipmentSummary shipment={shipment} entitlements={entitlements} movePath="" />);

  it('shows the delivery address', function() {
    const tr = wrapper.find('.delivery-address');
    expect(tr.contains(<Address address={shipment.delivery_address} />)).toBe(true);
  });
});
