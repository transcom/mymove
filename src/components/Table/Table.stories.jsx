import React from 'react';
import { withKnobs } from '@storybook/addon-knobs';

import { createHeader } from './utils';
import Table from './Table';

export default {
  title: 'TOO/TIO Components|Table',
  decorators: [
    withKnobs,
    (storyFn) => (
      <div style={{ margin: '10px', height: '80vh', display: 'flex', flexDirection: 'column', overflow: 'auto' }}>
        {storyFn()}
      </div>
    ),
  ],
};

const data = [
  {
    col1: 'Banks, Aaliyah',
    col2: '987654321',
    col3: 'New move',
    col4: 'LCKMAJ',
    col5: 'Navy',
    col6: '3',
    col7: 'NAS Jacksonville',
    col8: 'HAFC',
    col9: 'Garimundi, J (SW)',
  },
  {
    col1: 'Childers, Jamie',
    col2: '987654321',
    col3: 'New move',
    col4: 'XCQ5ZH',
    col5: 'Navy',
    col6: '3',
    col7: 'NAS Jacksonville',
    col8: 'HAFC',
    col9: 'Garimundi, J (SW)',
  },
  {
    col1: 'Clark-Nunez, Sofia',
    col2: '987654321',
    col3: 'New move',
    col4: 'UCAF8Q',
    col5: 'Navy',
    col6: '3',
    col7: 'NAS Jacksonville',
    col8: 'HAFC',
    col9: 'Garimundi, J (SW)',
  },
];

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

export const TXOTable = () => <Table data={data} columns={columns} />;
