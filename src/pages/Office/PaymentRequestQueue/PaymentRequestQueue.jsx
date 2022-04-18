import React from 'react';
import { withRouter } from 'react-router-dom';

import styles from './PaymentRequestQueue.module.scss';

import { usePaymentRequestQueueQueries, useUserQueries } from 'hooks/queries';
import { createHeader } from 'components/Table/utils';
import { HistoryShape } from 'types/router';
import { formatAgeToDays } from 'shared/formatters';
import { formatDateFromIso, serviceMemberAgencyLabel, paymentRequestStatusReadable } from 'utils/formatters';
import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import DateSelectFilter from 'components/Table/Filters/DateSelectFilter';
import { BRANCH_OPTIONS, GBLOC, PAYMENT_REQUEST_STATUS_OPTIONS } from 'constants/queues';
import TableQueue from 'components/Table/TableQueue';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const columns = (showBranchFilter = true) => [
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

const PaymentRequestQueue = ({ history }) => {
  const {
    // eslint-disable-next-line camelcase
    data: { office_user },
    isLoading,
    isError,
  } = useUserQueries();

  const showBranchFilter = office_user?.transportation_office?.gbloc !== GBLOC.USMC;

  const handleClick = (values) => {
    history.push(`/moves/${values.locator}/payment-requests`);
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

PaymentRequestQueue.propTypes = {
  history: HistoryShape.isRequired,
};

export default withRouter(PaymentRequestQueue);
