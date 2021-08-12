import React from 'react';
import { act } from 'react-dom/test-utils';
import { render, fireEvent, waitFor, screen } from '@testing-library/react';

import EditMaxBillableWeightModal from './EditMaxBillableWeightModal';

let onClose;
let onSubmit;
beforeEach(() => {
  onClose = jest.fn();
  onSubmit = jest.fn();
});

describe('EditMaxBillableWeightModal', () => {
  it('renders the component', () => {
    render(
      <EditMaxBillableWeightModal
        onSubmit={onSubmit}
        onClose={onClose}
        defaultWeight={7500}
        maxBillableWeight={8000}
      />,
    );
    expect(screen.getByRole('heading', { level: 4, name: 'Edit max billable weight' }));
    expect(screen.getByText('Default:').parentElement).toBeInstanceOf(HTMLDListElement);
    expect(screen.getByText('7,500 lbs').parentElement).toBeInstanceOf(HTMLDListElement);
    expect(screen.getByLabelText('New max billable weight'));
    expect(screen.getByDisplayValue('8,000 lbs'));
    expect(screen.getByRole('button', { name: 'Save' }));
    expect(screen.getByRole('button', { name: 'Back' }));
    expect(screen.getByLabelText('Close')).toBeInstanceOf(HTMLButtonElement);
  });

  it('closes the modal when close icon is clicked', () => {
    render(
      <EditMaxBillableWeightModal
        onSubmit={onSubmit}
        onClose={onClose}
        defaultWeight={7500}
        maxBillableWeight={8000}
      />,
    );

    act(() => {
      fireEvent.click(screen.getByLabelText('Close'));
    });

    expect(onClose.mock.calls.length).toBe(1);
  });

  it('closes the modal when the cancel button is clicked', () => {
    render(
      <EditMaxBillableWeightModal
        onSubmit={onSubmit}
        onClose={onClose}
        defaultWeight={7500}
        maxBillableWeight={8000}
      />,
    );

    act(() => {
      fireEvent.click(screen.getByRole('button', { name: 'Back' }));
    });

    expect(onClose.mock.calls.length).toBe(1);
  });

  it('calls the submit function when submit button is clicked', () => {
    render(
      <EditMaxBillableWeightModal
        onSubmit={onSubmit}
        onClose={onClose}
        defaultWeight={7500}
        maxBillableWeight={8000}
      />,
    );

    act(() => {
      fireEvent.click(screen.getByRole('button', { name: 'Save' }));
    });

    waitFor(() => {
      expect(onSubmit).toHaveBeenCalled();
    });
  });

  it('displays required validation error when max billable weight is empty', () => {
    render(
      <EditMaxBillableWeightModal
        onSubmit={onSubmit}
        onClose={onClose}
        defaultWeight={7500}
        maxBillableWeight={8000}
      />,
    );

    fireEvent.change(screen.getByDisplayValue('8,000 lbs'), { target: { value: '' } });
    waitFor(() => {
      expect(screen.getByRole('alert', { name: 'Required' }));
      expect(screen.getByRole('button', { name: 'Save' })).toHaveAttribute('disabled');
    });
  });

  it('displays minimum validation error when max billable weight value is less than 1', () => {
    render(
      <EditMaxBillableWeightModal
        onSubmit={onSubmit}
        onClose={onClose}
        defaultWeight={7500}
        maxBillableWeight={8000}
      />,
    );

    fireEvent.change(screen.getByDisplayValue('8,000 lbs'), { target: { value: '0' } });
    waitFor(() => {
      expect(screen.getByRole('alert', { name: 'Max billable weight must be greater than or equal to 1' }));
      expect(screen.getByRole('button', { name: 'Save' })).toHaveAttribute('disabled');
    });
  });
});
