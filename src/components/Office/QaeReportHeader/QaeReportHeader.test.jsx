import React from 'react';
import { render, screen } from '@testing-library/react';

import QaeReportHeader from './QaeReportHeader';

jest.mock('react-router', () => ({
  ...jest.requireActual('react-router'),
  useParams: jest.fn().mockReturnValue({ moveCode: 'LR4T8V' }),
}));

const getMockProps = (reportProps) => ({
  report: {
    id: 'b7305135-1d6d-4a0d-a6d5-f408474f6ee2',
    type: 'SHIPMENT',
    moveReferenceID: '7972-2874',
    ...reportProps,
  },
});

describe('QaeReportHeader', () => {
  it('renders correct content for a shipment report', () => {
    render(<QaeReportHeader {...getMockProps()} />);

    // h1
    expect(screen.getByRole('heading', { name: 'Shipment report', level: 1 })).toBeInTheDocument();

    // h6
    expect(screen.getByRole('heading', { name: 'REPORT ID #QA-B7305', level: 6 })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'MOVE CODE #LR4T8V', level: 6 })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'MTO REFERENCE ID #7972-2874', level: 6 })).toBeInTheDocument();
  });

  it('renders correct content for a counseling report', () => {
    render(<QaeReportHeader {...getMockProps({ type: 'COUNSELING' })} />);

    // h1
    expect(screen.getByRole('heading', { name: 'Counseling report', level: 1 })).toBeInTheDocument();

    // h6
    expect(screen.getByRole('heading', { name: 'REPORT ID #QA-B7305', level: 6 })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'MOVE CODE #LR4T8V', level: 6 })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'MTO REFERENCE ID #7972-2874', level: 6 })).toBeInTheDocument();
  });
});
