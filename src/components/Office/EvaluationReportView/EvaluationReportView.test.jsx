import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import userEvent from '@testing-library/user-event';

import EvaluationReportView from './EvaluationReportView';

import { useEvaluationReportShipmentListQueries } from 'hooks/queries';
import { qaeCSRRoutes } from 'constants/routes';
import { MockProviders, renderWithProviders } from 'testUtils';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

const mockReportId = 'db30c135-1d6d-4a0d-a6d5-f408474f6ee2';
const mockMoveId = '551dd01f-90cf-44d6-addb-ff919433dd61';
const mockViolationID = '9cdc8dc3-6cf4-46fb-b272-1468ef40796f';
const mockViolationID2 = '9cdc8dc3-6cf4-46fb-b272-1468ef40796g';
const mockShipmentID = '319e0751-1337-4ed9-b4b5-a15d4e6d272c';

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

const mockEvaluationReport = {
  id: mockReportId,
  createdAt: '2022-09-07T15:17:37.484Z',
  eTag: 'MjAyMi0wOS0wN1QxODowNjozNy44NjQxNDJa',
  evalEnd: '09:00',
  evalStart: '10:00',
  inspectionDate: '2022-09-08',
  inspectionType: 'VIRTUAL',
  location: 'ORIGIN',
  moveID: mockMoveId,
  moveReferenceID: '4118-8295',
  observedPickupDate: '2024-08-24',
  officeUser: {
    email: 'qae_role@office.mil',
    firstName: 'Leo',
    id: 'ef4f6d1f-4ac3-4159-a364-5403e7d958ff',
    lastName: 'Spaceman',
    phone: '415-555-1212',
  },
  remarks: 'test remarks',
  shipmentID: mockShipmentID,
  type: 'SHIPMENT',
  updatedAt: '2022-09-07T18:06:37.864Z',
  violationsObserved: true,
  seriousIncident: true,
  seriousIncidentDesc: 'there was a serious incident',
};

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

const mockViolation2 = {
  category: 'Category 1',
  displayOrder: 2,
  id: mockViolationID2,
  paragraphNumber: '1.2.4',
  requirementStatement: 'Test requirement statement for violation 2',
  requirementSummary: 'Test requirement summary for violation 2',
  subCategory: 'SubCategory 2',
  title: 'Title for violation 2',
};

const mockReportViolations = [
  {
    id: 'f3e2c135-336d-440d-a6d5-f404474f6ef3',
    reportId: mockReportId,
    violationId: mockViolationID,
    violation: mockViolation,
  },
];

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

const mockReportViolationsWithAppeals = [
  {
    id: 'superUniqueID',
    reportId: mockReportId,
    violationId: mockViolationID,
    violation: mockViolation,
    gsrAppeals: [
      {
        appealStatus: 'SUSTAINED',
        createdAt: '2024-10-23T16:41:20.514Z',
        id: '5b3fdda8-feb3-4b78-adb2-d164800aa1dc',
        officeUser: {
          firstName: 'Billy',
          lastName: 'Bob',
        },
        remarks: 'remarkable remarks',
        reportID: mockReportId,
        violationID: mockViolationID,
      },
    ],
  },
  {
    id: 'superUniqueID2',
    reportId: mockReportId,
    violationId: mockViolationID2,
    violation: mockViolation2,
  },
];

const mockReturnDataWithAppeals = {
  evaluationReport: mockEvaluationReport,
  isError: false,
  isLoading: false,
  isSuccess: true,
  mtoShipments: mockShipmentData,
  reportViolations: mockReportViolationsWithAppeals,
};

const mockReturnData = {
  evaluationReport: mockEvaluationReport,
  isError: false,
  isLoading: false,
  isSuccess: true,
  mtoShipments: mockShipmentData,
  reportViolations: mockReportViolations,
};

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

const mockRoutingParams = { moveCode: 'A12345', reportId: 'db30c135-1d6d-4a0d-a6d5-f408474f6ee2' };
const mockRoutingConfig = { path: qaeCSRRoutes.EVALUATION_REPORT_VIEW_PATH, params: mockRoutingParams };

jest.mock('hooks/queries', () => ({
  useEvaluationReportShipmentListQueries: jest.fn(),
}));

const renderForm = (props) => {
  useEvaluationReportShipmentListQueries.mockReturnValue(mockReturnData);
  const defaultProps = {
    customerInfo,
    grade: 'E_4',
    destinationDutyLocationPostalCode: '90210',
  };

  return renderWithProviders(<EvaluationReportView {...defaultProps} {...props} />, mockRoutingConfig);
};

const renderFormWithAppeals = (props) => {
  useEvaluationReportShipmentListQueries.mockReturnValue(mockReturnDataWithAppeals);
  const defaultProps = {
    customerInfo,
    grade: 'E_4',
    destinationDutyLocationPostalCode: '90210',
  };

  return renderWithProviders(<EvaluationReportView {...defaultProps} {...props} />, mockRoutingConfig);
};

