/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import ReactRouter from 'react-router';

import EvaluationReport from './EvaluationReport';

import { MockProviders } from 'testUtils';
import { useShipmentEvaluationReportQueries } from 'hooks/queries';

const mockReportId = 'db30c135-1d6d-4a0d-a6d5-f408474f6ee2';
const mockMtoRefId = '6789-1234';
const mockCustomer = {
  last_name: 'spaceman',
  first_name: 'leo',
  phone: '555-555-5555',
};
const mockOrders = {
  grade: 'E_1',
  agency: 'COAST_GUARD',
};

jest.mock('hooks/queries', () => ({
  useShipmentEvaluationReportQueries: jest.fn(),
}));

const mockShipmentData = {
  isLoading: false,
  isError: false,
  evaluationReport: {
    id: mockReportId,
    type: 'SHIPMENT',
    moveReferenceID: mockMtoRefId,
  },
  mtoShipment: {
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
    id: 'c37ccf04-637c-4afc-9ef6-dee1555e16ef',
    moveTaskOrderID: '35eb1c36-8916-46f4-a72a-32267c9cb070',
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
};

const mockCounselingData = {
  isLoading: false,
  isError: false,
  evaluationReport: {
    id: mockReportId,
    type: 'COUNSELING',
    moveReferenceID: mockMtoRefId,
  },
  mtoShipment: {
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
    id: 'c37ccf04-637c-4afc-9ef6-dee1555e16ef',
    moveTaskOrderID: '35eb1c36-8916-46f4-a72a-32267c9cb070',
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
};

describe('EvaluationReport', () => {
  it('renders the page components for shipment report', async () => {
    useShipmentEvaluationReportQueries.mockReturnValue(mockShipmentData);
    jest.spyOn(ReactRouter, 'useParams').mockReturnValue({ reportId: mockReportId });

    render(
      <MockProviders initialEntries={[`/moves/LR4T8V/evaluation-reports/${mockReportId}`]}>
        <EvaluationReport customerInfo={mockCustomer} orders={mockOrders} />
      </MockProviders>,
    );

    await waitFor(() => {
      // Page content sections
      expect(screen.getByRole('heading', { name: 'Shipment information', level: 2 })).toBeInTheDocument();
      expect(screen.getByText('Customer information')).toBeInTheDocument();
      expect(screen.getByText('QAE')).toBeInTheDocument();

      expect(screen.getByRole('heading', { name: 'Evaluation form', level: 2 })).toBeInTheDocument();
      expect(screen.getByRole('heading', { name: 'Evaluation information', level: 3 })).toBeInTheDocument();
    });

    // Buttons
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
    expect(screen.getByText('Save draft')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Review and submit' })).toBeInTheDocument();
  });

  it('renders the page components for counseling report', async () => {
    useShipmentEvaluationReportQueries.mockReturnValue(mockCounselingData);
    jest.spyOn(ReactRouter, 'useParams').mockReturnValue({ reportId: mockReportId });
    render(
      <MockProviders initialEntries={[`/moves/LR4T8V/evaluation-reports/${mockReportId}`]}>
        <EvaluationReport customerInfo={mockCustomer} orders={mockOrders} />
      </MockProviders>,
    );

    await waitFor(() => {
      // Page content sections
      expect(screen.getByRole('heading', { name: 'Move information', level: 2 })).toBeInTheDocument();
      expect(screen.getByText('Customer information')).toBeInTheDocument();
      expect(screen.getByText('QAE')).toBeInTheDocument();

      expect(screen.getByRole('heading', { name: 'Evaluation form', level: 2 })).toBeInTheDocument();
      expect(screen.getByRole('heading', { name: 'Evaluation information', level: 3 })).toBeInTheDocument();

      // Buttons
      expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save draft' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Review and submit' })).toBeInTheDocument();
    });
  });
});
