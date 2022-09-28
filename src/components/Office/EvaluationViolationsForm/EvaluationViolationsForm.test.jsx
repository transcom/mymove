/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import EvaluationViolationsForm from './EvaluationViolationsForm';

import { MockProviders } from 'testUtils';

const mockMoveCode = 'A12345';
const mockReportId = 'db30c135-1d6d-4a0d-a6d5-f408474f6ee2';
const mockViolationID = '9cdc8dc3-6cf4-46fb-b272-1468ef40796f';
const mockViolation = {
  category: 'Category 1',
  displayOrder: 1,
  id: mockViolationID,
  paragraphNumber: '1.2.3',
  requirementStatement: 'Test requirement statement for violation 1',
  requirementSummary: 'Test requirement summary for violation 1',
  subCategory: 'SubCategory 1',
  title: 'Title for violation 1',
};

const mockEvaluationReport = {
  createdAt: '2022-09-07T15:17:37.484Z',
  eTag: 'MjAyMi0wOS0wN1QxODowNjozNy44NjQxNDJa',
  evaluationLengthMinutes: 240,
  inspectionDate: '2022-09-08',
  inspectionType: 'DATA_REVIEW',
  location: mockMoveCode,
  moveID: '551dd01f-90cf-44d6-addb-ff919433dd61',
  moveReferenceID: '4118-8295',
  officeUser: {
    email: 'qae_csr_role@office.mil',
    firstName: 'Leo',
    id: 'ef4f6d1f-4ac3-4159-a364-5403e7d958ff',
    lastName: 'Spaceman',
    phone: '415-555-1212',
  },
  remarks: 'test',
  shipmentID: '319e0751-1337-4ed9-b4b5-a15d4e6d272c',
  type: 'SHIPMENT',
  updatedAt: '2022-09-07T18:06:37.864Z',
  violationsObserved: true,
};

const mockReportViolation = {
  id: 'f3e2c135-336d-440d-a6d5-f404474f6ef3',
  reportId: mockReportId,
  violationId: mockViolationID,
  violations: [mockViolation],
};

const mockPush = jest.fn();
jest.mock('react-router', () => ({
  ...jest.requireActual('react-router'),
  useHistory: () => ({
    push: mockPush,
  }),
  useParams: jest.fn().mockReturnValue({ moveCode: 'A12345', reportId: 'db30c135-1d6d-4a0d-a6d5-f408474f6ee2' }),
}));

const mockSaveEvaluationReport = jest.fn();
const mockAssociateReportViolations = jest.fn();
jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  saveEvaluationReport: (options) => mockSaveEvaluationReport(options),
  associateReportViolations: (options) => mockAssociateReportViolations(options),
}));

const renderForm = (props) => {
  const defaultProps = {
    evaluationReport: mockEvaluationReport,
    reportViolations: [mockReportViolation],
    violations: [mockViolation],
  };

  return render(
    <MockProviders initialEntries={[`/moves/${mockMoveCode}/evaluation-reports/${mockReportId}`]}>
      <EvaluationViolationsForm {...defaultProps} {...props} />
    </MockProviders>,
  );
};

beforeEach(() => {
  jest.clearAllMocks();
});

describe('EvaluationViolationsForm', () => {
  it('renders the form content', async () => {
    renderForm();

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

    renderForm({ violations: mockTwoCategoryViolations });

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

    renderForm({ violations: mockKpiViolation });

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
    renderForm();
    // Click back button
    await userEvent.click(await screen.findByRole('button', { name: '< Back to Evaluation form' }));

    // Verify that we re-route back to the eval report
    expect(mockPush).toHaveBeenCalledTimes(1);
    expect(mockPush).toHaveBeenCalledWith(`/moves/${mockMoveCode}/evaluation-reports/${mockReportId}`);
  });

  it('can save a draft and reroute back to the eval reports', async () => {
    renderForm();

    // Click save draft button
    await userEvent.click(await screen.findByRole('button', { name: 'Save draft' }));

    // Verify that report was saved, violations re-associated with report, and page rerouted
    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledTimes(1);
      expect(mockSaveEvaluationReport).toHaveBeenCalledTimes(1);
      expect(mockSaveEvaluationReport).toHaveBeenCalledWith({
        body: {
          evaluationLengthMinutes: 240,
          inspectionDate: '2022-09-08',
          inspectionType: 'DATA_REVIEW',
          location: 'A12345',
          observedClaimsResponseDate: undefined,
          observedPickupDate: undefined,
          observedPickupSpreadEndDate: undefined,
          observedPickupSpreadStartDate: undefined,
          remarks: 'test',
          seriousIncident: undefined,
          seriousIncidentDesc: undefined,
          violationsObserved: true,
        },
        ifMatchETag: 'MjAyMi0wOS0wN1QxODowNjozNy44NjQxNDJa',
        reportID: 'db30c135-1d6d-4a0d-a6d5-f408474f6ee2',
      });
      expect(mockAssociateReportViolations).toHaveBeenCalledTimes(1);
      expect(mockPush).toHaveBeenCalledWith(`/moves/${mockMoveCode}/evaluation-reports`, {
        showSaveDraftSuccess: true,
      });
    });
  });
});
