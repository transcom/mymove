import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Formik } from 'formik';

import ShipmentIncentiveAdvance from './ShipmentIncentiveAdvance';

describe('components/Office/ShipmentIncentiveAdvance', () => {
  it('should display content without props', () => {
    render(
      <Formik initialValues={{}}>
        <ShipmentIncentiveAdvance />
      </Formik>,
    );

    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Incentive & advance');
    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('Estimated incentive: $0');
    expect(screen.getByLabelText('No')).toBeChecked();
    expect(screen.getByLabelText('Yes')).not.toBeChecked();

    expect(screen.queryByLabelText('Amount requested')).not.toBeInTheDocument();
  });

  it('should respond to user radio button input', async () => {
    render(
      <Formik initialValues={{}}>
        <ShipmentIncentiveAdvance />
      </Formik>,
    );

    userEvent.click(screen.getByLabelText('Yes'));

    expect(await screen.findByLabelText('Amount requested')).toBeInTheDocument();
  });

  it('should respond to props and form values', async () => {
    render(
      <Formik initialValues={{ advanceRequested: true, amountRequested: '7000' }}>
        <ShipmentIncentiveAdvance estimatedIncentive={1111111} />
      </Formik>,
    );

    expect(await screen.findByLabelText('Amount requested')).toHaveValue('7,000');
    expect(screen.getByText('Reminder: your advance can not be more than $6,666.67')).toBeInTheDocument();
  });
});
