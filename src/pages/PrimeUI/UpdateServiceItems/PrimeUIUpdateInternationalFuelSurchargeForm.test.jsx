import React from 'react';
import { screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import PrimeUIUpdateInternationalFuelSurchargeForm from './PrimeUIUpdateInternationalFuelSurchargeForm';

import { renderWithProviders } from 'testUtils';
import { primeSimulatorRoutes } from 'constants/routes';

const mtoServiceItemID = '38569958-2889-41e5-8101-82c56ec48430';

const serviceItem = {
  id: mtoServiceItemID,
  reServiceCode: 'POEFSC',
  reServiceName: 'International POE fuel surcharge',
  status: 'APPROVED',
  mtoShipmentID: '38569958-2889-41e5-8102-82c56ec48430',
};

const portOfEmbarkation = {
  city: 'SEATTLE',
  country: 'UNITED STATES',
  county: 'KING',
  id: '38569958-2889-41e5-8101-82c56ec48430',
  portCode: 'SEA',
  portName: 'SEATTLE TACOMA INTL',
  portType: 'A',
  state: 'WASHINGTON',
  zip: '98158',
};

const mtoShipment = {
  id: '38569958-2889-41e5-8102-82c56ec48430',
  portOfEmbarkation,
};

const moveTaskOrder = {
  mtoShipments: [mtoShipment],
  mtoServiceItems: [serviceItem],
  serviceItem,
};

const onUpdateServiceItemMock = jest.fn();

// Mock the react-router-dom functions
const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
  useParams: jest.fn().mockReturnValue({ moveCodeOrID: ':moveCodeOrID' }),
}));

describe('PrimeUIUpdateInternationalFuelSurchargeForm', () => {
  it('renders the international fuel surcharge form', async () => {
    renderWithProviders(
      <PrimeUIUpdateInternationalFuelSurchargeForm
        moveTaskOrder={moveTaskOrder}
        mtoServiceItemId={mtoServiceItemID}
        onUpdateServiceItem={onUpdateServiceItemMock}
      />,
    );

    expect(
      screen.getByRole('heading', { name: 'Update International Fuel Surcharge Service Item', level: 2 }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole('heading', { name: 'POEFSC - International POE fuel surcharge', level: 3 }),
    ).toBeInTheDocument();
    expect(screen.getByText('Port:')).toBeInTheDocument();
    expect(screen.getByText('SEATTLE TACOMA INTL')).toBeInTheDocument();
    expect(screen.getByText('SEATTLE, WASHINGTON 98158')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Save' })).toBeEnabled();
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeEnabled();
  });

  it('fires off onUpdateServiceItemMock function when save button is clicked', async () => {
    renderWithProviders(
      <PrimeUIUpdateInternationalFuelSurchargeForm
        moveTaskOrder={moveTaskOrder}
        mtoServiceItemId={mtoServiceItemID}
        onUpdateServiceItem={onUpdateServiceItemMock}
      />,
    );
    const portCodeInput = screen.getByLabelText(/Port Code/);
    await userEvent.type(portCodeInput, 'SEA');
    const saveButton = await screen.findByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    expect(onUpdateServiceItemMock).toHaveBeenCalled();
  });

  it('port code value is set to uppercase when user types in a value', async () => {
    renderWithProviders(
      <PrimeUIUpdateInternationalFuelSurchargeForm
        moveTaskOrder={moveTaskOrder}
        mtoServiceItemId={mtoServiceItemID}
        onUpdateServiceItem={onUpdateServiceItemMock}
      />,
    );
    const portCodeInput = screen.getByLabelText(/Port Code/);
    expect(portCodeInput).toHaveValue('SEA');
    fireEvent.change(portCodeInput, { target: { value: 'pdx' } });

    expect(portCodeInput).toHaveValue('PDX');
  });

  it('does not fire off onUpdateServiceItemMock function when save button is clicked and port code is empty', async () => {
    renderWithProviders(
      <PrimeUIUpdateInternationalFuelSurchargeForm
        moveTaskOrder={moveTaskOrder}
        mtoServiceItemId={mtoServiceItemID}
        onUpdateServiceItem={onUpdateServiceItemMock}
      />,
    );
    const portCodeInput = screen.getByLabelText(/Port Code/);
    await userEvent.clear(portCodeInput, '');
    const saveButton = await screen.findByRole('button', { name: 'Save' });

    onUpdateServiceItemMock.mockClear();
    await userEvent.click(saveButton);
    expect(onUpdateServiceItemMock).not.toHaveBeenCalled();
  });

  it('does not fire off onUpdateServiceItemMock function when save button is clicked and port code is fewer than 3 characters', async () => {
    renderWithProviders(
      <PrimeUIUpdateInternationalFuelSurchargeForm
        moveTaskOrder={moveTaskOrder}
        mtoServiceItemId={mtoServiceItemID}
        onUpdateServiceItem={onUpdateServiceItemMock}
      />,
    );
    const portCodeInput = screen.getByLabelText(/Port Code/);
    await userEvent.clear(portCodeInput, '12');
    const saveButton = await screen.findByRole('button', { name: 'Save' });

    onUpdateServiceItemMock.mockClear();
    await userEvent.click(saveButton);
    expect(onUpdateServiceItemMock).not.toHaveBeenCalled();
  });

  it('directs the user back to the move page when cancel button is clicked', async () => {
    renderWithProviders(
      <PrimeUIUpdateInternationalFuelSurchargeForm
        moveTaskOrder={moveTaskOrder}
        mtoServiceItemId={mtoServiceItemID}
        onUpdateServiceItem={onUpdateServiceItemMock}
      />,
    );

    const cancelButton = await screen.findByRole('button', { name: 'Cancel' });

    await userEvent.click(cancelButton);

    expect(mockNavigate).toHaveBeenCalledWith(primeSimulatorRoutes.VIEW_MOVE_PATH);
  });

  it('renders asterisks for required fields', async () => {
    renderWithProviders(
      <PrimeUIUpdateInternationalFuelSurchargeForm
        moveTaskOrder={moveTaskOrder}
        mtoServiceItemId={mtoServiceItemID}
        onUpdateServiceItem={onUpdateServiceItemMock}
      />,
    );

    expect(document.querySelector('#reqAsteriskMsg')).toHaveTextContent('Fields marked with * are required.');
    expect(await screen.getByLabelText(/Port Code */)).toBeInTheDocument();
  });
});
