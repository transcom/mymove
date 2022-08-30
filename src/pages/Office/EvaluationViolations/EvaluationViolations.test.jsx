/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import EvaluationViolations from './EvaluationViolations';

import { MockProviders } from 'testUtils';

const mockMoveCode = 'A12345';
const mockReportId = 'db30c135-1d6d-4a0d-a6d5-f408474f6ee2';
const mockMtoRefId = '6789-1234';

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
    render(<EvaluationViolations />);

    // Check out heading
    expect(await screen.getByRole('heading', { name: 'Select violations', level: 2 })).toBeInTheDocument();

    // Check out buttons
    const buttons = await screen.getAllByRole('button');
    expect(buttons).toHaveLength(4);
    expect(buttons[0]).toHaveTextContent('< Back to Evaluation form');
    expect(buttons[1]).toHaveTextContent('Cancel');
    expect(buttons[2]).toHaveTextContent('Save draft');
    expect(buttons[3]).toHaveTextContent('Review and submit');
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
