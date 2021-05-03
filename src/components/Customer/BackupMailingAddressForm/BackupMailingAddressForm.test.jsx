import React from 'react';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import BackupMailingAddressForm from './BackupMailingAddressForm';

describe('BackupMailingAddressForm component', () => {
  const formFieldsName = 'backup_mailing_residence';

  const testProps = {
    formFieldsName,
    initialValues: {
      [formFieldsName]: {
        street_address_1: '',
        street_address_2: '',
        city: '',
        state: '',
        postal_code: '',
      },
    },
    onSubmit: jest.fn().mockImplementation(() => Promise.resolve()),
    onBack: jest.fn(),
  };

  const fakeAddress = {
    street_address_1: '235 Prospect Valley Road SE',
    street_address_2: '#125',
    city: 'El Paso',
    state: 'TX',
    postal_code: '79912',
  };

  it('renders the form inputs', async () => {
    const { getByLabelText } = render(<BackupMailingAddressForm {...testProps} />);

    await waitFor(() => {
      expect(getByLabelText('Address 1')).toBeInstanceOf(HTMLInputElement);

      expect(getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);

      expect(getByLabelText('City')).toBeInstanceOf(HTMLInputElement);

      expect(getByLabelText('State')).toBeInstanceOf(HTMLSelectElement);

      expect(getByLabelText('ZIP')).toBeInstanceOf(HTMLInputElement);
    });
  });

  it('shows an error message if trying to submit an invalid form', async () => {
    const { getByRole, findAllByRole } = render(<BackupMailingAddressForm {...testProps} />);
    const submitBtn = getByRole('button', { name: 'Next' });

    userEvent.click(submitBtn);

    const alerts = await findAllByRole('alert');

    expect(alerts.length).toBe(4);

    alerts.forEach((alert) => {
      expect(alert).toHaveTextContent('Required');
    });

    expect(testProps.onSubmit).not.toHaveBeenCalled();
  });

  it('submits the form when its valid', async () => {
    const { getByRole, getByLabelText } = render(<BackupMailingAddressForm {...testProps} />);
    const submitBtn = getByRole('button', { name: 'Next' });

    userEvent.type(getByLabelText('Address 1'), fakeAddress.street_address_1);
    userEvent.type(getByLabelText(/Address 2/), fakeAddress.street_address_2);
    userEvent.type(getByLabelText('City'), fakeAddress.city);
    userEvent.selectOptions(getByLabelText('State'), [fakeAddress.state]);
    userEvent.type(getByLabelText('ZIP'), fakeAddress.postal_code);

    userEvent.click(submitBtn);

    const expectedParams = {
      [formFieldsName]: fakeAddress,
    };

    await waitFor(() => {
      expect(testProps.onSubmit).toHaveBeenCalledWith(expectedParams, expect.anything());
    });
  });

  it('implements the onBack handler when the Back button is clicked', async () => {
    const { getByRole } = render(<BackupMailingAddressForm {...testProps} />);
    const backBtn = getByRole('button', { name: 'Back' });

    userEvent.click(backBtn);

    await waitFor(() => {
      expect(testProps.onBack).toHaveBeenCalled();
    });
  });

  afterEach(jest.resetAllMocks);
});
