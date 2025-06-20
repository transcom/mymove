import React from 'react';
import { act, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { createServiceItem } from '../../../services/primeApi';

import CreateServiceItem from './CreateServiceItem';

import { usePrimeSimulatorGetMove } from 'hooks/queries';
import { MockProviders } from 'testUtils';
import { primeSimulatorRoutes } from 'constants/routes';
import { createServiceItemModelTypes } from 'constants/prime';
import { SERVICE_ITEM_CODES } from 'constants/serviceItems';

const mockNavigate = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('hooks/queries', () => ({
  usePrimeSimulatorGetMove: jest.fn(),
}));

jest.mock('services/primeApi', () => ({
  ...jest.requireActual('services/primeApi'),
  createServiceItem: jest.fn(),
}));

const mockSetFlashMessage = jest.fn();

jest.mock('react-redux', () => ({
  ...jest.requireActual('react-redux'),
  connect: jest.fn(() => (component) => (props) => component({ ...props, setFlashMessage: mockSetFlashMessage })),
}));

const moveTaskOrder = {
  id: '1',
  moveCode: 'LN4T89',
  mtoShipments: [
    {
      id: '4',
      eTag: 'testEtag123',
      shipmentType: 'HHG',
      requestedPickupDate: '2021-12-01',
      marketCode: 'd',
      pickupAddress: {
        id: '1',
        streetAddress1: '800 Madison Avenue',
        city: 'New York',
        state: 'NY',
        postalCode: '10002',
      },
      destinationAddress: {
        id: '2',
        streetAddress1: '100 1st Avenue',
        city: 'New York',
        state: 'NY',
        postalCode: '10001',
      },
    },
  ],
};

const routingParams = { moveCodeOrID: 'LN4T89', shipmentId: '4' };

const moveReturnValue = {
  moveTaskOrder,
  isLoading: false,
  isError: false,
};

const renderWithProviders = () => {
  render(
    <MockProviders path={primeSimulatorRoutes.CREATE_SERVICE_ITEM_PATH} params={routingParams}>
      <CreateServiceItem setFlashMessage={jest.fn()} />
    </MockProviders>,
  );
};

describe('successful submission of form', () => {
  it('calls history router back to move details', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
    createServiceItem.mockReturnValue({
      id: '7',
      moveTaskOrderID: '1',
      paymentRequestNumber: '1111-1111-1',
      status: 'PENDING',
      paymentServiceItems: [],
    });

    renderWithProviders();

    const serviceItemsDropDown = screen.getByLabelText(/Service item type/);
    expect(serviceItemsDropDown).toBeInstanceOf(HTMLSelectElement);
    await userEvent.selectOptions(serviceItemsDropDown, createServiceItemModelTypes.MTOServiceItemDomesticShuttle);
    expect(serviceItemsDropDown).toHaveValue(createServiceItemModelTypes.MTOServiceItemDomesticShuttle);

    const serviceItemCode = screen.getByLabelText(/Service item code/);
    expect(serviceItemCode).toBeInstanceOf(HTMLSelectElement);
    await userEvent.selectOptions(serviceItemCode, SERVICE_ITEM_CODES.DOSHUT);
    expect(serviceItemCode).toHaveValue(SERVICE_ITEM_CODES.DOSHUT);

    await userEvent.type(screen.getByLabelText('Reason *'), 'Testing reason');

    const saveButton = screen.getByRole('button', { name: 'Create service item' });
    expect(saveButton).toBeEnabled();
    await act(async () => {
      await userEvent.click(screen.getByRole('button', { name: 'Create service item' }));
    });

    await waitFor(() => {
      expect(mockSetFlashMessage).toHaveBeenCalledWith(
        `MSG_CREATE_SERVICE_ITEM_SUCCESS${moveTaskOrder.moveCode}`,
        'success',
        'Successfully created service item',
        '',
        true,
      );
      expect(mockNavigate).toHaveBeenCalledWith('/simulator/moves/LN4T89/details');
    });
  });
});

describe('CreateServiceItem page', () => {
  describe('check loading and error component states', () => {
    const loadingReturnValue = {
      moveTaskOrder: undefined,
      isLoading: true,
      isError: false,
    };

    const errorReturnValue = {
      moveTaskOrder: undefined,
      isLoading: false,
      isError: true,
    };

    it('renders the loading placeholder when the query is still loading', () => {
      usePrimeSimulatorGetMove.mockReturnValue(loadingReturnValue);

      render(
        <MockProviders>
          <CreateServiceItem setFlashMessage={jest.fn()} />
        </MockProviders>,
      );

      expect(screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 }));
    });

    it('renders the Something Went Wrong component when the query has an error', () => {
      usePrimeSimulatorGetMove.mockReturnValue(errorReturnValue);

      render(
        <MockProviders>
          <CreateServiceItem setFlashMessage={jest.fn()} />
        </MockProviders>,
      );

      expect(screen.getByText(/Something went wrong./));
    });
  });
});
