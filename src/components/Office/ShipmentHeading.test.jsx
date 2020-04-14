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
  it('should render the text heading', () => {
    const wrapper = shallow(<ShipmentHeading shipmentInfo={headingInfo} />);
    expect(wrapper.find('h3').text()).toEqual('Household Goods');
  });
});
