import React, { useState, useEffect, useMemo } from 'react';
import { GridContainer } from '@trussworks/react-uswds';
import { withRouter } from 'react-router-dom';
import { useTable, useFilters } from 'react-table';

import styles from './PaymentRequestQueue.module.scss';

import { usePaymentRequestQueueQueries } from 'hooks/queries';
import Table from 'components/Table/Table';
import { createHeader } from 'components/Table/utils';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { HistoryShape } from 'types/router';
import {
  formatDateFromIso,
  formatAgeToDays,
  paymentRequestStatusReadable,
  serviceMemberAgencyLabel,
} from 'shared/formatters';
import TextBoxFilter from 'components/Table/Filters/TextBoxFilter';
import SelectFilter from 'components/Table/Filters/SelectFilter';

const branchFilterOptions = [
  { value: '', label: 'All' },
  { value: 'ARMY', label: 'Army' },
  { value: 'NAVY', label: 'Navy' },
  { value: 'MARINES', label: 'Marine Corps' },
  { value: 'AIR_FORCE', label: 'Air Force' },
  { value: 'COAST_GUARD', label: 'Coast Guard' },
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
  createHeader('Origin GBLOC', 'originGBLOC'),
];

const PaymentRequestQueue = ({ history }) => {
  const [paramFilters, setParamFilters] = useState([]);

  const {
    queuePaymentRequestsResult: { totalCount, queuePaymentRequests = [] },
    isLoading,
    isError,
  } = usePaymentRequestQueueQueries(paramFilters);

  // react-table setup below

  const defaultColumn = useMemo(
    () => ({
      // Let's set up our default Filter UI
      Filter: TextBoxFilter,
    }),
    [],
  );
  const tableData = useMemo(() => queuePaymentRequests, [queuePaymentRequests]);
  const tableColumns = useMemo(() => columns, []);
  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    rows,
    prepareRow,
    state: { filters },
  } = useTable(
    {
      columns: tableColumns,
      data: tableData,
      initialState: { hiddenColumns: ['id'] },
      defaultColumn, // Be sure to pass the defaultColumn option
      manualFilters: true,
    },
    useFilters,
  );

  // When these table states change, fetch new data!
  useEffect(() => {
    if (!isLoading && !isError) {
      setParamFilters(filters);
    }
  }, [filters, isLoading, isError]);

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
