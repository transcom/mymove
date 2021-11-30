import React from 'react';
import classnames from 'classnames';

import styles from './OfficeDefinitionLists.module.scss';

import { OrdersInfoShape } from 'types/order';
import { formatDate } from 'shared/dates';
import { departmentIndicatorReadable, ordersTypeReadable, ordersTypeDetailReadable } from 'shared/formatters';
import descriptionListStyles from 'styles/descriptionList.module.scss';

const OrdersList = ({ ordersInfo }) => {
  return (
    <div className={styles.OfficeDefinitionLists}>
      <dl className={descriptionListStyles.descriptionList}>
        <div className={descriptionListStyles.row}>
          <dt>Current duty location</dt>
          <dd data-testid="currentDutyLocation">{ordersInfo.currentDutyStation?.name}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>New duty location</dt>
          <dd data-testid="newDutyLocation">{ordersInfo.newDutyStation?.name}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Date issued</dt>
          <dd data-testid="issuedDate">{formatDate(ordersInfo.issuedDate, 'DD MMM YYYY')}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Report by date</dt>
          <dd data-testid="reportByDate">{formatDate(ordersInfo.reportByDate, 'DD MMM YYYY')}</dd>
        </div>
        <div
          className={classnames(descriptionListStyles.row, {
            [styles.missingInfoError]: !ordersInfo.departmentIndicator,
          })}
        >
          <dt>Department indicator</dt>
          <dd data-testid="departmentIndicator">{departmentIndicatorReadable(ordersInfo.departmentIndicator)}</dd>
        </div>
        <div className={classnames(descriptionListStyles.row, { [styles.missingInfoError]: !ordersInfo.ordersNumber })}>
          <dt>Orders number</dt>
          <dd data-testid="ordersNumber">{!ordersInfo.ordersNumber ? 'Missing' : ordersInfo.ordersNumber}</dd>
        </div>
        <div className={classnames(descriptionListStyles.row, { [styles.missingInfoError]: !ordersInfo.ordersType })}>
          <dt>Orders type</dt>
          <dd data-testid="ordersType">{ordersTypeReadable(ordersInfo.ordersType)}</dd>
        </div>
        <div
          className={classnames(descriptionListStyles.row, { [styles.missingInfoError]: !ordersInfo.ordersTypeDetail })}
        >
          <dt>Orders type detail</dt>
          <dd data-testid="ordersTypeDetail">{ordersTypeDetailReadable(ordersInfo.ordersTypeDetail)}</dd>
        </div>
        <div className={classnames(descriptionListStyles.row, { [styles.missingInfoError]: !ordersInfo.tacMDC })}>
          <dt>HHG TAC</dt>
          <dd data-testid="tacMDC">{!ordersInfo.tacMDC ? 'Missing' : ordersInfo.tacMDC}</dd>
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

OrdersList.propTypes = {
  ordersInfo: OrdersInfoShape.isRequired,
};

export default OrdersList;
