import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ShipmentEvaluationForm from './ShipmentEvaluationForm';

import { MockProviders } from 'testUtils';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: jest.fn().mockReturnValue({ moveCode: 'LR4T8V', reportID: '58350bae-8e87-4e83-bd75-74027fb4333a' }),
}));

afterEach(() => {
  jest.clearAllMocks();
});

describe('ShipmentEvaluationForm', () => {
  it('renders the form components', async () => {
    render(
      <MockProviders initialEntries={['/moves/LR4T8V/evaluation-reports/58350bae-8e87-4e83-bd75-74027fb4333a']}>
        <ShipmentEvaluationForm />
      </MockProviders>,
    );

    // Headers
    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Evaluation form');
      const h3s = screen.getAllByRole('heading', { level: 3 });
      expect(h3s).toHaveLength(3);

      expect(screen.getByText('Evaluation information')).toBeInTheDocument();
      expect(screen.getByText('Violations')).toBeInTheDocument();
      expect(screen.getByText('QAE remarks')).toBeInTheDocument();

      // // Form components
      expect(screen.getByTestId('form')).toBeInTheDocument();

      expect(screen.getByText('Date of inspection')).toBeInTheDocument();
      expect(screen.getByText('Evaluation type')).toBeInTheDocument();
      expect(screen.getByText('Evaluation location')).toBeInTheDocument();
      expect(screen.getByText('Evaluation length')).toBeInTheDocument();
      expect(screen.getByText('Violations observed')).toBeInTheDocument();
      expect(screen.getByText('Evaluation remarks')).toBeInTheDocument();

      // Conditionally shown fields should not be displayed initially
      expect(screen.queryByText('Travel time to evaluation')).not.toBeInTheDocument();
      expect(screen.queryByText('Observed pickup date')).not.toBeInTheDocument();
      expect(screen.queryByText('Observed delivery date')).not.toBeInTheDocument();

      // Form buttons
      expect(screen.getByText('Cancel')).toBeInTheDocument();
      expect(screen.getByText('Save draft')).toBeInTheDocument();
      expect(screen.getByText('Review and submit')).toBeInTheDocument();
    });
  });

  it('renders conditionally form components correctly', async () => {
    render(
      <MockProviders initialEntries={['/moves/LR4T8V/evaluation-reports/58350bae-8e87-4e83-bd75-74027fb4333a']}>
        <ShipmentEvaluationForm />
      </MockProviders>,
    );

    // Initially no conditional fields shown
    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Evaluation form');
      expect(screen.getAllByTestId('textarea')).toHaveLength(1);

      expect(screen.queryByText('Travel time to evaluation')).not.toBeInTheDocument();
      expect(screen.queryByText('Observed pickup date')).not.toBeInTheDocument();
      expect(screen.queryByText('Observed delivery date')).not.toBeInTheDocument();
    });

    // Select Physical Evaluation type, should show Travel time to evaluation picker
    await userEvent.click(screen.getByText('Physical'));
    await waitFor(() => {
      expect(screen.getByText('Travel time to evaluation')).toBeInTheDocument();
      expect(screen.queryByText('Observed delivery date')).not.toBeInTheDocument();
      expect(screen.queryByText('Observed pickup date')).not.toBeInTheDocument();
      expect(screen.getAllByTestId('textarea')).toHaveLength(1);
    });

    // Select Eval locations and validate correct fields are shown
    await userEvent.click(screen.getByText('Origin'));
    await waitFor(() => {
      expect(screen.getByText('Observed pickup date')).toBeInTheDocument();
      expect(screen.queryByText('Observed delivery date')).not.toBeInTheDocument();
      expect(screen.getAllByTestId('textarea')).toHaveLength(1);
    });

    await userEvent.click(screen.getByText('Destination'));
    await waitFor(() => {
      expect(screen.getByText('Observed delivery date')).toBeInTheDocument();
      expect(screen.queryByText('Observed pickup date')).not.toBeInTheDocument();
      expect(screen.getAllByTestId('textarea')).toHaveLength(1);
    });

    await userEvent.click(screen.getByText('Other'));
    await waitFor(() => {
      expect(screen.queryByText('Observed delivery date')).not.toBeInTheDocument();
      expect(screen.queryByText('Observed pickup date')).not.toBeInTheDocument();
      expect(screen.getAllByTestId('textarea')).toHaveLength(2);
    });

    // If not 'Physical' eval type, no conditional time fields should be shown
    await userEvent.click(screen.getByText('Virtual'));
    await waitFor(() => {
      expect(screen.queryByText('Travel time to evaluation')).not.toBeInTheDocument();
      expect(screen.queryByText('Observed delivery date')).not.toBeInTheDocument();
      expect(screen.queryByText('Observed pickup date')).not.toBeInTheDocument();
    });
  });

  it('displays the delete confirmation on cancel', async () => {
    render(
      <MockProviders initialEntries={['/moves/LR4T8V/evaluation-reports/58350bae-8e87-4e83-bd75-74027fb4333a']}>
        <ShipmentEvaluationForm />
      </MockProviders>,
    );

    expect(await screen.getByRole('heading', { level: 2 })).toHaveTextContent('Evaluation form');

    // Buttons
    await userEvent.click(await screen.getByRole('button', { name: 'Cancel' }));

    expect(
      await screen.findByRole('heading', { level: 3, name: 'Are you sure you want to cancel this report?' }),
    ).toBeInTheDocument();
  });
});
