import React from 'react';
import classnames from 'classnames';
import moment from 'moment';
import { PropTypes } from 'prop-types';
import { Tag } from '@trussworks/react-uswds';

import DataTableWrapper from '../../DataTableWrapper/index';
import DataTable from '../../DataTable/index';
import { SITExtensionShape } from '../../../types/sitExtensions';

import styles from './ShipmentSITDisplay.module.scss';

import { sitExtensionReasons, SIT_EXTENSION_STATUS } from 'constants/sitExtensions';
import { formatDateFromIso, formatDate } from 'utils/formatters';
import { utcDateFormat } from 'shared/dates';
import { SERVICE_ITEM_CODES } from 'constants/serviceItems';
import { ShipmentShape } from 'types/shipment';
import { SitStatusShape, LOCATION_TYPES } from 'types/sitStatusShape';

const SITHistoryItem = ({ sitItem }) => (
  <dl data-testid="sitHistoryItem">
    <div>
      <dt>Reason:</dt>
      <dd>{sitExtensionReasons[sitItem.requestReason]}</dd>
    </div>
    {sitItem.contractorRemarks && (
      <div>
        <dt>Contractor remarks:</dt>
        <dd>{sitItem.contractorRemarks}</dd>
      </div>
    )}
    {sitItem.officeRemarks && (
      <div>
        <dt>Office remarks:</dt>
        <dd>{sitItem.officeRemarks}</dd>
      </div>
    )}
  </dl>
);

const SITHistoryItemHeader = ({ sitItem }) => (
  <div className={styles.sitHistoryItemHeader}>
    Total days of SIT approved: {sitItem.approvedDays ?? '0'}{' '}
    <span>updated on {formatDateFromIso(sitItem.decisionDate, 'DD MMM YYYY')} </span>
  </div>
);

const SitHistoryList = ({ sitHistory, dayAllowance }) => {
  let approvedDays = dayAllowance;
  return (
    <div className={styles.tableContainer}>
      <p className={styles.sitHeader}>SIT history</p>
      {sitHistory.map((currentItem) => {
        const sitItem = {
          ...currentItem,
          approvedDays,
        };
        approvedDays -= currentItem.approvedDays;
        return (
          <DataTable
            key={sitItem.id}
            columnHeaders={[<SITHistoryItemHeader sitItem={sitItem} />]}
            dataRow={[<SITHistoryItem sitItem={sitItem} />]}
            custClass={styles.sitHistoryItem}
          />
        );
      })}
    </div>
  );
};

