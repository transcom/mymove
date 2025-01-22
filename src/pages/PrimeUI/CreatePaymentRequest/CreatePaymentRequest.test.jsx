import React from 'react';
import { act, render, screen, within, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import CreatePaymentRequest from './CreatePaymentRequest';

import { usePrimeSimulatorGetMove } from 'hooks/queries';
import { createPaymentRequest } from 'services/primeApi';
import { MockProviders } from 'testUtils';
import { primeSimulatorRoutes } from 'constants/routes';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));
const routingParams = { moveCodeOrID: 'LN4T89' };

jest.mock('hooks/queries', () => ({
  usePrimeSimulatorGetMove: jest.fn(),
  createPaymentRequest: jest.fn(),
}));

jest.mock('services/primeApi', () => ({
  ...jest.requireActual('services/primeApi'),
  createPaymentRequest: jest.fn(),
}));

const moveTaskOrder = {
  id: '1',
  moveCode: 'LN4T89',
  mtoShipments: [
    {
      id: '2',
      shipmentType: 'HHG',
      requestedPickupDate: '2021-11-26',
      pickupAddress: { streetAddress1: '100 1st Avenue', city: 'New York', state: 'NY', postalCode: '10001' },
      destinationAddress: { streetAddress1: '200 2nd Avenue', city: 'Buffalo', state: 'NY', postalCode: '1001' },
      marketCode: 'd',
    },
    {
      id: '3',
      shipmentType: 'HHG_INTO_NTS',
      requestedPickupDate: '2021-12-01',
      pickupAddress: { streetAddress1: '800 Madison Avenue', city: 'New York', state: 'NY', postalCode: '10002' },
      destinationAddress: { streetAddress1: '200 2nd Avenue', city: 'Buffalo', state: 'NY', postalCode: '1001' },
      marketCode: 'i',
    },
  ],
  mtoServiceItems: [
    { id: '4', reServiceCode: 'MS', reServiceName: 'Move management' },
    { id: '5', reServiceCode: 'DLH', mtoShipmentID: '2', reServiceName: 'Domestic linehaul' },
    { id: '6', reServiceCode: 'FSC', mtoShipmentID: '3', reServiceName: 'Fuel surcharge' },
  ],
};

const moveReturnValue = {
  moveTaskOrder,
  isLoading: false,
  isError: false,
};

const renderWithProviders = () => {
  render(
    <MockProviders path={primeSimulatorRoutes.CREATE_PAYMENT_REQUEST_PATH} params={routingParams}>
      <CreatePaymentRequest setFlashMessage={jest.fn()} />
    </MockProviders>,
  );
};

