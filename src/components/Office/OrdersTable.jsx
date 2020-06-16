import React from 'react';
import * as PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

function OrdersTable({ ordersInfo, handleViewEditOrdersClick }) {
  return (
    <div>
      <div className="stackedtable-header">
        <div>
          <h4>Orders</h4>
        </div>
        <div>
          <Button onClick={handleViewEditOrdersClick} secondary>
            <span>View & edit orders</span>
          </Button>
        </div>
      </div>
      <table className="table--stacked">
        <colgroup>
          <col style={{ width: '25%' }} />
          <col style={{ width: '75%' }} />
        </colgroup>
        <tbody>
          <tr>
            <th scope="row">Current duty Station</th>
            <td data-cy="currentDutyStation">{ordersInfo.currentDutyStation}</td>
          </tr>
          <tr>
            <th scope="row">New duty station</th>
            <td data-cy="newDutyStation">{ordersInfo.newDutyStation}</td>
          </tr>
          <tr>
            <th scope="row">Date issued</th>
            <td data-cy="issuedDate">{ordersInfo.issuedDate}</td>
          </tr>
          <tr>
            <th scope="row">Report by date</th>
            <td data-cy="reportByDate">{ordersInfo.reportByDate}</td>
          </tr>
          <tr>
            <th scope="row">Department indicator</th>
            <td data-cy="departmentIndicator">{ordersInfo.departmentIndicator}</td>
          </tr>
          <tr>
            <th scope="row">Orders number</th>
            <td data-cy="ordersNumber">{ordersInfo.ordersNumber}</td>
          </tr>
          <tr>
            <th scope="row">Orders type</th>
            <td data-cy="ordersType">{ordersInfo.ordersType}</td>
          </tr>
          <tr>
            <th scope="row">Orders type detail</th>
            <td data-cy="ordersTypeDetail">{ordersInfo.ordersTypeDetail}</td>
          </tr>
          <tr>
            <th scope="row">TAC / MDC</th>
            <td data-cy="tacMDC">{ordersInfo.tacMDC}</td>
          </tr>
          <tr>
            <th scope="row">SAC / SDN</th>
            <td data-cy="sacSDN">{ordersInfo.sacSDN}</td>
          </tr>
        </tbody>
      </table>
    </div>
  );
}

OrdersTable.propTypes = {
  ordersInfo: PropTypes.shape({
    currentDutyStation: PropTypes.string,
    newDutyStation: PropTypes.string,
    issuedDate: PropTypes.string,
    reportByDate: PropTypes.string,
    departmentIndicator: PropTypes.string,
    ordersNumber: PropTypes.string,
    ordersType: PropTypes.string,
    ordersTypeDetail: PropTypes.string,
    tacMDC: PropTypes.string,
    sacSDN: PropTypes.string,
  }).isRequired,
  handleViewEditOrdersClick: PropTypes.func,
};

export default OrdersTable;
