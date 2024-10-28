import React from 'react';
import { render, screen } from '@testing-library/react';

import PrimeUIUpdateSitServiceItem from './PrimeUIUpdateSitServiceItem';

import { MockProviders, ReactQueryWrapper } from 'testUtils';
import { primeSimulatorRoutes } from 'constants/routes';
import { usePrimeSimulatorGetMove } from 'hooks/queries';

// Mock the usePrimeSimulatorGetMove hook
jest.mock('hooks/queries', () => ({
  usePrimeSimulatorGetMove: jest.fn(),
}));

describe('PrimeUIUpdateSitServiceItems page', () => {
  it('renders the destination sit service item form', async () => {
    const routingParams = {
      moveCodeOrID: 'bf2fc98f-3cb5-40a0-a125-4c222096c35b',
      mtoServiceItemId: '45fe9475-d592-48f5-896a-45d4d6eb7e76',
    };

    const renderComponent = () => {
      render(
        <ReactQueryWrapper>
          <MockProviders path={primeSimulatorRoutes.UPDATE_SIT_SERVICE_ITEM_PATH} params={routingParams}>
            <PrimeUIUpdateSitServiceItem />
          </MockProviders>
        </ReactQueryWrapper>,
      );
    };

    const moveTaskOrder = {
      id: '1',
      moveCode: 'LN4T89',
      mtoShipments: [
        {
          id: '2',
          shipmentType: 'HHG',
          requestedPickupDate: '2021-11-26',
          pickupAddress: { streetAddress1: '100 1st Avenue', city: 'New York', state: 'NY', postalCode: '10001' },
        },
      ],
      mtoServiceItems: [
        {
          reServiceCode: 'DDDSIT',
          modelType: 'MTOServiceItemDestSIT',
          reason: 'Holiday break',
          sitDestinationFinalAddress: {
            streetAddress1: '444 Main Ave',
            streetAddress2: 'Apartment 9000',
            streetAddress3: 'c/o Some Person',
            city: 'Anytown',
            state: 'AL',
            postalCode: '90210',
          },
          id: '45fe9475-d592-48f5-896a-45d4d6eb7e76',
        },
      ],
    };

    const moveReturnValue = {
      moveTaskOrder,
      isLoading: false,
      isError: false,
    };
    usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);

    renderComponent();

    expect(screen.getByRole('heading', { name: 'Update Destination SIT Service Item', level: 2 })).toBeInTheDocument();
    expect(screen.getByRole('textbox', { name: 'SIT Departure Date' })).toBeInTheDocument();
    expect(screen.getByRole('textbox', { name: 'SIT Requested Delivery' })).toBeInTheDocument();
    expect(screen.getByRole('textbox', { name: 'SIT Customer Contacted' })).toBeInTheDocument();
    expect(screen.getByRole('textbox', { name: 'Update Reason' })).toBeInTheDocument();
  });

  it('renders the origin sit service item form', async () => {
    const routingParams = {
      moveCodeOrID: 'bf2fc98f-3cb5-40a0-a125-4c222096c35b',
      mtoServiceItemId: '45fe9475-d592-48f5-896a-45d4d6eb7e76',
    };

    const renderComponent = () => {
      render(
        <ReactQueryWrapper>
          <MockProviders path={primeSimulatorRoutes.UPDATE_SIT_SERVICE_ITEM_PATH} params={routingParams}>
            <PrimeUIUpdateSitServiceItem />
          </MockProviders>
        </ReactQueryWrapper>,
      );
    };

    const moveTaskOrder = {
      id: '1',
      moveCode: 'LN4T89',
      mtoShipments: [
        {
          id: '2',
          shipmentType: 'HHG',
          requestedPickupDate: '2021-11-26',
          pickupAddress: { streetAddress1: '100 1st Avenue', city: 'New York', state: 'NY', postalCode: '10001' },
        },
      ],
      mtoServiceItems: [
        {
          reServiceCode: 'DDDSIT',
          modelType: 'MTOServiceItemOriginSIT',
          reason: 'Holiday break',
          sitDestinationFinalAddress: {
            streetAddress1: '444 Main Ave',
            streetAddress2: 'Apartment 9000',
            streetAddress3: 'c/o Some Person',
            city: 'Anytown',
            state: 'AL',
            postalCode: '90210',
          },
          id: '45fe9475-d592-48f5-896a-45d4d6eb7e76',
        },
      ],
    };

    const moveReturnValue = {
      moveTaskOrder,
      isLoading: false,
      isError: false,
    };
    usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);

    renderComponent();

    expect(screen.getByRole('heading', { name: 'Update Origin SIT Service Item', level: 2 })).toBeInTheDocument();
    expect(screen.getByRole('textbox', { name: 'SIT Departure Date' })).toBeInTheDocument();
    expect(screen.getByRole('textbox', { name: 'SIT Requested Delivery' })).toBeInTheDocument();
    expect(screen.getByRole('textbox', { name: 'SIT Customer Contacted' })).toBeInTheDocument();
    expect(screen.getByRole('textbox', { name: 'Update Reason' })).toBeInTheDocument();
  });
});
