import React from 'react';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import PrimeUIUpdateInternationalFuelSurchargeForm from './PrimeUIUpdateInternationalFuelSurchargeForm';

import { renderWithProviders } from 'testUtils';
import { primeSimulatorRoutes } from 'constants/routes';

const intlFscInitialValues = {
  mtoServiceItemID: '48569958-2889-41e5-8101-82c56ec48430',
  reServiceCode: 'POEFSC',
  portCode: 'SEA',
};

const serviceItem = {
  reServiceCode: 'POEFSC',
  reServiceName: 'International POE Fuel Surcharge',
  status: 'APPROVED',
};

const port = {
  city: 'SEATTLE',
  country: 'UNITED STATES',
  county: 'KING',
  id: '48569958-2889-41e5-8101-82c56ec48430',
  portCode: 'SEA',
  portName: 'SEATTLE TACOMA INTL',
  portType: 'A',
  state: 'WASHINGTON',
  zip: '98158',
};

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
        initialValues={intlFscInitialValues}
        serviceItem={serviceItem}
        port={port}
        onSubmit={jest.fn()}
      />,
    );

    expect(
      screen.getByRole('heading', { name: 'Update International Fuel Surcharge Service Item', level: 2 }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole('heading', { name: 'POEFSC - International POE Fuel Surcharge', level: 3 }),
    ).toBeInTheDocument();
    expect(screen.getByText('Port:')).toBeInTheDocument();
    expect(screen.getByText('SEATTLE TACOMA INTL')).toBeInTheDocument();
    expect(screen.getByText('SEATTLE, WASHINGTON 98158')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Save' })).toBeEnabled();
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeEnabled();
  });

  it('fires off onSubmit function when save button is clicked', async () => {
    const onSubmitMock = jest.fn();
    renderWithProviders(
      <PrimeUIUpdateInternationalFuelSurchargeForm
        initialValues={intlFscInitialValues}
        serviceItem={serviceItem}
        port={port}
        onSubmit={onSubmitMock}
      />,
    );
    const portCodeInput = screen.getByLabelText(/Port Code/);
    await userEvent.type(portCodeInput, 'SEA');
    const saveButton = await screen.findByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    expect(onSubmitMock).toHaveBeenCalled();
  });

  it('does not fire off onSubmit function when save button is clicked and port code is empty', async () => {
    const onSubmitMock = jest.fn();
    renderWithProviders(
      <PrimeUIUpdateInternationalFuelSurchargeForm
        initialValues={intlFscInitialValues}
        serviceItem={serviceItem}
        port={port}
        onSubmit={onSubmitMock}
      />,
    );
    const portCodeInput = screen.getByLabelText(/Port Code/);
    await userEvent.clear(portCodeInput, '');
    const saveButton = await screen.findByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    expect(onSubmitMock).not.toHaveBeenCalled();
  });

  it('does not fire off onSubmit function when save button is clicked and port code is fewer than 3 characters', async () => {
    const onSubmitMock = jest.fn();
    renderWithProviders(
      <PrimeUIUpdateInternationalFuelSurchargeForm
        initialValues={intlFscInitialValues}
        serviceItem={serviceItem}
        port={port}
        onSubmit={onSubmitMock}
      />,
    );
    const portCodeInput = screen.getByLabelText(/Port Code/);
    await userEvent.clear(portCodeInput, '12');
    const saveButton = await screen.findByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    expect(onSubmitMock).not.toHaveBeenCalled();
  });

  it('directs the user back to the move page when cancel button is clicked', async () => {
    renderWithProviders(
      <PrimeUIUpdateInternationalFuelSurchargeForm
        initialValues={intlFscInitialValues}
        serviceItem={serviceItem}
        port={port}
        onSubmit={jest.fn()}
      />,
    );

    const cancelButton = await screen.findByRole('button', { name: 'Cancel' });

    await userEvent.click(cancelButton);

    expect(mockNavigate).toHaveBeenCalledWith(primeSimulatorRoutes.VIEW_MOVE_PATH);
  });
});
