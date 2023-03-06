import React from 'react';
import { useTable, useFilters, usePagination, useSortBy } from 'react-table';
import classnames from 'classnames';

import { createHeader } from '../../../Table/utils';
import { formatWeight, formatReviewShipmentWeightsDate } from '../../../../utils/formatters';
import { calculateTotalNetWeightForWeightTickets } from '../../../../utils/ppmCloseout';
import { shipmentTypes } from '../../../../constants/shipments';
import { SHIPMENT_OPTIONS } from '../../../../shared/constants';

import styles from './ReviewShipmentWeightsTable.module.scss';

import Table from 'components/Table/Table';

export const NoRowsMessages = {
  PPM: 'No PPM shipments have been created for this move.',
  NonPPM: 'No HHG, NTS, or NTS-Release shipments have been created for this move.',
};

export const PPMReviewWeightsTableColumns = [
  createHeader(
    '',
    (row) => {
      return (
        <div
          className={`${styles['review-shipment-weights-table-row']} ${styles['review-shipment-weights-table-row-PPM']}`}
        >
          <strong>
            {shipmentTypes[row.shipmentType]}
            {row.showNumber && ` ${row.shipmentNumber}`}
          </strong>{' '}
        </div>
      );
    },
    {
      id: 'shipmentType',
      isFilterable: false,
    },
  ),
  createHeader('Weight ticket', (row) => <a href={row.ppmShipment.reviewURL}> Review Documents </a>, {
    id: 'weightTicket',
    isFilterable: false,
  }),
  createHeader(
    'Pro-gear (lbs)',
    (row) => (row.ppmShipment.proGearWeight > 0 ? formatWeight(row.ppmShipment.proGearWeight) : '-'),
    {
      id: 'proGear',
      isFilterable: false,
    },
  ),
  createHeader(
    'Spouse pro-gear',
    (row) => (row.ppmShipment.spouseProGearWeight > 0 ? formatWeight(row.ppmShipment.spouseProGearWeight) : '-'),
    {
      id: 'spouseProGear',
      isFilterable: false,
    },
  ),
  createHeader(
    'Estimated Weight',
    (row) => (row.ppmShipment.estimatedWeight > 0 ? formatWeight(row.ppmShipment.estimatedWeight) : '-'),
    {
      id: 'estimatedWeight',
      isFilterable: false,
    },
  ),
  createHeader(
    'Net Weight',
    (row) => {
      const calculatedNetWeight = calculateTotalNetWeightForWeightTickets(row.ppmShipment?.weightTickets);
      return calculatedNetWeight > 0 ? formatWeight(calculatedNetWeight) : '-';
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
        : '-',
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
    (row) => (row.entitlement.proGearWeight > 0 ? formatWeight(row.entitlement.proGearWeight) : '-'),
    {
      id: 'proGear',
      isFilterable: false,
    },
  ),
  createHeader(
    'Spouse pro-gear (lbs)',
    (row) => (row.entitlement.spouseProGearWeight > 0 ? formatWeight(row.entitlement.spouseProGearWeight) : '-'),
    {
      id: 'spouseProGear',
      isFilterable: false,
    },
  ),
];

export const NonPPMTableColumns = [
  createHeader(
    '',
    (row) => {
      const shipmentClassName = classnames({
        [styles[`review-shipment-weights-table-row-NTS-release`]]: row.shipmentType === SHIPMENT_OPTIONS.NTSR,
        [styles[`review-shipment-weights-table-row-NTS`]]: row.shipmentType === SHIPMENT_OPTIONS.NTS,
        [styles[`review-shipment-weights-table-row-HHG`]]:
          row.shipmentType === SHIPMENT_OPTIONS.HHG ||
          row.shipmentType === SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC ||
          row.shipmentType === SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
      });
      return (
        <div className={`${styles['review-shipment-weights-table-row']} ${shipmentClassName}`}>
          <strong>
            {shipmentTypes[row.shipmentType]}
            {row.showNumber && ` ${row.shipmentNumber}`}
          </strong>{' '}
        </div>
      );
    },
    {
      id: 'shipmentType',
      isFilterable: false,
    },
  ),
  createHeader(
    'Estimated Weight',
    (row) => {
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
      return estimatedWeight ? formatWeight(estimatedWeight) : '-';
    },
    {
      id: 'estimatedWeight',
      isFilterable: false,
    },
  ),
  createHeader('Reweigh requested', (row) => (row.reweigh ? 'Yes' : 'No'), {
    id: 'reweighRequested',
    isFilterable: false,
  }),
  createHeader(
    'Billable weight',
    (row) => (row.calculatedBillableWeight > 0 ? formatWeight(row.calculatedBillableWeight) : '-'),
    {
      id: 'billableWeight',
      isFilterable: false,
    },
  ),
  createHeader(
    'Actual weight',
    (row) => {
      if (!row?.reweigh?.weight && !row?.primeActualWeight) {
        return '-';
      }
      let actualWeight;
      if (!row?.reweigh?.weight) {
        actualWeight = row.primeActualWeight;
      } else if (!row?.primeActualWeight) {
        actualWeight = row.reweigh.weight;
      } else {
        actualWeight = Math.min(row.primeActualWeight, row.reweigh.weight);
      }
      return actualWeight > 0 ? formatWeight(actualWeight) : '-';
    },
    {
      id: 'actualWeight',
      isFilterable: false,
    },
  ),
  createHeader(
    'Delivery date',
    (row) => (row?.actualDeliveryDate ? formatReviewShipmentWeightsDate(row.actualDeliveryDate) : '-'),
    {
      id: 'deliveryDate',
      isFilterable: false,
    },
  ),
];

const ReviewShipmentWeightsTable = (props) => {
  const { tableColumns, tableData, noRowsMsg } = props;

  const { getTableProps, getTableBodyProps, headerGroups, rows, prepareRow } = useTable(
    {
      columns: tableColumns,
      data: tableData,
      manualFilters: false,
      manualPagination: false,
      manualSortBy: false,
      disableMultiSort: true,
      defaultCanSort: false,
      disableSortBy: true,
      autoResetSortBy: false,
      // If this option is true, the filters we get back from this hook
      // will not be memoized, which makes it easy to get into infinite render loops
      autoResetFilters: false,
    },
    useFilters,
    useSortBy,
    usePagination,
  );
  return (
    <div data-testid="table-queue" className={styles.ReviewShipmentWeightsTable}>
      {rows.length > 0 ? (
        <div className={styles.tableContainer}>
          <Table
            getTableProps={getTableProps}
            getTableBodyProps={getTableBodyProps}
            headerGroups={headerGroups}
            rows={rows}
            prepareRow={prepareRow}
            handleClick={() => {}}
          />
        </div>
      ) : (
        <p>{noRowsMsg || 'No results found.'}</p>
      )}
    </div>
  );
};

export default ReviewShipmentWeightsTable;
