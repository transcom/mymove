import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ValidationCode from './ValidationCode';

import { MockProviders, renderWithProviders } from 'testUtils';
import { validateCode } from 'services/internalApi';

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  validateCode: jest.fn().mockImplementation(() => Promise.resolve()),
}));

afterEach(() => {
  jest.resetAllMocks();
});

const mockSubmit = jest.fn();

describe('ValidationCode', () => {
  test('it should render all text for the component and asterisks for required fields', async () => {
    renderWithProviders(<ValidationCode onSuccess={jest.fn()} />);

    expect(document.querySelector('#reqAsteriskMsg')).toHaveTextContent('Fields marked with * are required.');

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
        <ValidationCode onSuccess={mockSubmit} />
      </MockProviders>,
    );

    const nextBtn = await screen.findByRole('button', { name: 'Next' });
    expect(nextBtn).toBeDisabled();
    await userEvent.type(screen.getByLabelText('Validation code *'), 'TestCode123123');
    expect(screen.getByLabelText('Validation code *')).toBeRequired();

    expect(nextBtn).toBeEnabled();
    await userEvent.click(nextBtn);

    expect(mockSubmit).toHaveBeenCalled();
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
        <ValidationCode onSuccess={mockSubmit} />
      </MockProviders>,
    );

    const nextBtn = await screen.findByRole('button', { name: 'Next' });
    expect(nextBtn).toBeDisabled();
    await userEvent.type(screen.getByLabelText('Validation code *'), 'TestCode123123');
    expect(nextBtn).toBeEnabled();
    await userEvent.click(nextBtn);

    expect(mockSubmit).not.toHaveBeenCalled();
    expect(screen.getByText('Incorrect validation code')).toBeInTheDocument();
    expect(screen.getByText('Please try again')).toBeInTheDocument();
  });
});
