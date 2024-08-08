import React, { useCallback, useState } from 'react';
import { useNavigate } from 'react-router-dom';

import styles from './QAECSRMoveSearch.module.scss';

import { useMoveSearchQueries } from 'hooks/queries';
import SearchResultsTable from 'components/Table/SearchResultsTable';
import MoveSearchForm from 'components/MoveSearchForm/MoveSearchForm';
import { isNullUndefinedOrWhitespace } from 'shared/utils';

const QAECSRMoveSearch = () => {
  const navigate = useNavigate();
  const [search, setSearch] = useState({ moveCode: null, dodID: null, customerName: null, paymentRequestCode: null });
  const [searchHappened, setSearchHappened] = useState(false);

  const handleClick = useCallback(
    (values) => {
      navigate(`/moves/${values.locator}/details`);
    },
    [navigate],
  );
  const onSubmit = useCallback((values) => {
    const payload = {
      moveCode: null,
      dodID: null,
      customerName: null,
      paymentRequestCode: null,
    };
    if (!isNullUndefinedOrWhitespace(values.searchText)) {
      if (values.searchType === 'moveCode') {
        payload.moveCode = values.searchText;
      } else if (values.searchType === 'dodID') {
        payload.dodID = values.searchText;
      } else if (values.searchType === 'customerName') {
        payload.customerName = values.searchText;
      } else if (values.searchType === 'paymentRequestCode') {
        payload.paymentRequestCode = values.searchText.trim();
      }
    }

    setSearch(payload);
    setSearchHappened(true);
  }, []);

  return (
    <div className={styles.QAECSRMoveSearchWrapper}>
      <div data-testid="move-search" className={styles.QAECSRMoveSearchPage}>
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
            useQueries={useMoveSearchQueries}
            moveCode={search.moveCode}
            dodID={search.dodID}
            customerName={search.customerName}
            paymentRequestCode={search.paymentRequestCode}
          />
        )}
      </div>
    </div>
  );
};

QAECSRMoveSearch.propTypes = {};

export default QAECSRMoveSearch;
