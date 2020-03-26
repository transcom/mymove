import React from 'react';
// import { action } from '@storybook/addon-actions';
import { storiesOf } from '@storybook/react';
// import { ReactComponent as EditIcon } from 'shared/images/edit-24px.svg';
import { StackedTable, StackedTableRow, StackedTableHeader, StackedTableData } from '../components/StackedTable';

const StackedTableExample = () => (
  <div>
    <StackedTable fullWidth>
      <StackedTableRow>
        <StackedTableHeader>Orders Number</StackedTableHeader>
        <StackedTableData>999999999</StackedTableData>
      </StackedTableRow>
      <StackedTableRow>
        <StackedTableHeader>Authorized entitlement</StackedTableHeader>
        <StackedTableData>999999999</StackedTableData>
      </StackedTableRow>
    </StackedTable>
  </div>
);

storiesOf('Components|StackedTable', module).add('default', () => <StackedTableExample />);
