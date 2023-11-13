/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { fireEvent, render, screen } from '@testing-library/react';

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
    handleShowEditSitEntryDateModal: jest.fn(),
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
      <MockProviders>
        <ServiceItemsTable
          {...defaultProps}
          statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
          serviceItems={serviceItems}
        />
      </MockProviders>,
    );
    expect(wrapper.find('td').at(1).text()).toBe('—');
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

    expect(wrapper.find('dt').at(0).text()).toBe('SIT entry date:');
    expect(wrapper.find('dd').at(0).text()).toBe('31 Dec 2020');
    expect(wrapper.find('dt').at(1).text()).toBe('First available delivery date 1:');
    expect(wrapper.find('dd').at(1).text()).toBe('31 Dec 2020');
    expect(wrapper.find('dt').at(2).text()).toBe('Customer contact attempt 1:');
    expect(wrapper.find('dd').at(2).text()).toBe('31 Dec 2020, 0400Z');

    expect(wrapper.find('dt').at(3).text()).toBe('First available delivery date 2:');
    expect(wrapper.find('dd').at(3).text()).toBe('01 Jan 2021');
    expect(wrapper.find('dt').at(4).text()).toBe('Customer contact attempt 2:');
    expect(wrapper.find('dd').at(4).text()).toBe('01 Jan 2021, 0800Z');
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
    expect(wrapper.find('dt').at(0).contains('SIT entry date')).toBe(true);
    expect(wrapper.find('dd').at(0).contains('25 Dec 2023')).toBe(true);
    expect(wrapper.find('dt').at(1).contains('ZIP')).toBe(true);
    expect(wrapper.find('dd').at(1).contains('12345')).toBe(true);
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

  it('shows update requested tag when service item contains a requested sit address update', () => {
    const serviceItems = [
      {
        id: 'abc123',
        mtoShipmentID: 'xyz789',
        createdAt: '2020-11-20',
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
        createdAt: '2020-11-20',
        approvedAt: '2020-11-21',
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
    expect(wrapper.find('td').at(1).text()).toContain('Date approved: 21 Nov 2020');
  });

  it('properly displays service item table tag for rejected address update', () => {
    const serviceItems = [
      {
        id: 'abc123',
        mtoShipmentID: 'xyz789',
        createdAt: '2020-11-20',
        rejectedAt: '2020-11-21',
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
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
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
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
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
        createdAt: '2020-11-20',
        rejectedAt: '2020-11-21',
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

  jest.mock('react-query', () => ({
    ...jest.requireActual('react-query'),
    useQuery: jest.fn(),
  }));

  jest.mock('hooks/queries', () => ({
    ...jest.requireActual('hooks/queries'),
    useGHCGetMoveHistory: jest.fn(),
  }));

  it('renders a tooltip if resubmitted service item', async () => {
    // eslint-disable-next-line global-require
    const { useQuery } = require('react-query');
    // eslint-disable-next-line global-require
    const { useGHCGetMoveHistory } = require('hooks/queries');

    const history = {
      isLoading: false,
      isError: false,
      queueResult: {
        data: [
          {
            action: 'UPDATE',
            eventName: 'updateMTOServiceItem',
            actionTstampTx: '2022-03-09T15:33:38.579Z',
            changedValues: {
              status: 'SUBMITTED',
            },
            objectId: 'historyDATA',
            oldValues: {
              status: 'REJECTED',
            },
          },
        ],
      },
      isSuccess: true,
    };
    useGHCGetMoveHistory.mockReturnValue(history);

    const mockHistoryObject = {
      queueResult: {
        data: [
          {
            action: 'UPDATE',
            actionTstampClk: '2023-10-31T19:16:51.304Z',
            actionTstampStm: '2023-10-31T19:16:51.301Z',
            actionTstampTx: '2023-10-31T19:16:51.281Z',
            changedValues: {
              reason: 'Test reason 82',
              sit_departure_date: null,
              status: 'SUBMITTED',
            },
            context: [
              {
                name: 'Domestic origin 1st day SIT',
                shipment_id_abbr: '659a9',
                shipment_type: 'HHG',
              },
            ],
            eventName: 'updateMTOServiceItem',
            id: '2924e0d3-d63a-4586-ac1a-2f6a633d0ae8',
            objectId: '6eba127d-cbc8-4bf7-b276-b174a71a91cd',
            oldValues: {
              actual_weight: null,
              approved_at: null,
              description: null,
              estimated_weight: null,
              id: '6eba127d-cbc8-4bf7-b276-b174a71a91cd',
              pickup_postal_code: null,
              reason: 'Test reason 81',
              rejected_at: null,
              rejection_reason: null,
              requested_approvals_requested_status: true,
              sit_customer_contacted: null,
              sit_departure_date: '2023-08-31',
              sit_destination_original_address_id: null,
              sit_entry_date: '2023-08-31',
              sit_postal_code: '90211',
              sit_requested_delivery: '2019-08-31',
              status: 'REJECTED',
            },
            relId: 19112,
            schemaName: 'public',
            sessionUserId: '3ce06fa9-590a-48e5-9e30-6ad1e82b528c',
            tableName: 'mto_service_items',
            transactionId: 1478,
          },
          {
            action: 'UPDATE',
            actionTstampClk: '2023-10-31T19:16:36.041Z',
            actionTstampStm: '2023-10-31T19:16:36.041Z',
            actionTstampTx: '2023-10-31T19:16:36.040Z',
            changedValues: {
              sit_departure_date: '2023-08-31',
              status: 'REJECTED',
            },
            context: [
              {
                name: 'Domestic origin 1st day SIT',
                shipment_id_abbr: '659a9',
                shipment_type: 'HHG',
              },
            ],
            id: '60fc09ca-f11f-458a-a40c-e5ba347b68c2',
            objectId: '6eba127d-cbc8-4bf7-b276-b174a71a91cd',
            oldValues: {
              actual_weight: null,
              approved_at: null,
              description: null,
              estimated_weight: null,
              id: '6eba127d-cbc8-4bf7-b276-b174a71a91cd',
              pickup_postal_code: null,
              reason: 'Test reason 81',
              rejected_at: null,
              rejection_reason: null,
              requested_approvals_requested_status: true,
              sit_customer_contacted: null,
              sit_departure_date: null,
              sit_destination_original_address_id: null,
              sit_entry_date: '2023-08-31',
              sit_postal_code: '90211',
              sit_requested_delivery: '2019-08-31',
              status: 'SUBMITTED',
            },
            relId: 19112,
            schemaName: 'public',
            tableName: 'mto_service_items',
            transactionId: 1477,
          },
        ],
        id: '78f7f149-0c00-4b20-83ef-ce6aabadeaef',
        locator: 'SITEXT',
        page: 1,
        perPage: 20,
        referenceId: '1455-1734',
        totalCount: 46,
      },
      isLoading: false,
      isError: false,
      isSuccess: true,
    };
    useQuery.mockReturnValue(mockHistoryObject);

    const serviceItems = [
      {
        id: 'abc123',
        actionTstampTx: '2022-03-09T15:33:38.579Z',
        createdAt: '2023-11-5',
        serviceItem: 'Domestic origin 1st day SIT',
        code: 'DOFSIT',
        reServiceCode: 'DOFSIT',
        reServiceID: '987654321',
        reServiceName: 'Domestic origin 1st day SIT',
        mtoShipmentID: '123456789',
        status: 'SUBMITTED',
        sitRequestedDelivery: '2019-08-31',
        moveTaskOrderID: '78f7f149-0c00-4b20-83ef-ce6aabadeaef',
        SITPostalCode: '90211',
        details: {
          SITPostalCode: '90211',
          reason: 'Test reason 82',
          sitEntryDate: '2020-12-31',
        },
        sitAddressUpdates: [
          {
            contractorRemarks: 'contractor remarks',
            distance: 140,
            status: SIT_ADDRESS_UPDATE_STATUS.REQUESTED,
            officeRemarks: null,
          },
        ],
      },
      {
        id: '123456',
        actionTstampTx: '2022-03-09T15:33:38.579Z',
        createdAt: '2023-11-5',
        serviceItem: 'Domestic origin 1st day SIT',
        code: 'DOFSIT',
        reServiceCode: 'DOFSIT',
        reServiceID: '987654321',
        reServiceName: 'Domestic origin 1st day SIT',
        mtoShipmentID: '123456789',
        status: 'SUBMITTED',
        sitRequestedDelivery: '2019-08-31',
        moveTaskOrderID: '78f7f149-0c00-4b20-83ef-ce6aabadeaef',
        SITPostalCode: '90211',
        details: {
          SITPostalCode: '90211',
          reason: 'Test reason 82',
          sitEntryDate: '2020-12-31',
        },
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

    // TESTING WITH RENDER
    // IF MOCK DATA IS GOOD, TESTS WORK
    const { container } = render(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        />
      </MockProviders>,
    );

    screen.debug();
    const tooltipContainer = container.querySelector('.tooltipContainer');
    expect(tooltipContainer).toBeInTheDocument();
    fireEvent.mouseEnter(tooltipContainer);
    const tooltipText = tooltipContainer.textContent;
    expect(tooltipText).toContain('Previous:');
    // TESTING WITH RENDER
    // IF MOCK DATA IS GOOD, TESTS WORK

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        />
      </MockProviders>,
    );

    // WHEN MOCK DATA IS GOOD, THIS TEST WORKS
    // find the first instance of the tooltip
    // console.log(wrapper.debug());
    const toolTip = wrapper.find('ToolTip').at(0);
    expect(toolTip.exists()).toBe(true);
    expect(toolTip.props().text).toContain('Status\nNew: SUBMITTED \nPrevious: REJECTED\n\n');
    // END WHEN MOCK DATA IS GOOD, THIS TEST WORKS

    // TESTING OF MOCK FUNCTIONS
    // when the tooltip works, this is what it looks like:
    // <ToolTip text="Status\nNew: SUBMITTED \nPrevious: REJECTED\n\n" position="bottom" color="#0050d8">
  });
});
