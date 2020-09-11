/* eslint-ignore */
import React from 'react';
import { shape, string } from 'prop-types';
import { Link } from 'react-router-dom';

import reviewStyles from '../Review.module.scss';

const OrdersTable = ({ orders }) => (
  <div className="review-container">
    <div className="stackedtable-header">
      <h3>Orders</h3>
      <Link>Edit</Link>
    </div>
    <table className={`table--stacked ${reviewStyles['review-table']}`}>
      <colgroup>
        <col style={{ width: '40%' }} />
        <col style={{ width: '60%' }} />
      </colgroup>
      <tbody>
        <tr>
          <th scope="row">Orders type</th>
          <td style={{ wordBreak: 'break-word' }}>{orders.orders_type}</td>
        </tr>
        <tr>
          <th scope="row">Orders date</th>
          <td>{orders.issue_date}</td>
        </tr>
        <tr>
          <th scope="row">Report by date</th>
          <td>{orders.report_by_date}</td>
        </tr>
        <tr>
          <th scope="row">New duty station</th>
          <td>{orders.new_duty_station.name}</td>
        </tr>
        <tr>
          <th scope="row">Dependents</th>
          <td>{orders.has_dependents ? 'Yes' : 'No'}</td>
        </tr>
        <tr>
          <th scope="row">Orders</th>
          <td>
            {orders.uploaded_orders.uploads.length} file{orders.uploaded_orders.uploads.length > 1 ? 's' : ''}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
);

OrdersTable.propTypes = {
  orders: shape({}).isRequired,
};

export default OrdersTable;
