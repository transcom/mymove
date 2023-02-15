import React from 'react';
import { render, waitFor, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Formik } from 'formik';

import { UnsupportedZipCodePPMErrorMsg } from 'utils/validation';
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
        {() => {
          return <OriginZIPInfo currentZip="90210" postalCodeValidator={() => {}} />;
        }}
      </Formik>,
    );
    expect(await screen.getByRole('heading', { level: 2, name: 'Origin info' })).toBeInTheDocument();
    expect(screen.getByLabelText('Planned departure date')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Origin ZIP')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Second origin ZIP')).toBeInstanceOf(HTMLInputElement);
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
        {() => {
          return <OriginZIPInfo currentZip="90210" postalCodeValidator={() => {}} />;
        }}
      </Formik>,
    );
    const useCurrentZip = screen.getByText('Use current ZIP');
    const originZip = screen.getByLabelText('Origin ZIP');
    expect(originZip.value).toBe('');
    await userEvent.click(useCurrentZip);
    await waitFor(() => {
      expect(originZip.value).toBe('90210');
    });
  });

  describe('validates form fields and displays error messages', () => {
    it('marks required inputs when left empty', async () => {
      render(
        <Formik
          initialValues={{
            expectedDepartureDate: '',
            pickupPostalCode: '',
            secondPickupPostalCode: '',
          }}
        >
          {() => {
            return <OriginZIPInfo currentZip="90210" postalCodeValidator={() => UnsupportedZipCodePPMErrorMsg} />;
          }}
        </Formik>,
      );

      const wrapper = screen.getByTestId('originZIP');

      await userEvent.type(within(wrapper).getByTestId('textInput'), '88888');

      within(wrapper).getByTestId('textInput').blur();

      await waitFor(() => {
        expect(screen.getByRole('alert')).toHaveTextContent("We don't have rates for this ZIP code.");
        expect(screen.getByRole('alert').nextElementSibling).toHaveAttribute('name', 'pickupPostalCode');
      });
    });
  });
});
