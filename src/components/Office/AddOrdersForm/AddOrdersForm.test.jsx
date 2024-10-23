import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import AddOrdersForm from './AddOrdersForm';

import { dropdownInputOptions } from 'utils/formatters';
import { ORDERS_PAY_GRADE_OPTIONS } from 'constants/orders';

describe('CreateMoveCustomerInfo Component', () => {
  const initialValues = {
    ordersType: '',
    issueDate: '',
    reportByDate: '',
    hasDependents: '',
    newDutyLocation: '',
    grade: '',
    originDutyLocation: '',
  };
  const testProps = {
    initialValues,
    ordersTypeOptions: dropdownInputOptions(ORDERS_PAY_GRADE_OPTIONS),
    onSubmit: jest.fn(),
    onBack: jest.fn(),
  };

  it('renders the form inputs', async () => {
    render(<AddOrdersForm {...testProps} />);

    await waitFor(() => {
      expect(screen.getByText('Tell us about the orders')).toBeInTheDocument();
      expect(screen.getByLabelText('Orders type')).toBeInTheDocument();
      expect(screen.getByLabelText('Orders date')).toBeInTheDocument();
      expect(screen.getByLabelText('Report by date')).toBeInTheDocument();
      expect(screen.getByText('Are dependents included in the orders?')).toBeInTheDocument();
      expect(screen.getByTestId('hasDependentsYes')).toBeInTheDocument();
      expect(screen.getByTestId('hasDependentsNo')).toBeInTheDocument();
      expect(screen.getByLabelText('Current duty location')).toBeInTheDocument();
      expect(screen.getByLabelText('New duty location')).toBeInTheDocument();
      expect(screen.getByLabelText('Pay grade')).toBeInTheDocument();
    });
  });

  it('shows an error message if trying to submit an invalid form', async () => {
    const { getByRole, findAllByRole, getByLabelText } = render(<AddOrdersForm {...testProps} />);
    await userEvent.click(getByLabelText('Orders type'));
    await userEvent.click(getByLabelText('Orders date'));
    await userEvent.click(getByLabelText('Report by date'));
    await userEvent.click(getByLabelText('Current duty location'));
    await userEvent.click(getByLabelText('New duty location'));
    await userEvent.click(getByLabelText('Pay grade'));

    const submitBtn = getByRole('button', { name: 'Next' });
    await userEvent.click(submitBtn);

    const alerts = await findAllByRole('alert');
    expect(alerts.length).toBe(4);

    alerts.forEach((alert) => {
      expect(alert).toHaveTextContent('Required');
    });

    expect(testProps.onSubmit).not.toHaveBeenCalled();
  });
});
