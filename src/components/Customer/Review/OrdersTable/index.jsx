/* eslint-ignore */
import React from 'react';
import { string, arrayOf, bool, shape } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import reviewStyles from '../Review.module.scss';

const OrdersTable = ({ orderType, issueDate, reportByDate, newDutyStationName, hasDependents, uploads }) => (
  <div className={reviewStyles['review-container']}>
    <div className={reviewStyles['review-header']}>
      <h3>Orders</h3>
      <Button unstyled className={reviewStyles['edit-btn']}>
        Edit
      </Button>
    </div>
    <table className={`table--stacked ${reviewStyles['review-table']}`}>
      <colgroup>
        <col style={{ width: '40%' }} />
        <col style={{ width: '60%' }} />
      </colgroup>
      <tbody>
        <tr>
          <th scope="row">Orders type</th>
          <td style={{ wordBreak: 'break-word' }}>{orderType}</td>
        </tr>
        <tr>
          <th scope="row">Orders date</th>
          <td>{issueDate}</td>
        </tr>
        <tr>
          <th scope="row">Report by date</th>
          <td>{reportByDate}</td>
        </tr>
        <tr>
          <th scope="row">New duty station</th>
          <td>{newDutyStationName}</td>
        </tr>
        <tr>
          <th scope="row">Dependents</th>
          <td>{hasDependents ? 'Yes' : 'No'}</td>
        </tr>
        <tr>
          <th scope="row">Orders</th>
          <td>
            {uploads.length} file{uploads.length > 1 ? 's' : ''}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
);

OrdersTable.propTypes = {
  orderType: string.isRequired,
  issueDate: string.isRequired,
  reportByDate: string.isRequired,
  newDutyStationName: string.isRequired,
  hasDependents: bool.isRequired,
  uploads: arrayOf(shape({})).isRequired,
};

export default OrdersTable;
