import React from 'react';
import { capitalize } from 'lodash';
import { formatDate, formatDateTimeWithTZ } from 'shared/formatters';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faClock from '@fortawesome/fontawesome-free-solid/faClock';
import faExclamationCircle from '@fortawesome/fontawesome-free-solid/faExclamationCircle';
import { formatDate4DigitYear } from 'shared/formatters';
import moment from 'moment';
import { getEntitlements } from 'shared/entitlements';
import { sitDaysUsed, sitTotalDaysUsed } from 'shared/StorageInTransit/calculator';

// Abstracting react table column creation
const CreateReactTableColumn = (header, accessor, options = {}) => ({
  Header: header,
  accessor: accessor,
  ...options,
});

const status = CreateReactTableColumn('Status', 'synthetic_status', {
  Cell: row => (
    <span className="status" data-cy="status">
      {capitalize(row.value && row.value.replace('_', ' '))}
    </span>
  ),
});

const hhgStatus = CreateReactTableColumn('HHG status', 'hhg_status', {
  Cell: row => (
    <span className="status" data-cy="status">
      {row.value &&
        row.value
          .replace('_', ' ')
          .split(' ')
          .map(word => capitalize(word))
          .join(' ')}
    </span>
  ),
});

const customerName = CreateReactTableColumn('Customer name', 'customer_name');

const dodId = CreateReactTableColumn('DoD ID', 'edipi');

const rank = CreateReactTableColumn('Rank', 'rank', {
  Cell: row => <span className="rank">{row.value && row.value.replace('_', '-')}</span>,
});

const shipments = CreateReactTableColumn('Shipments', 'shipments');

const locator = CreateReactTableColumn('Locator #', 'locator', {
  Cell: row => <span data-cy="locator">{row.value}</span>,
});

const gbl = CreateReactTableColumn('GBL #', 'gbl_number');

const moveDate = CreateReactTableColumn('PPM start', 'move_date', {
  Cell: row => <span className="move_date">{formatDate(row.value)}</span>,
});

const pickupDate = CreateReactTableColumn('Pickup', 'move_date', {
  Cell: row => <span className="move_date">{formatDate(row.value)}</span>,
});

const lastModifiedDate = CreateReactTableColumn('Last modified', 'last_modified_date', {
  Cell: row => <span className="updated_at">{formatDateTimeWithTZ(row.value)}</span>,
});

const submittedDate = CreateReactTableColumn('Submitted', 'submitted_date', {
  Cell: row => <span className="submitted_date">{formatDateTimeWithTZ(row.value)}</span>,
});

const origin = CreateReactTableColumn('Origin', 'origin_duty_station_name', {
  Cell: row => <span>{row.value}</span>,
});

const destination = CreateReactTableColumn('Destination', 'destination_duty_station_name', {
  Cell: row => <span>{row.value}</span>,
});

const originGBLOC = CreateReactTableColumn('Origin GBLOC', 'origin_gbloc', {
  Cell: row => <span>{row.value}</span>,
});

const destinationGBLOC = CreateReactTableColumn('Destination GBLOC', 'destination_gbloc', {
  Cell: row => <span>{row.value}</span>,
});

const deliveredDate = CreateReactTableColumn('Delivered', 'delivered_date', {
  Cell: row => <span>{formatDate(row.value)}</span>,
});

const invoiceApprovedDate = CreateReactTableColumn('Invoice Approved', 'invoice_approved_date', {
  Cell: row => <span>{formatDate(row.value)}</span>,
});

const sitExpires = CreateReactTableColumn(
  'SIT expires',
  row => {
    if (row.storage_in_transits && row.storage_in_transits.some(sit => sit.actual_start_date)) {
      return formatDate4DigitYear(
        moment.min(
          row.storage_in_transits.filter(sit => sit.actual_start_date).map(sit => {
            return moment(sit.actual_start_date).add(
              getEntitlements(row.rank).storage_in_transit +
                sitDaysUsed(sit) -
                sitTotalDaysUsed(row.storage_in_transits),
              'days',
            );
          }),
        ),
      );
    }
    return null;
  },
  {
    Cell: row => <span>{row.value}</span>,
    id: 'sit_expires',
  },
);

export const calculateNeedsAttention = row => {
  const attentions = [];
  if ((row.hhg_status && row.hhg_status === 'ACCEPTED') || row.status === 'SUBMITTED') {
    attentions.push('Awaiting review');
  }
  if (row.has_unapproved_shipment_line_items) {
    attentions.push('Pre-approval requested');
  }
  if (
    row.storage_in_transits &&
    row.storage_in_transits.some(sit => sit.status === 'REQUESTED' && sit.location === 'ORIGIN')
  ) {
    attentions.push('Origin SIT requested');
  }
  if (
    row.storage_in_transits &&
    row.storage_in_transits.some(sit => sit.status === 'REQUESTED' && sit.location === 'DESTINATION')
  ) {
    attentions.push('Dest SIT requested');
  }
  return attentions;
};

const clockCell = value => {
  if (value === 'CLOCK') {
    return (
      <span data-cy="ppm-queue-icon">
        <FontAwesomeIcon icon={faClock} className="clock-icon" />
      </span>
    );
  } else if (value === 'BANG') {
    return (
      <span data-cy="ppm-queue-icon">
        <FontAwesomeIcon icon={faExclamationCircle} className="bang-icon" />
      </span>
    );
  }
  return '';
};

const needsAttentionClockIcon = CreateReactTableColumn(
  <FontAwesomeIcon icon={faClock} />,
  row => {
    const attentions = calculateNeedsAttention(row);
    if (attentions.length > 0) {
      if (attentions.includes('Awaiting review') && row.pm_survey_conducted_date) {
        return 'BANG';
      }
      return 'CLOCK';
    }
    return 'NONE';
  },
  {
    id: 'clockIcon',
    Cell: row => clockCell(row.value),
    width: 50,
  },
);

const needsAttention = CreateReactTableColumn('Needs Attention', calculateNeedsAttention, {
  Cell: row => (
    <div>
      {row.value.map((attention, index) => (
        <span key={index} className="needs-attention-alert">
          {attention}
        </span>
      ))}
    </div>
  ),
  id: 'needs_attention',
});

// Columns used to display in react table
export const newColumns = [
  needsAttentionClockIcon,
  needsAttention,
  customerName,
  hhgStatus,
  shipments,
  origin,
  dodId,
  locator,
  pickupDate,
  submittedDate,
];

export const ppmColumns = [status, customerName, origin, destination, dodId, locator, moveDate, lastModifiedDate];

export const hhgActiveColumns = [
  needsAttentionClockIcon,
  needsAttention,
  customerName,
  hhgStatus,
  origin,
  destination,
  locator,
  gbl,
  pickupDate,
  sitExpires,
];

export const hhgDeliveredColumns = [
  needsAttention,
  customerName,
  originGBLOC,
  destinationGBLOC,
  locator,
  gbl,
  deliveredDate,
  invoiceApprovedDate,
];

export const defaultColumns = [status, customerName, dodId, rank, locator, gbl, moveDate, lastModifiedDate];
