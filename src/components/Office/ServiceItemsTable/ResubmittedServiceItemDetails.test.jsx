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
      handleRequestSITAddressUpdateModal: jest.fn(),
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
            },
            objectId: 'historyObjectInServiceItemsTableTest',
            oldValues: {
              reason: 'Old reason in test code',
              status: 'REJECTED',
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
    expect(toolTip.props().text).toBe(resultString);
  });
});
