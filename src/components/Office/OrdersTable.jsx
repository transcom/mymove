import React from 'react';
import * as PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

function OrdersTable({ ordersInfo }) {
  return (
    <div>
      <div className="stackedtable-header">
        <div>
          <h4>Orders</h4>
        </div>
        <div>
          <Button secondary>
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
            <th scope="row" className="text-bold">
              Current duty Station
            </th>
            <td data-cy="currentDutyStation">{ordersInfo.currentDutyStation}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              New duty station
            </th>
            <td data-cy="newDutyStation">{ordersInfo.newDutyStation}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Date issued
            </th>
            <td data-cy="issuedDate">{ordersInfo.issuedDate}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Report by date
            </th>
            <td data-cy="reportByDate">{ordersInfo.reportByDate}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Department indicator
            </th>
            <td data-cy="departmentIndicator">{ordersInfo.departmentIndicator}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Orders number
            </th>
            <td data-cy="ordersNumber">{ordersInfo.ordersNumber}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Orders type
            </th>
            <td data-cy="ordersType">{ordersInfo.ordersType}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Orders type detail
            </th>
            <td data-cy="ordersTypeDetail">{ordersInfo.ordersTypeDetail}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              TAC / MDC
            </th>
            <td data-cy="tacMDC">{ordersInfo.tacMDC}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              SAC / SDN
            </th>
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
};

export default OrdersTable;
