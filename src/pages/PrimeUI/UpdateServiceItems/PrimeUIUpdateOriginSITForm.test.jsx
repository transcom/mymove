import React from 'react';
import { screen } from '@testing-library/react';

import PrimeUIUpdateOriginSITForm from './PrimeUIUpdateOriginSITForm';

import { renderWithProviders } from 'testUtils';

const originSitInitialValues = {
  sitDepartureDate: '01 Nov 2023',
  sitRequestedDelivery: '01 Dec 2023',
  sitCustomerContacted: '15 Oct 2023',
  mtoServiceItemID: '45fe9475-d592-48f5-896a-45d4d6eb7e76',
};

describe('PrimeUIRequestSITDestAddressChangeForm', () => {
  it('renders the address change request form', async () => {
    renderWithProviders(<PrimeUIUpdateOriginSITForm initialValues={originSitInitialValues} onSubmit={jest.fn()} />);

    expect(screen.getByRole('heading', { name: 'Update Origin SIT Service Item', level: 2 })).toBeInTheDocument();
    expect(await screen.findByLabelText('SIT Departure Date')).toHaveValue('01 Nov 2023');
    expect(await screen.findByLabelText('SIT Requested Delivery')).toHaveValue('01 Dec 2023');
    expect(await screen.findByLabelText('SIT Customer Contacted')).toHaveValue('15 Oct 2023');
    expect(screen.getByRole('button', { name: 'Save' })).toBeEnabled();
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeEnabled();
  });
});
