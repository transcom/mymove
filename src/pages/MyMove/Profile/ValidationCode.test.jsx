import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ValidationCode from './ValidationCode';

import { MockProviders, renderWithProviders } from 'testUtils';
import { customerRoutes } from 'constants/routes';
import { validateCode } from 'services/internalApi';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  validateCode: jest.fn().mockImplementation(() => Promise.resolve()),
}));

afterEach(() => {
  jest.resetAllMocks();
});

describe('ValidationCode', () => {
  test('it should render all text for the component', async () => {
    renderWithProviders(<ValidationCode />);

    expect(screen.getByText('Please enter a validation code to begin creating a move')).toBeInTheDocument();
    const nextBtn = await screen.findByRole('button', { name: 'Next' });
    expect(nextBtn).toBeDisabled();
  });

  test('it navigates the user when entering a successful code', async () => {
    validateCode.mockImplementation(() =>
      Promise.resolve({
        body: {
          parameterValue: 'TestCode123123',
          parameterName: 'validation_code',
        },
      }),
    );

    render(
      <MockProviders>
        <ValidationCode />
      </MockProviders>,
    );

    const nextBtn = await screen.findByRole('button', { name: 'Next' });
    expect(nextBtn).toBeDisabled();
    await userEvent.type(screen.getByLabelText('Validation code'), 'TestCode123123');
    expect(nextBtn).toBeEnabled();
    await userEvent.click(nextBtn);

    expect(mockNavigate).toHaveBeenCalledWith(customerRoutes.DOD_INFO_PATH);
  });

  test('it displays error when code is not correct', async () => {
    validateCode.mockImplementation(() =>
      Promise.resolve({
        body: {
          parameterValue: '',
          parameterName: 'validation_code',
        },
      }),
    );

    render(
      <MockProviders>
        <ValidationCode />
      </MockProviders>,
    );

    const nextBtn = await screen.findByRole('button', { name: 'Next' });
    expect(nextBtn).toBeDisabled();
    await userEvent.type(screen.getByLabelText('Validation code'), 'TestCode123123');
    expect(nextBtn).toBeEnabled();
    await userEvent.click(nextBtn);

    expect(mockNavigate).not.toHaveBeenCalled();
    expect(screen.getByText('Incorrect validation code')).toBeInTheDocument();
    expect(screen.getByText('Please try again')).toBeInTheDocument();
  });
});
