import React from 'react';
import { mount } from 'enzyme';

import ImportantShipmentDates from './ImportantShipmentDates';

describe('ImportantShipmentDates', () => {
  const requestedPickupDate = 'Thursday, 26 Mar 2020';
  const scheduledPickupDate = 'Friday, 27 Mar 2020';
  const actualPickupDate = 'Saturday, 28 Mar 2020';
  const requiredDeliveryDate = 'Monday, 30 Mar 2020';
  const requestedDeliveryDate = 'Sunday, 29 Mar 2020';
  const scheduledDeliveryDate = 'Tuesday, 1 Apr 2020';
  const actualDeliveryDate = 'Wednesday, 2 Apr 2020';

  it('should render the shipment dates we pass in', () => {
    const wrapper = mount(
      <ImportantShipmentDates
        requestedPickupDate={requestedPickupDate}
        scheduledPickupDate={scheduledPickupDate}
        actualPickupDate={actualPickupDate}
        requestedDeliveryDate={requestedDeliveryDate}
        scheduledDeliveryDate={scheduledDeliveryDate}
        actualDeliveryDate={actualDeliveryDate}
        requiredDeliveryDate={requiredDeliveryDate}
      />,
    );
    expect(wrapper.find('td').at(0).text()).toEqual(requiredDeliveryDate);
    expect(wrapper.find('td').at(1).text()).toEqual(requestedPickupDate);
    expect(wrapper.find('td').at(2).text()).toEqual(scheduledPickupDate);
    expect(wrapper.find('td').at(3).text()).toEqual(actualPickupDate);
    expect(wrapper.find('td').at(4).text()).toEqual(requestedDeliveryDate);
    expect(wrapper.find('td').at(5).text()).toEqual(scheduledDeliveryDate);
    expect(wrapper.find('td').at(6).text()).toEqual(actualDeliveryDate);
  });

  it('should render an em-dash when no date is provided', () => {
    const emDash = '\u2014';
    const wrapper = mount(<ImportantShipmentDates />);
    expect(wrapper.find('td').at(0).text()).toEqual(emDash);
    expect(wrapper.find('td').at(1).text()).toEqual(emDash);
    expect(wrapper.find('td').at(2).text()).toEqual(emDash);
    expect(wrapper.find('td').at(3).text()).toEqual(emDash);
    expect(wrapper.find('td').at(4).text()).toEqual(emDash);
    expect(wrapper.find('td').at(5).text()).toEqual(emDash);
    expect(wrapper.find('td').at(6).text()).toEqual(emDash);
  });
});
