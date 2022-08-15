/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { v4 } from 'uuid';
import ReactRouter from 'react-router';

import ShipmentEvaluationReport from './ShipmentEvaluationReport';

import { MockProviders } from 'testUtils';

const mockMoveCode = 'LR4T8V';
const mockReportId = v4();
const mockCustomer = {
  last_name: 'spaceman',
  first_name: 'leo',
  phone: '555-555-5555',
};
const mockOrders = {
  grade: 'E_1',
  agency: 'COAST_GUARD',
};

describe('ShipmentEvaluationReport', () => {
  it('renders the page components', async () => {
    jest.spyOn(ReactRouter, 'useParams').mockReturnValue({ moveCode: mockMoveCode, reportId: mockReportId });

    render(
      <MockProviders initialEntries={[`/moves/LR4T8V/evaluation-reports/${mockReportId}`]}>
        <ShipmentEvaluationReport customerInfo={mockCustomer} orders={mockOrders} />
      </MockProviders>,
    );

    await waitFor(() => {
      const h1 = screen.getByRole('heading', { name: 'Shipment report', level: 1 });
      expect(h1).toBeInTheDocument();

      expect(screen.getByText(`REPORT ID #QA-${mockReportId.slice(0, 5).toUpperCase()}`)).toBeInTheDocument();
      expect(screen.getByText(`MOVE CODE #${mockMoveCode}`)).toBeInTheDocument();
      expect(screen.getByText('MTO REFERENCE ID #')).toBeInTheDocument();

      // Page content sections
      expect(screen.getByRole('heading', { name: 'Shipment information', level: 2 })).toBeInTheDocument();
      expect(screen.getByRole('heading', { name: 'Evaluation form', level: 2 })).toBeInTheDocument();

      // Buttons
      expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save draft' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Review and submit' })).toBeInTheDocument();
    });
  });
});
