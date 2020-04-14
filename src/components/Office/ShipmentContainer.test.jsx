import React from 'react';
import { shallow } from 'enzyme';
import ShipmentHeading from './ShipmentHeading';
import ShipmentContainer from './ShipmentContainer';

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

describe('Shipment Container', () => {
  it('should render the container successfully', () => {
    const wrapper = shallow(
      <ShipmentContainer>
        <ShipmentHeading shipmentInfo={headingInfo} />
      </ShipmentContainer>,
    );
    expect(wrapper.find('.container').exists()).toBe(true);
  });
});