const SitStatusTables = ({ shipment, sitExtensions, sitStatus, openModalButton }) => {
  const pendingSITExtension = sitExtensions.find((se) => se.status === SIT_EXTENSION_STATUS.PENDING);

  const currentDaysInSIT = sitStatus?.totalSITDaysUsed || 0;
  const currentDaysInSITElement = <p>{currentDaysInSIT}</p>;
  let sitStartDate = sitStatus?.sitEntryDate;
  if (!sitStartDate) {
    sitStartDate = shipment.mtoServiceItems?.reduce((item, acc) => {
      if (item.sitEntryDate < acc.sitEntryDate) {
        return item;
      }
      return acc;
    }).sitEntryDate;
  }
  const sitStartDateElement = <p>{formatDate(sitStartDate, utcDateFormat, 'DD MMM YYYY')}</p>;

  // Subract one day from total days remaining to account for the current day
  const sitEndDate = sitStatus
    ? moment()
        .add(sitStatus.totalDaysRemaining - 1, 'days')
        .format('DD MMM YYYY')
    : moment(sitStartDate).add(shipment.sitDaysAllowance, 'days').format('DD MMM YYYY');

  // Previous SIT calculations and date ranges
  const previousDaysUsed = sitStatus?.pastSITServiceItems?.map((pastSITItem) => {
    const sitDaysUsed = moment(pastSITItem.sitDepartureDate).diff(pastSITItem.sitEntryDate, 'days');
    const location = pastSITItem.reServiceCode === SERVICE_ITEM_CODES.DOPSIT ? 'origin' : 'destination';

    const start = formatDate(pastSITItem.sitEntryDate, utcDateFormat, 'DD MMM YYYY');
    const end = formatDate(pastSITItem.sitDepartureDate, utcDateFormat, 'DD MMM YYYY');
    const text = `${sitDaysUsed} days at ${location} (${start} - ${end})`;

    return <p key={pastSITItem.id}>{text}</p>;
  });

  // Currently active SIT
  const currentLocation = sitStatus?.location === LOCATION_TYPES.DESTINATION ? 'destination SIT' : 'origin SIT';

  const totalSITDaysUsed = sitStatus?.totalSITDaysUsed || 0;
  const totalDaysRemaining = () => {
    const daysRemaining = sitStatus ? sitStatus.totalDaysRemaining : shipment.sitDaysAllowance;
    if (!sitStatus && daysRemaining > 0) {
      return daysRemaining;
    }
    if (sitStatus && daysRemaining > 0) {
      // Subract one day from the remaining days on the current sit to account for the current day
      return daysRemaining - 1;
    }
    return 'Expired';
  };

  return (
    <>
      <div className={styles.title}>
        <p>SIT (STORAGE IN TRANSIT){pendingSITExtension && <Tag>Extension requested</Tag>}</p>
        {openModalButton}
      </div>
      <div className={styles.tableContainer} data-testid="sitStatusTable">
        {/* Sit Total days table */}
        <DataTable
          columnHeaders={['Total days of SIT approved', 'Total days used', 'Total days remaining']}
          dataRow={[shipment.sitDaysAllowance, totalSITDaysUsed, totalDaysRemaining()]}
        />
      </div>

      <div className={styles.tableContainer} data-testid="sitStartAndEndTable">
        {/* Sit Start and End table */}
        {currentDaysInSIT > 0 && <p className={styles.sitHeader}>Current location: {currentLocation}</p>}
        <DataTable
          columnHeaders={[`SIT start date`, 'SIT authorized end date']}
          dataRow={[sitStartDateElement, sitEndDate]}
          custClass={styles.currentLocation}
        />
      </div>
      <div className={styles.tableContainer} data-testid="sitDaysAtCurrentLocation">
        {/* Total days at current location */}
        <DataTable columnHeaders={[`Total days in ${currentLocation}`]} dataRow={[currentDaysInSITElement]} />
      </div>
      {/* Service Items */}
      {sitStatus?.pastSITServiceItems && (
        <div className={styles.tableContainer}>
          <DataTable columnHeaders={['Previously used SIT']} dataRow={[previousDaysUsed]} />
        </div>
      )}
    </>
  );
};

const ShipmentSITDisplay = ({ sitExtensions, sitStatus, shipment, className, openModalButton }) => {
  const sitHistory = React.useMemo(
    () => sitExtensions.reverse().filter((sitItem) => sitItem.status !== SIT_EXTENSION_STATUS.PENDING),
    [sitExtensions],
  );

  return (
    <DataTableWrapper
      className={classnames('maxw-tablet', styles.mtoshipmentSITDisplay, className)}
      testID="sitExtensions"
    >
      <SitStatusTables
        openModalButton={openModalButton}
        shipment={shipment}
        sitStatus={sitStatus}
        sitExtensions={sitExtensions}
      />
      {/* Sit History */}
      {sitExtensions && sitHistory.length > 0 && (
        <SitHistoryList sitHistory={sitHistory} dayAllowance={shipment.sitDaysAllowance} />
      )}
    </DataTableWrapper>
  );
};

ShipmentSITDisplay.propTypes = {
  sitExtensions: PropTypes.arrayOf(SITExtensionShape),
  sitStatus: SitStatusShape,
  shipment: ShipmentShape.isRequired,
  openModalButton: PropTypes.element,
  className: PropTypes.string,
};

ShipmentSITDisplay.defaultProps = {
  sitExtensions: [],
  sitStatus: null,
  openModalButton: undefined,
  className: '',
};

export default ShipmentSITDisplay;
