import React from 'react';
import { render, screen } from '@testing-library/react';

import EvaluationReportShipmentDisplay from './EvaluationReportShipmentDisplay';

import { hhgInfo, ntsInfo, ntsReleaseInfo, ordersLOA } from 'components/Office/ShipmentDisplay/ShipmentDisplayTestData';

describe('Evaluation report - HHG Shipment', () => {
  it('renders the HHG component successfully', () => {
    render(
      <EvaluationReportShipmentDisplay
        shipmentId="1"
        displayInfo={hhgInfo}
        ordersLOA={ordersLOA}
        onChange={jest.fn()}
        isSubmitted
      />,
    );
    expect(screen.getByText('HHG')).toBeInTheDocument();
  });
});

describe('Evaluation report - NTS Shipment', () => {
  it('renders the NTS component successfully', () => {
    render(
      <EvaluationReportShipmentDisplay
        shipmentId="1"
        displayInfo={ntsInfo}
        ordersLOA={ordersLOA}
        onChange={jest.fn()}
        isSubmitted
      />,
    );
    expect(screen.getByText('NTS')).toBeInTheDocument();
  });
});

describe('Evaluation report - NTSR Shipment', () => {
  it('renders the NTSR component successfully', () => {
    render(
      <EvaluationReportShipmentDisplay
        shipmentId="1"
        displayInfo={ntsReleaseInfo}
        ordersLOA={ordersLOA}
        onChange={jest.fn()}
        isSubmitted
      />,
    );
    expect(screen.getByText('NTS-release')).toBeInTheDocument();
  });
});
