/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ShipmentList from './ShipmentList';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatWeight } from 'utils/formatters';

beforeEach(() => {
  jest.clearAllMocks();
});

describe('ShipmentList component', () => {
  const shipments = [
    { id: 'ID-1', shipmentType: SHIPMENT_OPTIONS.PPM },
    { id: 'ID-2', shipmentType: SHIPMENT_OPTIONS.HHG },
    { id: 'ID-3', shipmentType: SHIPMENT_OPTIONS.NTS },
    { id: 'ID-4', shipmentType: SHIPMENT_OPTIONS.NTSR },
  ];
  const onShipmentClick = jest.fn();
  const onDeleteClick = jest.fn();
  const defaultProps = {
    shipments,
    onShipmentClick,
    onDeleteClick,
    moveSubmitted: false,
  };
  it('renders ShipmentList with shipments', () => {
    render(<ShipmentList {...defaultProps} />);

    screen.getByRole('button', { name: /^ppm #id-1/i });
    screen.getByRole('button', { name: /^hhg #id-2/i });
    screen.getByRole('button', { name: /^nts #id-3/i });
    screen.getByRole('button', { name: /^nts-release #id-4/i });
  });

  it.each([
    ['Shows', false],
    ['Hides', true],
  ])('%s the edit link if moveSubmitted is %s', async (showHideEditLink, moveSubmitted) => {
    const props = { ...defaultProps, moveSubmitted };

    render(<ShipmentList {...props} />);

    userEvent.click(screen.getByRole('button', { name: /^ppm /i }));

    const checkShipmentClick = (shipmentID, shipmentNumber, shipmentType) => {
      if (showHideEditLink === 'Shows') {
        expect(onShipmentClick).toHaveBeenCalledWith(shipmentID, shipmentNumber, shipmentType);
      } else {
        expect(onShipmentClick).not.toHaveBeenCalled();
      }
    };

    await waitFor(() => {
      checkShipmentClick('ID-1', 1, SHIPMENT_OPTIONS.PPM);
    });

    userEvent.click(screen.getByRole('button', { name: /^hhg /i }));

    await waitFor(() => {
      checkShipmentClick('ID-2', 1, SHIPMENT_OPTIONS.HHG);
    });

    userEvent.click(screen.getByRole('button', { name: /^nts /i }));

    await waitFor(() => {
      checkShipmentClick('ID-3', 1, SHIPMENT_OPTIONS.NTS);
    });

    userEvent.click(screen.getByRole('button', { name: /^nts-release /i }));

    await waitFor(() => {
      checkShipmentClick('ID-4', 1, SHIPMENT_OPTIONS.NTSR);
    });
  });
});

describe('Shipment List being used for billable weight', () => {
  it('renders maximum billable weight, total billable weight, weight requested and weight allowance with no flags', () => {
    const shipments = [
      {
        id: '0001',
        shipmentType: SHIPMENT_OPTIONS.HHG,
        calculatedBillableWeight: 1161,
        primeEstimatedWeight: 200,
        reweigh: { id: '1234', weight: 50 },
      },
      {
        id: '0002',
        shipmentType: SHIPMENT_OPTIONS.HHG,
        calculatedBillableWeight: 3200,
        primeEstimatedWeight: 3000,
        reweigh: { id: '1234' },
      },
      {
        id: '0003',
        shipmentType: SHIPMENT_OPTIONS.HHG,
        calculatedBillableWeight: 3000,
        primeEstimatedWeight: 3000,
        reweigh: { id: '1234', weight: 40 },
      },
    ];

    const defaultProps = {
      shipments,
      moveSubmitted: false,
      showShipmentWeight: true,
    };

    render(<ShipmentList {...defaultProps} />);

    // flags
    expect(screen.queryByText('Over weight')).toBeInTheDocument();
    expect(screen.queryByText('Missing weight')).toBeInTheDocument();

    // weights
    expect(screen.getByText(formatWeight(shipments[0].calculatedBillableWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(shipments[1].calculatedBillableWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(shipments[2].calculatedBillableWeight))).toBeInTheDocument();
  });

  it('does not display weight flags when not appropriate', () => {
    const shipments = [
      { id: '0001', shipmentType: SHIPMENT_OPTIONS.HHG, calculatedBillableWeight: 5666, primeEstimatedWeight: 5600 },
      {
        id: '0002',
        shipmentType: SHIPMENT_OPTIONS.HHG,
        calculatedBillableWeight: 3200,
        primeEstimatedWeight: 3000,
        reweigh: { id: '1234', weight: 3400 },
      },
      { id: '0003', shipmentType: SHIPMENT_OPTIONS.HHG, calculatedBillableWeight: 5400, primeEstimatedWeight: 5000 },
      // we don't show flags for ntsr shipments - if this was an hhg, it would show a missing weight warning
      { id: '0004', shipmentType: SHIPMENT_OPTIONS.NTSR },
    ];

    const defaultProps = {
      shipments,
      moveSubmitted: false,
      showShipmentWeight: true,
    };

    render(<ShipmentList {...defaultProps} />);

    // flags
    expect(screen.queryByText('Over weight')).not.toBeInTheDocument();
    expect(screen.queryByText('Missing weight')).not.toBeInTheDocument();
  });
});
