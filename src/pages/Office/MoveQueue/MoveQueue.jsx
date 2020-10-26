import React from 'react';
import { withRouter } from 'react-router-dom';
import { GridContainer } from '@trussworks/react-uswds';

import styles from './MoveQueue.module.scss';

import { HistoryShape } from 'types/router';
import Table from 'components/Table/Table';
import { createHeader } from 'components/Table/utils';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useMovesQueueQueries } from 'hooks/queries';
import { departmentIndicatorLabel } from 'shared/formatters';

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
  createHeader('Status', 'status'),
  createHeader('Move ID', 'locator'),
  createHeader(
    'Branch',
    (row) => {
      return departmentIndicatorLabel(row.departmentIndicator);
    },
    { id: 'branch' },
  ),
  createHeader('# of shipments', 'shipmentsCount'),
  createHeader('Destination duty station', 'destinationDutyStation.name'),
  createHeader('Origin GBLOC', 'originGBLOC'),
];

const MoveQueue = ({ history }) => {
  const { queueMovesResult, isLoading, isError } = useMovesQueueQueries();

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  //  no-unused-vars
  const { page, perPage, totalCount, queueMoves } = queueMovesResult[`${undefined}`];

  const handleClick = (values) => {
    history.push(`/moves/${values.id}/details`);
  };

  return (
    <GridContainer containerSize="widescreen" className={styles.MoveQueue}>
      <h1>{`All moves (${totalCount})`}</h1>
      <div className={styles.tableContainer}>
        <Table columns={columns} data={queueMoves} hiddenColumns={['id']} handleClick={handleClick} />
      </div>
    </GridContainer>
  );
};

MoveQueue.propTypes = {
  history: HistoryShape.isRequired,
};

export default withRouter(MoveQueue);
