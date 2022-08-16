import React from 'react';
import { render, screen } from '@testing-library/react';

import EvaluationReportMoveInfo from './EvaluationReportMoveInfo';

import { ORDERS_BRANCH_OPTIONS } from 'constants/orders';

describe('EvaluationReportMoveInfo', () => {
  it('renders the correct content', async () => {
    const mockCustomerInfo = {
      last_name: 'Smith',
      first_name: 'John',
      phone: '+441234567890',
      email: 'abc@123.com',
      agency: ORDERS_BRANCH_OPTIONS.NAVY,
    };

    const mockOrders = {
      grade: 'E_1',
    };

    render(<EvaluationReportMoveInfo customerInfo={mockCustomerInfo} orders={mockOrders} />);

    expect(screen.getByRole('heading', { name: 'Move information', level: 2 })).toBeInTheDocument();
    expect(screen.getByText('Customer information')).toBeInTheDocument();
    expect(screen.getByText('QAE')).toBeInTheDocument();
  });
});
