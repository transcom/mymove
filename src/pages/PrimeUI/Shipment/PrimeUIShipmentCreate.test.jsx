import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';

import { primeSimulatorRoutes } from 'constants/routes';
import { MockProviders } from 'testUtils';
import { createPrimeMTOShipmentV3 } from 'services/primeApi';
import PrimeUIShipmentCreate from 'pages/PrimeUI/Shipment/PrimeUIShipmentCreate';

const moveCode = 'LR4T8V';
const moveId = '9c7b255c-2981-4bf8-839f-61c7458e2b4d';
const shipmentId = 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee';
const routingParams = {
  moveCode,
  moveCodeOrID: moveId,
  shipmentId,
};

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/primeApi', () => ({
  ...jest.requireActual('services/primeApi'),
  createPrimeMTOShipmentV3: jest.fn().mockImplementation(() => Promise.resolve()),
}));

const moveDetailsURL = generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID: moveId });

const mockedComponent = (
  <MockProviders path={primeSimulatorRoutes.CREATE_SHIPMENT_PATH} params={routingParams}>
    <PrimeUIShipmentCreate setFlashMessage={jest.fn()} />
  </MockProviders>
);

describe('Create Shipment Page', () => {
  it('renders the page without errors', async () => {
    render(mockedComponent);

    expect(await screen.findByText('Shipment Type')).toBeInTheDocument();
  });

  it('navigates the user to the home page when the cancel button is clicked', async () => {
    render(mockedComponent);

    expect(await screen.findByText('Shipment Type')).toBeInTheDocument();

    const cancel = screen.getByRole('button', { name: 'Cancel' });
    await userEvent.click(cancel);

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(moveDetailsURL);
    });
  });
});

describe('successful submission of form', () => {
  it('calls history router back to move details', async () => {
    createPrimeMTOShipmentV3.mockReturnValue({});

    render(mockedComponent);

    await userEvent.selectOptions(screen.getByLabelText('Shipment type'), 'HHG');

    const saveButton = await screen.getByRole('button', { name: 'Save' });

    expect(saveButton).not.toBeDisabled();
    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(moveDetailsURL);
    });
  });
});
