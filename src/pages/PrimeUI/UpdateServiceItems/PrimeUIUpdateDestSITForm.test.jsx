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

const reformatPrimeApiSITDestinationAddress = fromPrimeAPIAddressFormat(shipmentAddress);

const destSitInitialValues = {
  address: reformatPrimeApiSITDestinationAddress,
  sitDepartureDate: '01 Nov 2023',
  sitRequestedDelivery: '01 Dec 2023',
  sitCustomerContacted: '15 Oct 2023',
  mtoServiceItemID: '45fe9475-d592-48f5-896a-45d4d6eb7e76',
};

// Mock the react-router-dom functions
const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
  useParams: jest.fn().mockReturnValue({ moveCodeOrID: ':moveCodeOrID' }),
}));

describe('PrimeUIRequestSITDestAddressChangeForm', () => {
  it('renders the address change request form', async () => {
    renderWithProviders(
      <PrimeUIUpdateDestSITForm name="address" initialValues={destSitInitialValues} onSubmit={jest.fn()} />,
    );

    expect(screen.getByRole('heading', { name: 'Update Destination SIT Service Item', level: 2 })).toBeInTheDocument();
    expect(screen.getByLabelText('Address 1')).toBeInTheDocument();
    expect(screen.getByLabelText(/Address 2/)).toBeInTheDocument();
    expect(screen.getByLabelText(/Address 3/)).toBeInTheDocument();
    expect(screen.getByLabelText('City')).toBeInTheDocument();
    expect(screen.getByLabelText('State')).toBeInTheDocument();
    expect(screen.getByLabelText('ZIP')).toBeInTheDocument();
    expect(await screen.findByLabelText('SIT Departure Date')).toHaveValue('01 Nov 2023');
    expect(await screen.findByLabelText('SIT Requested Delivery')).toHaveValue('01 Dec 2023');
    expect(await screen.findByLabelText('SIT Customer Contacted')).toHaveValue('15 Oct 2023');
    expect(screen.getByRole('button', { name: 'Save' })).toBeEnabled();
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeEnabled();
  });

  it('directs the user back to the move page when cancel button is clicked', async () => {
    renderWithProviders(
      <PrimeUIUpdateDestSITForm name="address" initialValues={destSitInitialValues} onSubmit={jest.fn()} />,
    );

    const cancelButton = await screen.findByRole('button', { name: 'Cancel' });

    await userEvent.click(cancelButton);

    expect(mockNavigate).toHaveBeenCalledWith(primeSimulatorRoutes.VIEW_MOVE_PATH);
  });
});
