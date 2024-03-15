import React from 'react';
import { useParams } from 'react-router-dom';

import QueueNav from 'components/QueueNav/QueueNav';

const Queue = () => {
  const { queueType } = useParams();

  if (queueType === 'Search') {

    return (
      <div data-testid="move-search" className={styles.ServicesCounselingQueue}>
        {renderNavBar()}
        <h1>Search for a move</h1>
        <MoveSearchForm onSubmit={onSubmit} role={roleTypes.SERVICES_COUNSELOR} />
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
            roleType={roleTypes.SERVICES_COUNSELOR}
          />
        )}
      </div>
    );
  }
  else{

  }

  return <QueueNav />;
};
