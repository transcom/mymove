import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import * as Yup from 'yup';
import userEvent from '@testing-library/user-event';

import { requiredAddressSchema } from '../../../utils/validation';
import { fromPrimeAPIAddressFormat } from '../../../utils/formatters';

import PrimeUIShipmentUpdateAddressForm from './PrimeUIShipmentUpdateAddressForm';

const mockUseHistoryPush = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: jest.fn().mockReturnValue({ moveCodeOrID: 'LN4T89', shipmentId: '4' }),
  useHistory: () => ({
    push: mockUseHistoryPush,
  }),
}));

describe('PrimeUIShipmentUpdateAddressForm', () => {
  const shipmentAddress = {
    id: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
    streetAddress1: '444 Main Ave',
    streetAddress2: 'Apartment 9000',
    streetAddress3: '',
    city: 'Anytown',
    state: 'AL',
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

  const updatePickupAddressSchema = Yup.object().shape({
    addressID: Yup.string(),
    pickupAddress: Yup.object().shape({
      address: requiredAddressSchema,
    }),
    eTag: Yup.string(),
  });

  it('renders the form', async () => {
    render(
      <PrimeUIShipmentUpdateAddressForm
        initialValues={initialValuesPickupAddress}
        updateShipmentAddressSchema={updatePickupAddressSchema}
        addressLocation="Pickup address"
        onSubmit={jest.fn()}
        name="pickupAddress.address"
      />,
    );
    expect(screen.getByRole('heading', { name: 'Pickup address', level: 2 })).toBeInTheDocument();
    expect(screen.getByLabelText('Address 1')).toBeInTheDocument();
    expect(screen.getByLabelText(/Address 2/)).toBeInTheDocument();
    expect(screen.getByLabelText('City')).toBeInTheDocument();
    expect(screen.getByLabelText('State')).toBeInTheDocument();
    expect(screen.getByLabelText('ZIP')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Save' })).toBeEnabled();
  });

  it('change text and button is enabled', async () => {
    render(
      <PrimeUIShipmentUpdateAddressForm
        initialValues={initialValuesPickupAddress}
        updateShipmentAddressSchema={updatePickupAddressSchema}
        addressLocation="Pickup address"
        onSubmit={jest.fn()}
        name="pickupAddress.address"
      />,
    );

    await userEvent.type(screen.getByLabelText('Address 1'), '23 City Str');
    await userEvent.type(screen.getByLabelText('City'), 'City');
    await userEvent.clear(screen.getByLabelText('ZIP'));
    await userEvent.type(screen.getByLabelText('ZIP'), '90210');
    await userEvent.selectOptions(screen.getByLabelText('State'), ['CA']);

    const submitBtn = screen.getByRole('button', { name: 'Save' });
    await waitFor(() => {
      expect(submitBtn).toBeEnabled();
    });
    await userEvent.click(submitBtn);
  });

  it('disables the submit button when the zip is bad', async () => {
    render(
      <PrimeUIShipmentUpdateAddressForm
        initialValues={initialValuesPickupAddress}
        updateShipmentAddressSchema={updatePickupAddressSchema}
        addressLocation="Pickup address"
        onSubmit={jest.fn()}
        name="pickupAddress.address"
      />,
    );
    await userEvent.clear(screen.getByLabelText('ZIP'));
    await userEvent.type(screen.getByLabelText('ZIP'), '1');
    (await screen.getByLabelText('ZIP')).blur();
    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Save' })).toBeDisabled();
      expect(screen.getByText('Must be valid zip code')).toBeInTheDocument();
    });
  });

  it('disables the submit button when the address 1 is missing', async () => {
    render(
      <PrimeUIShipmentUpdateAddressForm
        initialValues={initialValuesPickupAddress}
        updateShipmentAddressSchema={updatePickupAddressSchema}
        addressLocation="Pickup address"
        onSubmit={jest.fn()}
        name="pickupAddress.address"
      />,
    );
    await userEvent.clear(screen.getByLabelText('Address 1'));
    (await screen.getByLabelText('Address 1')).blur();
    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Save' })).toBeDisabled();
      expect(screen.getByText('Required')).toBeInTheDocument();
    });
  });

  it('disables the submit button when city is missing', async () => {
    render(
      <PrimeUIShipmentUpdateAddressForm
        initialValues={initialValuesPickupAddress}
        updateShipmentAddressSchema={updatePickupAddressSchema}
        addressLocation="Pickup address"
        onSubmit={jest.fn()}
        name="pickupAddress.address"
      />,
    );
    await userEvent.clear(screen.getByLabelText('City'));
    (await screen.getByLabelText('City')).blur();
    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Save' })).toBeDisabled();
      expect(screen.getByText('Required')).toBeInTheDocument();
    });
  });
});
