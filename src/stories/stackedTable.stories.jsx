import React from 'react';
import { action } from '@storybook/addon-actions';
import { storiesOf } from '@storybook/react';
import { Button } from '@trussworks/react-uswds';
import { ReactComponent as DocsIcon } from 'shared/icon/documents.svg';
import { EditButton } from '../components/form';

const StackedTableExample = () => (
  <div className="table--stacked">
    <div className="display-flex">
      <div>
        <h4>Orders</h4>
      </div>
      <div>
        <Button className="usa-button--icon" onClick={action('clicked')}>
          <span className="icon">
            <DocsIcon />
          </span>
          <span>View orders</span>
        </Button>
      </div>
    </div>
    <table>
      <colgroup>
        <col style={{ width: '25%' }} />
        <col style={{ width: '75%' }} />
      </colgroup>
      <tbody>
        <tr>
          <th scope="row">Orders number</th>
          <td>999999999</td>
        </tr>
        <tr>
          <th scope="row">Authorized Entitlement</th>
          <td>999999999</td>
        </tr>
      </tbody>
    </table>
  </div>
);

const StackedTableWithButtons = () => (
  <div className="table--stacked table--stacked-wbuttons">
    <div className="display-flex">
      <div>
        <h4>Orders</h4>
      </div>
      <div>
        <Button className="usa-button--icon" onClick={action('clicked')}>
          <span className="icon">
            <DocsIcon />
          </span>
          <span>View orders</span>
        </Button>
      </div>
    </div>
    <table>
      <colgroup>
        <col style={{ width: '25%' }} />
        <col style={{ width: '75%' }} />
      </colgroup>
      <tbody>
        <tr>
          <th scope="row">Orders number</th>
          <td>
            999999999
            <EditButton unstyled onClick={action('should open edit form')} />
          </td>
        </tr>
        <tr>
          <th scope="row">Orders number</th>
          <td>
            999999999
            <EditButton unstyled onClick={action('should open edit form')} />
          </td>
        </tr>
        <tr>
          <th scope="row">Orders number</th>
          <td>
            999999999
            <EditButton unstyled onClick={action('should open edit form')} />
          </td>
        </tr>
      </tbody>
    </table>
  </div>
);

storiesOf('Components|StackedTable', module)
  .add('default', () => (
    <div style={{ padding: `20px`, background: `#f0f0f0` }}>
      <StackedTableExample />
    </div>
  ))
  .add('with buttons to edit', () => (
    <div style={{ padding: `20px`, background: `#f0f0f0` }}>
      <StackedTableWithButtons />
    </div>
  ));
