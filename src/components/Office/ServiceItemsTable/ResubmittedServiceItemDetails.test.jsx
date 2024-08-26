/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

// eslint-disable-next-line import/extensions
import { multiplePaymentRequests } from './resubmittedServiceItemUnitTestData.js';
import ServiceItemsTable from './ServiceItemsTable';

import { useMovePaymentRequestsQueries, useGHCGetMoveHistory } from 'hooks/queries';
import { SERVICE_ITEM_STATUS } from 'shared/constants';
import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

jest.mock('hooks/queries', () => ({
  useMovePaymentRequestsQueries: jest.fn(),
  useGHCGetMoveHistory: jest.fn(),
}));

describe('ServiceItemsTable', () => {
  it('renders a tooltip with old details if resubmitted service item', () => {
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
              reason: 'New reason in test code',
              status: 'SUBMITTED',
              id: 'abc12345',
              pickup_postal_code: '54321',
              sit_entry_date: '2023-01-01',
              sit_postal_code: '09876',
            },
            objectId: 'historyObjectInServiceItemsTableTest',
            oldValues: {
              reason: 'Old reason in test code',
              status: 'REJECTED',
              id: 'def67890',
              pickup_postal_code: '12345',
              sit_entry_date: '2022-12-12',
              sit_postal_code: '67890',
            },
          },
        ],
      },
      isSuccess: true,
    };

    useMovePaymentRequestsQueries.mockReturnValue(multiplePaymentRequests);
    useGHCGetMoveHistory.mockReturnValue(history);

    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        />
      </MockProviders>,
    );

    const toolTip = wrapper.find('ToolTip').at(0);
    expect(toolTip.exists()).toBe(true);
    let resultString = 'Reason\nNew: New reason in test code \nPrevious: Old reason in test code\n\n';
    resultString += 'Status\nNew: SUBMITTED \nPrevious: REJECTED\n\n';
    resultString += 'ID\nNew: abc12345 \nPrevious: def67890\n\n';
    resultString += 'Pickup Postal Code\nNew: 54321 \nPrevious: 12345\n\n';
    resultString += 'SIT Entry Date\nNew: 2023-01-01 \nPrevious: 2022-12-12\n\n';
    resultString += 'SIT Postal Code\nNew: 09876 \nPrevious: 67890\n\n';
    expect(toolTip.props().text).toBe(resultString);
  });

  it('does not render a tooltip for a service item that has not been resubmitted', () => {
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
    const history = {
      isLoading: false,
      isError: false,
      queueResult: {
        data: [
          {
            action: 'INSERT',
            eventName: 'updateMTOServiceItem',
            actionTstampTx: '2022-03-09T15:33:38.579Z',
            changedValues: {
              reason: 'New reason in test code',
              status: 'SUBMITTED',
              id: 'abc12345',
              pickup_postal_code: '54321',
              sit_entry_date: '2023-01-01',
              sit_postal_code: '09876',
            },
            objectId: 'historyObjectInServiceItemsTableTest',
          },
        ],
      },
      isSuccess: true,
    };

    useMovePaymentRequestsQueries.mockReturnValue(multiplePaymentRequests);
    useGHCGetMoveHistory.mockReturnValue(history);

    const approvedWrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.APPROVED}
        />
      </MockProviders>,
    );

    const toolTipAccepted = approvedWrapper.find('ToolTip').at(0);
    expect(toolTipAccepted.exists()).toBe(false);

    const rejectedWrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
        <ServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.REJECTED}
        />
      </MockProviders>,
    );

    const toolTipRejected = rejectedWrapper.find('ToolTip').at(0);
    expect(toolTipRejected.exists()).toBe(false);
  });
});
