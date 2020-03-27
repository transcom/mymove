import React from 'react';
import { action } from '@storybook/addon-actions';
import { storiesOf } from '@storybook/react';
import { Button } from '@trussworks/react-uswds';
import { ReactComponent as EditIcon } from 'shared/images/edit-24px.svg';
import { StackedTable, StackedTableRow, StackedTableHeader, StackedTableData } from '../components/StackedTable';

const StackedTableExample = () => (
  <div>
    <StackedTable fullWidth>
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
      <StackedTableRow>
        <StackedTableHeader>Table Header</StackedTableHeader>
        <StackedTableData>
          Table Data
          <Button className="usa-button--icon" onClick={action('should open edit form')}>
            <span className="icon">
              <EditIcon />
            </span>
            <span>Icon button text</span>
          </Button>
        </StackedTableData>
      </StackedTableRow>
      <StackedTableRow>
        <StackedTableHeader>Table Header</StackedTableHeader>
        <StackedTableData>
          Table Data
          <Button className="usa-button--icon" onClick={action('should open edit form')}>
            <span className="icon">
              <EditIcon />
            </span>
            <span>Icon button text</span>
          </Button>
        </StackedTableData>
      </StackedTableRow>
    </StackedTable>
  </div>
);

storiesOf('Components|StackedTable', module)
  .add('default', () => <StackedTableExample />)
  .add('with buttons to edit', () => <StackedTableWithButtons />);
