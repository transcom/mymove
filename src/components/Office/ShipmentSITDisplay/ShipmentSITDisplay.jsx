import React from 'react';
import classnames from 'classnames';
import moment from 'moment';
import { PropTypes } from 'prop-types';
import { Button, Tag } from '@trussworks/react-uswds';

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

const ShipmentSITDisplay = ({
  sitExtensions,
  sitStatus,
  shipment,
  showReviewSITExtension,
  showSubmitSITExtension,
  hideSITExtensionAction,
  className,
}) => {
  const pendingSITExtension = sitExtensions.find((se) => se.status === SIT_EXTENSION_STATUS.PENDING);

  const sitEndDate = `Ends ${moment().utc().add(sitStatus.totalDaysRemaining, 'days').format('DD MMM YYYY')}`;

  const mappedSITExtensionList = sitExtensions
    .filter((sitExt) => sitExt.status !== SIT_EXTENSION_STATUS.PENDING)
    .map((sitExt) => {
      return (
        <dl key={sitExt.id}>
          {sitExt.status === SIT_EXTENSION_STATUS.APPROVED ? (
            <div>
              <dt>{sitExt.approvedDays} days added</dt>
              <dd>on {formatDateFromIso(sitExt.decisionDate, 'DD MMM YYYY')}</dd>
            </div>
          ) : (
            <div>
              <dt>0 days added</dt>
              <dd>on {formatDateFromIso(sitExt.decisionDate, 'DD MMM YYYY')} — request rejected</dd>
            </div>
          )}
          <div>
            <dt>Reason:</dt>
            <dd>{sitExtensionReasons[sitExt.requestReason]}</dd>
          </div>
          {sitExt.contractorRemarks && (
            <div>
              <dt>Contractor remarks:</dt>
              <dd>{sitExt.contractorRemarks}</dd>
            </div>
          )}
          {sitExt.officeRemarks && (
            <div>
              <dt>Office remarks:</dt>
              <dd>{sitExt.officeRemarks}</dd>
            </div>
          )}
        </dl>
      );
    });

  const totalDaysAuthorizedAndUsed = (
    <>
      <p>{shipment.sitDaysAllowance} authorized</p>
      <p>{sitStatus.totalSITDaysUsed} used</p>
    </>
  );

  // data-happo-hide is in place to compensate for mockDate being ignored in Storybook
  const daysRemainingAndEndDate = (
    <>
      <p>{sitStatus.totalDaysRemaining} remaining</p>
      <p data-happo-hide>{sitEndDate}</p>
    </>
  );

  // Currently active SIT
  const currentLocation = sitStatus.location === LOCATION_TYPES.ORIGIN ? 'origin' : 'destination';

  const currentDaysInSit = <p>{sitStatus.totalSITDaysUsed}</p>;
  const currentDateEnteredSit = <p>{formatDate(sitStatus.sitEntryDate, utcDateFormat, 'DD MMM YYYY')}</p>;

  // Previous SIT calculations and date ranges
  const previousDaysUsed = sitStatus.pastSITServiceItems?.map((pastSITItem) => {
    const sitDaysUsed = moment(pastSITItem.sitDepartureDate).utc().diff(pastSITItem.sitEntryDate, 'days');
    const location = pastSITItem.reServiceCode === SERVICE_ITEM_CODES.DOPSIT ? 'origin' : 'destination';

    const start = formatDate(pastSITItem.sitEntryDate, utcDateFormat, 'DD MMM YYYY');
    const end = formatDate(pastSITItem.sitDepartureDate, utcDateFormat, 'DD MMM YYYY');
    const text = `${sitDaysUsed} days at ${location} (${start} - ${end})`;

    return <p key={pastSITItem.id}>{text}</p>;
  });

  return (
    <DataTableWrapper
      className={classnames('maxw-tablet', styles.mtoshipmentSITDisplay, className)}
      testID="sitExtensions"
    >
      <div className={styles.title}>
        <p>SIT (STORAGE IN TRANSIT){pendingSITExtension && <Tag>Extension requested</Tag>}</p>
        {!hideSITExtensionAction &&
          (pendingSITExtension ? (
            <p>
              <Button type="button" onClick={() => showReviewSITExtension(true)} unstyled>
                View request
              </Button>
            </p>
          ) : (
            <Button
              type="button"
              onClick={() => showSubmitSITExtension(true)}
              unstyled
              className={styles.submitSITEXtensionLink}
            >
              Edit
            </Button>
          ))}
      </div>

      <DataTable
        columnHeaders={['Total days of SIT', 'Total days remaining']}
        dataRow={[totalDaysAuthorizedAndUsed, daysRemainingAndEndDate]}
      />
      <p>Current location: {currentLocation}</p>
      <DataTable
        columnHeaders={[`Days in ${currentLocation} SIT`, 'Date entered SIT']}
        dataRow={[currentDaysInSit, currentDateEnteredSit]}
      />
      {sitStatus.pastSITServiceItems && (
        <DataTable columnHeaders={['Previously used SIT']} dataRow={[previousDaysUsed]} />
      )}
      {sitExtensions && mappedSITExtensionList.length > 0 && (
        <DataTable columnHeaders={['SIT extensions']} dataRow={[mappedSITExtensionList]} />
      )}
    </DataTableWrapper>
  );
};

ShipmentSITDisplay.propTypes = {
  sitExtensions: PropTypes.arrayOf(SITExtensionShape),
  sitStatus: SitStatusShape.isRequired,
  shipment: ShipmentShape.isRequired,
  showReviewSITExtension: PropTypes.func,
  showSubmitSITExtension: PropTypes.func,
  hideSITExtensionAction: PropTypes.bool,
  className: PropTypes.string,
};

ShipmentSITDisplay.defaultProps = {
  sitExtensions: [],
  showReviewSITExtension: undefined,
  showSubmitSITExtension: undefined,
  hideSITExtensionAction: false,
  className: '',
};

export default ShipmentSITDisplay;
