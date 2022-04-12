import React from 'react';
import classnames from 'classnames';
import PropTypes from 'prop-types';

import styles from './OfficeDefinitionLists.module.scss';

import { OrdersInfoShape } from 'types/order';
import { formatDate } from 'shared/dates';
import { departmentIndicatorReadable, ordersTypeReadable, ordersTypeDetailReadable } from 'shared/formatters';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import { formatLabelReportByDate } from 'utils/formatters';

const OrdersList = ({ ordersInfo, showMissingWarnings }) => {
  const { ordersType } = ordersInfo;
  const isRetiree = ordersType === 'RETIREMENT';
  const isSeparatee = ordersType === 'SEPARATION';
  const missingText = showMissingWarnings ? 'Missing' : '—';

  const reportDateRowLabel = formatLabelReportByDate(ordersType);

  return (
    <div className={styles.OfficeDefinitionLists}>
      <dl className={descriptionListStyles.descriptionList}>
        <div
          className={classnames(descriptionListStyles.row, {
            [styles.missingInfoError]: showMissingWarnings && !ordersInfo.ordersType,
          })}
        >
          <dt>Orders type</dt>
          <dd data-testid="ordersType">{ordersTypeReadable(ordersInfo.ordersType, missingText)}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Date issued</dt>
          <dd data-testid="issuedDate">{formatDate(ordersInfo.issuedDate, 'DD MMM YYYY')}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt data-testid="reportByDateLabel">{reportDateRowLabel}</dt>
          <dd data-testid="reportByDate">{formatDate(ordersInfo.reportByDate, 'DD MMM YYYY')}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Current duty location</dt>
          <dd data-testid="currentDutyLocation">{ordersInfo.currentDutyLocation?.name}</dd>
        </div>
        <div
          className={classnames(descriptionListStyles.row, {
            [styles.missingInfoError]: !ordersInfo.newDutyLocation?.name,
          })}
        >
          <dt data-testid="newDutyLocationLabel">
            {isRetiree || isSeparatee ? 'HOR, HOS, or PLEAD' : 'New duty location'}
          </dt>
          <dd data-testid="newDutyLocation">
            {ordersInfo.newDutyLocation?.name ? ordersInfo.newDutyLocation?.name : '-'}
          </dd>
        </div>
        <div
          className={classnames(descriptionListStyles.row, {
            [styles.missingInfoError]: showMissingWarnings && !ordersInfo.departmentIndicator,
          })}
        >
          <dt>Department indicator</dt>
          <dd data-testid="departmentIndicator">
            {departmentIndicatorReadable(ordersInfo.departmentIndicator, missingText)}
          </dd>
        </div>
        <div
          className={classnames(descriptionListStyles.row, {
            [styles.missingInfoError]: showMissingWarnings && !ordersInfo.ordersNumber,
          })}
        >
          <dt>Orders number</dt>
          <dd data-testid="ordersNumber">{!ordersInfo.ordersNumber ? missingText : ordersInfo.ordersNumber}</dd>
        </div>

        <div
          className={classnames(descriptionListStyles.row, {
            [styles.missingInfoError]: showMissingWarnings && !ordersInfo.ordersTypeDetail,
          })}
        >
          <dt>Orders type detail</dt>
          <dd data-testid="ordersTypeDetail">{ordersTypeDetailReadable(ordersInfo.ordersTypeDetail, missingText)}</dd>
        </div>
        <div
          className={classnames(descriptionListStyles.row, {
            [styles.missingInfoError]: showMissingWarnings && !ordersInfo.tacMDC,
          })}
        >
          <dt>HHG TAC</dt>
          <dd data-testid="tacMDC">{!ordersInfo.tacMDC ? missingText : ordersInfo.tacMDC}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>HHG SAC</dt>
          <dd data-testid="sacSDN">{!ordersInfo.sacSDN ? '—' : ordersInfo.sacSDN}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>NTS TAC</dt>
          <dd data-testid="NTStac">{!ordersInfo.NTStac ? '—' : ordersInfo.NTStac}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>NTS SAC</dt>
          <dd data-testid="NTSsac">{!ordersInfo.NTSsac ? '—' : ordersInfo.NTSsac}</dd>
        </div>
      </dl>
    </div>
  );
};

OrdersList.defaultProps = {
  showMissingWarnings: true,
};

OrdersList.propTypes = {
  ordersInfo: OrdersInfoShape.isRequired,
  showMissingWarnings: PropTypes.bool,
};

export default OrdersList;
