import React from 'react';
import { shallow } from 'enzyme';

import ShipmentHeading from './ShipmentHeading';

const headingInfo = {
  shipmentType: 'Household Goods',
  originCity: 'San Antonio',
  originState: 'TX',
  originPostalCode: '98421',
  destinationCity: 'Tacoma',
  destinationState: 'WA',
  destinationPostalCode: '98421',
  scheduledPickupDate: '27 Mar 2020',
};

describe('Shipment Heading', () => {
  it('should render the data passed to it within the heading', () => {
    const wrapper = shallow(<ShipmentHeading shipmentInfo={headingInfo} />);
    expect(wrapper.find('h3').text()).toEqual('Household Goods');
    expect(wrapper.find('small').text()).toContain('San Antonio TX 98421');
    expect(wrapper.find('small').text()).toContain('Tacoma WA 98421');
    expect(wrapper.find('small').text()).toContain('27 Mar 2020');
  });
});
