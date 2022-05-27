import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Formik } from 'formik';
import * as Yup from 'yup';

import ShipmentIncentiveAdvance from './ShipmentIncentiveAdvance';

import { getFormattedMaxAdvancePercentage } from 'utils/incentives';

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
    const estimatedIncentive = 1111111;
    const validationSchema = Yup.object().shape({
      advance: Yup.number().max(
        (estimatedIncentive * 0.6) / 100,
        `Enter an amount that is less than or equal to the maximum advance (${getFormattedMaxAdvancePercentage()} of estimated incentive)`,
      ),
    });

    render(
      <Formik validationSchema={validationSchema} initialValues={{ advanceRequested: true, advance: '7000' }}>
        <ShipmentIncentiveAdvance estimatedIncentive={estimatedIncentive} />
      </Formik>,
    );

    expect(await screen.findByLabelText('Amount requested')).toHaveValue('7,000');
    expect(
      screen.getByText(
        'Enter an amount that is less than or equal to the maximum advance (60% of estimated incentive)',
      ),
    ).toBeInTheDocument();
    expect(screen.getByText('Maximum advance: $6,666')).toBeInTheDocument();
  });
});
