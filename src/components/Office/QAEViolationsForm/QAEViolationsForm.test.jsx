/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import QAEViolationsForm from './QAEViolationsForm';

import { MockProviders } from 'testUtils';
import { saveEvaluationReport, associateReportViolations, submitEvaluationReport } from 'services/ghcApi';
import { useEvaluationReportShipmentListQueries } from 'hooks/queries';

const mockMoveCode = 'A12345';
const mockReportId = 'db30c135-1d6d-4a0d-a6d5-f408474f6ee2';
const mockMoveId = '551dd01f-90cf-44d6-addb-ff919433dd61';
const mockViolationID = '9cdc8dc3-6cf4-46fb-b272-1468ef40796f';
const mockShipmentID = '319e0751-1337-4ed9-b4b5-a15d4e6d272c';

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
  id: mockReportId,
  createdAt: '2022-09-07T15:17:37.484Z',
  eTag: 'MjAyMi0wOS0wN1QxODowNjozNy44NjQxNDJa',
  evalEnd: '09:00',
  evalStart: '10:00',
  inspectionDate: '2022-09-08',
  inspectionType: 'DATA_REVIEW',
  location: mockMoveCode,
  moveID: mockMoveId,
  moveReferenceID: '4118-8295',
  officeUser: {
    email: 'qae_csr_role@office.mil',
    firstName: 'Leo',
    id: 'ef4f6d1f-4ac3-4159-a364-5403e7d958ff',
    lastName: 'Spaceman',
    phone: '415-555-1212',
  },
  remarks: 'test',
  shipmentID: mockShipmentID,
  type: 'SHIPMENT',
  updatedAt: '2022-09-07T18:06:37.864Z',
  violationsObserved: true,
  seriousIncident: false,
};

const mockReportViolation = {
  id: 'f3e2c135-336d-440d-a6d5-f404474f6ef3',
  reportId: mockReportId,
  violationId: mockViolationID,
  violations: [mockViolation],
};

const customerInfo = {
  agency: 'ARMY',
  backup_contact: { email: 'email@example.com', name: 'name', phone: '555-555-5555' },
  current_address: {
    city: 'Beverly Hills',
    country: 'US',
    eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi4zMzIwOTFa',
    id: '28f11990-7ced-4d01-87ad-b16f2c86ea83',
    postalCode: '90210',
    state: 'CA',
    streetAddress1: '123 Any Street',
    streetAddress2: 'P.O. Box 12345',
    streetAddress3: 'c/o Some Person',
  },
  dodID: '5052247544',
  eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi4zNTkzNFo=',
  email: 'leo_spaceman_sm@example.com',
  first_name: 'Leo',
  id: 'ea557b1f-2660-4d6b-89a0-fb1b5efd2113',
  last_name: 'Spacemen',
  phone: '555-555-5555',
  userID: 'f4bbfcdf-ef66-4ce7-92f8-4c1bf507d596',
  grade: 'E_4',
};

const mockShipmentData = [
  {
    actualPickupDate: '2020-03-16',
    approvedDate: '2022-08-16T00:00:00.000Z',
    billableWeightCap: 4000,
    billableWeightJustification: 'heavy',
    createdAt: '2022-08-16T00:00:22.316Z',
    customerRemarks: 'Please treat gently',
    destinationAddress: {
      city: 'Fairfield',
      country: 'US',
      eTag: 'MjAyMi0wOC0xNlQwMDowMDoyMi4zMTQ0MDha',
      id: '1cf638df-1c1e-4c03-916a-bd20f7ec58ce',
      postalCode: '94535',
      state: 'CA',
      streetAddress1: '987 Any Avenue',
      streetAddress2: 'P.O. Box 9876',
      streetAddress3: 'c/o Some Person',
    },
    eTag: 'MjAyMi0wOC0xNlQwMDowMDoyMi4zMTY2N1o=',
    id: mockShipmentID,
    moveTaskOrderID: mockMoveId,
    pickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMi0wOC0xNlQwMDowMDoyMi4zMTIzOTZa',
      id: 'c0cf15bb-96ee-443a-982e-0e9ef2b9a80d',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    primeActualWeight: 2000,
    primeEstimatedWeight: 1400,
    requestedDeliveryDate: '2020-03-15',
    requestedPickupDate: '2020-03-15',
    scheduledPickupDate: '2020-03-16',
    shipmentType: 'HHG',
    status: 'APPROVED',
    updatedAt: '2022-08-16T00:00:22.316Z',
  },
];

const savedReportBody = {
  evalStart: mockEvaluationReport.evalStart,
  evalEnd: mockEvaluationReport.evalEnd,
  inspectionDate: mockEvaluationReport.inspectionDate,
  inspectionType: mockEvaluationReport.inspectionType,
  location: mockMoveCode,
  observedClaimsResponseDate: undefined,
  observedPickupDate: undefined,
  observedPickupSpreadEndDate: undefined,
  observedPickupSpreadStartDate: undefined,
  remarks: mockEvaluationReport.remarks,
  seriousIncident: false,
  seriousIncidentDesc: null,
  violationsObserved: mockEvaluationReport.violationsObserved,
};

