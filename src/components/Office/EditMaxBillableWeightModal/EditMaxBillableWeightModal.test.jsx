import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import EditMaxBillableWeightModal from './EditMaxBillableWeightModal';

let onClose;
let onSubmit;
beforeEach(() => {
  onClose = jest.fn();
  onSubmit = jest.fn();
});

afterEach(() => {
  jest.clearAllMocks();
});

describe('EditMaxBillableWeightModal', () => {
  it('renders the component', async () => {
    render(
      <EditMaxBillableWeightModal
        onSubmit={onSubmit}
        onClose={onClose}
        defaultWeight={7500}
        maxBillableWeight={8000}
      />,
    );

    expect(await screen.findByRole('heading', { level: 4, name: 'Edit max billable weight' })).toBeInTheDocument();
    expect(screen.getByText('Default:').parentElement).toBeInstanceOf(HTMLDListElement);
    expect(screen.getByText('7,500 lbs').parentElement).toBeInstanceOf(HTMLDListElement);
    expect(screen.getByLabelText('New max billable weight')).toBeInTheDocument();
    expect(screen.getByDisplayValue('8,000 lbs')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Save' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Back' })).toBeInTheDocument();
    expect(screen.getByLabelText('Close')).toBeInstanceOf(HTMLButtonElement);
  });

  it('closes the modal when close icon is clicked', async () => {
    render(
      <EditMaxBillableWeightModal
        onSubmit={onSubmit}
        onClose={onClose}
        defaultWeight={7500}
        maxBillableWeight={8000}
      />,
    );

    await userEvent.click(await screen.getByLabelText('Close'));

    await waitFor(() => {
      expect(onClose.mock.calls.length).toBe(1);
    });
  });

  it('closes the modal when the cancel button is clicked', async () => {
    render(
      <EditMaxBillableWeightModal
        onSubmit={onSubmit}
        onClose={onClose}
        defaultWeight={7500}
        maxBillableWeight={8000}
      />,
    );

    await userEvent.click(await screen.getByRole('button', { name: 'Back' }));

    await waitFor(() => {
      expect(onClose.mock.calls.length).toBe(1);
    });
  });

  it('calls the submit function when submit button is clicked', async () => {
    render(
      <EditMaxBillableWeightModal
        onSubmit={onSubmit}
        onClose={onClose}
        defaultWeight={7500}
        maxBillableWeight={8000}
      />,
    );

    await userEvent.click(await screen.getByRole('button', { name: 'Save' }));

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalled();
    });
  });

  it('displays required validation error when max billable weight is empty', async () => {
    render(
      <EditMaxBillableWeightModal
        onSubmit={onSubmit}
        onClose={onClose}
        defaultWeight={7500}
        maxBillableWeight={8000}
      />,
    );

    await userEvent.type(await screen.getByDisplayValue('8,000 lbs'), '{selectall}{del}');

    expect(await screen.findByTestId('errorMessage')).toHaveTextContent('Required');
    expect(screen.getByRole('button', { name: 'Save' })).toHaveAttribute('disabled');
  });

  it('displays minimum validation error when max billable weight value is less than 1', async () => {
    render(
      <EditMaxBillableWeightModal
        onSubmit={onSubmit}
        onClose={onClose}
        defaultWeight={7500}
        maxBillableWeight={8000}
      />,
    );

    await userEvent.type(await screen.getByDisplayValue('8,000 lbs'), '{selectall}{del}0');

    expect(await screen.findByTestId('errorMessage')).toHaveTextContent(
      'Max billable weight must be greater than or equal to 1',
    );
    expect(screen.getByRole('button', { name: 'Save' })).toHaveAttribute('disabled');
  });
});
