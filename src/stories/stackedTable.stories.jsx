import React from 'react';
import { action } from '@storybook/addon-actions';
import { storiesOf } from '@storybook/react';
import { Button } from '@trussworks/react-uswds';
import { ReactComponent as EditIcon } from 'shared/icon/edit.svg';
import { StackedTable, StackedTableRow, StackedTableHeader, StackedTableData } from '../components/StackedTable';

const StackedTableExample = () => (
  <div style={{ background: 'F0F0F0' }}>
    <StackedTable fullWidth>
      <caption>
        <h4>Orders</h4>
      </caption>
      <StackedTableRow>
        <StackedTableHeader>Table Header</StackedTableHeader>
        <StackedTableData>Table Data</StackedTableData>
      </StackedTableRow>
      <StackedTableRow>
        <StackedTableHeader>Table Header</StackedTableHeader>
        <StackedTableData>Table Data</StackedTableData>
      </StackedTableRow>
    </StackedTable>
  </div>
);

const StackedTableWithButtons = () => (
  <div>
    <StackedTable fullWidth>
      <caption>
        <h4>Orders</h4>
      </caption>
      <StackedTableRow>
        <StackedTableHeader>Table Header</StackedTableHeader>
        <StackedTableData>
          Table Data
          <Button className="usa-button--unstyled" onClick={action('should open edit form')}>
            <span className="icon">
              <EditIcon />
            </span>
            <span>Edit</span>
          </Button>
        </StackedTableData>
      </StackedTableRow>
      <StackedTableRow>
        <StackedTableHeader>Table Header</StackedTableHeader>
        <StackedTableData>
          Table Data
          <Button className="usa-button--unstyled" onClick={action('should open edit form')}>
            <span className="icon">
              <EditIcon />
            </span>
            <span>Edit</span>
          </Button>
        </StackedTableData>
      </StackedTableRow>
      <StackedTableRow>
        <StackedTableHeader>Table Header</StackedTableHeader>
        <StackedTableData>
          Table Data
          <Button className="usa-button--unstyled" onClick={action('should open edit form')}>
            <span className="icon">
              <EditIcon />
            </span>
            <span>Edit</span>
          </Button>
        </StackedTableData>
      </StackedTableRow>
    </StackedTable>
  </div>
);

storiesOf('Components|StackedTable', module)
  .add('default', () => <StackedTableExample />)
  .add('with buttons to edit', () => <StackedTableWithButtons />);
