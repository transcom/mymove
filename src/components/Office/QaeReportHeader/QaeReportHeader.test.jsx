import React from 'react';
import { screen } from '@testing-library/react';

import QaeReportHeader from './QaeReportHeader';

import { qaeCSRRoutes } from 'constants/routes';
import { renderWithRouter } from 'testUtils';

const mockReportId = 'b7305135-1d6d-4a0d-a6d5-f408474f6ee2';
const mockRoutingConfig = {
  path: qaeCSRRoutes.BASE_EVALUATION_REPORT_PATH,
  params: { moveCode: 'LR4T8V', reportId: mockReportId },
};
const getMockProps = (reportProps) => ({
  report: {
    id: mockReportId,
    type: 'SHIPMENT',
    moveReferenceID: '7972-2874',
    ...reportProps,
  },
});

describe('QaeReportHeader', () => {
  it('renders correct content for a shipment report', () => {
    renderWithRouter(<QaeReportHeader {...getMockProps()} />, mockRoutingConfig);

    // h1
    expect(screen.getByRole('heading', { name: 'Shipment report', level: 1 })).toBeInTheDocument();

    // h6
    expect(screen.getByRole('heading', { name: 'REPORT ID #QA-B7305', level: 6 })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'MOVE CODE #LR4T8V', level: 6 })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'MTO REFERENCE ID #7972-2874', level: 6 })).toBeInTheDocument();
  });

  it('renders correct content for a counseling report', () => {
    renderWithRouter(<QaeReportHeader {...getMockProps({ type: 'COUNSELING' })} />, mockRoutingConfig);

    // h1
    expect(screen.getByRole('heading', { name: 'Counseling report', level: 1 })).toBeInTheDocument();

    // h6
    expect(screen.getByRole('heading', { name: 'REPORT ID #QA-B7305', level: 6 })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'MOVE CODE #LR4T8V', level: 6 })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'MTO REFERENCE ID #7972-2874', level: 6 })).toBeInTheDocument();
  });
});
