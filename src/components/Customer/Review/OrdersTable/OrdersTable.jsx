import React from 'react';
import { string, arrayOf, bool, shape, func } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import reviewStyles from '../Review.module.scss';

import { formatCustomerDate, formatLabelReportByDate } from 'utils/formatters';
import { ORDERS_TYPE_OPTIONS } from 'constants/orders';

const OrdersTable = ({
  hasDependents,
  issueDate,
  moveId,
  newDutyLocationName,
  onEditClick,
  orderType,
  reportByDate,
  uploads,
  originDutyLocationName,
  payGrade,
  orderId,
  counselingOfficeName,
  accompaniedTour,
  dependentsUnderTwelve,
  dependentsTwelveAndOver,
}) => {
  const isRetirementOrSeparation = ['RETIREMENT', 'SEPARATION'].includes(orderType);
  const editPath = `/move/${moveId}/review/edit-orders/${orderId}`;

  return (
    <div className={reviewStyles['review-container']}>
      <div className={reviewStyles['review-header']}>
        <h2>Orders</h2>
        <Button
          unstyled
          className={reviewStyles['edit-btn']}
          data-testid="edit-orders-table"
          onClick={() => onEditClick(editPath)}
        >
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
            <td style={{ wordBreak: 'break-word' }}>{ORDERS_TYPE_OPTIONS[orderType] || orderType}</td>
          </tr>
          <tr>
            <th scope="row">Orders date</th>
            <td>{formatCustomerDate(issueDate)}</td>
          </tr>
          <tr>
            <th scope="row">{formatLabelReportByDate(orderType)}</th>
            <td>{formatCustomerDate(reportByDate)}</td>
          </tr>
          <tr>
            <th className={reviewStyles['table-divider-top']} scope="row" style={{ borderBottom: 'none' }}>
              Current duty location
            </th>
            <td className={reviewStyles['table-divider-top']} style={{ borderBottom: 'none' }}>
              {originDutyLocationName}
            </td>
          </tr>
          <tr>
            <th className={reviewStyles['table-divider-top']} scope="row" style={{ borderBottom: 'none' }}>
              Counseling office
            </th>
            <td className={reviewStyles['table-divider-top']} style={{ borderBottom: 'none' }}>
              {counselingOfficeName}
            </td>
          </tr>
          <tr>
            <th scope="row">{isRetirementOrSeparation ? 'HOR, PLEAD or HOS' : 'New duty location'}</th>
            <td>{newDutyLocationName}</td>
          </tr>
          <tr>
            <th scope="row">Dependents</th>
            <td>{hasDependents ? 'Yes' : 'No'}</td>
          </tr>
          {/* Group conditionally rendered OCONUS fields */}
          {(accompaniedTour || dependentsUnderTwelve > 0 || dependentsTwelveAndOver > 0) && (
            <>
              <tr>
                <th scope="row">Accompanied tour</th>
                <td>{accompaniedTour ? 'Yes' : 'No'}</td>
              </tr>
              <tr>
                <th scope="row">Dependents under twelve</th>
                <td>{dependentsUnderTwelve || 0}</td>
              </tr>
              <tr>
                <th scope="row">Dependents twelve and over</th>
                <td>{dependentsTwelveAndOver || 0}</td>
              </tr>
            </>
          )}
          {/* End grouping of UB fields */}
          <tr>
            <th scope="row">Orders</th>
            <td>
              {uploads.length} file{uploads.length > 1 ? 's' : ''}
            </td>
          </tr>
          <tr>
            <th scope="row">Pay grade</th>
            <td>{payGrade}</td>
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
  newDutyLocationName: string.isRequired,
  onEditClick: func.isRequired,
  orderType: string.isRequired,
  reportByDate: string.isRequired,
  uploads: arrayOf(shape({})).isRequired,
  payGrade: string.isRequired,
  originDutyLocationName: string.isRequired,
};

export default OrdersTable;
