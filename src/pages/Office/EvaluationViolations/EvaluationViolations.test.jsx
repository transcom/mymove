/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import EvaluationViolations from './EvaluationViolations';

import { MockProviders } from 'testUtils';

const mockMoveCode = 'A12345';
const mockReportId = 'db30c135-1d6d-4a0d-a6d5-f408474f6ee2';
const mockMtoRefId = '6789-1234';
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

const mockUsePWSViolationsQueries = jest.fn();
jest.mock('hooks/queries', () => ({
  useShipmentEvaluationReportQueries: jest.fn(() => ({
    isLoading: false,
    isError: false,
    evaluationReport: {
      id: mockReportId,
      violationsObserved: true,
      type: 'SHIPMENT',
      moveReferenceID: mockMtoRefId,
    },
  })),
  usePWSViolationsQueries: () => mockUsePWSViolationsQueries(),
}));

const mockDeleteEvaluationReport = jest.fn();
jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  deleteEvaluationReport: (reportId) => mockDeleteEvaluationReport(reportId),
}));

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

describe('EvaluationViolations', () => {
  it('renders the component content', async () => {
    mockUsePWSViolationsQueries.mockImplementation(() => ({
      isLoading: false,
      isError: false,
      violations: [mockViolation],
    }));

    render(<EvaluationViolations />);

    // Check out heading
    expect(await screen.getByRole('heading', { name: 'Select violations', level: 2 })).toBeInTheDocument();

    // Violation Accordion is present
    expect(await screen.getByRole('heading', { name: 'Category 1', level: 3 })).toBeInTheDocument();
    expect(await screen.getByRole('heading', { name: 'SubCategory 1', level: 4 })).toBeInTheDocument();
    expect(await screen.getByRole('button', { name: 'SubCategory 1' })).toBeInTheDocument();

    // Verify Action/Naviation buttons
    expect(await screen.getByRole('button', { name: '< Back to Evaluation form' })).toBeInTheDocument();
    expect(await screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
    expect(await screen.getByRole('button', { name: 'Save draft' })).toBeInTheDocument();
    expect(await screen.getByRole('button', { name: 'Review and submit' })).toBeInTheDocument();
  });

  it('renders a violation accordion for each category', async () => {
    mockUsePWSViolationsQueries.mockImplementation(() => ({
      isLoading: false,
      isError: false,
      violations: [
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
      ],
    }));

    render(<EvaluationViolations />);

    // Content for each category is present
    expect(await screen.getByRole('heading', { name: 'Category 1', level: 3 })).toBeInTheDocument();
    expect(await screen.getByRole('heading', { name: 'SubCategory of Category 1', level: 4 })).toBeInTheDocument();
    expect(await screen.getByRole('heading', { name: 'Category 2', level: 3 })).toBeInTheDocument();
    expect(await screen.getByRole('heading', { name: 'SubCategory of Category 2', level: 4 })).toBeInTheDocument();
  });

  it('re-routes back to the eval report', async () => {
    render(
      <MockProviders initialEntries={[`/moves/${mockMoveCode}/evaluation-reports/${mockReportId}`]}>
        <EvaluationViolations />
      </MockProviders>,
    );

    // Click back button
    await userEvent.click(await screen.findByRole('button', { name: '< Back to Evaluation form' }));

    // Verify that we re-route back to the eval report
    expect(mockPush).toHaveBeenCalledTimes(1);
    expect(mockPush).toHaveBeenCalledWith(`/moves/${mockMoveCode}/evaluation-reports/${mockReportId}`);
  });
});
