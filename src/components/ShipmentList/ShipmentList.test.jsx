import React from 'react';
import { render, screen, queryByRole, getByRole } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ShipmentList from './ShipmentList';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatWeight } from 'utils/formatters';

beforeEach(() => {
  jest.clearAllMocks();
});

describe('ShipmentList component', () => {
  const shipments = [
    {
      id: 'ID-1',
      shipmentType: SHIPMENT_OPTIONS.PPM,
      ppmShipment: {
        id: 'ppm',
        hasRequestedAdvance: false,
      },
    },
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
  it('renders ShipmentList with shipments', async () => {
    render(<ShipmentList {...defaultProps} />);

    expect(screen.getAllByTestId('shipment-list-item-container').length).toBe(4);
    expect(screen.getAllByTestId('shipment-list-item-container')[0]).toHaveTextContent(/^ppm/i);
    expect(screen.getAllByTestId('shipment-list-item-container')[1]).toHaveTextContent(/^hhg/i);
    expect(screen.getAllByTestId('shipment-list-item-container')[2]).toHaveTextContent(/^nts/i);
    expect(screen.getAllByTestId('shipment-list-item-container')[3]).toHaveTextContent(/^nts-release/i);
  });

  it.each([
    ['Shows', false],
    ['Hides', true],
  ])('%s the edit link if moveSubmitted is %s', async (showHideEditLink, moveSubmitted) => {
    const props = { ...defaultProps, moveSubmitted };

    render(<ShipmentList {...props} />);

    let editBtn = queryByRole(screen.getAllByTestId('shipment-list-item-container')[0], 'button', { name: 'Edit' });

    const checkShipmentClick = async (shipmentID, shipmentNumber, shipmentType) => {
      if (showHideEditLink === 'Shows') {
        await userEvent.click(editBtn);
        expect(onShipmentClick).toHaveBeenCalledWith(shipmentID, shipmentNumber, shipmentType);
      } else {
        expect(editBtn).toBeNull();
      }
    };

    await checkShipmentClick('ID-1', 1, SHIPMENT_OPTIONS.PPM);

    editBtn = queryByRole(screen.getAllByTestId('shipment-list-item-container')[1], 'button', { name: 'Edit' });
    await checkShipmentClick('ID-2', 1, SHIPMENT_OPTIONS.HHG);

    editBtn = queryByRole(screen.getAllByTestId('shipment-list-item-container')[2], 'button', { name: 'Edit' });
    await checkShipmentClick('ID-3', 1, SHIPMENT_OPTIONS.NTS);

    editBtn = queryByRole(screen.getAllByTestId('shipment-list-item-container')[3], 'button', { name: 'Edit' });
    await checkShipmentClick('ID-4', 1, SHIPMENT_OPTIONS.NTSR);
  });

  it.each([
    [SHIPMENT_OPTIONS.PPM, 1],
    [SHIPMENT_OPTIONS.HHG, 2],
    [SHIPMENT_OPTIONS.NTS, 3],
    [SHIPMENT_OPTIONS.NTSR, 4],
  ])('calls onDeleteClick for shipment type %s when delete is clicked', async (_, id) => {
    render(<ShipmentList {...defaultProps} />);
    const deleteBtn = getByRole(screen.getAllByTestId('shipment-list-item-container')[id - 1], 'button', {
      name: 'Delete',
    });

    await userEvent.click(deleteBtn);
    expect(onDeleteClick).toHaveBeenCalledWith(`ID-${id}`);
    expect(onDeleteClick).toHaveBeenCalledTimes(1);
  });
});

describe('Shipment List being used for billable weight', () => {
  it('renders maximum billable weight, actual billable weight, actual weight and weight allowance with no flags', () => {
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

  describe('Shipment List correctly displays incomplete PPM card', () => {
    it('display incomplete badge on PPM shipment missing an hasRequestedAdvance value', () => {
      const shipments = [
        {
          id: '0001',
          shipmentType: SHIPMENT_OPTIONS.PPM,
          ppmShipment: { id: '1234', hasRequestedAdvance: null },
        },
      ];

      const defaultProps = {
        shipments,
        moveSubmitted: false,
      };

      render(<ShipmentList {...defaultProps} />);

      expect(screen.getByText('Incomplete')).toBeInTheDocument();
    });

    it('do not show incomplete badge when PPM is complete', () => {
      const shipments = [
        {
          id: '0001',
          shipmentType: SHIPMENT_OPTIONS.PPM,
          ppmShipment: { id: '1234', hasRequestedAdvance: false },
        },
      ];

      const defaultProps = {
        shipments,
        moveSubmitted: false,
      };

      render(<ShipmentList {...defaultProps} />);

      expect(screen.queryByText('Incomplete')).not.toBeInTheDocument();
    });
  });
});
describe('Shipment List with PPM', () => {
  it('displays both estimated weight on PPM shipment', () => {
    const shipments = [
      {
        id: '0001',
        shipmentType: SHIPMENT_OPTIONS.PPM,
        ppmShipment: {
          id: '1234',
          hasRequestedAdvance: null,
          primeEstimatedWeight: '1000',
          calculatedBillableWeight: '1500',
        },
      },
    ];

    const defaultProps = {
      shipments,
      moveSubmitted: true,
      showShipmentWeight: true,
    };

    render(<ShipmentList {...defaultProps} />);

    expect(screen.getByText('Estimated')).toBeInTheDocument();
    expect(screen.getByText('Actual')).toBeInTheDocument();
  });
  it('should contain actual weight as full weight minus empty weight', () => {
    const shipments = [
      {
        id: '0001',
        shipmentType: SHIPMENT_OPTIONS.PPM,
        ppmShipment: {
          id: '1234',
          hasRequestedAdvance: null,
          primeEstimatedWeight: '1000',
          weightTickets: [
            {
              id: '1',
              fullWeight: '25000',
              emptyWeight: '22500',
            },
          ],
        },
      },
    ];
    const defaultProps = {
      shipments,
      moveSubmitted: true,
      showShipmentWeight: true,
    };
    render(<ShipmentList {...defaultProps} />);

    expect(screen.getByText('2,500 lbs')).toBeInTheDocument();
  });
});
