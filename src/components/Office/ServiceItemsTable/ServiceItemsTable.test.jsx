/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import ServiceItemsTable from './ServiceItemsTable';

import { SERVICE_ITEM_STATUS } from 'shared/constants';
import { SIT_ADDRESS_UPDATE_STATUS } from 'constants/sitUpdates';
import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

describe('ServiceItemsTable', () => {
  const defaultProps = {
    handleUpdateMTOServiceItemStatus: jest.fn(),
    handleShowRejectionDialog: jest.fn(),
    handleRequestSITAddressUpdateModal: jest.fn(),
    handleShowEditSitAddressModal: jest.fn(),
    serviceItemAddressUpdateAlert: {
      makeVisible: false,
      alertMessage: '',
      alertType: '',
    },
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
          customerContacts: [
            { timeMilitary: '0400Z', firstAvailableDeliveryDate: '2020-12-31', dateOfContact: '2020-12-31' },
            { timeMilitary: '0800Z', firstAvailableDeliveryDate: '2021-01-01', dateOfContact: '2021-01-01' },
          ],
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
    expect(wrapper.find('dt').at(0).text()).toBe('First available delivery date 1:');
    expect(wrapper.find('dd').at(0).text()).toBe('31 Dec 2020');
    expect(wrapper.find('dt').at(1).text()).toBe('Customer contact attempt 1:');
    expect(wrapper.find('dd').at(1).text()).toBe('31 Dec 2020, 0400Z');

    expect(wrapper.find('dt').at(2).text()).toBe('First available delivery date 2:');
    expect(wrapper.find('dd').at(2).text()).toBe('01 Jan 2021');
    expect(wrapper.find('dt').at(3).text()).toBe('Customer contact attempt 2:');
    expect(wrapper.find('dd').at(3).text()).toBe('01 Jan 2021, 0800Z');
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

  it('shows update requested tag when service item contains a requested sit address update', () => {
    const serviceItems = [
      {
        id: 'abc123',
        mtoShipmentID: 'xyz789',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic destination SIT delivery',
        code: 'DDDSIT',
        sitAddressUpdates: [
          {
            contractorRemarks: 'contractor remarks',
            distance: 140,
            status: SIT_ADDRESS_UPDATE_STATUS.REQUESTED,
            officeRemarks: null,
          },
        ],
      },
    ];

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.APPROVED}
        />
      </MockProviders>,
    );

    expect(wrapper.find('[data-testid="sitAddressUpdateTag"]').exists()).toBe(true);
  });

  it('properly displays service item table tag for approved address update', () => {
    const serviceItems = [
      {
        id: 'abc123',
        mtoShipmentID: 'xyz789',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic destination SIT delivery',
        code: 'DDDSIT',
        sitAddressUpdates: [
          {
            contractorRemarks: 'contractor remarks',
            distance: 140,
            status: SIT_ADDRESS_UPDATE_STATUS.APPROVED,
            officeRemarks: 'I have approved',
          },
        ],
      },
    ];

    const propsForApprovedUpdate = {
      handleUpdateMTOServiceItemStatus: jest.fn(),
      handleShowRejectionDialog: jest.fn(),
      serviceItemAddressUpdateAlert: {
        makeVisible: true,
        alertMessage: 'warning',
        alertType: 'Address update over 50 miles approved.',
      },
    };

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
        <ServiceItemsTable
          {...propsForApprovedUpdate}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.APPROVED}
        />
      </MockProviders>,
    );

    expect(wrapper.find('[data-testid="serviceItemAddressUpdateAlert"]').exists()).toBe(true);
  });

  it('properly displays service item table tag for rejected address update', () => {
    const serviceItems = [
      {
        id: 'abc123',
        mtoShipmentID: 'xyz789',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic destination SIT delivery',
        code: 'DDDSIT',
        sitAddressUpdates: [
          {
            contractorRemarks: 'contractor remarks',
            distance: 140,
            status: SIT_ADDRESS_UPDATE_STATUS.REJECTED,
            officeRemarks: 'I have rejected',
          },
        ],
      },
    ];

    const propsForApprovedUpdate = {
      handleUpdateMTOServiceItemStatus: jest.fn(),
      handleShowRejectionDialog: jest.fn(),
      serviceItemAddressUpdateAlert: {
        makeVisible: true,
        alertMessage: 'info',
        alertType: 'Address update over 50 miles rejected.',
      },
    };

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
        <ServiceItemsTable
          {...propsForApprovedUpdate}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.APPROVED}
        />
      </MockProviders>,
    );

    expect(wrapper.find('[data-testid="serviceItemAddressUpdateAlert"]').exists()).toBe(true);
  });

  it('properly displays service item table tag for edited address update', () => {
    const serviceItems = [
      {
        id: 'abc123',
        mtoShipmentID: 'xyz789',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic destination SIT delivery',
        code: 'DDDSIT',
        sitAddressUpdates: [
          {
            contractorRemarks: null,
            distance: 49,
            status: SIT_ADDRESS_UPDATE_STATUS.APPROVED,
            officeRemarks: 'I have edited',
          },
        ],
      },
    ];

    const propsForApprovedUpdate = {
      handleUpdateMTOServiceItemStatus: jest.fn(),
      handleShowRejectionDialog: jest.fn(),
      serviceItemAddressUpdateAlert: {
        makeVisible: true,
        alertMessage: 'info',
        alertType: 'Address update within 50 miles.',
      },
    };

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
        <ServiceItemsTable
          {...propsForApprovedUpdate}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.APPROVED}
        />
      </MockProviders>,
    );

    expect(wrapper.find('[data-testid="serviceItemAddressUpdateAlert"]').exists()).toBe(true);
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
      <ServiceItemsTable
        {...defaultProps}
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.APPROVED}
      />,
    );

    expect(wrapper.find('button[data-testid="editTextButton"]').length).toBeFalsy();
  });

  it('shows edit button when service item does not have requested sit address updates', () => {
    const serviceItems = [
      {
        id: 'abc123',
        mtoShipmentID: 'xyz789',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic destination SIT delivery',
        code: 'DDDSIT',
        sitAddressUpdates: [
          {
            contractorRemarks: 'contractor remarks',
            distance: 140,
            status: SIT_ADDRESS_UPDATE_STATUS.APPROVED,
            officeRemarks: null,
          },
        ],
      },
    ];

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.APPROVED}
        />
      </MockProviders>,
    );

    const editButton = wrapper.find('button[data-testid="editTextButton"]');

    expect(editButton.length).toBeTruthy();

    expect(editButton.at(0).contains('Edit')).toBe(true);
  });

  it('calls the handleShowEditSitAddressModal handler when the edit button is clicked', () => {
    const serviceItems = [
      {
        id: 'abc123',
        mtoShipmentID: 'xyz789',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic destination SIT delivery',
        code: 'DDDSIT',
        sitAddressUpdates: [
          {
            contractorRemarks: 'contractor remarks',
            distance: 140,
            status: SIT_ADDRESS_UPDATE_STATUS.APPROVED,
            officeRemarks: null,
          },
        ],
      },
    ];

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.APPROVED}
        />
      </MockProviders>,
    );

    wrapper.find('button[data-testid="editTextButton"]').simulate('click');

    expect(defaultProps.handleShowEditSitAddressModal).toHaveBeenCalledWith('abc123', 'xyz789');
  });

  it('shows review request button when service item contains requested sit address update', () => {
    const serviceItems = [
      {
        id: 'abc123',
        mtoShipmentID: 'xyz789',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic destination SIT delivery',
        code: 'DDDSIT',
        sitAddressUpdates: [
          {
            contractorRemarks: 'contractor remarks',
            distance: 140,
            status: SIT_ADDRESS_UPDATE_STATUS.REQUESTED,
            officeRemarks: null,
          },
        ],
      },
    ];

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.APPROVED}
        />
      </MockProviders>,
    );

    const reviewRequestButton = wrapper.find('button[data-testid="reviewRequestTextButton"]');

    expect(reviewRequestButton.length).toBeTruthy();

    expect(reviewRequestButton.at(0).contains('Review Request')).toBe(true);
  });

  it('calls the handleRequestSITAddressUpdateModal handler when the Review Request button is clicked', () => {
    const serviceItems = [
      {
        id: 'abc123',
        mtoShipmentID: 'xyz789',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic destination SIT delivery',
        code: 'DDDSIT',
        sitAddressUpdates: [
          {
            contractorRemarks: 'contractor remarks',
            distance: 140,
            status: SIT_ADDRESS_UPDATE_STATUS.REQUESTED,
            officeRemarks: null,
          },
        ],
      },
    ];

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.APPROVED}
        />
      </MockProviders>,
    );

    wrapper.find('button[data-testid="reviewRequestTextButton"]').simulate('click');

    expect(defaultProps.handleRequestSITAddressUpdateModal).toHaveBeenCalledWith('abc123', 'xyz789');
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
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
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
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
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
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic destination SIT delivery',
        code: 'DDDSIT',
      },
    ];

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
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
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
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
});
