import React from 'react';
import { GridContainer } from '@trussworks/react-uswds';

import styles from './PaymentRequestQueue.module.scss';

import Table from 'components/Table/Table';
import { createHeader } from 'components/Table/utils';

const columns = [
  createHeader('Customer name', ''),
  createHeader('DoD ID', 'customer.dodID'),
  createHeader('Status', 'status'),
  createHeader('Age', 'age'),
  createHeader('Submitted', 'submittedAt'),
  createHeader('Move ID', 'locator'),
  createHeader('Branch', 'departmentIndicator'),
  createHeader('Destination duty station', 'destinationDutyStation.name'),
  createHeader('Origin GBLOC', 'originGBLOC'),
];

const PaymentRequestQueue = () => {
  return (
    <GridContainer containerSize="widescreen" className={styles.PaymentRequestQueue}>
      <h1>Payment requests (0)</h1>
      <div className={styles.tableContainer}>
        <Table columns={columns} />
      </div>
    </GridContainer>
  );
};

export default PaymentRequestQueue;
