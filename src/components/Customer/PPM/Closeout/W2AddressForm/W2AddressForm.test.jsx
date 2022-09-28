import React from 'react';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import W2AddressForm from './W2AddressForm';

describe('W2AddressForm component', () => {
  const formFieldsName = 'w2_address';

  const testProps = {
    formFieldsName,
    initialValues: {
      [formFieldsName]: {
        streetAddress1: '',
        streetAddress2: '',
        city: '',
        state: '',
        postalCode: '',
      },
    },
  };

  it('renders the form inputs', async () => {
    const { getByLabelText } = render(<W2AddressForm {...testProps} />);

    await waitFor(() => {
      expect(getByLabelText('Address 1')).toBeInstanceOf(HTMLInputElement);

      expect(getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);

      expect(getByLabelText('City')).toBeInstanceOf(HTMLInputElement);

      expect(getByLabelText('State')).toBeInstanceOf(HTMLSelectElement);

      expect(getByLabelText('ZIP')).toBeInstanceOf(HTMLInputElement);
    });
  });

  it('passes custom validators to fields', async () => {
    const postalCodeValidator = jest.fn().mockImplementation(() => undefined);

    const { findByLabelText } = render(
      <W2AddressForm {...testProps} validators={{ postalCode: postalCodeValidator }} />,
    );

    const postalCodeInput = await findByLabelText('ZIP');

    const postalCode = '99999';

    userEvent.type(postalCodeInput, postalCode);
    userEvent.tab();

    await waitFor(() => {
      expect(postalCodeValidator).toHaveBeenCalledWith(postalCode);
    });
  });

  afterEach(jest.resetAllMocks);
});
