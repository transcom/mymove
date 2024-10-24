import React from 'react';
import { string } from 'prop-types';

import styles from './index.module.scss';

import { OrderShape, CustomerShape } from 'types/order';
import { formatCustomerDate, formatLabelReportByDate } from 'utils/formatters';
import {
  CHECK_SPECIAL_ORDERS_TYPES,
  ORDERS_BRANCH_OPTIONS,
  ORDERS_PAY_GRADE_OPTIONS,
  SPECIAL_ORDERS_TYPES,
} from 'constants/orders.js';
import SERVICE_MEMBER_AGENCIES from 'content/serviceMemberAgencies';
import MOVE_STATUSES from 'constants/moves';
import { roleTypes } from 'constants/userRoles';
import departmentIndicators from 'constants/departmentIndicators';

const CustomerHeader = ({ customer, order, moveCode, move, userRole }) => {
  const isCoastGuard = customer.agency === departmentIndicators.COAST_GUARD;
  // eslint-disable-next-line camelcase
  const { order_type: orderType } = order;

  const isRetireeOrSeparatee = ['RETIREMENT', 'SEPARATION'].includes(orderType);
  const isSpecialMove = CHECK_SPECIAL_ORDERS_TYPES(orderType);

  /**
   * Depending on the order type, this row dt label can be either:
   * Report by date (PERMANENT_CHANGE_OF_STATION)
   * Date of retirement (RETIREMENT)
   * Date of separation (SEPARATION)
   */
  const reportDateLabel = formatLabelReportByDate(orderType);
  // This logic to show different originGLBOC is based on queue table's backend logic
  const originGBLOC =
    move?.status === MOVE_STATUSES.NEEDS_SERVICE_COUNSELING ||
    userRole === roleTypes.SERVICES_COUNSELOR ||
    !move?.shipmentGBLOC
      ? order.originDutyLocationGBLOC
      : move.shipmentGBLOC;
  const originGBLOCDisplay =
    order.agency === SERVICE_MEMBER_AGENCIES.MARINES ? `${order.originDutyLocationGBLOC} / USMC` : originGBLOC;

  return (
    <div className={styles.custHeader}>
      <div>
        <div data-testid="nameBlock" className={styles.nameBlock}>
          <h2>
            {customer.last_name}, {customer.first_name}
          </h2>
          <span className="usa-tag usa-tag--cyan usa-tag--large">#{moveCode}</span>
        </div>
        <div>
          <p>
            <span data-testid="deptPayGrade" className={styles.details}>
              {ORDERS_BRANCH_OPTIONS[`${order.agency}`]} {ORDERS_PAY_GRADE_OPTIONS[`${order.grade}`]}
            </span>
            <span className={styles.verticalBar}>|</span>
            <span data-testid="edipi" className={styles.details}>
              DoD ID {customer.edipi}
            </span>
            {isCoastGuard && (
              <>
                <span className={styles.verticalBar}>|</span>
                <span data-testid="emplid" className={styles.details}>
                  EMPLID {customer.emplid}
                </span>
              </>
            )}
          </p>
        </div>
      </div>
      {isSpecialMove ? (
        <div data-testid="specialMovesLabel" className={styles.specialMovesLabel}>
          <p>{SPECIAL_ORDERS_TYPES[`${orderType}`]}</p>
        </div>
      ) : null}
      <div data-testid="infoBlock" className={styles.infoBlock}>
        <div>
          <p>Authorized origin</p>
          <h4>{order.originDutyLocation.name}</h4>
        </div>
        {order.destinationDutyLocation.name && (
          <div>
            <p data-testid="destinationLabel">
              {isRetireeOrSeparatee ? 'HOR, HOS or PLEAD' : 'Authorized destination'}
            </p>
            <h4>{order.destinationDutyLocation.name}</h4>
          </div>
        )}
        <div>
          <p data-testid="reportDateLabel">{reportDateLabel}</p>
          <h4>{formatCustomerDate(order.report_by_date)}</h4>
        </div>
        <div>
          <p data-testid="originGBLOC">Origin GBLOC</p>
          <h4>{originGBLOCDisplay}</h4>
        </div>
      </div>
    </div>
  );
};

CustomerHeader.propTypes = {
  customer: CustomerShape.isRequired,
  order: OrderShape.isRequired,
  moveCode: string.isRequired,
};

export default CustomerHeader;
