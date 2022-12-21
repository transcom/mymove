import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router';

import { primeSimulatorRoutes } from 'constants/routes';
import { MockProviders } from 'testUtils';
import { createPrimeMTOShipment } from 'services/primeApi';
import PrimeUIShipmentCreate from 'pages/PrimeUI/Shipment/PrimeUIShipmentCreate';

const moveId = '9c7b255c-2981-4bf8-839f-61c7458e2b4d';

const mockPush = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: jest.fn().mockReturnValue({
    moveCode: 'LR4T8V',
    moveCodeOrID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    shipmentId: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
  }),
  useHistory: () => ({
    push: mockPush,
  }),
}));

jest.mock('services/primeApi', () => ({
  ...jest.requireActual('services/primeApi'),
  createPrimeMTOShipment: jest.fn().mockImplementation(() => Promise.resolve()),
}));

const createShipmentURL = generatePath(primeSimulatorRoutes.CREATE_SHIPMENT_PATH, {
  moveCodeOrID: moveId,
});
const moveDetailsURL = generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID: moveId });

const mockedComponent = (
  <MockProviders initialEntries={[createShipmentURL]}>
    <PrimeUIShipmentCreate setFlashMessage={jest.fn()} />
  </MockProviders>
);

describe('Create Shipment Page', () => {
  it('renders the page without errors', async () => {
    render(mockedComponent);

    expect(await screen.findByText('Shipment Type')).toBeInTheDocument();
  });

  it('navigates the user to the home page when the cancel button is clicked', async () => {
    render(
      <MockProviders>
        <PrimeUIShipmentCreate setFlashMessage={jest.fn()} />
      </MockProviders>,
    );

    const cancel = screen.getByRole('button', { name: 'Cancel' });
    userEvent.click(cancel);

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith(moveDetailsURL);
    });
  });
});

describe('successful submission of form', () => {
  it('calls history router back to move details', async () => {
    createPrimeMTOShipment.mockReturnValue({});

    render(
      <MockProviders>
        <PrimeUIShipmentCreate setFlashMessage={jest.fn()} />
      </MockProviders>,
    );

    userEvent.selectOptions(screen.getByLabelText('Shipment type'), 'HHG');

    const saveButton = await screen.getByRole('button', { name: 'Save' });

    expect(saveButton).not.toBeDisabled();
    userEvent.click(saveButton);

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith(moveDetailsURL);
    });
  });
});
