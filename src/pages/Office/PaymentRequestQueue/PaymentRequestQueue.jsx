import React, { useMemo } from 'react';
import { GridContainer } from '@trussworks/react-uswds';
import { withRouter } from 'react-router-dom';
import { useTable } from 'react-table';

import styles from './PaymentRequestQueue.module.scss';

import { usePaymentRequestQueueQueries } from 'hooks/queries';
import Table from 'components/Table/Table';
import { createHeader } from 'components/Table/utils';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { HistoryShape } from 'types/router';
import {
  departmentIndicatorLabel,
  formatDateFromIso,
  formatAgeToDays,
  paymentRequestStatusReadable,
} from 'shared/formatters';

const columns = [
  createHeader('ID', 'id'),
  createHeader(
    'Customer name',
    (row) => {
      return `${row.customer.last_name}, ${row.customer.first_name}`;
    },
    { id: 'name' },
  ),
  createHeader('DoD ID', 'customer.dodID'),
  createHeader(
    'Status',
    (row) => {
      return paymentRequestStatusReadable(row.status);
    },
    'status',
  ),
  createHeader(
    'Age',
    (row) => {
      return formatAgeToDays(row.age);
    },
    'age',
  ),
  createHeader(
    'Submitted',
    (row) => {
      return formatDateFromIso(row.submittedAt, 'DD MMM YYYY');
    },
    'submittedAt',
  ),
  createHeader('Move ID', 'locator'),
  createHeader(
    'Branch',
    (row) => {
      return departmentIndicatorLabel(row.departmentIndicator);
    },
    { id: 'branch' },
  ),
  createHeader('Origin GBLOC', 'originGBLOC'),
];

const PaymentRequestQueue = ({ history }) => {
  const {
    queuePaymentRequestsResult: { totalCount, queuePaymentRequests = [] },
    isLoading,
    isError,
  } = usePaymentRequestQueueQueries();

  // react-table setup below
  const tableData = useMemo(() => queuePaymentRequests, [queuePaymentRequests]);
  const tableColumns = useMemo(() => columns, []);
  const { getTableProps, getTableBodyProps, headerGroups, rows, prepareRow } = useTable({
    columns: tableColumns,
    data: tableData,
    initialState: { hiddenColumns: ['id'] },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const handleClick = (values) => {
    history.push(`/moves/MOVE_CODE/payment-requests/${values.id}`);
  };

  return (
    <GridContainer containerSize="widescreen" className={styles.PaymentRequestQueue}>
      <h1>{`Payment requests (${totalCount})`}</h1>
      <div className={styles.tableContainer}>
        <Table
          handleClick={handleClick}
          getTableProps={getTableProps}
          getTableBodyProps={getTableBodyProps}
          headerGroups={headerGroups}
          rows={rows}
          prepareRow={prepareRow}
        />
      </div>
    </GridContainer>
  );
};

PaymentRequestQueue.propTypes = {
  history: HistoryShape.isRequired,
};

export default withRouter(PaymentRequestQueue);
