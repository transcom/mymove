import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';

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
      expect(screen.getByLabelText('Current duty location')).toBeInTheDocument();
      expect(screen.getByLabelText('New duty location')).toBeInTheDocument();
      expect(screen.getByLabelText('Pay grade')).toBeInTheDocument();
    });
  });
});
