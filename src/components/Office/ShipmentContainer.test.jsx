import React from 'react';
import { shallow } from 'enzyme';

import ShipmentHeading from './ShipmentHeading';
import ShipmentContainer from './ShipmentContainer';

import { SHIPMENT_OPTIONS } from 'shared/constants';

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
    expect(wrapper.find('[data-testid="ShipmentContainer"]').exists()).toBe(true);
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
      <ShipmentContainer shipmentType={SHIPMENT_OPTIONS.HHG}>
        <ShipmentHeading shipmentInfo={headingInfo} />
      </ShipmentContainer>,
    );
    expect(wrapper.find('.container--accent--hhg').length).toBe(1);

    wrapper = shallow(
      <ShipmentContainer shipmentType={SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC}>
        <ShipmentHeading shipmentInfo={headingInfo} />
      </ShipmentContainer>,
    );
    expect(wrapper.find('.container--accent--hhg').length).toBe(1);

    wrapper = shallow(
      <ShipmentContainer shipmentType={SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC}>
        <ShipmentHeading shipmentInfo={headingInfo} />
      </ShipmentContainer>,
    );
    expect(wrapper.find('.container--accent--hhg').length).toBe(1);
  });
  it('renders a container with className container--accent--nts', () => {
    const wrapper = shallow(
      <ShipmentContainer shipmentType={SHIPMENT_OPTIONS.NTS}>
        <ShipmentHeading shipmentInfo={headingInfo} />
      </ShipmentContainer>,
    );
    expect(wrapper.find('.container--accent--nts').length).toBe(1);
  });
  it('renders a container with className container--accent--ntsr', () => {
    const wrapper = shallow(
      <ShipmentContainer shipmentType={SHIPMENT_OPTIONS.NTSR}>
        <ShipmentHeading shipmentInfo={headingInfo} />
      </ShipmentContainer>,
    );
    expect(wrapper.find('.container--accent--ntsr').length).toBe(1);
  });
});
