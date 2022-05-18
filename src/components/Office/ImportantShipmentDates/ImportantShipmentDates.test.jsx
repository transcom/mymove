import React from 'react';
import { mount } from 'enzyme';

import ImportantShipmentDates from './ImportantShipmentDates';

describe('ImportantShipmentDates', () => {
  const requestedPickupDate = 'Thursday, 26 Mar 2020';
  const scheduledPickupDate = 'Friday, 27 Mar 2020';
  const requiredDeliveryDate = 'Monday, 30 Mar 2020';

  it('should render the shipment dates we pass in', () => {
    const wrapper = mount(
      <ImportantShipmentDates
        requestedPickupDate={requestedPickupDate}
        scheduledPickupDate={scheduledPickupDate}
        requiredDeliveryDate={requiredDeliveryDate}
      />,
    );
    expect(wrapper.find('td').at(0).text()).toEqual(requestedPickupDate);
    expect(wrapper.find('td').at(1).text()).toEqual(scheduledPickupDate);
    expect(wrapper.find('td').at(2).text()).toEqual(requiredDeliveryDate);
  });

  it('should render an em-dash when no date is provided', () => {
    const emDash = '\u2014';
    const wrapper = mount(<ImportantShipmentDates />);
    expect(wrapper.find('td').at(0).text()).toEqual(emDash);
    expect(wrapper.find('td').at(1).text()).toEqual(emDash);
    expect(wrapper.find('td').at(2).text()).toEqual(emDash);
  });
});
