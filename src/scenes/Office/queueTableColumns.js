import React from 'react';
import { capitalize } from 'lodash';
import { formatDate, formatDateTimeWithTZ } from 'shared/formatters';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faClock from '@fortawesome/fontawesome-free-solid/faClock';

// Abstracting react table column creation
const CreateReactTableColumn = (header, accessor, options = {}) => ({
  Header: header,
  accessor: accessor,
  ...options,
});

const clockIcon = CreateReactTableColumn(<FontAwesomeIcon icon={faClock} />, 'synthetic_status', {
  id: 'clockIcon',
  Cell: row =>
    row.value === 'PAYMENT_REQUESTED' || row.value === 'SUBMITTED' ? (
      <span data-cy="ppm-queue-icon">
        <FontAwesomeIcon icon={faClock} style={{ color: 'orange' }} />
      </span>
    ) : (
      ''
    ),
  width: 50,
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
      {(row.value && row.value.replace('_', ' '))
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

const gbl = CreateReactTableColumn('GBL', 'gbl_number');

const moveDate = CreateReactTableColumn('Move date', 'move_date', {
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

const sitExpires = CreateReactTableColumn('SIT expires', 'sit_expires', {
  Cell: row => <span>{row.value}</span>,
});

// Columns used to display in react table

export const newColumns = [clockIcon, customerName, locator, dodId, rank, shipments, moveDate, submittedDate];

export const ppmColumns = [clockIcon, status, customerName, dodId, rank, locator, moveDate, lastModifiedDate];

export const hhgActiveColumns = [
  clockIcon,
  customerName,
  hhgStatus,
  origin,
  destination,
  locator,
  gbl,
  moveDate,
  sitExpires,
];

export const defaultColumns = [status, customerName, dodId, rank, locator, gbl, moveDate, lastModifiedDate];
