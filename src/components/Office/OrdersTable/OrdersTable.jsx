import React from 'react';
import { Link } from 'react-router-dom';
import classnames from 'classnames';

import styles from './OrdersTable.module.scss';

import { OrdersInfoShape } from 'types/order';
import { formatDate } from 'shared/dates';
import { departmentIndicatorReadable, ordersTypeReadable, ordersTypeDetailReadable } from 'shared/formatters';

function OrdersTable({ ordersInfo }) {
  return (
    <div className={styles.OrdersTable}>
      <div className="stackedtable-header">
        <div>
          <h4>Orders</h4>
        </div>
        <div>
          <Link className="usa-button usa-button--secondary" data-testid="edit-orders" to="orders">
            Edit orders
          </Link>
        </div>
      </div>
      <table className="table--stacked">
        <colgroup>
          <col style={{ width: '25%' }} />
          <col style={{ width: '75%' }} />
        </colgroup>
        <tbody>
          <tr>
            <th scope="row">Current duty station</th>
            <td data-testid="currentDutyStation">{ordersInfo.currentDutyStation?.name}</td>
          </tr>
          <tr>
            <th scope="row">New duty station</th>
            <td data-testid="newDutyStation">{ordersInfo.newDutyStation?.name}</td>
          </tr>
          <tr>
            <th scope="row">Date issued</th>
            <td data-testid="issuedDate">{formatDate(ordersInfo.issuedDate, 'DD MMM YYYY')}</td>
          </tr>
          <tr>
            <th scope="row">Report by date</th>
            <td data-testid="reportByDate">{formatDate(ordersInfo.reportByDate, 'DD MMM YYYY')}</td>
          </tr>
          <tr className={classnames({ [styles.missingInfoError]: !ordersInfo.departmentIndicator })}>
            <th scope="row">Department indicator</th>
            <td data-testid="departmentIndicator">{departmentIndicatorReadable(ordersInfo.departmentIndicator)}</td>
          </tr>
          <tr className={classnames({ [styles.missingInfoError]: !ordersInfo.ordersNumber })}>
            <th scope="row">Orders number</th>
            <td data-testid="ordersNumber">{!ordersInfo.ordersNumber ? 'Missing' : ordersInfo.ordersNumber}</td>
          </tr>
          <tr className={classnames({ [styles.missingInfoError]: !ordersInfo.ordersType })}>
            <th scope="row">Orders type</th>
            <td data-testid="ordersType">{ordersTypeReadable(ordersInfo.ordersType)}</td>
          </tr>
          <tr className={classnames({ [styles.missingInfoError]: !ordersInfo.ordersTypeDetail })}>
            <th scope="row">Orders type detail</th>
            <td data-testid="ordersTypeDetail">{ordersTypeDetailReadable(ordersInfo.ordersTypeDetail)}</td>
          </tr>
          <tr className={classnames({ [styles.missingInfoError]: !ordersInfo.tacMDC })}>
            <th scope="row">TAC / MDC</th>
            <td data-testid="tacMDC">{!ordersInfo.tacMDC ? 'Missing' : ordersInfo.tacMDC}</td>
          </tr>
          <tr>
            <th scope="row">SAC / SDN</th>
            <td data-testid="sacSDN">{ordersInfo.sacSDN}</td>
          </tr>
        </tbody>
      </table>
    </div>
  );
}

OrdersTable.propTypes = {
  ordersInfo: OrdersInfoShape.isRequired,
};

export default OrdersTable;
