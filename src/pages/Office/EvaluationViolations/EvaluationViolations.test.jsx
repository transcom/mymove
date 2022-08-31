/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';

import EvaluationViolations from './EvaluationViolations';

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
  useEvaluationReportQueries: jest.fn(() => ({
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

    await waitFor(() => {
      // Displays heading
      expect(screen.getByRole('heading', { name: 'Shipment report', level: 1 })).toBeInTheDocument();

      // Displays Evalutaion Violations Form
      expect(screen.getByRole('heading', { name: 'Select violations', level: 2 })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: '< Back to Evaluation form' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save draft' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Review and submit' })).toBeInTheDocument();
    });
  });
});
