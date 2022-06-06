import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Formik } from 'formik';

import DestinationZIPInfo from 'components/Office/DestinationZIPInfo/DestinationZIPInfo';

describe('DestinationZIPInfo component', () => {
  it('renders blank form on load', async () => {
    render(
      <Formik
        initialValues={{
          destinationPostalCode: '',
          secondDestinationPostalCode: '',
        }}
      >
        <DestinationZIPInfo dutyZip="90210" postalCodeValidator={() => {}} />
      </Formik>,
    );
    expect(await screen.getByRole('heading', { level: 2, name: 'Destination info' })).toBeInTheDocument();
    expect(screen.getByLabelText('Destination ZIP')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Second destination ZIP')).toBeInstanceOf(HTMLInputElement);
  });

  it('fills in duty ZIP when use duty ZIP checkbox is checked', async () => {
    render(
      <Formik
        initialValues={{
          destinationPostalCode: '',
          secondDestinationPostalCode: '',
        }}
      >
        {() => {
          return <DestinationZIPInfo dutyZip="90210" postalCodeValidator={() => {}} />;
        }}
      </Formik>,
    );
    const useDutyZip = screen.getByText('Use ZIP for new duty location');
    const destinationZip = screen.getByLabelText('Destination ZIP');
    expect(destinationZip.value).toBe('');
    userEvent.click(useDutyZip);
    await waitFor(() => {
      expect(destinationZip.value).toBe('90210');
    });
  });
});
