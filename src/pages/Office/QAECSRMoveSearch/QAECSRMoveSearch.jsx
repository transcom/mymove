import React, { useCallback, useMemo, useState } from 'react';
import { withRouter } from 'react-router-dom';
import { GridContainer } from '@trussworks/react-uswds';

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
    isFilterable: false,
  }),
  createHeader('DOD ID', 'dodID', {
    id: 'dodID',
    isFilterable: false,
  }),
  createHeader(
    'Customer name',
    (row) => {
      return `${row.lastName}, ${row.firstName}`;
    },
    {
      id: 'customerName',
      isFilterable: false,
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
      return row.originDutyLocationPostalCode;
    },
    {
      id: 'originPostalCode',
      isFilterable: true,
    },
  ),
  createHeader(
    'Destination ZIP',
    (row) => {
      return row.destinationDutyLocationPostalCode;
    },
    {
      id: 'destinationPostalCode',
      isFilterable: true,
    },
  ),
  createHeader(
    'Branch',
    (row) => {
      return serviceMemberAgencyLabel(row.branch);
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
  createHeader(
    'Number of Shipments',
    (row) => {
      return Number(row.shipmentsCount || 0);
    },
    { id: 'shipmentsCount', isFilterable: true },
  ),
];

const QAECSRMoveSearch = ({ history }) => {
  const [search, setSearch] = useState({ moveCode: null, dodID: null, customerName: null });
  const [searchHappened, setSearchHappened] = useState(false);

  const handleClick = useCallback(
    (values) => {
      history.push(`/moves/${values.locator}/details`);
    },
    [history],
  );
  const onSubmit = useCallback((values) => {
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
  }, []);

  const isLoading = false;
  const isError = false;

  const tableColumns = useMemo(() => columns(true), []);

  return (
    <div className={styles.QAECSRMoveSearchWrapper}>
      <GridContainer data-testid="move-search" containerSize="widescreen" className={styles.QAECSRMoveSearchPage}>
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
            disableSortBy={false}
            columns={tableColumns}
            title="Results"
            handleClick={handleClick}
            useQueries={useQAECSRMoveSearchQueries}
            moveCode={search.moveCode}
            dodID={search.dodID}
            customerName={search.customerName}
          />
        )}
      </GridContainer>
    </div>
  );
};

QAECSRMoveSearch.propTypes = {
  history: HistoryShape.isRequired,
};

export default withRouter(QAECSRMoveSearch);
