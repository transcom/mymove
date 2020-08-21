import React from 'react';
import { mount } from 'enzyme';

import ImportantShipmentDates from './ImportantShipmentDates';

describe('ImportantShipmentDates', () => {
  const requestedPickupDate = 'Thursday, 26 Mar 2020';
  const scheduledPickupDate = 'Friday, 27 Mar 2020';

  it('should render the shipment dates we pass in', () => {
    const wrapper = mount(
      <ImportantShipmentDates requestedPickupDate={requestedPickupDate} scheduledPickupDate={scheduledPickupDate} />,
    );
    expect(wrapper.find('td').at(0).text()).toEqual(requestedPickupDate);
    expect(wrapper.find('td').at(1).text()).toEqual(scheduledPickupDate);
  });
});