jest.mock('hooks/queries', () => ({
  useEvaluationReportShipmentListQueries: jest.fn(),
}));

const mockPush = jest.fn();
jest.mock('react-router', () => ({
  ...jest.requireActual('react-router'),
  useHistory: () => ({
    push: mockPush,
  }),
  useParams: jest.fn().mockReturnValue({ moveCode: 'A12345', reportId: 'db30c135-1d6d-4a0d-a6d5-f408474f6ee2' }),
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  saveEvaluationReport: jest.fn(),
  associateReportViolations: jest.fn(),
  submitEvaluationReport: jest.fn(),
}));

const renderForm = (props) => {
  useEvaluationReportShipmentListQueries.mockReturnValue(mockShipmentData);
  const defaultProps = {
    evaluationReport: mockEvaluationReport,
    reportViolations: [mockReportViolation],
    violations: [mockViolation],
    customerInfo,
    mtoShipments: mockShipmentData,
  };

  return render(
    <MockProviders initialEntries={[`/moves/A12345/evaluation-reports/db30c135-1d6d-4a0d-a6d5-f408474f6ee2`]}>
      <QAEViolationsForm {...defaultProps} {...props} />
    </MockProviders>,
  );
};

beforeEach(() => {
  jest.clearAllMocks();
});

describe('QAEViolationsForm', () => {
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
      mockViolation,
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
      expect(screen.getByRole('heading', { name: mockViolation.category, level: 3 })).toBeInTheDocument();
      expect(screen.getByRole('heading', { name: mockViolation.subCategory, level: 4 })).toBeInTheDocument();
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
    await userEvent.click(checkbox);

    await waitFor(() => {
      // Date picker should be shown if corresponding item is checked
      expect(screen.getByTestId('violation-checkbox')).toBeInTheDocument();
      expect(screen.getByText('Observed pickup spread start date')).toBeInTheDocument();
      expect(screen.getByText('Observed pickup spread end date')).toBeInTheDocument();
    });
  });
});

describe('QAEViolationsForm Buttons', () => {
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
      expect(saveEvaluationReport).toHaveBeenCalledTimes(1);
      expect(saveEvaluationReport).toHaveBeenCalledWith({
        body: savedReportBody,
        ifMatchETag: mockEvaluationReport.eTag,
        reportID: mockReportId,
      });
      expect(associateReportViolations).toHaveBeenCalledTimes(1);
      expect(mockPush).toHaveBeenCalledWith(`/moves/${mockMoveCode}/evaluation-reports`, {
        showSaveDraftSuccess: true,
      });
    });
  });

  it('click the review and submit button, and see the report preivew with violations', async () => {
    renderForm();

    // Click save draft button
    await userEvent.click(await screen.findByTestId('reviewAndSubmit'));

    // Verify that report was saved, violations re-associated with report, and submission preview modal is rendered
    await waitFor(() => {
      expect(saveEvaluationReport).toHaveBeenCalledTimes(1);
      expect(saveEvaluationReport).toHaveBeenCalledWith({
        body: savedReportBody,
        ifMatchETag: mockEvaluationReport.eTag,
        reportID: mockReportId,
      });
      expect(associateReportViolations).toHaveBeenCalledTimes(1);

      expect(screen.getByText('Preview and submit shipment report')).toBeInTheDocument();

      // check that violations render in the preview
      expect(screen.getByRole('heading', { name: mockViolation.category, level: 3 })).toBeInTheDocument();
      expect(screen.getByRole('heading', { name: mockViolation.subCategory, level: 4 })).toBeInTheDocument();

      // check that back and submission buttons render
      expect(screen.getByRole('button', { name: '< Back to Evaluation form' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Submit' })).toBeInTheDocument();
    });
  });

  it('click the back button from the preview page and close the preview modal', async () => {
    renderForm();

    // Click save draft button
    await userEvent.click(await screen.findByTestId('reviewAndSubmit'));

    // Click back button
    await userEvent.click(await screen.findByTestId('backToEvalFromSubmit'));

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Select violations', level: 2 })).toBeInTheDocument();
    });
  });

  it('click the submit button from the preview page and close the preview modal', async () => {
    renderForm();

    // Click save draft button
    await userEvent.click(await screen.findByRole('button', { name: 'Review and submit' }));

    await userEvent.click(await screen.findByRole('button', { name: 'Submit' }));

    await waitFor(() => {
      expect(submitEvaluationReport).toHaveBeenCalledTimes(1);
      // Verify that we re-route back to the eval report
      expect(mockPush).toHaveBeenCalledTimes(1);
      expect(mockPush).toHaveBeenCalledWith(`/moves/${mockMoveCode}/evaluation-reports`, { showSubmitSuccess: true });
    });
  });
});
