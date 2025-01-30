import React from 'react';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import PrimeUIUpdateInternationalShuttleForm from './PrimeUIUpdateInternationalShuttleForm';

import { renderWithProviders } from 'testUtils';
import { primeSimulatorRoutes } from 'constants/routes';

const internationalShuttleInitialValues = {
  estimatedWeight: 500,
  actualWeight: 600,
  mtoServiceItemID: '45fe9475-d592-48f5-896a-45d4d6eb7e76',
  reServiceCode: 'IDSHUT',
};

const serviceItem = {
  reServiceCode: 'IDSHUT',
  reServiceName: 'International Shuttle',
};

// Mock the react-router-dom functions
const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
  useParams: jest.fn().mockReturnValue({ moveCodeOrID: ':moveCodeOrID' }),
}));

describe('PrimeUIUpdateInternationalShuttleForm', () => {
  it('renders the shuttle change request form', async () => {
    renderWithProviders(
      <PrimeUIUpdateInternationalShuttleForm
        initialValues={internationalShuttleInitialValues}
        serviceItem={serviceItem}
        onSubmit={jest.fn()}
      />,
    );

    expect(
      screen.getByRole('heading', { name: 'Update International Shuttle Service Item', level: 2 }),
    ).toBeInTheDocument();
    expect(await screen.findByLabelText('Estimated Weight')).toHaveValue(500);
    expect(await screen.findByLabelText('Actual Weight')).toHaveValue(600);
    expect(screen.getByRole('button', { name: 'Save' })).toBeEnabled();
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeEnabled();
  });

  it('fires off onSubmit function when save button is clicked', async () => {
    const onSubmitMock = jest.fn();
    renderWithProviders(
      <PrimeUIUpdateInternationalShuttleForm
        initialValues={internationalShuttleInitialValues}
        serviceItem={serviceItem}
        onSubmit={jest.fn()}
      />,
    );

    const saveButton = await screen.findByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    expect(onSubmitMock).toHaveBeenCalled();
  });

  it('directs the user back to the move page when cancel button is clicked', async () => {
    renderWithProviders(
      <PrimeUIUpdateInternationalShuttleForm
        initialValues={internationalShuttleInitialValues}
        serviceItem={serviceItem}
        onSubmit={jest.fn()}
      />,
    );

    const cancelButton = await screen.findByRole('button', { name: 'Cancel' });

    await userEvent.click(cancelButton);

    expect(mockNavigate).toHaveBeenCalledWith(primeSimulatorRoutes.VIEW_MOVE_PATH);
  });
});
