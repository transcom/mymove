import React, { useMemo, useState } from 'react';
import { withRouter } from 'react-router-dom';

import styles from './QAECSRMoveSearch.module.scss';

import { HistoryShape } from 'types/router';
import { createHeader } from 'components/Table/utils';
import { useQAECSRMoveSearchQueries } from 'hooks/queries';
import { serviceMemberAgencyLabel } from 'utils/formatters';
import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import { BRANCH_OPTIONS, MOVE_STATUS_OPTIONS, MOVE_STATUS_LABELS } from 'constants/queues';
import SearchResultsTable from 'components/Table/SearchResultsTable';
import MoveSearchForm from 'components/MoveSearchForm/MoveSearchForm';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const columns = (showBranchFilter = true) => [
  createHeader('Move code', 'locator', {
    id: 'locator',
    isFilterable: true,
  }),
  createHeader('DOD ID', 'customer.dodID', {
    id: 'dodID',
    isFilterable: true,
  }),
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
  createHeader(
    'Status',
    (row) => {
      return MOVE_STATUS_LABELS[`${row.status}`];
    },
    {
      id: 'status',
      isFilterable: true,
      // eslint-disable-next-line react/jsx-props-no-spreading
      Filter: (props) => <MultiSelectCheckBoxFilter options={MOVE_STATUS_OPTIONS} {...props} />,
    },
  ),
  createHeader(
    'Origin ZIP',
    (row) => {
      return `${row.originDutyLocation?.address?.postalCode}`;
    },
    {
      id: 'originZIP',
      isFilterable: true,
    },
  ),
  createHeader(
    'Destination ZIP',
    (row) => {
      return `${row.destinationDutyLocation?.address?.postalCode}`;
    },
    {
      id: 'destinationZIP',
      isFilterable: true,
    },
  ),
  createHeader(
    'Branch',
    (row) => {
      return serviceMemberAgencyLabel(row.customer?.agency);
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
  createHeader('# of shipments', 'shipmentsCount', { disableSortBy: true }),
];

const QAECSRMoveSearch = ({ history }) => {
  const [search, setSearch] = useState({ moveCode: null, dodID: null, customerName: null });
  const [searchHappened, setSearchHappened] = useState(false);

  const handleClick = (values) => {
    history.push(`/moves/${values.locator}/details`);
  };
  const onSubmit = (values) => {
    const payload = {
      moveCode: null,
      dodID: null,
      customerName: null,
    };
    if (values.searchType === 'moveCode') {
      payload.moveCode = values.searchText;
    } else if (values.searchType === 'dodID') {
      payload.dodID = values.searchText;
    } else if (values.searchType === 'customerName') {
      payload.customerName = values.searchText;
    }
    setSearch(payload);
    setSearchHappened(true);
  };

  const { searchResult, isLoading, isError } = useQAECSRMoveSearchQueries({
    moveCode: search.moveCode,
    dodID: search.dodID,
    customerName: search.customerName,
  });

  const { data = [] } = searchResult;
  const tableColumns = useMemo(() => columns(true), []);
  const tableData = useMemo(() => data, [data]);
  return (
    <div className={styles.QAECSRMoveSearchPage}>
      <h1>Search for a move</h1>
      <MoveSearchForm onSubmit={onSubmit} />
      {isLoading && <LoadingPlaceholder />}
      {!isLoading && isError && <SomethingWentWrong />}
      {!isLoading && !isError && searchHappened && (
        <SearchResultsTable
          showFilters
          showPagination
          defaultCanSort
          disableMultiSort
          manualFilters={false}
          disableSortBy={false}
          columns={tableColumns}
          title="Results"
          handleClick={handleClick}
          useQueries={useQAECSRMoveSearchQueries}
          data={tableData}
        />
      )}
    </div>
  );
};

QAECSRMoveSearch.propTypes = {
  history: HistoryShape.isRequired,
};

export default withRouter(QAECSRMoveSearch);
