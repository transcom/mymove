import React from 'react';
import { action } from '@storybook/addon-actions';

import ServiceItemContainer from './ServiceItemContainer';

import RequestedServiceItemsTable from 'components/Office/RequestedServiceItemsTable/RequestedServiceItemsTable';
import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';
import { MTO_SERVICE_ITEM_STATUS } from 'shared/constants';

export default {
  title: 'Office Components/ServiceItemContainer',
  component: ServiceItemContainer,
};

const handleUpdateMTOServiceItemStatusMock = action('handleUpdateMTOServiceItemStatus');
const handleShowRejectionDialogMock = action('handleShowRejectionDialog');

const requestedServiceItemInfo = {
  serviceItems: [
    {
      id: 'abc123',
      submittedAt: '2020-11-20',
      approvedAt: '2021-01-24',
      serviceItem: 'Move management',
      code: 'FSC',
      details: {},
    },
  ],
  handleUpdateMTOServiceItemStatus: handleUpdateMTOServiceItemStatusMock,
  handleShowRejectionDialog: handleShowRejectionDialogMock,
  statusForTableType: MTO_SERVICE_ITEM_STATUS.APPROVED,
};

export const MTOServiceItem = () => (
  <MockProviders permissions={[permissionTypes.createShipmentCancellation]}>
    <ServiceItemContainer>
      <RequestedServiceItemsTable {...requestedServiceItemInfo} />
    </ServiceItemContainer>
  </MockProviders>
);
