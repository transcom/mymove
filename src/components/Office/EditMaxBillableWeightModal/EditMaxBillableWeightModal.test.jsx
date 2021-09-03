import React from 'react';
import { act } from 'react-dom/test-utils';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';

import EditMaxBillableWeightModal from './EditMaxBillableWeightModal';

let onClose;
let onSubmit;
beforeEach(() => {
  onClose = jest.fn();
  onSubmit = jest.fn();
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
    expect(await screen.getByRole('heading', { level: 4, name: 'Edit max billable weight' })).toBeInTheDocument();
    expect(await screen.getByText('Default:').parentElement).toBeInstanceOf(HTMLDListElement);
    expect(await screen.getByText('7,500 lbs').parentElement).toBeInstanceOf(HTMLDListElement);
    expect(await screen.getByLabelText('New max billable weight')).toBeInTheDocument();
    expect(await screen.getByDisplayValue('8,000 lbs')).toBeInTheDocument();
    expect(await screen.getByRole('button', { name: 'Save' })).toBeInTheDocument();
    expect(await screen.getByRole('button', { name: 'Back' })).toBeInTheDocument();
    expect(await screen.getByLabelText('Close')).toBeInstanceOf(HTMLButtonElement);
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

    await act(async () => {
      fireEvent.click(await screen.getByLabelText('Close'));
    });

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

    await act(async () => {
      fireEvent.click(await screen.getByRole('button', { name: 'Back' }));
    });

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

    await act(async () => {
      fireEvent.click(await screen.getByRole('button', { name: 'Save' }));
    });

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

    await act(async () => {
      fireEvent.change(await screen.getByDisplayValue('8,000 lbs'), { target: { value: '' } });
    });

    waitFor(() => {
      expect(screen.getByRole('alert', { name: 'Required' }));
      expect(screen.getByRole('button', { name: 'Save' })).toHaveAttribute('disabled');
    });
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

    await act(async () => {
      fireEvent.change(await screen.getByDisplayValue('8,000 lbs'), { target: { value: '0' } });
    });

    waitFor(() => {
      expect(screen.getByRole('alert', { name: 'Max billable weight must be greater than or equal to 1' }));
      expect(screen.getByRole('button', { name: 'Save' })).toHaveAttribute('disabled');
    });
  });
});
