import React from 'react';
import { render, screen, queryByRole, getByRole } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ShipmentList from './ShipmentList';

import { SHIPMENT_OPTIONS, SHIPMENT_TYPES } from 'shared/constants';
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
    { id: 'ID-5', shipmentType: SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE },
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

    expect(screen.getAllByTestId('shipment-list-item-container').length).toBe(5);
    expect(screen.getAllByTestId('shipment-list-item-container')[0]).toHaveTextContent(/^ppm/i);
    expect(screen.getAllByTestId('shipment-list-item-container')[1]).toHaveTextContent(/^hhg/i);
    expect(screen.getAllByTestId('shipment-list-item-container')[2]).toHaveTextContent(/^nts/i);
    expect(screen.getAllByTestId('shipment-list-item-container')[3]).toHaveTextContent(/^nts-release/i);
    expect(screen.getAllByTestId('shipment-list-item-container')[4]).toHaveTextContent(/^UB/i);
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

    editBtn = queryByRole(screen.getAllByTestId('shipment-list-item-container')[4], 'button', { name: 'Edit' });
    await checkShipmentClick('ID-5', 1, SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE);
  });

  it.each([
    [SHIPMENT_OPTIONS.PPM, 1],
    [SHIPMENT_OPTIONS.HHG, 2],
    [SHIPMENT_OPTIONS.NTS, 3],
    [SHIPMENT_OPTIONS.NTSR, 4],
    [SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE, 5],
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

describe('ShipmentList shipment weight tooltip', () => {
  const defaultProps = {
    moveSubmitted: false,
  };

  it.each([
    [SHIPMENT_OPTIONS.HHG, 'ID-2', '110% Prime Estimated Weight'],
    [SHIPMENT_OPTIONS.NTS, 'ID-3', '110% Prime Estimated Weight'],
    [SHIPMENT_OPTIONS.NTSR, 'ID-4', '110% Previously Recorded Weight'],
    [SHIPMENT_TYPES.BOAT_HAUL_AWAY, 'ID-5', '110% Prime Estimated Weight'],
    [SHIPMENT_TYPES.BOAT_TOW_AWAY, 'ID-6', '110% Prime Estimated Weight'],
    [SHIPMENT_OPTIONS.MOBILE_HOME, 'ID-7', '110% Prime Estimated Weight'],
  ])('shipment weight tooltip, show is true. %s', async (type, id, expectedTooltipText) => {
    // Render component
    const props = { ...defaultProps, showShipmentTooltip: true, shipments: [{ id, shipmentType: type }] };
    render(<ShipmentList {...props} />);

    // Verify tooltip exists
    const tooltipIcon = screen.getByTestId('tooltip-container');
    expect(tooltipIcon).toBeInTheDocument();

    // Click the tooltip
    await userEvent.hover(tooltipIcon);

    // Verify tooltip text
    const tooltipText = await screen.findByText(expectedTooltipText);
    expect(tooltipText).toBeInTheDocument();
  });

  it.each([
    [SHIPMENT_OPTIONS.HHG, 'ID-2'],
    [SHIPMENT_OPTIONS.NTS, 'ID-3'],
    [SHIPMENT_OPTIONS.NTSR, 'ID-4'],
    [SHIPMENT_TYPES.BOAT_HAUL_AWAY, 'ID-5'],
    [SHIPMENT_TYPES.BOAT_TOW_AWAY, 'ID-6'],
    [SHIPMENT_OPTIONS.MOBILE_HOME, 'ID-7'],
  ])('shipment weight tooltip, show is false. %s', async (type, id) => {
    // Render component
    const props = { ...defaultProps, showShipmentTooltip: false, shipments: [{ id, shipmentType: type }] };
    render(<ShipmentList {...props} />);

    // Verify tooltip doesn't exist
    expect(screen.queryByTestId('tooltip-container')).not.toBeInTheDocument();
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
