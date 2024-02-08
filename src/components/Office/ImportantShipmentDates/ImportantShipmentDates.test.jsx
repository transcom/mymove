import React from 'react';
import { mount } from 'enzyme';

import ImportantShipmentDates from './ImportantShipmentDates';

describe('ImportantShipmentDates', () => {
  const requiredDeliveryDate = 'Wednesday, 25 Mar 2020';
  const requestedPickupDate = 'Thursday, 26 Mar 2020';
  const scheduledPickupDate = 'Friday, 27 Mar 2020';
  const plannedMoveDate = 'Saturday, 28 Mar 2020';
  const actualMoveDate = 'Sunday, 29 Mar 2020';
  const requestedDeliveryDate = 'Monday, 30 Mar 2020';
  const scheduledDeliveryDate = 'Tuesday, 1 Apr 2020';
  const actualDeliveryDate = 'Wednesday, 2 Apr 2020';
  const actualPickupDate = 'Thursday, 3 Apr 2020';

  it('should render the shipment dates we pass in', () => {
    const wrapper = mount(
      <ImportantShipmentDates
        requestedPickupDate={requestedPickupDate}
        requiredDeliveryDate={requiredDeliveryDate}
        scheduledPickupDate={scheduledPickupDate}
        plannedMoveDate={plannedMoveDate}
        actualMoveDate={actualMoveDate}
        actualPickupDate={actualPickupDate}
        requestedDeliveryDate={requestedDeliveryDate}
        scheduledDeliveryDate={scheduledDeliveryDate}
        actualDeliveryDate={actualDeliveryDate}
        isPPM={false}
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
  });

  it('should show relevant PPM fields when it is a PPM', () => {
    const wrapper = mount(
      <ImportantShipmentDates
        requestedPickupDate={requestedPickupDate}
        requiredDeliveryDate={requiredDeliveryDate}
        scheduledPickupDate={scheduledPickupDate}
        plannedMoveDate={plannedMoveDate}
        actualMoveDate={actualMoveDate}
        actualPickupDate={actualPickupDate}
        requestedDeliveryDate={requestedDeliveryDate}
        scheduledDeliveryDate={scheduledDeliveryDate}
        actualDeliveryDate={actualDeliveryDate}
        isPPM
      />,
    );
    expect(wrapper.find('td').at(0).text()).toEqual(plannedMoveDate);
    expect(wrapper.find('td').at(1).text()).toEqual(actualMoveDate);
  });

  it('should not show irrelevant fields when it is a PPM', () => {
    const wrapper = mount(
      <ImportantShipmentDates
        requestedPickupDate={requestedPickupDate}
        scheduledPickupDate={scheduledPickupDate}
        plannedMoveDate={plannedMoveDate}
        actualMoveDate={actualMoveDate}
        requestedDeliveryDate={requestedDeliveryDate}
        scheduledDeliveryDate={scheduledDeliveryDate}
        actualDeliveryDate={actualDeliveryDate}
        isPPM
      />,
    );
    expect(wrapper.find('td').at(0).text()).not.toEqual(requestedPickupDate);
    expect(wrapper.find('td').at(1).text()).not.toEqual(requestedPickupDate);
    expect(wrapper.find('td').at(0).text()).not.toEqual(scheduledPickupDate);
    expect(wrapper.find('td').at(1).text()).not.toEqual(scheduledPickupDate);
    expect(wrapper.find('td').at(0).text()).not.toEqual(requestedDeliveryDate);
    expect(wrapper.find('td').at(1).text()).not.toEqual(requestedDeliveryDate);
    expect(wrapper.find('td').at(0).text()).not.toEqual(scheduledDeliveryDate);
    expect(wrapper.find('td').at(1).text()).not.toEqual(scheduledDeliveryDate);
    expect(wrapper.find('td').at(0).text()).not.toEqual(actualDeliveryDate);
    expect(wrapper.find('td').at(1).text()).not.toEqual(actualDeliveryDate);
  });

  it('should show relevant fields when it is a PPM', () => {
    const wrapper = mount(
      <ImportantShipmentDates
        requestedPickupDate={requestedPickupDate}
        scheduledPickupDate={scheduledPickupDate}
        plannedMoveDate={plannedMoveDate}
        actualMoveDate={actualMoveDate}
        requestedDeliveryDate={requestedDeliveryDate}
        scheduledDeliveryDate={scheduledDeliveryDate}
        actualDeliveryDate={actualDeliveryDate}
        isPPM
      />,
    );
    expect(wrapper.find('td').at(0).text()).toEqual(plannedMoveDate);
    expect(wrapper.find('td').at(1).text()).toEqual(actualMoveDate);
  });

  it('should not show irrelevant fields when it is a PPM', () => {
    const wrapper = mount(
      <ImportantShipmentDates
        requestedPickupDate={requestedPickupDate}
        scheduledPickupDate={scheduledPickupDate}
        plannedMoveDate={plannedMoveDate}
        actualMoveDate={actualMoveDate}
        requestedDeliveryDate={requestedDeliveryDate}
        scheduledDeliveryDate={scheduledDeliveryDate}
        actualDeliveryDate={actualDeliveryDate}
        isPPM
      />,
    );
    expect(wrapper.find('td').at(0).text()).not.toEqual(requestedPickupDate);
    expect(wrapper.find('td').at(1).text()).not.toEqual(requestedPickupDate);
    expect(wrapper.find('td').at(0).text()).not.toEqual(scheduledPickupDate);
    expect(wrapper.find('td').at(1).text()).not.toEqual(scheduledPickupDate);
    expect(wrapper.find('td').at(0).text()).not.toEqual(requestedDeliveryDate);
    expect(wrapper.find('td').at(1).text()).not.toEqual(requestedDeliveryDate);
    expect(wrapper.find('td').at(0).text()).not.toEqual(scheduledDeliveryDate);
    expect(wrapper.find('td').at(1).text()).not.toEqual(scheduledDeliveryDate);
    expect(wrapper.find('td').at(0).text()).not.toEqual(actualDeliveryDate);
    expect(wrapper.find('td').at(1).text()).not.toEqual(actualDeliveryDate);
  });
});