describe('CreatePaymentRequest page', () => {
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

    it('renders the loading placeholder when the query is still loading', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(loadingReturnValue);

      renderWithProviders();

      expect(await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 }));
    });

    it('renders the Something Went Wrong component when the query has an error', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(errorReturnValue);

      renderWithProviders();

      expect(await screen.getByText(/Something went wrong./));
    });
  });

  describe('displaying move, shipment, and service item information', () => {
    it('displays the move information and basic service items', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);

      renderWithProviders();

      const moveHeading = screen.getByRole('heading', { name: 'Move', level: 2 });
      expect(moveHeading).toBeInTheDocument();

      const moveContainer = moveHeading.parentElement;
      expect(within(moveContainer).getByText('Move Code:')).toBeInTheDocument();
      expect(within(moveContainer).getByText(moveTaskOrder.moveCode)).toBeInTheDocument();
      expect(within(moveContainer).getByText('Move Id:')).toBeInTheDocument();
      expect(within(moveContainer).getByText('1')).toBeInTheDocument();

      const moveServiceItems = screen.getByRole('heading', { name: 'Move Service Items', level: 2 });
      expect(moveServiceItems).toBeInTheDocument();

      const moveServiceItemsContainer = moveServiceItems.parentElement;
      expect(
        within(moveServiceItemsContainer).getByRole('heading', { name: 'Move management', level: 3 }),
      ).toBeInTheDocument();
      expect(within(moveServiceItemsContainer).getByText('Service Code:')).toBeInTheDocument();
      expect(within(moveServiceItemsContainer).getByText('MS')).toBeInTheDocument();
      expect(
        within(moveServiceItemsContainer).getByRole('checkbox', { name: 'Add to payment request', checked: false }),
      ).toBeInTheDocument();
    });

    it('displays the shipment information and shipment service items', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);

      renderWithProviders();

      // Verify the main shipments heading
      const shipmentsHeading = screen.getByRole('heading', { name: 'Shipments', level: 2 });
      expect(shipmentsHeading).toBeInTheDocument();

      // Locate the container holding all shipment headings
      const shipmentsContainer = shipmentsHeading.parentElement;

      // Get all shipment headings (level 3) within the shipments container
      const shipmentHeadings = within(shipmentsContainer).getAllByRole('heading', { level: 3 });
      const headingTexts = shipmentHeadings.map((heading) => heading.textContent);

      // Expected shipment headings (for both HHG and NTS)
      const expectedHeadings = [`dHHG shipment`, `iNTS shipment`];

      // Check that all expected headings are present, regardless of order
      expectedHeadings.forEach((expectedHeading) => {
        expect(headingTexts).toContain(expectedHeading);
      });

      // Check service items and checkboxes within the HHG container
      const hhgContainer = shipmentHeadings.find((heading) => heading.textContent === expectedHeadings[0]).parentElement
        .parentElement;

      expect(
        within(hhgContainer).getByRole('checkbox', { name: 'Add all service items', checked: false }),
      ).toBeInTheDocument();

      expect(within(hhgContainer).getByRole('heading', { name: 'Domestic linehaul', level: 3 }));

      expect(
        within(hhgContainer).getByRole('checkbox', { name: 'Add to payment request', checked: false }),
      ).toBeInTheDocument();

      // Check service items and checkboxes within the NTS container
      const ntsContainer = shipmentHeadings.find((heading) => heading.textContent === expectedHeadings[1]).parentElement
        .parentElement;

      expect(
        within(ntsContainer).getByRole('checkbox', { name: 'Add all service items', checked: false }),
      ).toBeInTheDocument();

      expect(within(ntsContainer).getByRole('heading', { name: 'Fuel surcharge', level: 3 }));

      expect(
        within(ntsContainer).getByRole('checkbox', { name: 'Add to payment request', checked: false }),
      ).toBeInTheDocument();
    });

    it('displays the submit button and hint text', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);

      renderWithProviders();

      expect(screen.getByRole('button', { name: 'Submit Payment Request' })).toBeDisabled();
      expect(
        screen.getByText(
          'At least one basic service item or shipment service item is required to create a payment request',
        ),
      ).toBeInTheDocument();
    });
  });

  describe('error alert display', () => {
    it('displays the error alert when the api submission returns an error', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      createPaymentRequest.mockRejectedValue({ response: { body: { title: 'Error title', detail: 'Error detail' } } });

      renderWithProviders();

      const serviceItemInputs = screen.getAllByRole('checkbox', { name: 'Add to payment request' });
      // avoiding linter pitfalls with async for loops
      await userEvent.click(serviceItemInputs[0]);
      await userEvent.click(serviceItemInputs[1]);
      await userEvent.click(serviceItemInputs[2]);

      await act(async () => {
        await userEvent.click(screen.getByRole('button', { name: 'Submit Payment Request' }));
      });

      await waitFor(() => {
        expect(screen.getByText('Prime API: Error title')).toBeInTheDocument();
        expect(screen.getByText('Error detail')).toBeInTheDocument();
      });
    });

    it('displays the unknown error when none is provided', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      createPaymentRequest.mockRejectedValue('malformed api error response');

      renderWithProviders();

      const serviceItemInputs = screen.getAllByRole('checkbox', { name: 'Add to payment request' });
      await userEvent.click(serviceItemInputs[0]);
      await userEvent.click(serviceItemInputs[1]);
      await userEvent.click(serviceItemInputs[2]);

      await act(async () => {
        await userEvent.click(screen.getByRole('button', { name: 'Submit Payment Request' }));
      });

      await waitFor(() => {
        expect(screen.getByText('Unexpected error')).toBeInTheDocument();
        expect(
          screen.getByText(
            'An unknown error has occurred, please check the state of the shipment and service items data for this move',
          ),
        ).toBeInTheDocument();
      });
    });
  });

  describe('successful submission of form', () => {
    it('calls history router back to move details', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      createPaymentRequest.mockReturnValue({
        id: '7',
        moveTaskOrderID: '1',
        paymentRequestNumber: '1111-1111-1',
        status: 'PENDING',
        paymentServiceItems: [],
      });

      renderWithProviders();

      const serviceItemInputs = screen.getAllByRole('checkbox', { name: 'Add to payment request' });
      await userEvent.click(serviceItemInputs[0]);
      await userEvent.click(serviceItemInputs[1]);
      await userEvent.click(serviceItemInputs[2]);

      await act(async () => {
        await userEvent.click(screen.getByRole('button', { name: 'Submit Payment Request' }));
      });

      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalledWith('/simulator/moves/LN4T89/details');
      });
    });
  });
});
