import React from 'react';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Provider } from 'react-redux';

import ResidentialAddressForm from './ResidentialAddressForm';

import { configureStore } from 'shared/store';

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

  it('renders the form inputs and help text', async () => {
    const mockStore = configureStore({});
    const { getByLabelText, getByText } = render(
      <Provider store={mockStore.store}>
        <ResidentialAddressForm {...testProps} />
      </Provider>,
    );

    await waitFor(() => {
      expect(getByLabelText(/Address 1/)).toBeInstanceOf(HTMLInputElement);

      expect(getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);

      expect(getByLabelText(/City/)).toBeInstanceOf(HTMLInputElement);

      expect(getByLabelText('State')).toBeInstanceOf(HTMLInputElement);

      expect(getByLabelText(/ZIP/)).toBeInstanceOf(HTMLInputElement);

      expect(getByText('Must be a physical address.')).toBeInTheDocument();
    });
  });

  it('passes custom validators to fields', async () => {
    const postalCodeValidator = jest.fn().mockImplementation(() => undefined);

    const { findByLabelText } = render(
      <ResidentialAddressForm {...testProps} validators={{ postalCode: postalCodeValidator }} />,
    );

    const postalCodeInput = await findByLabelText(/ZIP/);

    const postalCode = '99999';

    await userEvent.type(postalCodeInput, postalCode);

    await waitFor(() => {
      // We expect this to be called 6 times.
      // 1 - validate on mount
      // 5 - once for each 9 that was typed, since we are validating on change
      expect(postalCodeValidator).toHaveBeenCalledTimes(6);
      expect(postalCodeValidator).toHaveBeenCalledWith(postalCode);
    });
  });

  it('shows an error message if trying to submit an invalid form', async () => {
    const mockStore = configureStore({});
    const { getByRole, findAllByRole, getByLabelText } = render(
      <Provider store={mockStore.store}>
        <ResidentialAddressForm {...testProps} />
      </Provider>,
    );
    await userEvent.click(getByLabelText('Address 1'));
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
        <ResidentialAddressForm {...dataProps} />
      </Provider>,
    );
    const submitBtn = getByRole('button', { name: 'Next' });

    await userEvent.type(getByLabelText(/Address 1/), fakeAddress.streetAddress1);
    await userEvent.type(getByLabelText(/Address 2/), fakeAddress.streetAddress2);
    await userEvent.type(getByLabelText(/City/), fakeAddress.city);
    await userEvent.selectOptions(getByLabelText(/State/), [fakeAddress.state]);
    await userEvent.type(getByLabelText(/ZIP/), fakeAddress.postalCode);
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
        <ResidentialAddressForm {...testProps} />
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
