import React from 'react';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { act } from 'react-dom/test-utils';
import selectEvent from 'react-select-event';
import { renderWithRouter } from 'testUtils';

import ResidentialAddressForm from './ResidentialAddressForm';
import { searchLocationByZipCity } from 'services/internalApi';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  searchLocationByZipCity: jest.fn(),
}));

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

  it('renders the form inputs and help text', async () => {
    const { getByRole, getByLabelText, getByText } = render(<ResidentialAddressForm {...testProps} />);

    await waitFor(() => {
      expect(getByLabelText('Address 1')).toBeInstanceOf(HTMLInputElement);

      expect(getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);

      expect(getByRole('combobox', { id: 'zipCity-input' })).toBeInstanceOf(HTMLInputElement);

      expect(getByLabelText('City')).toBeInstanceOf(HTMLInputElement);

      expect(getByLabelText('State')).toBeInstanceOf(HTMLInputElement);

      expect(getByLabelText('County')).toBeInstanceOf(HTMLInputElement);

      expect(getByLabelText('ZIP')).toBeInstanceOf(HTMLInputElement);

      expect(getByText('Must be a physical address.')).toBeInTheDocument();
    });
  });

  it('shows an error message if trying to submit an invalid form', async () => {
    const mockAddress = {
      city: 'El Paso',
      state: 'TX',
      postalCode: '79912',
      county: 'El Paso',
    };
    const mockSearchLocationByZipCity= () => Promise.resolve(mockAddress);
    searchLocationByZipCity.mockImplementation(mockSearchLocationByZipCity);

    const { getByRole, findAllByRole, getByLabelText } = render(<ResidentialAddressForm {...testProps} />);
    await userEvent.click(getByLabelText('Address 1'));
    await userEvent.click(getByLabelText(/Address 2/));
    const input = getByRole('combobox', { id: 'zipCity-input' });
    fireEvent.change(input, {target: {value: mockAddress.postalCode} });
    await act(() => selectEvent.select(input, 79912));

    await waitFor(() => {
      expect(screen.getByText(mockAddress.postalCode)).toBeInTheDocument();
    });

    // fireEvent.keyPress(input, { key: 'Enter', code: 13 });

    await waitFor(() => {
      expect(searchLocationByZipCity).toHaveBeenCalledTimes(1);
    });
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
    const mockAddress = {
      city: 'El Paso',
      state: 'TX',
      postalCode: '79912',
      county: 'El Paso',
    };
    const mockSearchLocationByZipCity= () => Promise.resolve(mockAddress);
    searchLocationByZipCity.mockImplementation(mockSearchLocationByZipCity);

    // const { getByRole, getByLabelText } = render(<ResidentialAddressForm {...testProps} />);
    renderWithRouter(<ResidentialAddressForm {...testProps} />);
    const submitBtn = screen.getByRole('button', { name: 'Next' });

    await userEvent.type(screen.getByLabelText('Address 1'), fakeAddress.streetAddress1);
    await userEvent.type(screen.getByLabelText(/Address 2/), fakeAddress.streetAddress2);
    const input = screen.getByLabelText('Zip/City Lookup');
    await userEvent.type(input, '79912');
    fireEvent.change(input, {target: {value: mockAddress.postalCode} });
    // await waitFor(() => {
    //   expect(screen.getByText(mockAddress.state)).toBeInTheDocument();
    // });
    await act(() => selectEvent.select(input, mockAddress.postalCode));
    expect(screen.getByLabelText('State')).toHaveValue(mockAddress.state);
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
