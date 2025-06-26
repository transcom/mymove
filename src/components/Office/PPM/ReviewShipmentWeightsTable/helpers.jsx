import React from 'react';

import styles from './ReviewShipmentWeightsTable.module.scss';

import { formatReviewShipmentWeightsDate, formatWeight } from 'utils/formatters';
import { shipmentTypes } from 'constants/shipments';
import { createHeader } from 'components/Table/utils';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { getTotalNetWeightForWeightTickets } from 'utils/shipmentWeights';

export const DASH = '—';

export function addShipmentNumbersToTableData(tableData, determineShipmentNumbers) {
  if (!determineShipmentNumbers) {
    return tableData;
  }
  const shipments = tableData;
  const shipmentNumbersByType = {};
  const shipmentCountByType = {};

  // Count up each shipment type in the table data.
  shipments.forEach((shipment) => {
    const { shipmentType } = shipment;
    if (shipmentCountByType[shipmentType]) {
      shipmentCountByType[shipmentType] += 1;
    } else {
      shipmentCountByType[shipmentType] = 1;
    }
  });

  // Add the shipmentNumber and showNumber vars
  // to each shipment in the table data
  const modifiedTableData = shipments.map((shipment) => {
    const { shipmentType } = shipment;

    // Determine the number for each shipment.
    if (shipmentNumbersByType[shipmentType]) {
      shipmentNumbersByType[shipmentType] += 1;
    } else {
      shipmentNumbersByType[shipmentType] = 1;
    }
    const shipmentNumber = shipmentNumbersByType[shipmentType];

    // We only show the number for the shipment
    // if there are more than one of that type in the table data.
    const showNumber = shipmentCountByType[shipmentType] > 1;

    return { ...shipment, showNumber, shipmentNumber };
  });
  return modifiedTableData;
}

export function determineTableRowClassname(shipmentType) {
  let shipmentClassname = '';
  switch (shipmentType) {
    case SHIPMENT_OPTIONS.NTSR:
      shipmentClassname = styles[`review-shipment-weights-table-row-NTS-release`];
      break;
    case SHIPMENT_OPTIONS.NTS:
      shipmentClassname = styles[`review-shipment-weights-table-row-NTS`];
      break;
    case SHIPMENT_OPTIONS.HHG:
      shipmentClassname = styles[`review-shipment-weights-table-row-HHG`];
      break;
    case SHIPMENT_OPTIONS.PPM:
      shipmentClassname = styles[`review-shipment-weights-table-row-PPM`];
      break;
    default:
      break;
  }
  return shipmentClassname;
}

export const ShipmentTypeCell = (props) => {
  const { row } = props;
  const shipmentClassName = determineTableRowClassname(row.shipmentType);
  return (
    <div className={`${styles['review-shipment-weights-table-row']} ${shipmentClassName}`}>
      <strong>
        {shipmentTypes[row.shipmentType]}
        {row.showNumber && ` ${row.shipmentNumber}`}
      </strong>{' '}
    </div>
  );
};

export function estimatedWeightDisplayHelper(row) {
  let estimatedWeight;
  switch (row.shipmentType) {
    // Estimated weight doesn't apply to NTSR shipments
    case SHIPMENT_OPTIONS.NTSR:
      return 'N/A';
    case SHIPMENT_OPTIONS.NTS:
      estimatedWeight = row.ntsRecordedWeight;
      break;
    default:
      estimatedWeight = row.primeEstimatedWeight;
      break;
  }
  return estimatedWeight ? formatWeight(estimatedWeight) : DASH;
}

export function actualWeightDisplayHelper(row) {
  if (!row?.reweigh?.weight && !row?.primeActualWeight) {
    return DASH;
  }
  let actualWeight;
  if (!row?.reweigh?.weight) {
    actualWeight = row.primeActualWeight;
  } else if (!row?.primeActualWeight) {
    actualWeight = row.reweigh.weight;
  } else {
    actualWeight = Math.min(row.primeActualWeight, row.reweigh.weight);
  }
  return actualWeight > 0 ? formatWeight(actualWeight) : DASH;
}

export const NO_ROWS_MESSAGES = {
  PPM: 'No PPM shipments have been created for this move.',
  NonPPM: 'No HHG, NTS, or NTS-Release shipments have been created for this move.',
};

export const PPMReviewWeightsTableColumns = [
  createHeader('', (row) => <ShipmentTypeCell row={row} />, {
    id: 'shipmentType',
    isFilterable: false,
  }),
  createHeader(
    'Weight ticket',
    (row) =>
      row.ppmShipment.weightTickets.length > 0 ? (
        <a href={row.ppmShipment.reviewShipmentWeightsURL}> Review Documents </a>
      ) : (
        DASH
      ),
    {
      id: 'weightTicket',
      isFilterable: false,
    },
  ),
  createHeader(
    'Pro-gear (lbs)',
    (row) => (row.actualProGearWeight > 0 ? formatWeight(row.actualProGearWeight) : DASH),
    {
      id: 'proGear',
      isFilterable: false,
    },
  ),
  createHeader(
    'Spouse pro-gear',
    (row) => (row.actualSpouseProGearWeight > 0 ? formatWeight(row.actualSpouseProGearWeight) : DASH),
    {
      id: 'spouseProGear',
      isFilterable: false,
    },
  ),
  createHeader('Gun safe', (row) => (row.actualGunSafeWeight > 0 ? formatWeight(row.actualGunSafeWeight) : DASH), {
    id: 'gunSafe',
    isFilterable: false,
  }),
  createHeader(
    'Estimated weight',
    (row) => (row.ppmShipment.estimatedWeight > 0 ? formatWeight(row.ppmShipment.estimatedWeight) : DASH),
    {
      id: 'estimatedWeight',
      isFilterable: false,
    },
  ),
  createHeader(
    'Net weight',
    (row) => {
      const calculatedNetWeight = getTotalNetWeightForWeightTickets(row.ppmShipment?.weightTickets);
      return calculatedNetWeight > 0 ? formatWeight(calculatedNetWeight) : DASH;
    },
    {
      id: 'netWeight',
      isFilterable: false,
    },
  ),
  createHeader(
    'Actual Departure date',
    (row) => (row.ppmShipment.actualMoveDate ? formatReviewShipmentWeightsDate(row.ppmShipment.actualMoveDate) : DASH),
    {
      id: 'departureDate',
      isFilterable: false,
    },
  ),
];

