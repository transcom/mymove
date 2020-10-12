import React from 'react';
import { GridContainer } from '@trussworks/react-uswds';

import styles from './MoveQueue.module.scss';

import Table from 'components/Table/Table';
import { createHeader } from 'components/Table/utils';

const columns = [
  createHeader('Customer name', 'col1'),
  createHeader('DoD ID', 'col2'),
  createHeader('Status', 'col3'),
  createHeader('Move ID', 'col4'),
  createHeader('Branch', 'col5'),
  createHeader('# of shipments', 'col6'),
  createHeader('Destination duty station', 'col7'),
  createHeader('Origin GBLOC', 'col8'),
  createHeader('Last modified by', 'col9'),
];

const MoveQueue = () => {
  return (
    <GridContainer containerSize="widescreen" className={styles.MoveQueue}>
      <h1>All moves</h1>
      <div className={styles.tableContainer}>
        <Table columns={columns} />
      </div>
    </GridContainer>
  );
};

export default MoveQueue;
