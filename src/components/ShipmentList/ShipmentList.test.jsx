import React from 'react';
import { render, screen, queryByRole, getByRole } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ShipmentList from './ShipmentList';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatWeight } from 'utils/formatters';

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

    expect(screen.getAllByTestId('shipment-list-item-container').length).toBe(4);
    expect(screen.getAllByTestId('shipment-list-item-container')[0]).toHaveTextContent(/^ppm #id-1/i);
    expect(screen.getAllByTestId('shipment-list-item-container')[1]).toHaveTextContent(/^hhg #id-2/i);
    expect(screen.getAllByTestId('shipment-list-item-container')[2]).toHaveTextContent(/^nts #id-3/i);
    expect(screen.getAllByTestId('shipment-list-item-container')[3]).toHaveTextContent(/^nts-release #id-4/i);
  });

  it.each([
    ['Shows', false],
    ['Hides', true],
  ])('%s the edit link if moveSubmitted is %s', (showHideEditLink, moveSubmitted) => {
    const props = { ...defaultProps, moveSubmitted };

    render(<ShipmentList {...props} />);

    let editBtn = queryByRole(screen.getAllByTestId('shipment-list-item-container')[0], 'button', { name: 'Edit' });

    const checkShipmentClick = (shipmentID, shipmentNumber, shipmentType) => {
      if (showHideEditLink === 'Shows') {
        userEvent.click(editBtn);
        expect(onShipmentClick).toHaveBeenCalledWith(shipmentID, shipmentNumber, shipmentType);
      } else {
        expect(editBtn).toBeNull();
      }
    };

    checkShipmentClick('ID-1', 1, SHIPMENT_OPTIONS.PPM);

    editBtn = queryByRole(screen.getAllByTestId('shipment-list-item-container')[1], 'button', { name: 'Edit' });
    checkShipmentClick('ID-2', 1, SHIPMENT_OPTIONS.HHG);

    editBtn = queryByRole(screen.getAllByTestId('shipment-list-item-container')[2], 'button', { name: 'Edit' });
    checkShipmentClick('ID-3', 1, SHIPMENT_OPTIONS.NTS);

    editBtn = queryByRole(screen.getAllByTestId('shipment-list-item-container')[3], 'button', { name: 'Edit' });
    checkShipmentClick('ID-4', 1, SHIPMENT_OPTIONS.NTSR);
  });

  it('calls onDeleteClick when delete is clicked', () => {
    render(<ShipmentList {...defaultProps} />);
    for (let i = 0; i < defaultProps.shipments.length; i += 1) {
      const deleteBtn = getByRole(screen.getAllByTestId('shipment-list-item-container')[i], 'button', {
        name: 'Delete',
      });
      userEvent.click(deleteBtn);
      expect(onDeleteClick).toHaveBeenCalledWith(`ID-${i + 1}`);
      expect(onDeleteClick).toHaveBeenCalledTimes(i + 1);
    }
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
