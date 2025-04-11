import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ConnectedContactInfo from './ContactInfo';

import { MockProviders } from 'testUtils';
import { patchOfficeUser } from 'services/ghcApi';
import { officeRoutes } from 'constants/routes';

const mockNavigate = jest.fn();
jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  patchOfficeUser: jest.fn(),
}));
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

afterEach(jest.resetAllMocks);

describe('ContactInfo Component', () => {
  const props = {
    setFlashMessage: jest.fn(),
  };
  const mockState = {
    entities: {
      user: {
        userId123: {
          id: 'userId123',
          office_user: {
            id: '123',
            first_name: 'John',
            middle_name: 'M',
            last_name: 'Doe',
            telephone: '804-456-7890',
            email: 'john.doe@example.com',
          },
        },
      },
    },
  };

  it('renders the ContactInfo component correctly', () => {
    render(
      <MockProviders initialState={mockState}>
        <ConnectedContactInfo {...props} />
      </MockProviders>,
    );

    expect(screen.getByRole('heading', { name: 'Edit contact info' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Your contact info' })).toBeInTheDocument();
    expect(screen.getByLabelText('First name')).toBeInTheDocument();
    expect(screen.getByLabelText('Middle name')).toBeInTheDocument();
    expect(screen.getByLabelText('Last name')).toBeInTheDocument();
    expect(screen.getByLabelText('Email')).toBeInTheDocument();
    expect(screen.getByLabelText('Phone *')).toBeInTheDocument();
  });

  it('navigates back to profile when the cancel button is clicked', () => {
    render(
      <MockProviders initialState={mockState}>
        <ConnectedContactInfo {...props} />
      </MockProviders>,
    );

    const cancelButton = screen.getByRole('button', { name: 'Cancel' });
    fireEvent.click(cancelButton);

    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();

    expect(mockNavigate).toHaveBeenCalledWith(officeRoutes.PROFILE_PATH);
  });

  it('redirects to profile page when submission is successful', async () => {
    patchOfficeUser.mockResolvedValueOnce({});
    render(
      <MockProviders initialState={mockState}>
        <ConnectedContactInfo {...props} />
      </MockProviders>,
    );

    expect(screen.getByRole('button', { name: 'Save' })).toBeEnabled();

    await userEvent.click(screen.getByRole('button', { name: 'Save' }));

    await waitFor(() => {
      expect(patchOfficeUser).toHaveBeenCalledWith('123', { telephone: '804-456-7890' });
      expect(mockNavigate).toHaveBeenCalledWith(officeRoutes.PROFILE_PATH);
    });
  });

  it('displays an error message when submission fails', async () => {
    const errorMessage = 'Server error';
    patchOfficeUser.mockRejectedValueOnce({ response: { body: { detail: errorMessage } } });

    render(
      <MockProviders initialState={mockState}>
        <ConnectedContactInfo {...props} />
      </MockProviders>,
    );

    const submitButton = screen.getByRole('button', { name: /save/i });
    await userEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: /an error occurred/i })).toBeInTheDocument();
      expect(
        screen.getByText(`Failed to update contact info due to server error: ${errorMessage}`),
      ).toBeInTheDocument();
    });
  });
});
