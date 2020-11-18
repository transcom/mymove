import React from 'react';
import { withRouter } from 'react-router-dom';

import styles from './PaymentRequestQueue.module.scss';

import { usePaymentRequestQueueQueries } from 'hooks/queries';
import { createHeader } from 'components/Table/utils';
import { HistoryShape } from 'types/router';
import {
  formatDateFromIso,
  formatAgeToDays,
  paymentRequestStatusReadable,
  serviceMemberAgencyLabel,
} from 'shared/formatters';
import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import DateSelectFilter from 'components/Table/Filters/DateSelectFilter';
import { BRANCH_OPTIONS, PAYMENT_REQUEST_STATUS_OPTIONS } from 'constants/queues';
import TableQueue from 'components/Table/TableQueue';

const paymentRequestStatusOptions = Object.keys(PAYMENT_REQUEST_STATUS_OPTIONS).map((key) => ({
  value: key,
  label: PAYMENT_REQUEST_STATUS_OPTIONS[`${key}`],
}));

const branchFilterOptions = [
  { value: '', label: 'All' },
  ...Object.keys(BRANCH_OPTIONS).map((key) => ({
    value: key,
    label: BRANCH_OPTIONS[`${key}`],
  })),
];

const columns = [
  createHeader('ID', 'id'),
  createHeader(
    'Customer name',
    (row) => {
      return `${row.customer.last_name}, ${row.customer.first_name}`;
    },
    {
      id: 'lastName',
      isFilterable: true,
    },
  ),
  createHeader('DoD ID', 'customer.dodID', {
    id: 'dodID',
    isFilterable: true,
  }),
  createHeader(
    'Status',
    (row) => {
      return paymentRequestStatusReadable(row.status);
    },
    {
      id: 'status',
      isFilterable: true,
      // eslint-disable-next-line react/jsx-props-no-spreading
      Filter: (props) => <MultiSelectCheckBoxFilter options={paymentRequestStatusOptions} {...props} />,
    },
  ),
  createHeader(
    'Age',
    (row) => {
      return formatAgeToDays(row.age);
    },
    { id: 'age' },
  ),
  createHeader(
    'Submitted',
    (row) => {
      return formatDateFromIso(row.submittedAt, 'DD MMM YYYY');
    },
    {
      id: 'submittedAt',
      isFilterable: true,
      Filter: DateSelectFilter,
    },
  ),
  createHeader('Move Code', 'locator', {
    id: 'moveID',
    isFilterable: true,
  }),
  createHeader(
    'Branch',
    (row) => {
      return serviceMemberAgencyLabel(row.customer.agency);
    },
    {
      id: 'branch',
      isFilterable: true,
      // eslint-disable-next-line react/jsx-props-no-spreading
      Filter: (props) => <SelectFilter options={branchFilterOptions} {...props} />,
    },
  ),
  createHeader('Origin GBLOC', 'originGBLOC', { disableSortBy: true }),
];

const PaymentRequestQueue = ({ history }) => {
  const handleClick = (values) => {
    history.push(`/moves/MOVE_CODE/payment-requests/${values.id}`);
  };

  return (
    <div className={styles.PaymentRequestQueue}>
      <TableQueue
        showFilters
        showPagination
        manualSortBy
        defaultCanSort
        defaultSortedColumns={[{ id: 'status', desc: false }]}
        columns={columns}
        title="Payment requests"
        handleClick={handleClick}
        useQueries={usePaymentRequestQueueQueries}
      />
    </div>
  );
};

PaymentRequestQueue.propTypes = {
  history: HistoryShape.isRequired,
};

export default withRouter(PaymentRequestQueue);
