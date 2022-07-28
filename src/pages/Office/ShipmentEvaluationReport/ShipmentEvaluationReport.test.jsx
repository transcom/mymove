/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen } from '@testing-library/react';
import { v4 } from 'uuid';
import ReactRouter from 'react-router';

import ShipmentEvaluationReport from './ShipmentEvaluationReport';

import { MockProviders } from 'testUtils';

const mockMoveCode = 'LR4T8V';
const mockReportId = v4();
const mockMtoRefId = 'TODO';

describe('ShipmentEvaluationReport', () => {
  it('renders the page components', async () => {
    jest.spyOn(ReactRouter, 'useParams').mockReturnValue({ moveCode: mockMoveCode, reportId: mockReportId });

    render(
      <MockProviders initialEntries={[`/moves/LR4T8V/evaluation-reports/${mockReportId}`]}>
        <ShipmentEvaluationReport />
      </MockProviders>,
    );

    const h1 = await screen.getByRole('heading', { name: 'Shipment report', level: 1 });
    expect(h1).toBeInTheDocument();

    expect(await screen.getByText(`REPORT ID #${mockReportId}`)).toBeInTheDocument();
    expect(await screen.getByText(`MOVE CODE ${mockMoveCode}`)).toBeInTheDocument();
    expect(await screen.getByText(`MTO REFERENCE ID ${mockMtoRefId}`)).toBeInTheDocument();

    // Page content sections
    expect(await screen.getByRole('heading', { name: 'Shipment information', level: 2 })).toBeInTheDocument();
    expect(await screen.getByRole('heading', { name: 'Evaluation form', level: 2 })).toBeInTheDocument();

    // Buttons
    expect(await screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
    expect(await screen.getByRole('button', { name: 'Save draft' })).toBeInTheDocument();
    expect(await screen.getByRole('button', { name: 'Submit' })).toBeInTheDocument();
  });

  it('displays the delete confirmation on cancel', async () => {
    jest.spyOn(ReactRouter, 'useParams').mockReturnValue({ moveCode: mockMoveCode, reportId: mockReportId });

    render(
      <MockProviders initialEntries={[`/moves/LR4T8V/evaluation-reports/${mockReportId}`]}>
        <ShipmentEvaluationReport />
      </MockProviders>,
    );

    const h1 = await screen.getByRole('heading', { name: 'Shipment report', level: 1 });
    expect(h1).toBeInTheDocument();

    // Buttons
    await screen.getByRole('button', { name: 'Cancel' }).click();

    expect(
      await screen.findByRole('heading', { level: 3, name: 'Are you sure you want to cancel this report?' }),
    ).toBeInTheDocument();
  });
});
