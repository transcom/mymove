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
import { swaggerDateFormat } from 'shared/dates';
import { SERVICE_ITEM_CODES } from 'constants/serviceItems';
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

  const pendingSITExtension = sitExtensions.find((se) => se.status === SIT_EXTENSION_STATUS.PENDING);
  const currentDaysInSIT = sitStatus.currentSIT?.daysInSIT || 0;
  const sitDepartureDate =
    formatDate(sitStatus.currentSIT?.sitDepartureDate, swaggerDateFormat, 'DD MMM YYYY') || DEFAULT_EMPTY_VALUE;
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
    formatDate(sitStatus.currentSIT?.sitAllowanceEndDate, swaggerDateFormat, 'DD MMM YYYY') || '\u2014';

  // Previous SIT calculations and date ranges
  const previousDaysUsed = sitStatus.pastSITServiceItems?.map((pastSITItem) => {
    const sitDaysUsed = moment(pastSITItem.sitDepartureDate).diff(pastSITItem.sitEntryDate, 'days');
    const location = pastSITItem.reServiceCode === SERVICE_ITEM_CODES.DOFSIT ? 'origin' : 'destination';

    const start = formatDate(pastSITItem.sitEntryDate, swaggerDateFormat, 'DD MMM YYYY');
    const end = formatDate(pastSITItem.sitDepartureDate, swaggerDateFormat, 'DD MMM YYYY');
    const text = `${sitDaysUsed} days at ${location} (${start} - ${end})`;

    return <p key={pastSITItem.id}>{text}</p>;
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

  // Customer delivery request
  const customerContactDate =
    formatDate(sitStatus?.currentSIT?.sitCustomerContacted, swaggerDateFormat, 'DD MMM YYYY') || DEFAULT_EMPTY_VALUE;
  const sitRequestedDelivery =
    formatDate(sitStatus?.currentSIT?.sitRequestedDelivery, swaggerDateFormat, 'DD MMM YYYY') || DEFAULT_EMPTY_VALUE;

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
        <p>SIT (STORAGE IN TRANSIT){pendingSITExtension && <Tag>Additional Days Requested</Tag>}</p>
        {!pendingSITExtension && isConvertedToCustomerExpense && <Tag>Converted To Customer Expense</Tag>}
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
            {/* Sit Start and End table */}
            {currentDaysInSIT > 0 && <p className={styles.sitHeader}>Current location: {currentLocation}</p>}
            <DataTable
              columnHeaders={[`SIT start date`, 'SIT authorized end date', 'Calculated total SIT days']}
              dataRow={[sitStartDateElement, sitEndDate, sitStatus.calculatedTotalDaysInSIT]}
              custClass={styles.currentLocation}
            />
          </div>
          <div className={styles.tableContainer} data-testid="sitDaysAtCurrentLocation">
            {/* Total days at current location */}
            <DataTable
              testID="currentSITDateData"
              columnHeaders={[`Total days in ${currentLocation}`, `SIT departure date`]}
              dataRow={[currentDaysInSITElement, sitDepartureDate]}
            />
          </div>
        </>
      )}

      {/* Service Items */}
      {sitStatus.pastSITServiceItems && (
        <div className={styles.tableContainer}>
          <DataTable columnHeaders={['Previously used SIT']} dataRow={[previousDaysUsed]} />
        </div>
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
