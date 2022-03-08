/* eslint-ignore */
import React from 'react';
import { string, arrayOf, bool, shape, func } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import reviewStyles from '../Review.module.scss';

import { formatOrderType, formatLabelReportByDate } from 'utils/formatters';

const OrdersTable = ({
  hasDependents,
  issueDate,
  moveId,
  newDutyStationName,
  onEditClick,
  orderType,
  reportByDate,
  uploads,
}) => {
  const editPath = `/moves/${moveId}/review/edit-orders`;
  return (
    <div className={reviewStyles['review-container']}>
      <div className={reviewStyles['review-header']}>
        <h2>Orders</h2>
        <Button unstyled className={reviewStyles['edit-btn']} onClick={() => onEditClick(editPath)}>
          Edit
        </Button>
      </div>
      <table className={`table--stacked ${reviewStyles['review-table']}`}>
        <colgroup>
          <col />
          <col />
        </colgroup>
        <tbody>
          <tr>
            <th scope="row">Orders type</th>
            <td style={{ wordBreak: 'break-word' }}>{formatOrderType(orderType)}</td>
          </tr>
          <tr>
            <th scope="row">Orders date</th>
            <td>{issueDate}</td>
          </tr>
          <tr>
            <th scope="row">{formatLabelReportByDate(orderType)}</th>
            <td>{reportByDate}</td>
          </tr>
          <tr>
            <th scope="row">New duty location</th>
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
};

OrdersTable.propTypes = {
  hasDependents: bool.isRequired,
  issueDate: string.isRequired,
  moveId: string.isRequired,
  newDutyStationName: string.isRequired,
  onEditClick: func.isRequired,
  orderType: string.isRequired,
  reportByDate: string.isRequired,
  uploads: arrayOf(shape({})).isRequired,
};

export default OrdersTable;
