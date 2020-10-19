import React from 'react';
import { GridContainer } from '@trussworks/react-uswds';
import { withRouter } from 'react-router-dom';

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
  const { queuePaymentRequestsResult, isLoading, isError } = usePaymentRequestQueueQueries();

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  // eslint-disable-next-line no-unused-vars
  const { page, perPage, totalCount, queuePaymentRequests } = queuePaymentRequestsResult[`${undefined}`];

  const handleClick = (values) => {
    history.push(`/moves/MOVE_CODE/payment-requests/${values.id}`);
  };

  return (
    <GridContainer containerSize="widescreen" className={styles.PaymentRequestQueue}>
      <h1>{`Payment requests (${totalCount})`}</h1>
      <div className={styles.tableContainer}>
        <Table columns={columns} data={queuePaymentRequests} hiddenColumns={['id']} handleClick={handleClick} />
      </div>
    </GridContainer>
  );
};

PaymentRequestQueue.propTypes = {
  history: HistoryShape.isRequired,
};

export default withRouter(PaymentRequestQueue);
