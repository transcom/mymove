import React from 'react';
import { render, screen } from '@testing-library/react';

import RequestedServiceItemsTable from '../RequestedServiceItemsTable/RequestedServiceItemsTable';

import ServiceItemContainer from './ServiceItemContainer';

import { MockProviders } from 'testUtils';
import { MTO_SERVICE_ITEM_STATUS } from 'shared/constants';

const requestedServiceItemInfo = {
  serviceItems: [
    {
      id: 'abc123',
      submittedAt: '2020-11-20',
      serviceItem: 'Move management',
      code: 'FSC',
      details: {},
    },
  ],
  handleUpdateMTOServiceItemStatus: jest.fn(),
  handleShowRejectionDialog: jest.fn(),
  statusForTableType: MTO_SERVICE_ITEM_STATUS.APPROVED,
};

describe('Service Item Container', () => {
  it('renders the container successfully', async () => {
    render(
      <MockProviders>
        <ServiceItemContainer>
          <RequestedServiceItemsTable {...requestedServiceItemInfo} />
        </ServiceItemContainer>{' '}
      </MockProviders>,
    );

    const serviceItemContainer = await screen.findByTestId('ServiceItemContainer');

    expect(serviceItemContainer).toBeInTheDocument();

    expect(serviceItemContainer.className).toContain('container--accent--default');
  });

  it('renders a child component passed to it', async () => {
    render(
      <MockProviders>
        <ServiceItemContainer>
          <RequestedServiceItemsTable {...requestedServiceItemInfo} />
        </ServiceItemContainer>
      </MockProviders>,
    );

    const childHeading = await screen.findByText('Move Task Order Approved Service Items');
    expect(childHeading).toBeInTheDocument();
  });
});
