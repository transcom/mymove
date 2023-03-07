import React from 'react';
import classnames from 'classnames';

import { formatReviewShipmentWeightsDate, formatWeight } from '../../../../utils/formatters';
import { shipmentTypes } from '../../../../constants/shipments';
import { createHeader } from '../../../Table/utils';
import { SHIPMENT_OPTIONS } from '../../../../shared/constants';
import { calculateTotalNetWeightForWeightTickets } from '../../../../utils/shipmentWeights';

import styles from './ReviewShipmentWeightsTable.module.scss';

const DASH = '—';

function shipmentTypeCellDisplayHelper(row) {
  const shipmentClassName = classnames({
    [styles[`review-shipment-weights-table-row-NTS-release`]]: row.shipmentType === SHIPMENT_OPTIONS.NTSR,
    [styles[`review-shipment-weights-table-row-NTS`]]: row.shipmentType === SHIPMENT_OPTIONS.NTS,
    [styles[`review-shipment-weights-table-row-HHG`]]:
      row.shipmentType === SHIPMENT_OPTIONS.HHG ||
      row.shipmentType === SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC ||
      row.shipmentType === SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
    [styles[`review-shipment-weights-table-row-PPM`]]: row.shipmentType === SHIPMENT_OPTIONS.PPM,
  });
  return (
    <div className={`${styles['review-shipment-weights-table-row']} ${shipmentClassName}`}>
      <strong>
        {shipmentTypes[row.shipmentType]}
        {row.showNumber && ` ${row.shipmentNumber}`}
      </strong>{' '}
    </div>
  );
}

function estimatedWeightCellDisplayHelper(row) {
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

function actualWeightCellDisplayHelper(row) {
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

export const NoRowsMessages = {
  PPM: 'No PPM shipments have been created for this move.',
  NonPPM: 'No HHG, NTS, or NTS-Release shipments have been created for this move.',
};

export const PPMReviewWeightsTableColumns = [
  createHeader('', (row) => shipmentTypeCellDisplayHelper(row), {
    id: 'shipmentType',
    isFilterable: false,
  }),
  createHeader('Weight ticket', (row) => <a href={row.ppmShipment.reviewURL}> Review Documents </a>, {
    id: 'weightTicket',
    isFilterable: false,
  }),
  createHeader(
    'Pro-gear (lbs)',
    (row) => (row.ppmShipment.proGearWeight > 0 ? formatWeight(row.ppmShipment.proGearWeight) : DASH),
    {
      id: 'proGear',
      isFilterable: false,
    },
  ),
  createHeader(
    'Spouse pro-gear',
    (row) => (row.ppmShipment.spouseProGearWeight > 0 ? formatWeight(row.ppmShipment.spouseProGearWeight) : DASH),
    {
      id: 'spouseProGear',
      isFilterable: false,
    },
  ),
  createHeader(
    'Estimated Weight',
    (row) => (row.ppmShipment.estimatedWeight > 0 ? formatWeight(row.ppmShipment.estimatedWeight) : DASH),
    {
      id: 'estimatedWeight',
      isFilterable: false,
    },
  ),
  createHeader(
    'Net Weight',
    (row) => {
      const calculatedNetWeight = calculateTotalNetWeightForWeightTickets(row.ppmShipment?.weightTickets);
      return calculatedNetWeight > 0 ? formatWeight(calculatedNetWeight) : DASH;
    },
    {
      id: 'netWeight',
      isFilterable: false,
    },
  ),
  createHeader(
    'Departure Date',
    (row) =>
      row.ppmShipment.expectedDepartureDate
        ? formatReviewShipmentWeightsDate(row.ppmShipment.expectedDepartureDate)
        : DASH,
    {
      id: 'departureDate',
      isFilterable: false,
    },
  ),
];

export const ProGearTableColumns = [
  createHeader(
    '',
    () => {
      return (
        <div
          className={`${styles['review-shipment-weights-table-row']} ${styles['review-shipment-weights-table-row-pro-gear']}`}
        >
          <strong>Pro-gear</strong>{' '}
        </div>
      );
    },
    {
      id: 'shipmentType',
      isFilterable: false,
    },
  ),
  createHeader(
    'Pro-gear (lbs)',
    (row) => (row.entitlement.proGearWeight > 0 ? formatWeight(row.entitlement.proGearWeight) : DASH),
    {
      id: 'proGear',
      isFilterable: false,
    },
  ),
  createHeader(
    'Spouse pro-gear (lbs)',
    (row) => (row.entitlement.spouseProGearWeight > 0 ? formatWeight(row.entitlement.spouseProGearWeight) : DASH),
    {
      id: 'spouseProGear',
      isFilterable: false,
    },
  ),
];

export const NonPPMTableColumns = [
  createHeader('', (row) => shipmentTypeCellDisplayHelper(row), {
    id: 'shipmentType',
    isFilterable: false,
  }),
  createHeader('Estimated Weight', (row) => estimatedWeightCellDisplayHelper(row), {
    id: 'estimatedWeight',
    isFilterable: false,
  }),
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
  createHeader('Actual weight', (row) => actualWeightCellDisplayHelper(row), {
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
