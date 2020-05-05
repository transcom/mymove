import React from 'react';
import * as PropTypes from 'prop-types';
import { DocsButton } from '../form';

function OrdersTable({ ordersInfo }) {
  return (
    <div>
      <div className="stackedtable-header">
        <div>
          <h4>Orders</h4>
        </div>
        <div>
          <DocsButton label="View orders" />
        </div>
      </div>
      <table className="table--stacked">
        <colgroup>
          <col style={{ width: '25%' }} />
          <col style={{ width: '75%' }} />
        </colgroup>
        <tbody>
          <tr>
            <th scope="row">Current Duty Station</th>
            <td>{ordersInfo.currentDutyStation ? ordersInfo.currentDutyStation : ''}</td>
          </tr>
          <tr>
            <th scope="row">New duty station</th>
            <td>{ordersInfo.newDutyStation ? ordersInfo.newDutyStation : ''}</td>
          </tr>
          <tr>
            <th scope="row">Date issuedc</th>
            <td>{ordersInfo.issuedDate ? ordersInfo.issuedDate : ''}</td>
          </tr>
          <tr>
            <th scope="row">Report by date</th>
            <td>{ordersInfo.reportByDate ? ordersInfo.reportByDate : ''}</td>
          </tr>
          <tr>
            <th scope="row">Department indicator</th>
            <td>{ordersInfo.departmentIndicator ? ordersInfo.departmentIndicator : ''}</td>
          </tr>
          <tr>
            <th scope="row">Orders number</th>
            <td>{ordersInfo.ordersNumber ? ordersInfo.ordersNumber : ''}</td>
          </tr>
          <tr>
            <th scope="row">Orders type</th>
            <td>{ordersInfo.ordersType ? ordersInfo.ordersType : ''}</td>
          </tr>
          <tr>
            <th scope="row">Orders type detail</th>
            <td>{ordersInfo.ordersTypeDetail ? ordersInfo.ordersTypeDetail : ''}</td>
          </tr>
          <tr>
            <th scope="row">TAC / MDC</th>
            <td>{ordersInfo.tacMDC ? ordersInfo.tacMDC : ''}</td>
          </tr>
          <tr>
            <th scope="row">SAC / SDN</th>
            <td>{ordersInfo.sacSDN ? ordersInfo.sacSDN : ''}</td>
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
};

export default OrdersTable;
