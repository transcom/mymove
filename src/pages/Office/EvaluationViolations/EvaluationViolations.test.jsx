/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';

import EvaluationViolations from './EvaluationViolations';

import { useEvaluationReportShipmentListQueries, usePWSViolationsQueries } from 'hooks/queries';

const mockMoveCode = 'A12345';
const mockMoveId = '551dd01f-90cf-44d6-addb-ff919433dd61';
const mockReportId = 'db30c135-1d6d-4a0d-a6d5-f408474f6ee2';
const mockViolationID = '9cdc8dc3-6cf4-46fb-b272-1468ef40796f';
const mockShipmentID = '319e0751-1337-4ed9-b4b5-a15d4e6d272c';

const mockEvaluationReport = {
  id: mockReportId,
  createdAt: '2022-09-07T15:17:37.484Z',
  eTag: 'MjAyMi0wOS0wN1QxODowNjozNy44NjQxNDJa',
  evalStart: '09:00',
  evalEnd: '13:00',
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

const mockShipmentData = {
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
};

jest.mock('hooks/queries', () => ({
  useEvaluationReportShipmentListQueries: jest.fn(),
  usePWSViolationsQueries: jest.fn(),
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
  // note: not using mockMoveCode or reportId here because of execution order
  useParams: jest.fn().mockReturnValue({ moveCode: 'A12345', reportId: 'db30c135-1d6d-4a0d-a6d5-f408474f6ee2' }),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

describe('EvaluationViolations', () => {
  it('renders the component content', async () => {
    usePWSViolationsQueries.mockImplementation(() => ({
      isLoading: false,
      isError: false,
      violations: [mockViolation],
    }));

    useEvaluationReportShipmentListQueries.mockImplementation(() => ({
      isLoading: false,
      isError: false,
      mtoShipments: [mockShipmentData],
      evaluationReport: mockEvaluationReport,
      reportViolations: [mockViolation],
    }));

    render(<EvaluationViolations {...{ customerInfo }} />);

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
