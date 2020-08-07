import React from 'react';
import { Button } from '@trussworks/react-uswds';

import { OrdersInfoShape } from '../../types/moveOrder';

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
            <td data-testid="currentDutyStation">{ordersInfo.currentDutyStation?.name}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              New duty station
            </th>
            <td data-testid="newDutyStation">{ordersInfo.newDutyStation?.name}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Date issued
            </th>
            <td data-testid="issuedDate">{ordersInfo.issuedDate}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Report by date
            </th>
            <td data-testid="reportByDate">{ordersInfo.reportByDate}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Department indicator
            </th>
            <td data-testid="departmentIndicator">{ordersInfo.departmentIndicator}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Orders number
            </th>
            <td data-testid="ordersNumber">{ordersInfo.ordersNumber}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Orders type
            </th>
            <td data-testid="ordersType">{ordersInfo.ordersType}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Orders type detail
            </th>
            <td data-testid="ordersTypeDetail">{ordersInfo.ordersTypeDetail}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              TAC / MDC
            </th>
            <td data-testid="tacMDC">{ordersInfo.tacMDC}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              SAC / SDN
            </th>
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
