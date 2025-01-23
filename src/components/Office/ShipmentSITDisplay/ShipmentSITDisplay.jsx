import React, { useEffect, useState } from 'react';
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
import { formatDateForDatePicker, swaggerDateFormat } from 'shared/dates';
import { ShipmentShape } from 'types/shipment';
import { SitStatusShape, LOCATION_TYPES } from 'types/sitStatusShape';
import { DEFAULT_EMPTY_VALUE } from 'shared/constants';

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

const SitStatusTables = ({ shipment, sitExtensions, sitStatus, openModalButton, openConvertModalButton }) => {
  const [isConvertedToCustomerExpense, setIsConvertedToCustomerExpense] = useState(false);
  // Descending sort of past SIT service item groups by SIT Departure Date
  const sortedPastSITGroups = sitStatus.pastSITServiceItemGroupings?.sort(
    (a, b) => new Date(b.summary.sitDepartureDate) - new Date(a.summary.sitDepartureDate),
  );

  // Get the most recent past SIT group if Current SIT doesn't exist
  // This allows the support of the following PO requirement:
  // - "If a new SIT hasn't replaced the old one, then the old SIT's departure date, requested delivery date, and contact date should still show on the dashboard"
  // TODO: Eventually refactor this out when we receive the feature request to list multiple SITs at the same time
  const mostRecentPastSITGroup = sortedPastSITGroups?.[0];

  const pendingSITExtension = sitExtensions.find((se) => se.status === SIT_EXTENSION_STATUS.PENDING);
  const currentDaysInSIT = sitStatus.currentSIT?.daysInSIT || 0;

  const sitDepartureDate = sitStatus.currentSIT
    ? formatDate(sitStatus.currentSIT.sitDepartureDate, swaggerDateFormat, 'DD MMM YYYY') || DEFAULT_EMPTY_VALUE
    : formatDate(mostRecentPastSITGroup?.summary.sitDepartureDate, swaggerDateFormat, 'DD MMM YYYY') ||
      DEFAULT_EMPTY_VALUE;

  const currentDaysInSITElement = <p>{currentDaysInSIT}</p>;
  let sitEntryDate = sitStatus.currentSIT?.sitEntryDate;
  if (!sitEntryDate) {
    sitEntryDate = shipment.mtoServiceItems?.reduce((item, acc) => {
      // Check if the current
      if (item.sitEntryDate < acc.sitEntryDate) {
        return item;
      }
      return acc;
    }).sitEntryDate;
  }

  sitEntryDate = moment(sitEntryDate, swaggerDateFormat);
  const sitStartDateElement = <p>{formatDate(sitEntryDate, swaggerDateFormat, 'DD MMM YYYY')}</p>;
  const sitEndDate =
    formatDateForDatePicker(moment(sitStatus.currentSIT?.sitAuthorizedEndDate, swaggerDateFormat)) || '\u2014';

  // Previous SIT calculations and date ranges
  const previousDaysUsed = sitStatus.pastSITServiceItemGroupings?.map((sitGroup) => {
    // Build the past SIT text based off the past sit group summary rather than individual service items
    // The server provides sitDaysUsed
    const sitDaysUsed = sitGroup.summary.daysInSIT || DEFAULT_EMPTY_VALUE;
    const location = sitGroup.summary.location === LOCATION_TYPES.ORIGIN ? 'origin' : 'destination';

    // Display the dates the server used to calculate sitDaysUsed
    const start = formatDate(sitGroup.summary.sitEntryDate, swaggerDateFormat, 'DD MMM YYYY');
    const end = formatDate(sitGroup.summary.sitDepartureDate, swaggerDateFormat, 'DD MMM YYYY');
    const authorizedEndDate = sitGroup.summary.sitAuthorizedEndDate
      ? formatDate(sitGroup.summary.sitAuthorizedEndDate, swaggerDateFormat, 'DD MMM YYYY')
      : DEFAULT_EMPTY_VALUE;

    const text = `${sitDaysUsed} days at ${location} (${start} - ${end}),\nAuthorized End Date: ${authorizedEndDate}`;

    return <p key={sitGroup.summary.firstDaySITServiceItemID}>{text}</p>;
  });

  // Currently active SIT
  const currentLocation =
    sitStatus.currentSIT?.location === LOCATION_TYPES.DESTINATION ? 'destination SIT' : 'origin SIT';
  const totalSITDaysUsed = sitStatus.totalSITDaysUsed || 0;
  const daysRemaining = sitStatus ? sitStatus.totalDaysRemaining : shipment.sitDaysAllowance;
  const totalDaysRemaining = () => {
    if (daysRemaining > 0) {
      return daysRemaining;
    }
    return 'Expired';
  };

  const showConvertToCustomerExpense = daysRemaining <= 30;

  const customerContactDate = sitStatus.currentSIT
    ? formatDate(sitStatus.currentSIT.sitCustomerContacted, swaggerDateFormat, 'DD MMM YYYY') || DEFAULT_EMPTY_VALUE
    : formatDate(mostRecentPastSITGroup?.summary.sitCustomerContacted, swaggerDateFormat, 'DD MMM YYYY') ||
      DEFAULT_EMPTY_VALUE;

  const sitRequestedDelivery = sitStatus.currentSIT
    ? formatDate(sitStatus.currentSIT.sitRequestedDelivery, swaggerDateFormat, 'DD MMM YYYY') || DEFAULT_EMPTY_VALUE
    : formatDate(mostRecentPastSITGroup?.summary.sitRequestedDelivery, swaggerDateFormat, 'DD MMM YYYY') ||
      DEFAULT_EMPTY_VALUE;

  useEffect(() => {
    if (shipment.mtoServiceItems) {
      const itemsArray = Object.values(shipment.mtoServiceItems);
      const currentSIT = itemsArray.find((item) => item.id === sitStatus?.currentSIT?.serviceItemID);
      if (currentSIT?.convertToCustomerExpense) setIsConvertedToCustomerExpense(true);
      else setIsConvertedToCustomerExpense(false);
    }
  }, [shipment.mtoServiceItems, sitStatus?.currentSIT?.serviceItemID]);

  return (
    <>
      <div className={styles.title}>
        <p>
          SIT (STORAGE IN TRANSIT){pendingSITExtension && <Tag>SIT EXTENSION REQUESTED</Tag>}
          {!pendingSITExtension && isConvertedToCustomerExpense && <Tag>Converted To Customer Expense</Tag>}
        </p>

        {sitStatus.currentSIT &&
          !pendingSITExtension &&
          showConvertToCustomerExpense &&
          !isConvertedToCustomerExpense &&
          openConvertModalButton}
        {sitStatus.currentSIT && openModalButton}
      </div>
      <div className={styles.tableContainer} data-testid="sitStatusTable">
        {/* Sit Total days table */}
        <DataTable
          columnHeaders={['Total days of SIT approved', 'Total days used', 'Total days remaining']}
          dataRow={[shipment.sitDaysAllowance, totalSITDaysUsed, totalDaysRemaining()]}
        />
      </div>

      {/* Current SIT Info Section */}
      {sitStatus.currentSIT && (
        <>
          <div className={styles.tableContainer} data-testid="sitStartAndEndTable">
            {/* Sit Start and End table with total days at current location */}
            {currentDaysInSIT > 0 && <p className={styles.sitHeader}>Current location: {currentLocation}</p>}
            <DataTable
              columnHeaders={[`SIT start date`, 'SIT authorized end date', `Total days in ${currentLocation}`]}
              dataRow={[sitStartDateElement, sitEndDate, currentDaysInSITElement]}
              custClass={styles.currentLocation}
            />
          </div>
          <div className={styles.tableContainer} data-testid="currentSitDepartureDate">
            {/* Current SIT departure date */}
            <DataTable
              testID="currentSITDepartureDate"
              columnHeaders={[`SIT departure date`]}
              dataRow={[sitDepartureDate]}
            />
          </div>
        </>
      )}

      {/* Past SIT Service Items Info Section */}
      {sitStatus.pastSITServiceItemGroupings && (
        <>
          <div className={styles.tableContainer} data-testid="previouslyUsedSitTable">
            <DataTable columnHeaders={['Previously used SIT']} dataRow={[previousDaysUsed]} />
          </div>
          {!sitStatus.currentSIT && (
            <div className={styles.tableContainer} data-testid="pastSitDepartureDateTable">
              {/*
              SIT departure date row for if there is no current SIT.
              The customer wants the most recent SIT departure date to show as an independent entry
              similar to how Current SIT works.
              */}
              <DataTable
                testID="currentSITDateData"
                columnHeaders={[`SIT departure date`]}
                dataRow={[sitDepartureDate]}
              />
            </div>
          )}
        </>
      )}
      <div className={styles.tableContainer}>
        <p className={styles.sitHeader}>Customer delivery request</p>
        <DataTable
          columnHeaders={['Customer contact date', 'Requested delivery date']}
          dataRow={[customerContactDate, sitRequestedDelivery]}
          custClass={styles.currentLocation}
        />
      </div>
    </>
  );
};

const ShipmentSITDisplay = ({
  sitExtensions,
  sitStatus,
  shipment,
  className,
  openModalButton,
  openConvertModalButton,
}) => {
  const sitHistory = React.useMemo(
    () => sitExtensions.filter((sitItem) => sitItem.status !== SIT_EXTENSION_STATUS.PENDING),
    [sitExtensions],
  );

  return (
    <DataTableWrapper
      className={classnames('maxw-tablet', styles.mtoshipmentSITDisplay, className)}
      testID="sitExtensions"
    >
      <SitStatusTables
        openConvertModalButton={openConvertModalButton}
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
  openConvertModalButton: PropTypes.element,
  openModalButton: PropTypes.element,
  className: PropTypes.string,
};

ShipmentSITDisplay.defaultProps = {
  sitExtensions: [],
  sitStatus: undefined,
  openConvertModalButton: undefined,
  openModalButton: undefined,
  className: '',
};

export default ShipmentSITDisplay;
