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
  it('renders the container successfully', () => {
    const wrapper = shallow(
      <ShipmentContainer>
        <ShipmentHeading shipmentInfo={headingInfo} />
      </ShipmentContainer>,
    );
    expect(wrapper.find('.container').exists()).toBe(true);
  });
  it('renders a child component passed to it', () => {
    const wrapper = shallow(
      <ShipmentContainer>
        <ShipmentHeading shipmentInfo={headingInfo} />
      </ShipmentContainer>,
    );
    expect(wrapper.find(ShipmentHeading).length).toBe(1);
  });
  it('renders a container with className container--accent--hhg', () => {
    let wrapper = shallow(
      <ShipmentContainer shipmentType="HHG">
        <ShipmentHeading shipmentInfo={headingInfo} />
      </ShipmentContainer>,
    );
    expect(wrapper.find('.container--accent--hhg').length).toBe(1);

    wrapper = shallow(
      <ShipmentContainer shipmentType="HHG_SHORTHAUL_DOMESTIC">
        <ShipmentHeading shipmentInfo={headingInfo} />
      </ShipmentContainer>,
    );
    expect(wrapper.find('.container--accent--hhg').length).toBe(1);

    wrapper = shallow(
      <ShipmentContainer shipmentType="HHG_LONGHAUL_DOMESTIC">
        <ShipmentHeading shipmentInfo={headingInfo} />
      </ShipmentContainer>,
    );
    expect(wrapper.find('.container--accent--hhg').length).toBe(1);
  });
  it('renders a container with className container--accent--nts', () => {
    const wrapper = shallow(
      <ShipmentContainer shipmentType="NTS">
        <ShipmentHeading shipmentInfo={headingInfo} />
      </ShipmentContainer>,
    );
    expect(wrapper.find('.container--accent--nts').length).toBe(1);
  });
});
