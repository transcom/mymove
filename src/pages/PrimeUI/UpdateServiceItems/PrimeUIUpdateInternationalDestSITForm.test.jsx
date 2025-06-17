import React from 'react';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import PrimeUIUpdateDestSITForm from './PrimeUIUpdateDestSITForm';

import { fromPrimeAPIAddressFormat } from 'utils/formatters';
import { renderWithProviders } from 'testUtils';
import { primeSimulatorRoutes } from 'constants/routes';

const shipmentAddress = {
  streetAddress1: '444 Main Ave',
  streetAddress2: 'Apartment 9000',
  streetAddress3: 'Something else',
  city: 'Anytown',
  state: 'AL',
  postalCode: '90210',
};

const serviceItem = {
  reServiceCode: 'IDDSIT',
  reServiceName: 'International destination SIT delivery',
  status: 'REJECTED',
};

const reformatPrimeApiSITDestinationAddress = fromPrimeAPIAddressFormat(shipmentAddress);

const destSitInitialValues = {
  sitDestinationFinalAddress: reformatPrimeApiSITDestinationAddress,
  sitDepartureDate: '01 Nov 2023',
  sitRequestedDelivery: '01 Dec 2023',
  sitCustomerContacted: '15 Oct 2023',
  mtoServiceItemID: '45fe9475-d592-48f5-896a-45d4d6eb7e76',
  reServiceCode: 'IDDSIT',
};

// Mock the react-router-dom functions
const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
  useParams: jest.fn().mockReturnValue({ moveCodeOrID: ':moveCodeOrID' }),
}));

describe('PrimeUIRequestInternationalSITDestAddressChangeForm', () => {
  it('renders the address change request form', async () => {
    renderWithProviders(
      <PrimeUIUpdateDestSITForm
        name="sitDestinationFinalAddress"
        initialValues={destSitInitialValues}
        serviceItem={serviceItem}
        onSubmit={jest.fn()}
      />,
    );

    expect(screen.getByRole('heading', { name: 'Update Destination SIT Service Item', level: 2 })).toBeInTheDocument();
    expect(
      screen.getByRole('heading', { name: 'IDDSIT - International destination SIT delivery', level: 3 }),
    ).toBeInTheDocument();
    expect(await screen.findByLabelText('SIT Departure Date')).toHaveValue('01 Nov 2023');
    expect(await screen.findByLabelText('SIT Requested Delivery')).toHaveValue('01 Dec 2023');
    expect(await screen.findByLabelText('SIT Customer Contacted')).toHaveValue('15 Oct 2023');
    expect(screen.getByRole('button', { name: 'Save' })).toBeEnabled();
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeEnabled();
  });

  it('fires off onSubmit function when save button is clicked', async () => {
    const onSubmitMock = jest.fn();
    renderWithProviders(
      <PrimeUIUpdateDestSITForm
        initialValues={destSitInitialValues}
        serviceItem={serviceItem}
        onSubmit={onSubmitMock}
      />,
    );

    const saveButton = await screen.findByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    expect(onSubmitMock).toHaveBeenCalled();
  });

  it('directs the user back to the move page when cancel button is clicked', async () => {
    renderWithProviders(
      <PrimeUIUpdateDestSITForm
        name="sitDestinationFinalAddress"
        initialValues={destSitInitialValues}
        serviceItem={serviceItem}
        onSubmit={jest.fn()}
      />,
    );

    const cancelButton = await screen.findByRole('button', { name: 'Cancel' });

    await userEvent.click(cancelButton);

    expect(mockNavigate).toHaveBeenCalledWith(primeSimulatorRoutes.VIEW_MOVE_PATH);
  });

  it('renders asterisks for required fields', async () => {
    renderWithProviders(
      <PrimeUIUpdateDestSITForm
        name="sitDestinationFinalAddress"
        initialValues={destSitInitialValues}
        serviceItem={serviceItem}
        onSubmit={jest.fn()}
      />,
    );

    expect(document.querySelector('#reqAsteriskMsg')).toHaveTextContent('Fields marked with * are required.');
    expect(await screen.findByLabelText('Update Reason *')).toBeInTheDocument();
  });
});
