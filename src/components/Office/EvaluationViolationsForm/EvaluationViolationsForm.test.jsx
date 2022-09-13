/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import EvaluationViolationsForm from './EvaluationViolationsForm';

import { MockProviders } from 'testUtils';

const mockMoveCode = 'A12345';
const mockReportId = 'db30c135-1d6d-4a0d-a6d5-f408474f6ee2';

const mockPush = jest.fn();
jest.mock('react-router', () => ({
  ...jest.requireActual('react-router'),
  useHistory: () => ({
    push: mockPush,
  }),
  useParams: jest.fn().mockReturnValue({ moveCode: 'A12345', reportId: 'db30c135-1d6d-4a0d-a6d5-f408474f6ee2' }),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

describe('EvaluationViolationsForm', () => {
  it('renders the form content', async () => {
    const mockViolation = {
      category: 'Category 1',
      displayOrder: 1,
      id: '9cdc8dc3-6cf4-46fb-b272-1468ef40796f',
      paragraphNumber: '1.2.3',
      requirementStatement: 'Test requirement statement for violation 1',
      requirementSummary: 'Test requirement summary for violation 1',
      subCategory: 'SubCategory 1',
      title: 'Title for violation 1',
    };

    render(<EvaluationViolationsForm violations={[mockViolation]} />);

    await waitFor(() => {
      // Check out headings
      expect(screen.getByRole('heading', { name: 'Select violations', level: 2 })).toBeInTheDocument();
      expect(screen.getByRole('heading', { name: 'Serious incident', level: 3 })).toBeInTheDocument();
      expect(screen.getByTestId('seriousIncidentLegend')).toHaveTextContent('Serious incident');

      // Violation Accordion is present
      expect(screen.getByRole('heading', { name: 'Category 1', level: 3 })).toBeInTheDocument();
      expect(screen.getByRole('heading', { name: 'SubCategory 1', level: 4 })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'SubCategory 1' })).toBeInTheDocument();

      // Verify Action/Naviation buttons
      expect(screen.getByRole('button', { name: '< Back to Evaluation form' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save draft' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Review and submit' })).toBeInTheDocument();

      // Conditionally shown textarea not shown by default
      expect(screen.queryByText('Serious incident description')).not.toBeInTheDocument();

      // Date pickers should not be shown by default
      expect(screen.queryByText('Observed claims response date')).not.toBeInTheDocument();
      expect(screen.queryByText('Observed pickup date')).not.toBeInTheDocument();
      expect(screen.queryByText('Observed pickup spread start date')).not.toBeInTheDocument();
      expect(screen.queryByText('Observed pickup spread end date')).not.toBeInTheDocument();
    });
  });

  it('renders a violation accordion for each category', async () => {
    const mockTwoCategoryViolations = [
      {
        category: 'Category 1',
        displayOrder: 1,
        id: '9cdc8dc3-6cf4-46fb-b272-1468ef40796f',
        paragraphNumber: '1.2.2',
        requirementStatement: 'Test requirement statement for violation 1',
        requirementSummary: 'Test requirement summary for violation 1',
        subCategory: 'SubCategory of Category 1',
        title: 'Title for violation 1',
      },
      {
        category: 'Category 2',
        displayOrder: 2,
        id: '4fdc8dc3-6cf4-46fb-b272-1468ef4079ab',
        paragraphNumber: '1.2.3',
        requirementStatement: 'Test requirement statement for violation 2',
        requirementSummary: 'Test requirement summary for violation 2',
        subCategory: 'SubCategory of Category 2',
        title: 'Title for violation 2',
      },
    ];

    render(<EvaluationViolationsForm violations={mockTwoCategoryViolations} />);

    // Content for each category is present
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Category 1', level: 3 })).toBeInTheDocument();
      expect(screen.getByRole('heading', { name: 'SubCategory of Category 1', level: 4 })).toBeInTheDocument();
      expect(screen.getByRole('heading', { name: 'Category 2', level: 3 })).toBeInTheDocument();
      expect(screen.getByRole('heading', { name: 'SubCategory of Category 2', level: 4 })).toBeInTheDocument();
    });
  });

  it('renders conditional datepicker when kpi violation is selected', async () => {
    const mockKpiViolation = [
      {
        additionalDataElem: 'observedPickupSpreadDates',
        category: 'Pre-Move Services',
        displayOrder: 7,
        isKpi: true,
        id: 'e1ee1719-a6d5-49b0-ad3b-c4dac0a3f16f',
        paragraphNumber: '1.2.5.3.1',
        requirementStatement: 'requirement statement 1',
        requirementSummary: 'Schedule relocation using pickup spread rules',
        subCategory: 'Counseling',
        title: 'Scheduling Requirements',
      },
    ];

    render(<EvaluationViolationsForm violations={mockKpiViolation} />);

    await waitFor(() => {
      expect(screen.queryByText('Observed pickup spread start date')).not.toBeInTheDocument();
      expect(screen.queryByText('Observed pickup spread end date')).not.toBeInTheDocument();
    });
    const checkbox = screen.getByTestId('violation-checkbox');
    userEvent.click(checkbox);

    await waitFor(() => {
      // Date picker should be shown if corresponding item is checked
      expect(screen.getByTestId('violation-checkbox')).toBeInTheDocument();
      expect(screen.getByText('Observed pickup spread start date')).toBeInTheDocument();
      expect(screen.getByText('Observed pickup spread end date')).toBeInTheDocument();
    });
  });

  it('re-routes back to the eval report', async () => {
    render(
      <MockProviders initialEntries={[`/moves/${mockMoveCode}/evaluation-reports/${mockReportId}`]}>
        <EvaluationViolationsForm violations={[]} />
      </MockProviders>,
    );

    // Click back button
    await userEvent.click(await screen.findByRole('button', { name: '< Back to Evaluation form' }));

    // Verify that we re-route back to the eval report
    expect(mockPush).toHaveBeenCalledTimes(1);
    expect(mockPush).toHaveBeenCalledWith(`/moves/${mockMoveCode}/evaluation-reports/${mockReportId}`);
  });
});
