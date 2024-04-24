import React from 'react';
import { useNavigate } from 'react-router-dom';

import styles from './PaymentRequestQueue.module.scss';

import { usePaymentRequestQueueQueries, useUserQueries } from 'hooks/queries';
import { createHeader } from 'components/Table/utils';
import {
  formatDateFromIso,
  serviceMemberAgencyLabel,
  paymentRequestStatusReadable,
  formatAgeToDays,
} from 'utils/formatters';
import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import DateSelectFilter from 'components/Table/Filters/DateSelectFilter';
import { BRANCH_OPTIONS, GBLOC, PAYMENT_REQUEST_STATUS_OPTIONS } from 'constants/queues';
import TableQueue from 'components/Table/TableQueue';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { CHECK_SPECIAL_ORDERS_TYPES, SPECIAL_ORDERS_TYPES } from 'constants/orders';

const columns = (showBranchFilter = true) => [
  createHeader('ID', 'id'),
  createHeader(
    'Customer name',
    (row) => {
      return (
        <div>
          {CHECK_SPECIAL_ORDERS_TYPES(row.orderType) ? (
            <span className={styles.specialMoves}>{SPECIAL_ORDERS_TYPES[`${row.orderType}`]}</span>
          ) : null}
          {`${row.customer.last_name}, ${row.customer.first_name}`}
        </div>
      );
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
      Filter: (props) => (
        <div data-testid="statusFilter">
          {/* eslint-disable-next-line react/jsx-props-no-spreading */}
          <MultiSelectCheckBoxFilter options={PAYMENT_REQUEST_STATUS_OPTIONS} {...props} />
        </div>
      ),
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
      // eslint-disable-next-line react/jsx-props-no-spreading
      Filter: (props) => <DateSelectFilter dateTime {...props} />,
    },
  ),
  createHeader('Move Code', 'locator', {
    id: 'locator',
    isFilterable: true,
  }),
  createHeader(
    'Branch',
    (row) => {
      return serviceMemberAgencyLabel(row.customer.agency);
    },
    {
      id: 'branch',
      isFilterable: showBranchFilter,
      Filter: (props) => (
        // eslint-disable-next-line react/jsx-props-no-spreading
        <SelectFilter options={BRANCH_OPTIONS} {...props} />
      ),
    },
  ),
  createHeader('Origin GBLOC', 'originGBLOC', { disableSortBy: true }),
  createHeader(
    'Origin Duty Location',
    (row) => {
      return row.originDutyLocation.name;
    },
    {
      id: 'originDutyLocation',
      isFilterable: true,
    },
  ),
];

const PaymentRequestQueue = () => {
  const navigate = useNavigate();
  const {
    // eslint-disable-next-line camelcase
    data: { office_user },
    isLoading,
    isError,
  } = useUserQueries();

  // eslint-disable-next-line camelcase
  const showBranchFilter = office_user?.transportation_office?.gbloc !== GBLOC.USMC;

  const handleClick = (values) => {
    navigate(`/moves/${values.locator}/payment-requests`);
  };

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <div className={styles.PaymentRequestQueue}>
      <TableQueue
        showFilters
        showPagination
        manualSortBy
        defaultCanSort
        defaultSortedColumns={[{ id: 'age', desc: true }]}
        disableMultiSort
        disableSortBy={false}
        columns={columns(showBranchFilter)}
        title="Payment requests"
        handleClick={handleClick}
        useQueries={usePaymentRequestQueueQueries}
      />
    </div>
  );
};

export default PaymentRequestQueue;
