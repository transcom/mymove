import React from 'react';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ResidentialAddressForm from './ResidentialAddressForm';
import addressFactory from 'utils/test/factories/address';

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
        county: '',
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
    county: 'El Paso',
  };

  const mockAddress = addressFactory();

  it('renders the form inputs and help text', async () => {
    const { getByLabelText, getByText } = render(<ResidentialAddressForm {...testProps} />);

    await waitFor(() => {
      expect(getByLabelText('Address 1')).toBeInstanceOf(HTMLInputElement);

      expect(getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);

      expect(getByLabelText('City')).toBeInstanceOf(HTMLInputElement);

      expect(getByLabelText('State')).toBeInstanceOf(HTMLSelectElement);

      expect(getByLabelText('County')).toBeInstanceOf(HTMLSelectElement);

      expect(getByLabelText('ZIP')).toBeInstanceOf(HTMLInputElement);

      expect(getByText('Must be a physical address.')).toBeInTheDocument();
    });
  });

  it('shows an error message if trying to submit an invalid form', async () => {
    const { getByRole, findAllByRole, getByLabelText } = render(<ResidentialAddressForm {...testProps} />);
    await userEvent.click(getByLabelText('Address 1'));
    await userEvent.click(getByLabelText(/Address 2/));
    const postalCodeInput = await screen.findByLabelText('Zip/City Lookup');

    const postalCode = '79912';

    await userEvent.type(postalCodeInput, postalCode);
    await userEvent.click(await screen.findByText('79912'));
    const submitBtn = getByRole('button', { name: 'Next' });
    await userEvent.click(submitBtn);

    const alerts = await findAllByRole('alert');

    expect(alerts.length).toBe(1);

    alerts.forEach((alert) => {
      expect(alert).toHaveTextContent('Required');
    });

    expect(testProps.onSubmit).not.toHaveBeenCalled();
  });

  it('submits the form when its valid', async () => {
    const { getByRole, getByLabelText } = render(<ResidentialAddressForm {...testProps} />);
    const submitBtn = getByRole('button', { name: 'Next' });

    await userEvent.type(getByLabelText('Address 1'), mockAddress.streetAddress1);
    await userEvent.type(getByLabelText(/Address 2/), mockAddress.streetAddress2);
    const input = getByRole('combobox', { id: 'zipCity-input' });
    await userEvent.type(input, mockAddress.postalCode);
    await waitFor(() => {
      expect(screen.getByText(mockAddress.city)).toBeInTheDocument();
    });

    fireEvent.keyPress(input, { key: 'Enter', code: 13 });

    await waitFor(() => {
      expect(submitBtn).toBeEnabled();
    });
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
