import React from 'react';
import { action } from '@storybook/addon-actions';
import { storiesOf } from '@storybook/react';
import { Button } from '@trussworks/react-uswds';
import { ReactComponent as EditIcon } from 'shared/icon/edit.svg';
import { StackedTable, StackedTableRow, StackedTableHeader, StackedTableData } from '../components/StackedTable';

const StackedTableExample = () => (
  <div className="table--stacked">
    <StackedTable fullWidth>
      <caption>
        <h4>Orders</h4>
      </caption>
      <col style={{ width: '25%' }} />
      <col style={{ width: '75%' }} />
      <StackedTableRow>
        <StackedTableHeader>Orders number</StackedTableHeader>
        <StackedTableData>999999999</StackedTableData>
      </StackedTableRow>
      <StackedTableRow>
        <StackedTableHeader>Authorized Entitlement</StackedTableHeader>
        <StackedTableData>999999999</StackedTableData>
      </StackedTableRow>
    </StackedTable>
  </div>
);

const StackedTableWithButtons = () => (
  <div className="table--stacked table--stacked-wbuttons">
    <StackedTable fullWidth>
      <caption>
        <h4>Orders</h4>
      </caption>
      <col style={{ width: '25%' }} />
      <col style={{ width: '75%' }} />
      <StackedTableRow>
        <StackedTableHeader>Orders number</StackedTableHeader>
        <StackedTableData>
          999999999
          <Button className="usa-button--unstyled" onClick={action('should open edit form')}>
            <span className="icon">
              <EditIcon />
            </span>
            <span>Edit</span>
          </Button>
        </StackedTableData>
      </StackedTableRow>
      <StackedTableRow>
        <StackedTableHeader>Orders number</StackedTableHeader>
        <StackedTableData>
          999999999
          <Button className="usa-button--unstyled" onClick={action('should open edit form')}>
            <span className="icon">
              <EditIcon />
            </span>
            <span>Edit</span>
          </Button>
        </StackedTableData>
      </StackedTableRow>
      <StackedTableRow>
        <StackedTableHeader>Orders number</StackedTableHeader>
        <StackedTableData>
          999999999
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