describe('EvaluationReportView', () => {
  it('renders loading placeholder', () => {
    useEvaluationReportShipmentListQueries.mockReturnValue({
      isLoading: true,
      isError: false,
      isSuccess: false,
      evaluationReport: {},
      reportViolations: [],
      mtoShipments: [],
    });

    render(
      <MockProviders>
        <EvaluationReportView customerInfo={{}} grade="E-5" destinationDutyLocationPostalCode="12345" />
      </MockProviders>,
    );

    expect(screen.getByText('Loading, please wait...')).toBeInTheDocument();
  });

  it('renders Something Went Wrong page', () => {
    useEvaluationReportShipmentListQueries.mockReturnValue({
      isLoading: false,
      isError: true,
      isSuccess: false,
      evaluationReport: {},
      reportViolations: [],
      mtoShipments: [],
    });

    render(
      <MockProviders>
        <EvaluationReportView customerInfo={{}} grade="E-5" destinationDutyLocationPostalCode="12345" />
      </MockProviders>,
    );

    const errorMessage = screen.getByText(/Something went wrong./);
    expect(errorMessage).toBeInTheDocument();
  });

  it('renders the evaluation report with all details', () => {
    renderForm();

    // Check for basic report details
    expect(screen.getByText('Evaluation report')).toBeInTheDocument();
    expect(screen.getByText('Information')).toBeInTheDocument();
    expect(screen.getByText('Scheduled pickup')).toBeInTheDocument();
    expect(screen.getByText('Observed pickup')).toBeInTheDocument();
    expect(screen.getByText('Inspection date')).toBeInTheDocument();
    expect(screen.getByText('Report submission')).toBeInTheDocument();
    expect(screen.getByText('Inspection date')).toBeInTheDocument();
    expect(screen.getByText('Time evaluation started')).toBeInTheDocument();
    expect(screen.getByText('Time evaluation ended')).toBeInTheDocument();

    // Check for violations
    expect(screen.getByText('Violations observed')).toBeInTheDocument();
    expect(screen.getByText('1.2.3 Title for violation 1')).toBeInTheDocument();
    expect(screen.getByText('Test requirement summary for violation 1')).toBeInTheDocument();
    expect(screen.getByText('Observed Pickup Date')).toBeInTheDocument();

    // Check for serious incident details
    expect(screen.getByText('Serious Incident')).toBeInTheDocument();
    expect(screen.getByTestId('seriousIncidentYesNo')).toHaveTextContent('Yes');
    expect(screen.getByText('there was a serious incident')).toBeInTheDocument();

    // Check for QAE remarks
    expect(screen.getByText('QAE remarks')).toBeInTheDocument();
    expect(screen.getByText('Evaluation remarks')).toBeInTheDocument();
    expect(screen.getByText('test remarks')).toBeInTheDocument();

    expect(screen.getByTestId('backBtn')).toBeInTheDocument();
  });

  it('displays no violations observed when there are none', () => {
    useEvaluationReportShipmentListQueries.mockReturnValue({
      isLoading: false,
      isError: false,
      isSuccess: true,
      evaluationReport: {
        violationsObserved: false,
      },
      reportViolations: [],
      mtoShipments: [],
    });

    render(
      <MockProviders>
        <EvaluationReportView customerInfo={{}} grade="E-5" destinationDutyLocationPostalCode="12345" />
      </MockProviders>,
    );

    expect(screen.getByTestId('noViolationsObserved')).toHaveTextContent('No');
  });

  it('displays no serious incident data when there is not one', () => {
    useEvaluationReportShipmentListQueries.mockReturnValue({
      isLoading: false,
      isError: false,
      isSuccess: true,
      evaluationReport: {
        seriousIncident: false,
      },
      reportViolations: [],
      mtoShipments: [],
    });

    render(
      <MockProviders>
        <EvaluationReportView customerInfo={{}} grade="E-5" destinationDutyLocationPostalCode="12345" />
      </MockProviders>,
    );

    expect(screen.getByTestId('seriousIncidentYesNo')).toHaveTextContent('No');
  });

  it('allows a GSR user to see appeals attached to a violation', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    renderFormWithAppeals();

    expect(screen.queryByText('Billy Bob')).not.toBeInTheDocument();
    expect(screen.queryByText('remarkable remarks')).not.toBeInTheDocument();

    expect(screen.getByText('Show appeals')).toBeInTheDocument();
    const expandBtn = await screen.findByTestId('showAppealBtn');
    await userEvent.click(expandBtn);

    expect(screen.queryByText('Billy Bob')).toBeInTheDocument();
    expect(screen.queryByText('remarkable remarks')).toBeInTheDocument();
  });

  it('allows a GSR user to add an appeal to a violation', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    renderFormWithAppeals();

    const addAppealBtn = await screen.findByTestId('addAppealBtn');
    expect(addAppealBtn).toBeInTheDocument();
    await userEvent.click(addAppealBtn);

    // modal heading should appear
    expect(screen.getByText('Add Violation Appeal')).toBeInTheDocument();
  });
});
