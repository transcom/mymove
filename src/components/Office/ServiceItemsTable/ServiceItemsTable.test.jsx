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
    handleShowEditSitAddressModal: jest.fn(),
    handleShowEditSitEntryDateModal: jest.fn(),
    serviceItemAddressUpdateAlert: {
      makeVisible: false,
      alertMessage: '',
      alertType: '',
    },
  };

  it('renders with no estimated price when no estimated price exists', () => {
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
      <MockProviders>
        <ServiceItemsTable
          {...defaultProps}
          statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
          serviceItems={serviceItems}
        />
      </MockProviders>,
    );
    expect(wrapper.find('td').at(1).text()).toBe('Estimated Price: -');
  });

  it('renders with estimated price shown when estimated price', () => {
    const serviceItems = [
      {
        id: 'abc123',
        submittedAt: '2020-11-20',
        serviceItem: 'Fuel Surcharge',
        code: 'FSC',
        details: {
          estimatedPrice: 2314,
        },
      },
    ];
    const wrapper = mount(
      <MockProviders>
        <ServiceItemsTable
          {...defaultProps}
          statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
          serviceItems={serviceItems}
        />
      </MockProviders>,
    );
    expect(wrapper.find('td').at(1).text()).toBe('Estimated Price: $23.14');
  });

  it('renders a thumbnail image with dimensions for item and crating', () => {
    const serviceItems = [
      {
        id: 'abc123',
        createdAt: '2020-11-20',
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
      <MockProviders>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        />
      </MockProviders>,
    );

    expect(wrapper.find('td').at(0).text()).toContain('Date requested: 20 Nov 2020');
    expect(wrapper.find('dt').at(0).text()).toBe('Description:');
    expect(wrapper.find('dd').at(0).text()).toBe('grandfather clock');
    expect(wrapper.find('dt').at(1).text()).toBe('Item size:');
    expect(wrapper.find('dd').at(1).text()).toBe('7"x2"x3.5"');
    expect(wrapper.find('dt').at(2).text()).toBe('Crate size:');
    expect(wrapper.find('dd').at(2).text()).toBe('10"x2.5"x5"');
  });

  it('renders details for international crating (ICRT)', () => {
    const serviceItems = [
      {
        id: 'abc123',
        createdAt: '2020-11-20',
        serviceItem: 'International Crating',
        code: 'ICRT',
        details: {
          description: 'grandfather clock',
          itemDimensions: { length: 7000, width: 2000, height: 3500 },
          crateDimensions: { length: 10000, width: 2500, height: 5000 },
          market: 'OCONUS',
          externalCrate: true,
          standaloneCrate: true,
        },
      },
    ];

    const wrapper = mount(
      <MockProviders>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        />
      </MockProviders>,
    );

    expect(wrapper.find('td').at(0).text()).toContain('International Crating - Standalone');
    expect(wrapper.find('td').at(0).text()).toContain('Date requested: 20 Nov 2020');
    expect(wrapper.find('dt').at(0).text()).toBe('Description:');
    expect(wrapper.find('dd').at(0).text()).toBe('grandfather clock');
    expect(wrapper.find('dt').at(1).text()).toBe('Item size:');
    expect(wrapper.find('dd').at(1).text()).toBe('7"x2"x3.5"');
    expect(wrapper.find('dt').at(2).text()).toBe('Crate size:');
    expect(wrapper.find('dd').at(2).text()).toBe('10"x2.5"x5"');
    expect(wrapper.find('dt').at(3).text()).toBe('External crate:');
    expect(wrapper.find('dd').at(3).text()).toBe('Yes');
    expect(wrapper.find('dt').at(4).text()).toBe('Market:');
    expect(wrapper.find('dd').at(4).text()).toBe('OCONUS');
    expect(wrapper.find('dt').at(5).text()).toBe('Reason:');
    expect(wrapper.find('dd').at(5).text()).toBe('-');
  });

  it('renders details for international crating (IUCRT)', () => {
    const serviceItems = [
      {
        id: 'abc123',
        createdAt: '2020-11-20',
        serviceItem: 'International Crating',
        code: 'ICRT',
        details: {
          description: 'grandfather clock',
          itemDimensions: { length: 7000, width: 2000, height: 3500 },
          crateDimensions: { length: 10000, width: 2500, height: 5000 },
          market: 'OCONUS',
          externalCrate: null,
        },
      },
    ];

    const wrapper = mount(
      <MockProviders>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        />
      </MockProviders>,
    );

    expect(wrapper.find('td').at(0).text()).toContain('Date requested: 20 Nov 2020');
    expect(wrapper.find('dt').at(0).text()).toBe('Description:');
    expect(wrapper.find('dd').at(0).text()).toBe('grandfather clock');
    expect(wrapper.find('dt').at(1).text()).toBe('Item size:');
    expect(wrapper.find('dd').at(1).text()).toBe('7"x2"x3.5"');
    expect(wrapper.find('dt').at(2).text()).toBe('Crate size:');
    expect(wrapper.find('dd').at(2).text()).toBe('10"x2.5"x5"');
    expect(wrapper.find('dt').at(3).text()).toBe('Market:');
    expect(wrapper.find('dd').at(3).text()).toBe('OCONUS');
    expect(wrapper.find('dt').at(4).text()).toBe('Reason:');
    expect(wrapper.find('dd').at(4).text()).toBe('-');
  });

  it('renders with authorized price for MS item', () => {
    const serviceItems = [
      {
        id: 'abc123',
        createdAt: '2020-11-20',
        serviceItem: 'Domestic Crating',
        code: 'MS',
        details: {
          estimatedPrice: 100000,
        },
      },
    ];

    const wrapper = mount(
      <MockProviders>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        />
      </MockProviders>,
    );

    expect(wrapper.find('td').at(0).text()).toContain('Date requested: 20 Nov 2020');
    expect(wrapper.find('dt').at(0).text()).toBe('Price:');
    expect(wrapper.find('dd').at(0).text()).toBe('$1,000.00');
  });

  it('renders with authorized price for CS item', () => {
    const serviceItems = [
      {
        id: 'abc123',
        createdAt: '2020-11-20',
        serviceItem: 'Domestic Crating',
        code: 'CS',
        details: {
          estimatedPrice: 100000,
        },
      },
    ];

    const wrapper = mount(
      <MockProviders>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        />
      </MockProviders>,
    );

    expect(wrapper.find('td').at(0).text()).toContain('Date requested: 20 Nov 2020');
    expect(wrapper.find('dt').at(0).text()).toBe('Price:');
    expect(wrapper.find('dd').at(0).text()).toBe('$1,000.00');
  });

  it('renders the customer contacts for DDFSIT service item', () => {
    const serviceItems = [
      {
        id: 'abc123',
        createdAt: '2020-11-20',
        serviceItem: 'Domestic Crating',
        code: 'DDFSIT',
        details: {
          sitEntryDate: '2020-12-31',
          customerContacts: [
            {
              timeMilitary: '0400Z',
              firstAvailableDeliveryDate: '2020-12-31',
              dateOfContact: '2020-12-31',
            },
            { timeMilitary: '0800Z', firstAvailableDeliveryDate: '2021-01-01', dateOfContact: '2021-01-01' },
          ],
          sitDestinationOriginalAddress: {
            city: 'Destination Original Tampa',
            eTag: 'MjAyNC0wMy0xMlQxOTo1OTowOC41NjkxMzla',
            id: '7fd6cb90-54cd-44d8-8735-102e28734d84',
            postalCode: '33621',
            state: 'FL',
            streetAddress1: 'MacDill',
          },
        },
      },
    ];

    const wrapper = mount(
      <MockProviders>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        />
      </MockProviders>,
    );

    expect(wrapper.find('table').exists()).toBe(true);
    expect(wrapper.find('dt').at(0).text()).toBe('Original Delivery Address:');
    expect(wrapper.find('dd').at(0).text()).toBe('Destination Original Tampa, FL 33621');

    expect(wrapper.find('dt').at(1).text()).toBe('SIT entry date:');
    expect(wrapper.find('dd').at(1).text()).toBe('31 Dec 2020');

    expect(wrapper.find('dt').at(2).text()).toBe('First available delivery date 1:');
    expect(wrapper.find('dd').at(2).text()).toBe('31 Dec 2020');
    expect(wrapper.find('dt').at(3).text()).toBe('Customer contact attempt 1:');
    expect(wrapper.find('dd').at(3).text()).toBe('31 Dec 2020, 0400Z');

    expect(wrapper.find('dt').at(4).text()).toBe('First available delivery date 2:');
    expect(wrapper.find('dd').at(4).text()).toBe('01 Jan 2021');
    expect(wrapper.find('dt').at(5).text()).toBe('Customer contact attempt 2:');
    expect(wrapper.find('dd').at(5).text()).toBe('01 Jan 2021, 0800Z');
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
          sitEntryDate: '2023-12-25T00:00:00.000Z',
          sitOriginHHGOriginalAddress: {
            city: 'Origin Original Tampa',
            eTag: 'MjAyNC0wMy0xMlQxOTo1OTowOC41NjkxMzla',
            id: '7fd6cb90-54cd-44d8-8735-102e28734d84',
            postalCode: '33621',
            state: 'FL',
            streetAddress1: 'MacDill',
          },
        },
      },
    ];

    const wrapper = mount(
      <MockProviders>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        />
      </MockProviders>,
    );
    expect(wrapper.find('dt').at(0).contains('Original Pickup Address')).toBe(true);
    expect(wrapper.find('dd').at(0).contains('Origin Original Tampa, FL 33621')).toBe(true);

    expect(wrapper.find('dt').at(1).contains('SIT entry date')).toBe(true);
    expect(wrapper.find('dd').at(1).contains('25 Dec 2023')).toBe(true);
    expect(wrapper.find('dt').at(2).contains('Reason')).toBe(true);
    expect(wrapper.find('dd').at(2).contains('This is the reason')).toBe(true);
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
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem, permissionTypes.updateMTOPage]}>
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
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem, permissionTypes.updateMTOPage]}>
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
      <MockProviders>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        />
      </MockProviders>,
    );

    expect(wrapper.find('button[data-testid="acceptButton"]').length).toBeFalsy();
    expect(wrapper.find('button[data-testid="rejectButton"]').length).toBeFalsy();
    expect(wrapper.find('button[data-testid="approveTextButton"]').length).toBeFalsy();
    expect(wrapper.find('button[data-testid="rejectTextButton"]').length).toBeFalsy();
  });

  it('does not show accept or reject buttons when updateMTOPage permission is missing', () => {
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

    expect(wrapper.find('button[data-testid="acceptButton"]').length).toBeFalsy();
    expect(wrapper.find('button[data-testid="rejectButton"]').length).toBeFalsy();
    expect(wrapper.find('button[data-testid="approveTextButton"]').length).toBeFalsy();
    expect(wrapper.find('button[data-testid="rejectTextButton"]').length).toBeFalsy();
  });

  it('does not show accept button when DSH is rejected as a result of delivery address change', () => {
    const serviceItems = [
      {
        id: 'dsh123',
        mtoShipmentID: 'xyz789',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic shorthaul',
        code: 'DSH',
        details: {
          rejectionReason:
            'Automatically rejected due to change in delivery address affecting the ZIP code qualification for short haul / line haul.',
        },
      },
    ];

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem, permissionTypes.updateMTOPage]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.REJECTED}
        />
      </MockProviders>,
    );

    const approveTextButton = wrapper.find('button[data-testid="approveTextButton"]');

    expect(approveTextButton.length).toBeFalsy();

    expect(approveTextButton.at(0).find('svg[data-icon="check"]').length).toBe(0);
    expect(approveTextButton.at(0).contains('Approve')).toBe(false);
  });

  it('does not show accept button when DLH is rejected as a result of delivery address change', () => {
    const serviceItems = [
      {
        id: 'dlh123',
        mtoShipmentID: 'xyz789',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic linehaul',
        code: 'DLH',
        details: {
          rejectionReason:
            'Automatically rejected due to change in delivery address affecting the ZIP code qualification for short haul / line haul.',
        },
      },
    ];

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem, permissionTypes.updateMTOPage]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.REJECTED}
        />
      </MockProviders>,
    );

    const approveTextButton = wrapper.find('button[data-testid="approveTextButton"]');

    expect(approveTextButton.length).toBeFalsy();

    expect(approveTextButton.at(0).find('svg[data-icon="check"]').length).toBe(0);
    expect(approveTextButton.at(0).contains('Approve')).toBe(false);
  });

  it('shows accept button when DSH is rejected but NOT as a result of delivery address change', () => {
    const serviceItems = [
      {
        id: 'dsh123',
        mtoShipmentID: 'xyz789',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic shorthaul',
        code: 'DSH',
        details: {
          rejectionReason:
            'Any reason other than "Automatically rejected due to change in delivery address affecting the ZIP code qualification for short haul / line haul."',
        },
      },
    ];

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem, permissionTypes.updateMTOPage]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.REJECTED}
        />
      </MockProviders>,
    );

    const approveTextButton = wrapper.find('button[data-testid="approveTextButton"]');

    expect(approveTextButton.length).toBeTruthy();

    expect(approveTextButton.at(0).find('svg[data-icon="check"]').length).toBe(1);
    expect(approveTextButton.at(0).contains('Approve')).toBe(true);
  });

  it('shows accept button when DLH is rejected but NOT as a result of delivery address change', () => {
    const serviceItems = [
      {
        id: 'dlh123',
        mtoShipmentID: 'xyz789',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic linehaul',
        code: 'DLH',
        details: {
          rejectionReason:
            'Any reason other than "Automatically rejected due to change in delivery address affecting the ZIP code qualification for short haul / line haul."',
        },
      },
    ];

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem, permissionTypes.updateMTOPage]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.REJECTED}
        />
      </MockProviders>,
    );

    const approveTextButton = wrapper.find('button[data-testid="approveTextButton"]');

    expect(approveTextButton.length).toBeTruthy();

    expect(approveTextButton.at(0).find('svg[data-icon="check"]').length).toBe(1);
    expect(approveTextButton.at(0).contains('Approve')).toBe(true);
  });

  it('does not show edit/review request button when service item code is not DDDSIT', () => {
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
      <MockProviders>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.APPROVED}
        />
      </MockProviders>,
    );

    expect(wrapper.find('button[data-testid="editTextButton"]').length).toBeFalsy();
  });

  it('calls the handleShowEditSitEntryDateModal handler when the edit button is clicked for DOFSIT service item', () => {
    const serviceItems = [
      {
        id: 'abc123',
        mtoShipmentID: 'xyz789',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic origin 1st day SIT',
        code: 'DOFSIT',
      },
    ];

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem, permissionTypes.updateMTOPage]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.APPROVED}
        />
      </MockProviders>,
    );

    wrapper.find('button[data-testid="editTextButton"]').simulate('click');

    expect(defaultProps.handleShowEditSitEntryDateModal).toHaveBeenCalledWith('abc123', 'xyz789');
  });

  it('calls the handleShowEditSitEntryDateModal handler when the edit button is clicked for DDFSIT service item', () => {
    const serviceItems = [
      {
        id: 'abc123',
        mtoShipmentID: 'xyz789',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic destination 1st day SIT',
        code: 'DDFSIT',
      },
    ];

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem, permissionTypes.updateMTOPage]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.APPROVED}
        />
      </MockProviders>,
    );

    wrapper.find('button[data-testid="editTextButton"]').simulate('click');

    expect(defaultProps.handleShowEditSitEntryDateModal).toHaveBeenCalledWith('abc123', 'xyz789');
  });

  it('shows Reject text button for approved service item', () => {
    const serviceItems = [
      {
        id: 'abc123',
        mtoShipmentID: 'xyz789',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic destination SIT delivery',
        code: 'DDDSIT',
      },
    ];

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem, permissionTypes.updateMTOPage]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.APPROVED}
        />
      </MockProviders>,
    );

    const rejectTextButton = wrapper.find('button[data-testid="rejectTextButton"]');

    expect(rejectTextButton.length).toBeTruthy();

    expect(rejectTextButton.at(0).contains('Reject')).toBe(true);
  });

  it('calls the handleShowRejectionDialog handler when the Reject text button is clicked', () => {
    const serviceItems = [
      {
        id: 'abc123',
        mtoShipmentID: 'xyz789',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic destination SIT delivery',
        code: 'DDDSIT',
      },
    ];

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem, permissionTypes.updateMTOPage]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.APPROVED}
        />
      </MockProviders>,
    );

    wrapper.find('button[data-testid="rejectTextButton"]').simulate('click');

    expect(defaultProps.handleShowRejectionDialog).toHaveBeenCalledWith('abc123', 'xyz789');
  });

  it('shows Approve text button for rejected service item', () => {
    const serviceItems = [
      {
        id: 'abc123',
        mtoShipmentID: 'xyz789',
        createdAt: '2020-11-20',
        rejectedAt: '2020-11-21',
        serviceItem: 'Domestic destination SIT delivery',
        code: 'DDDSIT',
      },
    ];

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem, permissionTypes.updateMTOPage]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.REJECTED}
        />
      </MockProviders>,
    );

    const approveTextButton = wrapper.find('button[data-testid="approveTextButton"]');

    expect(approveTextButton.length).toBeTruthy();

    expect(approveTextButton.at(0).find('svg[data-icon="check"]').length).toBe(1);
    expect(approveTextButton.at(0).contains('Approve')).toBe(true);
    expect(wrapper.find('td').at(0).text()).toContain('Date rejected: 21 Nov 2020');
  });

  it('calls the handleUpdateMTOServiceItemStatus handler when the Approve text button is clicked', () => {
    const serviceItems = [
      {
        id: 'abc123',
        mtoShipmentID: 'xyz789',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic destination SIT delivery',
        code: 'DDDSIT',
      },
    ];

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem, permissionTypes.updateMTOPage]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.REJECTED}
        />
      </MockProviders>,
    );

    wrapper.find('button[data-testid="approveTextButton"]').simulate('click');

    expect(defaultProps.handleUpdateMTOServiceItemStatus).toHaveBeenCalledWith('abc123', 'xyz789', 'APPROVED');
  });

  it('disables the accept button when the move is locked', () => {
    const serviceItems = [
      {
        id: 'dlh123',
        mtoShipmentID: 'xyz789',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic linehaul',
        code: 'DLH',
        details: {
          rejectionReason:
            'Any reason other than "Automatically rejected due to change in delivery address affecting the ZIP code qualification for short haul / line haul."',
        },
      },
    ];

    const isMoveLocked = true;

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem, permissionTypes.updateMTOPage]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.REJECTED}
          isMoveLocked={isMoveLocked}
        />
      </MockProviders>,
    );

    // approve button shows up but is disabled
    const approveTextButton = wrapper.find('button[data-testid="approveTextButton"]');
    expect(approveTextButton.length).toBeTruthy();
    expect(wrapper.find('button[data-testid="approveTextButton"]').prop('disabled')).toBe(true);
  });
});
