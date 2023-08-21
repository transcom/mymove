import React from 'react';
import { screen } from '@testing-library/react';
import * as Yup from 'yup';

import PrimeUIRequestSITDestAddressChangeForm from './PrimeUIRequestSITDestAddressChangeForm';

import { addressSchema } from 'utils/validation';
import { fromPrimeAPIAddressFormat } from 'utils/formatters';
import { renderWithProviders } from 'testUtils';

const shipmentAddress = {
  streetAddress1: '444 Main Ave',
  streetAddress2: 'Apartment 9000',
  city: 'Anytown',
  state: 'AL',
  postalCode: '90210',
};

const reformatPrimeApiSITDestinationAddress = fromPrimeAPIAddressFormat(shipmentAddress);

const initialValues = {
  address: reformatPrimeApiSITDestinationAddress,
  contractorRemarks: '',
  mtoServiceItemID: '45fe9475-d592-48f5-896a-45d4d6eb7e76',
};

const destAddressChangeRequestSchema = Yup.object().shape({
  requestedAddress: addressSchema,
  contractorRemarks: Yup.string().required(),
  mtoServiceItemID: Yup.string(),
});

describe('PrimeUIRequestSITDestAddressChangeForm', () => {
  it('renders the address change request form', async () => {
    renderWithProviders(
      <PrimeUIRequestSITDestAddressChangeForm
        name="address"
        destAddressChangeRequestSchema={destAddressChangeRequestSchema}
        initialValues={initialValues}
        onSubmit={jest.fn()}
      />,
    );

    expect(
      screen.getByRole('heading', { name: 'Request Destination SIT Address Change', level: 2 }),
    ).toBeInTheDocument();
    expect(screen.getByLabelText('Address 1')).toBeInTheDocument();
    expect(screen.getByLabelText(/Address 2/)).toBeInTheDocument();
    expect(screen.getByLabelText('City')).toBeInTheDocument();
    expect(screen.getByLabelText('State')).toBeInTheDocument();
    expect(screen.getByLabelText('ZIP')).toBeInTheDocument();
    expect(screen.getByLabelText('Contractor Remarks')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Save' })).toBeEnabled();
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeEnabled();
  });

  it('Prepopulates the form with the current destination address of the DDDSIT service item', async () => {
    renderWithProviders(
      <PrimeUIRequestSITDestAddressChangeForm
        name="address"
        destAddressChangeRequestSchema={destAddressChangeRequestSchema}
        initialValues={initialValues}
        onSubmit={jest.fn()}
      />,
    );

    expect(screen.getByLabelText('Address 1')).toHaveValue('444 Main Ave');
    expect(screen.getByLabelText(/Address 2/)).toHaveValue('Apartment 9000');
    expect(screen.getByLabelText('City')).toHaveValue('Anytown');
    expect(screen.getByLabelText('State')).toHaveValue('AL');
    expect(screen.getByLabelText('Contractor Remarks')).toHaveValue('');
    expect(screen.getByLabelText('ZIP')).toHaveValue('90210');
  });
});
