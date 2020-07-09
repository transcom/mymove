import React from 'react';
import { action } from '@storybook/addon-actions';
import * as Yup from 'yup';

import { StackedTableRowForm, DocsButton } from '../components/form';

export default {
  title: 'Components|Stacked Table',
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/55411d0f-417d-48d0-9964-2d575ba4b390?mode=design',
    },
  },
};

const StackedTableExample = () => (
  <div>
    <div className="stackedtable-header">
      <div>
        <h4>Orders</h4>
      </div>
      <div>
        <DocsButton label="View orders" onClick={action('would open orders document viewer')} />
      </div>
    </div>
    <table className="table--stacked">
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
  <div>
    <div className="stackedtable-header">
      <div>
        <h4>Orders</h4>
      </div>
      <div>
        <DocsButton label="View orders" onClick={action('would open orders document viewer')} />
      </div>
    </div>
    <table className="table--stacked table--stacked-wbuttons">
      <colgroup>
        <col style={{ width: '25%' }} />
        <col style={{ width: '75%' }} />
      </colgroup>
      <tbody>
        <StackedTableRowForm
          initialValues={{ ordersNumber: '99999999' }}
          validationSchema={Yup.object({
            ordersNumber: Yup.string().max(15, 'Must be 15 characters or less').required('Required'),
          })}
          onSubmit={action('Orders Number Submit')}
          onReset={action('Orders Number Cancel')}
          id="ordersNumber"
          name="ordersNumber"
          type="text"
          label="Orders Number"
        />
        <StackedTableRowForm
          initialValues={{ madeUpField: '' }}
          validationSchema={Yup.object({
            madeUpField: Yup.string().max(15, 'Must be 15 characters or less').required('Required'),
          })}
          onSubmit={action('Made Up Field Submit')}
          onReset={action('Made Up Field Cancel')}
          id="madeUpField"
          name="madeUpField"
          type="text"
          label="Made Up Field"
        />
        <StackedTableRowForm
          initialValues={{ madeUpField2: 'Other data' }}
          validationSchema={Yup.object({
            madeUpField2: Yup.string().max(15, 'Must be 15 characters or less').required('Required'),
          })}
          onSubmit={action('Made Up Field 2 Submit')}
          onReset={action('Made Up Field 2 Cancel')}
          id="madeUpField2"
          name="madeUpField2"
          type="text"
          label="Made Up Field 2"
        />
        <StackedTableRowForm
          initialValues={{ madeUpField3: 'More Data' }}
          validationSchema={Yup.object({
            madeUpField3: Yup.string().max(15, 'Must be 15 characters or less').required('Required'),
          })}
          onSubmit={action('Made Up Field 3 Submit')}
          onReset={action('Made Up Field 3 Cancel')}
          id="madeUpField3"
          name="madeUpField3"
          type="text"
          label="Made Up Field 3"
        />
      </tbody>
    </table>
  </div>
);
const StackedTableWithSomeButtons = () => (
  <div>
    <div className="stackedtable-header">
      <div>
        <h4>Title Goes here</h4>
      </div>
    </div>
    <table className="table--stacked table--stacked-wbuttons">
      <colgroup>
        <col style={{ width: '25%' }} />
        <col style={{ width: '75%' }} />
      </colgroup>
      <tbody>
        <tr>
          <th scope="row">Made Up Read Only Row</th>
          <td>999999999</td>
        </tr>
        <tr>
          <th scope="row">Made Up Read Only Row</th>
          <td>999999999</td>
        </tr>
        <StackedTableRowForm
          initialValues={{ madeUpField: '' }}
          validationSchema={Yup.object({
            madeUpField: Yup.string().max(15, 'Must be 15 characters or less').required('Required'),
          })}
          onSubmit={action('Made Up Field Submit')}
          onReset={action('Made Up Field Cancel')}
          id="madeUpField"
          name="madeUpField"
          type="text"
          label="Made Up Field"
        />
        <StackedTableRowForm
          initialValues={{ madeUpField2: 'Other data' }}
          validationSchema={Yup.object({
            madeUpField2: Yup.string().max(15, 'Must be 15 characters or less').required('Required'),
          })}
          onSubmit={action('Made Up Field 2 Submit')}
          onReset={action('Made Up Field 2 Cancel')}
          id="madeUpField2"
          name="madeUpField2"
          type="text"
          label="Made Up Field 2"
        />
        <tr>
          <th scope="row">Made Up Read Only Row</th>
          <td>999999999</td>
        </tr>
        <tr>
          <th scope="row">Made Up Read Only Row</th>
          <td>999999999</td>
        </tr>
        <tr>
          <th scope="row">Made Up Read Only Row</th>
          <td>999999999</td>
        </tr>
        <StackedTableRowForm
          initialValues={{ madeUpField3: 'More Data' }}
          validationSchema={Yup.object({
            madeUpField3: Yup.string().max(15, 'Must be 15 characters or less').required('Required'),
          })}
          onSubmit={action('Made Up Field 3 Submit')}
          onReset={action('Made Up Field 3 Cancel')}
          id="madeUpField3"
          name="madeUpField3"
          type="text"
          label="Made Up Field 3"
        />
        <tr>
          <th scope="row">Made Up Read Only Row</th>
          <td>999999999</td>
        </tr>
        <StackedTableRowForm
          initialValues={{ madeUpField4: 'More Data' }}
          validationSchema={Yup.object({
            madeUpField4: Yup.string().max(15, 'Must be 15 characters or less').required('Required'),
          })}
          onSubmit={action('Made Up Field 4 Submit')}
          onReset={action('Made Up Field 4 Cancel')}
          id="madeUpField4"
          name="madeUpField4"
          type="text"
          label="Made Up Field 4"
        />
        <tr>
          <th scope="row">Made Up Read Only Row</th>
          <td>999999999</td>
        </tr>
      </tbody>
    </table>
  </div>
);

export const Default = () => (
  <div style={{ padding: `20px`, background: `#f0f0f0` }}>
    <StackedTableExample />
  </div>
);

export const withButtonsToEdit = () => (
  <div style={{ padding: `20px`, background: `#f0f0f0` }}>
    <StackedTableWithButtons />
  </div>
);

export const withButtonsSomeToEdit = () => (
  <div style={{ padding: `20px`, background: `#f0f0f0` }}>
    <StackedTableWithSomeButtons />
  </div>
);
