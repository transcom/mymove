import React, { useCallback, useState } from 'react';
import { withRouter } from 'react-router-dom';
import { GridContainer } from '@trussworks/react-uswds';

import styles from './QAECSRMoveSearch.module.scss';

import { HistoryShape } from 'types/router';
import { useQAECSRMoveSearchQueries } from 'hooks/queries';
import SearchResultsTable from 'components/Table/SearchResultsTable';
import MoveSearchForm from 'components/MoveSearchForm/MoveSearchForm';

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

  return (
    <div className={styles.QAECSRMoveSearchWrapper}>
      <GridContainer data-testid="move-search" containerSize="widescreen" className={styles.QAECSRMoveSearchPage}>
        <h1>Search for a move</h1>
        <MoveSearchForm onSubmit={onSubmit} />
        {searchHappened && (
          <SearchResultsTable
            showFilters
            showPagination
            defaultCanSort
            disableMultiSort
            disableSortBy={false}
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
