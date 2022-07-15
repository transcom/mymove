import React from 'react';
import { render, screen } from '@testing-library/react';

import EvaluationReportTable from './EvaluationReportTable';

const submittedReport = {
  createdAt: '2022-07-14T19:21:27.573Z',
  evaluationLengthMinutes: 45,
  id: 'a7fdb0b3-f876-4686-b94f-ad20a2c9a63d',
  inspectionDate: '2022-07-14T00:00:00.000Z',
  inspectionType: 'VIRTUAL',
  location: 'DESTINATION',
  moveID: 'bd1bbbdc-1710-4831-aa41-e35ebedff5cd',
  remarks: 'this is a submitted shipment report',
  shipmentID: '38e87840-d385-413e-9746-b2da2c8245bb',
  submittedAt: '2022-07-14T19:21:27.565Z',
  type: 'SHIPMENT',
  violationsObserved: true,
};
const draftReport = {
  createdAt: '2022-07-14T19:21:27.579Z',
  evaluationLengthMinutes: 45,
  id: '1f9d18a8-7688-428d-be8e-3f3c59a0ae5e',
  inspectionDate: '2022-07-14T00:00:00.000Z',
  inspectionType: 'PHYSICAL',
  location: null,
  locationDescription: 'Route 66 at crash inspection site 3',
  moveID: 'bd1bbbdc-1710-4831-aa41-e35ebedff5cd',
  remarks: 'this is a submitted NTS shipment report',
  shipmentID: 'd46825dd-cf90-442b-96de-c5075ea2f1bf',
  submittedAt: null,
  travelTimeMinutes: 60,
  type: 'SHIPMENT',
  violationsObserved: true,
};

describe('EvaluationReportTable', () => {
  it('renders empty table', () => {
    render(<EvaluationReportTable reports={[]} />);
    expect(screen.getByText('Report ID')).toBeInTheDocument();
    expect(screen.getByText('Date submitted')).toBeInTheDocument();
    expect(screen.getByText('Location')).toBeInTheDocument();
    expect(screen.getByText('Violations')).toBeInTheDocument();
    expect(screen.getByText('Serious Incident')).toBeInTheDocument();
    expect(screen.getByText('No QAE reports have been submitted for this shipment')).toBeInTheDocument();
  });
  it('renders table with a draft report', () => {
    render(<EvaluationReportTable reports={[draftReport]} />);
    expect(screen.getByText('Report ID')).toBeInTheDocument();
    expect(screen.getByText('Date submitted')).toBeInTheDocument();
    expect(screen.getByText('Location')).toBeInTheDocument();
    expect(screen.getByText('Violations')).toBeInTheDocument();
    expect(screen.getByText('Serious Incident')).toBeInTheDocument();
    expect(screen.queryByText('No QAE reports have been submitted for this shipment')).not.toBeInTheDocument();

    expect(screen.getByTestId('tag')).toHaveTextContent('DRAFT');

    expect(screen.getByText('#QA-1F9D1')).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'View report' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Download' })).toBeInTheDocument();
  });
  it('renders table with a submitted report', () => {
    render(<EvaluationReportTable reports={[submittedReport]} />);
    expect(screen.getByText('Report ID')).toBeInTheDocument();
    expect(screen.getByText('Date submitted')).toBeInTheDocument();
    expect(screen.getByText('Location')).toBeInTheDocument();
    expect(screen.getByText('Violations')).toBeInTheDocument();
    expect(screen.getByText('Serious Incident')).toBeInTheDocument();

    expect(screen.queryByText('No QAE reports have been submitted for this shipment')).not.toBeInTheDocument();
    expect(screen.queryByTestId('tag')).not.toBeInTheDocument();

    expect(screen.getByText('#QA-A7FDB')).toBeInTheDocument();
    expect(screen.getByText('14 Jul 2022')).toBeInTheDocument();
    expect(screen.getByText('Destination')).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'View report' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Download' })).toBeInTheDocument();
  });
});
