/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import ServiceItemsTable from './ServiceItemsTable';

import { SERVICE_ITEM_STATUS } from 'shared/constants';
import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

describe('ServiceItemsTable', () => {
  const defaultProps = {
    handleUpdateMTOServiceItemStatus: jest.fn(),
    handleShowRejectionDialog: jest.fn(),
  };

  it('renders with no details', () => {
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
      <ServiceItemsTable
        {...defaultProps}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        serviceItems={serviceItems}
      />,
    );
    expect(wrapper.find('td').at(1).text()).toBe('—');
  });

  it('renders a thumbnail image with dimensions for item and crating', () => {
    const serviceItems = [
      {
        id: 'abc123',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic Crating',
        code: 'DCRT',
        details: {
          description: 'grandfather clock',
          itemDimensions: { length: 7000, width: 2000, height: 3500 },
          crateDimensions: { length: 10000, width: 2500, height: 5000 },
        },
      },
    ];

    const wrapper = mount(
      <ServiceItemsTable
        {...defaultProps}
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
      />,
    );

    expect(wrapper.find('dt').at(0).text()).toBe('Description:');
    expect(wrapper.find('dd').at(0).text()).toBe('grandfather clock');
    expect(wrapper.find('dt').at(1).text()).toBe('Item size:');
    expect(wrapper.find('dd').at(1).text()).toBe('7"x2"x3.5"');
    expect(wrapper.find('dt').at(2).text()).toBe('Crate size:');
    expect(wrapper.find('dd').at(2).text()).toBe('10"x2.5"x5"');
  });

  it('renders the customer contacts for DDFSIT service item', () => {
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
      <ServiceItemsTable
        {...defaultProps}
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
      />,
    );

    expect(wrapper.find('table').exists()).toBe(true);
    expect(wrapper.find('dt').at(0).text()).toBe('First Customer Contact:');
    expect(wrapper.find('dd').at(0).text()).toBe('0400Z');
    expect(wrapper.find('dt').at(1).text()).toBe('First Available Delivery Date:');
    expect(wrapper.find('dd').at(1).text()).toBe('31 Dec 2020');

    expect(wrapper.find('dt').at(2).text()).toBe('Second Customer Contact:');
    expect(wrapper.find('dd').at(2).text()).toBe('0800Z');
    expect(wrapper.find('dt').at(3).text()).toBe('Second Available Delivery Date:');
    expect(wrapper.find('dd').at(3).text()).toBe('01 Jan 2021');
  });

  it('should render the SITPostalCode ZIP, and reason for DOFSIT service item', () => {
    const serviceItems = [
      {
        id: 'abc123',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic Origin 1st Day SIT',
        code: 'DOFSIT',
        details: {
          pickupPostalCode: '11111',
          SITPostalCode: '12345',
          reason: 'This is the reason',
        },
      },
    ];

    const wrapper = mount(
      <ServiceItemsTable
        {...defaultProps}
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
      />,
    );
    expect(wrapper.find('dt').at(0).contains('ZIP')).toBe(true);
    expect(wrapper.find('dd').at(0).contains('12345')).toBe(true);
    expect(wrapper.find('dt').at(1).contains('Reason')).toBe(true);
    expect(wrapper.find('dd').at(1).contains('This is the reason')).toBe(true);
  });

  it('calls the update service item status handler when the accept button is clicked', () => {
    const serviceItems = [
      {
        id: 'abc123',
        mtoShipmentID: 'xyz789',
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
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        />
      </MockProviders>,
    );

    wrapper.find('button[data-testid="acceptButton"]').simulate('click');

    expect(defaultProps.handleUpdateMTOServiceItemStatus).toHaveBeenCalledWith(
      'abc123',
      'xyz789',
      SERVICE_ITEM_STATUS.APPROVED,
    );
  });

  it('calls the show rejection handler when the reject button is clicked', () => {
    const serviceItems = [
      {
        id: 'abc123',
        mtoShipmentID: 'xyz789',
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
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        />
      </MockProviders>,
    );

    wrapper.find('button[data-testid="rejectButton"]').simulate('click');

    expect(defaultProps.handleShowRejectionDialog).toHaveBeenCalledWith('abc123', 'xyz789');
  });

  it('does not show accept or reject buttons when permissions are missing', () => {
    const serviceItems = [
      {
        id: 'abc123',
        mtoShipmentID: 'xyz789',
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
      <ServiceItemsTable
        {...defaultProps}
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
      />,
    );

    expect(wrapper.find('button[data-testid="acceptButton"]').length).toBeFalsy();
    expect(wrapper.find('button[data-testid="rejectButton"]').length).toBeFalsy();
    expect(wrapper.find('button[data-testid="approveTextButton"]').length).toBeFalsy();
    expect(wrapper.find('button[data-testid="rejectTextButton"]').length).toBeFalsy();
  });
});