export const PPMReviewWeightsTableColumnsWithoutGunSafe = [
  createHeader('', (row) => <ShipmentTypeCell row={row} />, {
    id: 'shipmentType',
    isFilterable: false,
  }),
  createHeader(
    'Weight ticket',
    (row) =>
      row.ppmShipment.weightTickets.length > 0 ? (
        <a href={row.ppmShipment.reviewShipmentWeightsURL}> Review Documents </a>
      ) : (
        DASH
      ),
    {
      id: 'weightTicket',
      isFilterable: false,
    },
  ),
  createHeader(
    'Pro-gear (lbs)',
    (row) => (row.actualProGearWeight > 0 ? formatWeight(row.actualProGearWeight) : DASH),
    {
      id: 'proGear',
      isFilterable: false,
    },
  ),
  createHeader(
    'Spouse pro-gear',
    (row) => (row.actualSpouseProGearWeight > 0 ? formatWeight(row.actualSpouseProGearWeight) : DASH),
    {
      id: 'spouseProGear',
      isFilterable: false,
    },
  ),
  createHeader(
    'Estimated weight',
    (row) => (row.ppmShipment.estimatedWeight > 0 ? formatWeight(row.ppmShipment.estimatedWeight) : DASH),
    {
      id: 'estimatedWeight',
      isFilterable: false,
    },
  ),
  createHeader(
    'Net weight',
    (row) => {
      const calculatedNetWeight = getTotalNetWeightForWeightTickets(row.ppmShipment?.weightTickets);
      return calculatedNetWeight > 0 ? formatWeight(calculatedNetWeight) : DASH;
    },
    {
      id: 'netWeight',
      isFilterable: false,
    },
  ),
  createHeader(
    'Actual Departure date',
    (row) => (row.ppmShipment.actualMoveDate ? formatReviewShipmentWeightsDate(row.ppmShipment.actualMoveDate) : DASH),
    {
      id: 'departureDate',
      isFilterable: false,
    },
  ),
];

export const nonPPMTableColumns = [
  createHeader('', (row) => <ShipmentTypeCell row={row} />, {
    id: 'shipmentType',
    isFilterable: false,
  }),
  createHeader('Estimated weight', (row) => estimatedWeightDisplayHelper(row), {
    id: 'estimatedWeight',
    isFilterable: false,
  }),
  createHeader(
    'Pro-gear \n(lbs)',
    (row) => (row.actualProGearWeight > 0 ? formatWeight(row.actualProGearWeight) : DASH),
    {
      id: 'actualProGear',
      isFilterable: false,
    },
  ),
  createHeader(
    'Spouse pro-gear \n(lbs)',
    (row) => (row.actualSpouseProGearWeight > 0 ? formatWeight(row.actualSpouseProGearWeight) : DASH),
    {
      id: 'actualSpouseProGear',
      isFilterable: false,
    },
  ),
  createHeader('Reweigh requested', (row) => (row.reweigh ? 'Yes' : 'No'), {
    id: 'reweighRequested',
    isFilterable: false,
  }),
  createHeader(
    'Billable weight',
    (row) => (row.calculatedBillableWeight > 0 ? formatWeight(row.calculatedBillableWeight) : DASH),
    {
      id: 'billableWeight',
      isFilterable: false,
    },
  ),
  createHeader('Actual weight', (row) => actualWeightDisplayHelper(row), {
    id: 'actualWeight',
    isFilterable: false,
  }),
  createHeader(
    'Delivery date',
    (row) => (row?.actualDeliveryDate ? formatReviewShipmentWeightsDate(row.actualDeliveryDate) : DASH),
    {
      id: 'deliveryDate',
      isFilterable: false,
    },
  ),
];

export const PPMReviewWeightsTableConfig = {
  tableColumns: PPMReviewWeightsTableColumns,
  noRowsMsg: NO_ROWS_MESSAGES.PPM,
  determineShipmentNumbers: true,
};

export const PPMReviewWeightsTableConfigWithoutGunSafe = {
  tableColumns: PPMReviewWeightsTableColumnsWithoutGunSafe,
  noRowsMsg: NO_ROWS_MESSAGES.PPM,
  determineShipmentNumbers: true,
};

export const nonPPMReviewWeightsTableConfig = {
  tableColumns: nonPPMTableColumns,
  noRowsMsg: NO_ROWS_MESSAGES.NonPPM,
  determineShipmentNumbers: true,
};
