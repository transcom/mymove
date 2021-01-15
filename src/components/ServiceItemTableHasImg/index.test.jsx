/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import { SERVICE_ITEM_STATUS } from '../../shared/constants';

import ServiceItemTableHasImg from './index';

describe('ServiceItemTableHasImg', () => {
  const defaultProps = {
    handleUpdateMTOServiceItemStatus: jest.fn(),
    handleShowRejectionDialog: jest.fn(),
  };

  it('should render no details', () => {
    const serviceItems = [
      {
        id: 'abc123',
        submittedAt: '2020-11-20',
        serviceItem: 'Fuel Surcharge',
        code: 'FSC',
        details: {},
      },
    ];
    const wrapper = mount(
      <ServiceItemTableHasImg
        {...defaultProps}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        serviceItems={serviceItems}
      />,
    );
    expect(wrapper.find('td').at(1).text()).toBe('—');
  });

  it('should render a thumbnail image with dimensions for item and crating', () => {
    const serviceItems = [
      {
        id: 'abc123',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic Crating',
        code: 'DCRT',
        details: {
          description: 'grandfather clock',
          imgURL: 'https://live.staticflickr.com/4735/24289917967_27840ed1af_b.jpg',
          itemDimensions: { length: 7000, width: 2000, height: 3500 },
          crateDimensions: { length: 10000, width: 2500, height: 5000 },
        },
      },
    ];

    const wrapper = mount(
      <ServiceItemTableHasImg
        {...defaultProps}
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
      />,
    );

    expect(wrapper.find('.siThumbnail').exists()).toBe(true);
    expect(wrapper.find('img').prop('src')).toBe(serviceItems[0].details.imgURL);
    expect(wrapper.find('dt').at(0).text()).toBe('Item Dimensions:');
    expect(wrapper.find('dd').at(0).text()).toBe('7"x2"x3.5"');
    expect(wrapper.find('dt').at(1).text()).toBe('Crate Dimensions:');
    expect(wrapper.find('dd').at(1).text()).toBe('10"x2.5"x5"');
  });

  it('should render customer contacts for DDFSIT service item', () => {
    const serviceItems = [
      {
        id: 'abc123',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic Crating',
        code: 'DDFSIT',
        details: {
          firstCustomerContact: { timeMilitary: '0400Z', firstAvailableDeliveryDate: '2020-12-31' },
          secondCustomerContact: { timeMilitary: '0800Z', firstAvailableDeliveryDate: '2021-01-01' },
        },
      },
    ];

    const wrapper = mount(
      <ServiceItemTableHasImg
        {...defaultProps}
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
      />,
    );

    expect(wrapper.find('table').exists()).toBe(true);
    expect(wrapper.find('.siThumbnail').exists()).toBe(false);
    expect(wrapper.find('dt').at(0).text()).toBe('First Customer Contact:');
    expect(wrapper.find('dd').at(0).text()).toBe('0400Z');
    expect(wrapper.find('dt').at(1).text()).toBe('First Available Delivery Date:');
    expect(wrapper.find('dd').at(1).text()).toBe('31 Dec 2020');

    expect(wrapper.find('dt').at(2).text()).toBe('Second Customer Contact:');
    expect(wrapper.find('dd').at(2).text()).toBe('0800Z');
    expect(wrapper.find('dt').at(3).text()).toBe('Second Available Delivery Date:');
    expect(wrapper.find('dd').at(3).text()).toBe('01 Jan 2021');
  });

  it('should render a zip and reason for DOFSIT service item', () => {
    const serviceItems = [
      {
        id: 'abc123',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic Origin 1st Day SIT',
        code: 'DOFSIT',
        details: {
          pickupPostalCode: '11111',
          reason: 'This is the reason',
        },
      },
    ];

    const wrapper = mount(
      <ServiceItemTableHasImg
        {...defaultProps}
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
      />,
    );
    expect(wrapper.find('dt').at(0).contains('ZIP')).toBe(true);
    expect(wrapper.find('dd').at(0).contains('11111')).toBe(true);
    expect(wrapper.find('dt').at(1).contains('Reason')).toBe(true);
    expect(wrapper.find('dd').at(1).contains('This is the reason')).toBe(true);
  });
});
