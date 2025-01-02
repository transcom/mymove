import React from 'react';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Provider } from 'react-redux';

import BackupAddressForm from './BackupAddressForm';

import { configureStore } from 'shared/store';

describe('BackupAddressForm component', () => {
  const formFieldsName = 'backup_mailing_residence';

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
    county: 'El Paso',
  };

  const dataProps = {
    formFieldsName,
    initialValues: {
      [formFieldsName]: {
        streetAddress1: '',
        streetAddress2: '',
        city: fakeAddress.city,
        state: fakeAddress.state,
        postalCode: fakeAddress.postalCode,
        county: fakeAddress.county,
      },
    },
    onSubmit: jest.fn().mockImplementation(() => Promise.resolve()),
    onBack: jest.fn(),
  };

  it('renders the form inputs', async () => {
    const mockStore = configureStore({});

    const { getByLabelText, getByTestId } = render(
      <Provider store={mockStore.store}>
        <BackupAddressForm {...testProps} />
      </Provider>,
    );

    await waitFor(() => {
      expect(getByLabelText(/Address 1/)).toBeInstanceOf(HTMLInputElement);

      expect(getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);

      expect(getByTestId('City')).toBeInstanceOf(HTMLLabelElement);

      expect(getByTestId('State')).toBeInstanceOf(HTMLLabelElement);

      expect(getByTestId('ZIP')).toBeInstanceOf(HTMLLabelElement);
    });
  });

  it('shows an error message if trying to submit an invalid form', async () => {
    const mockStore = configureStore({});
    const { getByRole, findAllByRole, getByLabelText } = render(
      <Provider store={mockStore.store}>
        <BackupAddressForm {...testProps} />
      </Provider>,
    );
    await userEvent.click(getByLabelText(/Address 1/));
    await userEvent.click(getByLabelText(/Address 2/));

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
    const mockStore = configureStore({});
    const { getByRole, getByLabelText } = render(
      <Provider store={mockStore.store}>
        <BackupAddressForm {...dataProps} />
      </Provider>,
    );
    const submitBtn = getByRole('button', { name: 'Next' });

    await userEvent.type(getByLabelText(/Address 1/), fakeAddress.streetAddress1);
    await userEvent.type(getByLabelText(/Address 2/), fakeAddress.streetAddress2);

    await userEvent.tab();

    await waitFor(() => {
      expect(submitBtn).toBeEnabled();
    });
    await userEvent.click(submitBtn);

    const expectedParams = {
      [formFieldsName]: fakeAddress,
    };

    await waitFor(() => {
      expect(dataProps.onSubmit).toHaveBeenCalledWith(expectedParams, expect.anything());
    });
  });

  it('implements the onBack handler when the Back button is clicked', async () => {
    const mockStore = configureStore({});
    const { getByRole } = render(
      <Provider store={mockStore.store}>
        <BackupAddressForm {...testProps} />
      </Provider>,
    );
    const backBtn = getByRole('button', { name: 'Back' });

    await userEvent.click(backBtn);

    await waitFor(() => {
      expect(testProps.onBack).toHaveBeenCalled();
    });
  });

  afterEach(jest.resetAllMocks);
});
