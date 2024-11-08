import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import * as Yup from 'yup';
import userEvent from '@testing-library/user-event';

import PrimeUIShipmentUpdateAddressForm from './PrimeUIShipmentUpdateAddressForm';

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

describe('PrimeUIShipmentUpdateAddressForm', () => {
  const shipmentAddress = {
    id: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
    streetAddress1: '444 Main Ave',
    streetAddress2: 'Apartment 9000',
    streetAddress3: 'c/o Anyone',
    city: 'Anytown',
    state: 'AL',
    county: 'Los Angeles',
    postalCode: '90210',
    country: 'USA',
    eTag: '1234567890',
  };

  const reformatPrimeApiShipmentAddress = fromPrimeAPIAddressFormat(shipmentAddress);

  const initialValuesPickupAddress = {
    addressID: shipmentAddress.id,
    pickupAddress: {
      address: reformatPrimeApiShipmentAddress,
    },
    eTag: shipmentAddress.eTag,
  };

  const initialValuesDestinationAddress = {
    addressID: shipmentAddress.id,
    destinationAddress: {
      address: reformatPrimeApiShipmentAddress,
    },
    eTag: shipmentAddress.eTag,
  };

  const updateAddressSchema = Yup.object().shape({
    addressID: Yup.string(),
    pickupAddress: Yup.object().shape({
      address: requiredAddressSchema,
    }),
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
      <PrimeUIShipmentUpdateAddressForm
        initialValues={initialValuesPickupAddress}
        updateShipmentAddressSchema={updateAddressSchema}
        addressLocation="Pickup address"
        onSubmit={jest.fn()}
        name="pickupAddress.address"
      />,
    );
    expect(screen.getByRole('heading', { name: 'Pickup address', level: 2 })).toBeInTheDocument();
    expect(screen.getByLabelText(/Address 1/)).toBeInTheDocument();
    expect(screen.getByLabelText(/Address 2/)).toBeInTheDocument();
    expect(screen.getByLabelText(/Address 3/)).toBeInTheDocument();
    expect(screen.getByLabelText('City')).toBeInTheDocument();
    expect(screen.getByLabelText('City')).toHaveValue(shipmentAddress.city);
    expect(screen.getByLabelText('County')).toBeInTheDocument();
    expect(screen.getByLabelText('County')).toHaveValue(shipmentAddress.county);
    expect(screen.getByLabelText('State')).toBeInTheDocument();
    expect(screen.getByLabelText('State')).toHaveValue(shipmentAddress.state);
    expect(screen.getByLabelText('ZIP')).toBeInTheDocument();
    expect(screen.getByLabelText('ZIP')).toHaveValue(shipmentAddress.postalCode);

    expect(screen.getByRole('button', { name: 'Save' })).toBeEnabled();
  });

  it('change text and button is enabled', async () => {
    renderWithProviders(
      <PrimeUIShipmentUpdateAddressForm
        initialValues={initialValuesPickupAddress}
        updateShipmentAddressSchema={updateAddressSchema}
        addressLocation="Pickup address"
        onSubmit={jest.fn()}
        name="pickupAddress.address"
      />,
    );

    await userEvent.type(screen.getByLabelText(/Address 1/), '23 City Str');
    await userEvent.type(screen.getByLabelText(/Address 2/), 'Apt 23');
    await userEvent.type(screen.getByLabelText(/Address 3/), 'C/O Twenty Three');

    const submitBtn = screen.getByRole('button', { name: 'Save' });
    await waitFor(() => {
      expect(submitBtn).toBeEnabled();
    });
    await userEvent.click(submitBtn);
  });

  it('does not disable the submit button when address lines 2 or 3 are blank', async () => {
    renderWithProviders(
      <PrimeUIShipmentUpdateAddressForm
        initialValues={initialValuesDestinationAddress}
        updateShipmentAddressSchema={updateAddressSchema}
        addressLocation="Destination address"
        onSubmit={jest.fn()}
        name="destinationAddress.address"
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

  it('disables the submit button when the address 1 is missing - pickup', async () => {
    renderWithProviders(
      <PrimeUIShipmentUpdateAddressForm
        initialValues={initialValuesPickupAddress}
        updateShipmentAddressSchema={updateAddressSchema}
        addressLocation="Pickup address"
        onSubmit={jest.fn()}
        name="pickupAddress.address"
      />,
    );
    await userEvent.clear(screen.getByLabelText(/Address 1/));
    (await screen.getByLabelText(/Address 1/)).blur();
    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Save' })).toBeDisabled();
      expect(screen.getByText('Required')).toBeInTheDocument();
    });
  });

  it('disables the submit button when the address 1 is missing - desination', async () => {
    renderWithProviders(
      <PrimeUIShipmentUpdateAddressForm
        initialValues={initialValuesDestinationAddress}
        updateShipmentAddressSchema={updateAddressSchema}
        addressLocation="Destination address"
        onSubmit={jest.fn()}
        name="destinationAddress.address"
      />,
    );
    await userEvent.clear(screen.getByLabelText(/Address 1/));
    (await screen.getByLabelText(/Address 1/)).blur();
    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Save' })).toBeDisabled();
    });
  });
});
