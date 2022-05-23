import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Formik } from 'formik';

import OriginZIPInfo from 'components/Office/OriginZIPInfo/OriginZIPInfo';

describe('OriginZIPInfo component', () => {
  it('renders blank form on load', async () => {
    render(
      <Formik
        initialValues={{
          expectedDepartureDate: '',
          pickupPostalCode: '',
          secondPickupPostalCode: '',
        }}
      >
        {({ setFieldValue }) => {
          return <OriginZIPInfo currentZip="90210" setFieldValue={setFieldValue} />;
        }}
      </Formik>,
    );
    expect(await screen.getByRole('heading', { level: 2, name: 'Origin info' })).toBeInTheDocument();
    expect(screen.getByLabelText('Planned departure date')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Origin ZIP')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Second origin ZIP (optional)')).toBeInstanceOf(HTMLInputElement);
  });

  it('fills in current ZIP when use current ZIP checkbox is checked', async () => {
    render(
      <Formik
        initialValues={{
          expectedDepartureDate: '',
          pickupPostalCode: '',
          secondPickupPostalCode: '',
        }}
      >
        {({ setFieldValue }) => {
          return <OriginZIPInfo currentZip="90210" setFieldValue={setFieldValue} />;
        }}
      </Formik>,
    );
    const useCurrentZip = screen.getByText('Use current ZIP');
    const originZip = screen.getByLabelText('Origin ZIP');
    expect(originZip.value).toBe('');
    userEvent.click(useCurrentZip);
    await waitFor(() => {
      expect(originZip.value).toBe('90210');
    });
  });
});
