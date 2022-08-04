import React from 'react';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ResidentialAddressForm from './ResidentialAddressForm';

describe('ResidentialAddressForm component', () => {
  const formFieldsName = 'current_residence';

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
    onSubmit: jest.fn().mockImplementation(() => Promise.resolve()),
    onBack: jest.fn(),
  };

  const fakeAddress = {
    streetAddress1: '235 Prospect Valley Road SE',
    streetAddress2: '#125',
    city: 'El Paso',
    state: 'TX',
    postalCode: '79912',
  };

  it('renders the form inputs', async () => {
    const { getByLabelText } = render(<ResidentialAddressForm {...testProps} />);

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
      <ResidentialAddressForm {...testProps} validators={{ postalCode: postalCodeValidator }} />,
    );

    const postalCodeInput = await findByLabelText('ZIP');

    const postalCode = '99999';

    await userEvent.type(postalCodeInput, postalCode);
    await userEvent.tab();

    await waitFor(() => {
      expect(postalCodeValidator).toHaveBeenCalledWith(postalCode);
    });
  });

  it('shows an error message if trying to submit an invalid form', async () => {
    const { getByRole, findAllByRole } = render(<ResidentialAddressForm {...testProps} />);
    const submitBtn = getByRole('button', { name: 'Next' });

    await userEvent.click(submitBtn);

    const alerts = await findAllByRole('alert');

    expect(alerts.length).toBe(4);

    alerts.forEach((alert) => {
      expect(alert).toHaveTextContent('Required');
    });

    expect(testProps.onSubmit).not.toHaveBeenCalled();
  });

  it('submits the form when its valid', async () => {
    const { getByRole, getByLabelText } = render(<ResidentialAddressForm {...testProps} />);
    const submitBtn = getByRole('button', { name: 'Next' });

    await userEvent.type(getByLabelText('Address 1'), fakeAddress.streetAddress1);
    await userEvent.type(getByLabelText(/Address 2/), fakeAddress.streetAddress2);
    await userEvent.type(getByLabelText('City'), fakeAddress.city);
    await userEvent.selectOptions(getByLabelText('State'), [fakeAddress.state]);
    await userEvent.type(getByLabelText('ZIP'), fakeAddress.postalCode);

    await userEvent.click(submitBtn);

    const expectedParams = {
      [formFieldsName]: fakeAddress,
    };

    await waitFor(() => {
      expect(testProps.onSubmit).toHaveBeenCalledWith(expectedParams, expect.anything());
    });
  });

  it('implements the onBack handler when the Back button is clicked', async () => {
    const { getByRole } = render(<ResidentialAddressForm {...testProps} />);
    const backBtn = getByRole('button', { name: 'Back' });

    await userEvent.click(backBtn);

    await waitFor(() => {
      expect(testProps.onBack).toHaveBeenCalled();
    });
  });

  afterEach(jest.resetAllMocks);
});
