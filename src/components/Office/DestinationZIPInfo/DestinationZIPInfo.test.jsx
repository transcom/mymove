import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Formik } from 'formik';

import DestinationZIPInfo from 'components/Office/DestinationZIPInfo/DestinationZIPInfo';

const defaultProps = {
  setFieldValue: jest.fn(),
  dutyZip: '90210',
  isUseDutyZIPChecked: false,
  postalCodeValidator: jest.fn(),
};

describe('DestinationZIPInfo component', () => {
  it('renders blank form on load', async () => {
    render(
      <Formik
        initialValues={{
          destinationPostalCode: '',
          useDutyZIP: false,
          secondDestinationPostalCode: '',
        }}
      >
        <DestinationZIPInfo {...defaultProps} />
      </Formik>,
    );
    expect(await screen.getByRole('heading', { level: 2, name: 'Destination info' })).toBeInTheDocument();
    expect(screen.getByLabelText('Destination ZIP')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Second destination ZIP (optional)')).toBeInstanceOf(HTMLInputElement);
  });

  it('fills in duty ZIP when use duty ZIP checkbox is checked', async () => {
    render(
      <Formik
        initialValues={{
          destinationPostalCode: '',
          useDutyZIP: false,
          secondDestinationPostalCode: '',
        }}
      >
        {({ setFieldValue }) => {
          return <DestinationZIPInfo {...defaultProps} setFieldValue={setFieldValue} />;
        }}
      </Formik>,
    );
    const useDutyZip = screen.getByText('Use ZIP for new duty location');
    const destinationZip = screen.getByLabelText('Destination ZIP');
    expect(destinationZip.value).toBe('');
    userEvent.click(useDutyZip);
    await waitFor(() => {
      expect(destinationZip.value).toBe(defaultProps.dutyZip);
    });
  });
});
