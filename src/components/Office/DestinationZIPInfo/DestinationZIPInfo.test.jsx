import React from 'react';
import { render, waitFor, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Formik } from 'formik';

import { UnsupportedZipCodePPMErrorMsg } from 'utils/validation';
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
    await userEvent.click(useDutyZip);
    await waitFor(() => {
      expect(destinationZip.value).toBe('90210');
    });
  });

  describe('validates form fields and displays error messages', () => {
    it('marks required inputs when left empty', async () => {
      render(
        <Formik
          initialValues={{
            destinationPostalCode: '',
            secondDestinationPostalCode: '',
          }}
        >
          {() => {
            return <DestinationZIPInfo dutyZip="90210" postalCodeValidator={() => UnsupportedZipCodePPMErrorMsg} />;
          }}
        </Formik>,
      );

      const wrapper = screen.getByTestId('destinationZIP');

      await userEvent.type(within(wrapper).getByTestId('textInput'), '88888');

      within(wrapper).getByTestId('textInput').blur();

      await waitFor(() => {
        expect(screen.getByRole('alert')).toHaveTextContent("We don't have rates for this ZIP code.");
        expect(screen.getByRole('alert').nextElementSibling).toHaveAttribute('name', 'destinationPostalCode');
      });
    });
  });
});
