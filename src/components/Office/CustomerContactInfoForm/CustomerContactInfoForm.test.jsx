import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import { Provider } from 'react-redux';

import CustomerContactInfoForm from './CustomerContactInfoForm';

import { roleTypes } from 'constants/userRoles';
import { configureStore } from 'shared/store';

describe('CustomerContactInfoForm Component', () => {
  const loggedInTOOState = {
    auth: {
      activeRole: roleTypes.TOO,
      isLoading: false,
      isLoggedIn: true,
    },
    entities: {
      user: {
        userId234: {
          id: 'userId234',
          roles: [{ roleType: roleTypes.TIO }],
        },
      },
    },
  };
  const initialValues = {
    firstName: '',
    middleName: '',
    lastName: '',
    suffix: '',
    customerTelephone: '',
    customerEmail: '',
    customerAddress: {
      streetAddress1: '123 Happy St',
      streetAddress2: 'Unit 4',
      city: 'Missoula',
      state: 'MT',
      postalCode: '59802',
    },
    name: '',
    telephone: '',
    email: '',
    cacUser: true,
  };
  const testProps = {
    initialValues,
    onSubmit: jest.fn(),
    onBack: jest.fn(),
  };

  const initialValuesCacValidated = {
    firstName: 'joe',
    middleName: 'bob',
    lastName: 'bob',
    suffix: 'jr',
    customerTelephone: '855-222-1111',
    customerEmail: 'joebob@gmail.com',
    customerAddress: {
      streetAddress1: '123 Happy St',
      streetAddress2: 'Unit 4',
      city: 'Missoula',
      state: 'MT',
      postalCode: '59802',
      county: 'MISSOULA',
    },
    name: 'joe bob',
    telephone: '855-222-1111',
    email: 'joebob@gmail.com',
    cacUser: null,
  };
  const testPropsCacValidated = {
    initialValuesCacValidated,
    onSubmit: jest.fn(),
    onBack: jest.fn(),
  };

  it('renders the form inputs', async () => {
    const mockStore = configureStore({
      ...loggedInTOOState,
    });
    render(
      <Provider store={mockStore.store}>
        <CustomerContactInfoForm {...testProps} />
      </Provider>,
    );

    await waitFor(() => {
      expect(screen.getByText('Contact info')).toBeInstanceOf(HTMLHeadingElement);
      expect(screen.getByLabelText('This is not the person named on the orders.')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('First name')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('First name')).toBeRequired();

      expect(screen.getByLabelText(/Middle name/)).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByLabelText('Last name')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Last name')).toBeRequired();

      expect(screen.getByLabelText(/Suffix/)).toBeInstanceOf(HTMLInputElement);

      expect(screen.getAllByLabelText('Phone')[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('Phone')[0]).toBeRequired();

      expect(screen.getAllByLabelText(/Alternate Phone/)[0]).toBeInstanceOf(HTMLInputElement);

      expect(screen.getAllByLabelText('Email')[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('Email')[0]).toBeRequired();

      expect(screen.getByText('Current Address')).toBeInstanceOf(HTMLHeadingElement);
      expect(screen.getByDisplayValue('123 Happy St')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByDisplayValue('Unit 4')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByDisplayValue('Missoula')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByDisplayValue('MT')).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getByDisplayValue('59802')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByLabelText('Name')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('Phone')[1]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('Email')[1]).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText('CAC Validation')).toBeInstanceOf(HTMLHeadingElement);
      expect(
        screen.getByText('Is the customer a non-CAC user or do they need to bypass CAC validation?'),
      ).toBeInTheDocument();
      expect(
        screen.getByText(
          'If this is checked yes, then the customer has already validated their CAC or their identity has been validated by a trusted office user.',
        ),
      ).toBeInTheDocument();
      expect(screen.getByTestId('cac-user-yes')).toBeInTheDocument();
      expect(screen.getByTestId('cac-user-no')).toBeInTheDocument();
    });
  });

  it('does not allow submission without cac_validated value', async () => {
    const mockStore = configureStore({
      ...loggedInTOOState,
    });
    render(
      <Provider store={mockStore.store}>
        <CustomerContactInfoForm {...testPropsCacValidated} />
      </Provider>,
    );

    await waitFor(() => {
      expect(screen.getByText('CAC Validation')).toBeInstanceOf(HTMLHeadingElement);
      expect(
        screen.getByText('Is the customer a non-CAC user or do they need to bypass CAC validation?'),
      ).toBeInTheDocument();
      expect(
        screen.getByText(
          'If this is checked yes, then the customer has already validated their CAC or their identity has been validated by a trusted office user.',
        ),
      ).toBeInTheDocument();
      expect(screen.getByTestId('cac-user-yes')).toBeInTheDocument();
      expect(screen.getByTestId('cac-user-no')).toBeInTheDocument();
    });
    expect(screen.getByRole('button', { name: 'Save' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Save' })).toBeDisabled();
  });
});
