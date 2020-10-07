import React from 'react';
import { withKnobs } from '@storybook/addon-knobs';

import { createStringHeader } from './utils';
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
  createStringHeader('Customer name', 'col1'),
  createStringHeader('DoD ID', 'col2'),
  createStringHeader('Status', 'col3'),
  createStringHeader('Move ID', 'col4'),
  createStringHeader('Branch', 'col5'),
  createStringHeader('# of shipments', 'col6'),
  createStringHeader('Destination duty station', 'col7'),
  createStringHeader('Origin GBLOC', 'col8'),
  createStringHeader('Last modified by', 'col9'),
];

export const Default = () => <Table data={data} columns={columns} />;
