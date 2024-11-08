import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import * as Yup from 'yup';
import userEvent from '@testing-library/user-event';

import PrimeUIShipmentUpdateDestinationAddressForm from './PrimeUIShipmentUpdateDestinationAddressForm';

import { requiredAddressSchema } from 'utils/validation';
import { fromPrimeAPIAddressFormat } from 'utils/formatters';
import { MockProviders, ReactQueryWrapper } from 'testUtils';
import { primeSimulatorRoutes } from 'constants/routes';

const mockNavigate = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));
const routingParams = { moveCodeOrID: 'LN4T89', shipmentId: '4' };

describe('PrimeUIShipmentUpdateDestinationAddressForm', () => {
  const shipmentAddress = {
    id: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
    streetAddress1: '444 Main Ave',
    streetAddress2: 'Apartment 9000',
    streetAddress3: 'J. Jonah Jameson',
    city: 'Anytown',
    county: 'SOLANO',
    state: 'AL',
    postalCode: '90210',
    country: 'USA',
    eTag: '1234567890',
  };

  const reformatPrimeApiDestinationAddress = fromPrimeAPIAddressFormat(shipmentAddress);

  const initialValuesDestinationAddress = {
    mtoShipmentID: shipmentAddress.id,
    contractorRemarks: '',
    newAddress: {
      address: reformatPrimeApiDestinationAddress,
    },
    eTag: shipmentAddress.eTag,
  };

  const updateDestinationAddressSchema = Yup.object().shape({
    mtoShipmentID: Yup.string(),
    newAddress: Yup.object().shape({
      address: requiredAddressSchema,
    }),
    contractorRemarks: Yup.string().required('Contractor remarks are required to make these changes'),
    eTag: Yup.string(),
  });

  const renderWithProviders = (component) => {
    render(
      <ReactQueryWrapper>
        <MockProviders path={primeSimulatorRoutes.SHIPMENT_UPDATE_ADDRESS_PATH} params={routingParams}>
          {component}
        </MockProviders>
      </ReactQueryWrapper>,
    );
  };

  it('renders the form', async () => {
    renderWithProviders(
      <PrimeUIShipmentUpdateDestinationAddressForm
        initialValues={initialValuesDestinationAddress}
        updateDestinationAddressSchema={updateDestinationAddressSchema}
        onSubmit={jest.fn()}
        name="newAddress.address"
      />,
    );
    expect(screen.getByLabelText(/Address 1/)).toBeInTheDocument();
    expect(screen.getByLabelText(/Address 2/)).toBeInTheDocument();
    expect(screen.getByLabelText('City')).toBeInTheDocument();
    expect(screen.getByLabelText('State')).toBeInTheDocument();
    expect(screen.getByLabelText('ZIP')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Save' })).toBeEnabled();
  });

  it('change text and button is enabled', async () => {
    renderWithProviders(
      <PrimeUIShipmentUpdateDestinationAddressForm
        initialValues={initialValuesDestinationAddress}
        updateDestinationAddressSchema={updateDestinationAddressSchema}
        onSubmit={jest.fn()}
        name="newAddress.address"
      />,
    );

    await userEvent.type(screen.getByLabelText(/Address 1/), '23 City Str');
    await userEvent.type(screen.getByLabelText(/Address 2/), 'Apt 23');
    await userEvent.type(screen.getByLabelText(/Address 3/), 'C/O Twenty Three');
    await userEvent.type(screen.getByLabelText('Contractor Remarks'), 'Test remarks');

    const submitBtn = screen.getByRole('button', { name: 'Save' });
    await waitFor(() => {
      expect(submitBtn).toBeEnabled();
    });
    await userEvent.click(submitBtn);
  });

  it('disables the submit button when the address 1 is missing', async () => {
    renderWithProviders(
      <PrimeUIShipmentUpdateDestinationAddressForm
        initialValues={initialValuesDestinationAddress}
        updateDestinationAddressSchema={updateDestinationAddressSchema}
        onSubmit={jest.fn()}
        name="newAddress.address"
      />,
    );
    await userEvent.clear(screen.getByLabelText(/Address 1/));
    (await screen.getByLabelText(/Address 1/)).blur();
    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Save' })).toBeDisabled();
      expect(screen.getByText('Required')).toBeInTheDocument();
    });
  });

  it('does not disable the submit button when address lines 2 or 3 are blank', async () => {
    renderWithProviders(
      <PrimeUIShipmentUpdateDestinationAddressForm
        initialValues={initialValuesDestinationAddress}
        updateDestinationAddressSchema={updateDestinationAddressSchema}
        addressLocation="Pickup address"
        onSubmit={jest.fn()}
        name="newAddress.address"
      />,
    );

    await userEvent.clear(screen.getByLabelText(/Address 3/));
    (await screen.getByLabelText(/Address 3/)).blur();
    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Save' }).getAttribute('disabled')).toBeFalsy();
    });

    await userEvent.clear(screen.getByLabelText(/Address 2/));
    (await screen.getByLabelText(/Address 2/)).blur();
    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Save' }).getAttribute('disabled')).toBeFalsy();
    });
  });
});
