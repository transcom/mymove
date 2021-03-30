import React from 'react';
import { shallow } from 'enzyme';

import ShipmentHeading from './ShipmentHeading';

const shipmentDestinationAddressWithPostalOnly = {
  postal_code: '98421',
};

const shipmentDestinationAddress = {
  street_1: '123 Main St',
  city: 'Tacoma',
  state: 'WA',
  postal_code: '98421',
};

const headingInfo = {
  shipmentID: '1',
  shipmentType: 'Household Goods',
  originCity: 'San Antonio',
  originState: 'TX',
  originPostalCode: '98421',
  destinationAddress: shipmentDestinationAddress,
  scheduledPickupDate: '27 Mar 2020',
  shipmentStatus: 'SUBMITTED',
};

describe('Shipment Heading with full destination address', () => {
  it('should render the data passed to it within the heading', () => {
    const wrapper = shallow(<ShipmentHeading shipmentInfo={headingInfo} handleUpdateMTOShipmentStatus={jest.fn()} />);
    expect(wrapper.find('h3').text()).toEqual('Household Goods');
    expect(wrapper.find('small').text()).toContain('San Antonio, TX 98421');
    expect(wrapper.find('small').text()).toContain('Tacoma, WA 98421');
    expect(wrapper.find('small').text()).toContain('27 Mar 2020');
  });
});

describe('Shipment Heading with missing destination address', () => {
  it("only renders the postal_code of the order's new duty station", () => {
    headingInfo.destinationAddress = shipmentDestinationAddressWithPostalOnly;
    const wrapper = shallow(<ShipmentHeading shipmentInfo={headingInfo} handleUpdateMTOShipmentStatus={jest.fn()} />);
    expect(wrapper.find('h3').text()).toEqual('Household Goods');
    expect(wrapper.find('small').text()).toContain('San Antonio, TX 98421');
    expect(wrapper.find('small').text()).toContain('98421');
    expect(wrapper.find('small').text()).toContain('27 Mar 2020');
  });
});
