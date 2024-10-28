import React from 'react';
import { render, screen } from '@testing-library/react';

import EvaluationReportShipmentDisplay from './EvaluationReportShipmentDisplay';

import {
  hhgInfo,
  ntsInfo,
  ntsReleaseInfo,
  canceledInfo,
  ordersLOA,
  ppmInfo,
} from 'components/Office/ShipmentDisplay/ShipmentDisplayTestData';

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
    expect(screen.getByTestId('shipment-display')).toHaveTextContent('HHG');
    expect(screen.getByTestId('ShipmentContainer')).toHaveTextContent('EVLRPT-01');
  });
  it('renders with canceled tag', () => {
    render(
      <EvaluationReportShipmentDisplay shipmentId="1" displayInfo={canceledInfo} onChange={jest.fn()} isSubmitted />,
    );
    expect(screen.getByText('canceled')).toBeInTheDocument();
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
    expect(screen.getByTestId('shipment-display')).toHaveTextContent('NTS');
    expect(screen.getByTestId('ShipmentContainer')).toHaveTextContent('EVLRPT-02');
  });
  it('renders with canceled tag', () => {
    render(
      <EvaluationReportShipmentDisplay shipmentId="1" displayInfo={canceledInfo} onChange={jest.fn()} isSubmitted />,
    );
    expect(screen.getByText('canceled')).toBeInTheDocument();
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
    expect(screen.getByTestId('shipment-display')).toHaveTextContent('NTS-release');
  });
  it('renders with canceled tag', () => {
    render(
      <EvaluationReportShipmentDisplay shipmentId="1" displayInfo={canceledInfo} onChange={jest.fn()} isSubmitted />,
    );
    expect(screen.getByText('canceled')).toBeInTheDocument();
  });
});

describe('Evaluation report - PPM Shipment', () => {
  it('renders the PPM component successfully', () => {
    render(
      <EvaluationReportShipmentDisplay
        shipmentId="3"
        displayInfo={ppmInfo}
        ordersLOA={ordersLOA}
        onChange={jest.fn()}
        isSubmitted
      />,
    );
    expect(screen.getByTestId('shipment-display')).toHaveTextContent('PPM');
    expect(screen.getByTestId('ShipmentContainer')).toHaveTextContent('EVLRPT-03');
  });
  it('renders with canceled tag', () => {
    render(
      <EvaluationReportShipmentDisplay shipmentId="1" displayInfo={canceledInfo} onChange={jest.fn()} isSubmitted />,
    );
    expect(screen.getByText('canceled')).toBeInTheDocument();
  });
});
