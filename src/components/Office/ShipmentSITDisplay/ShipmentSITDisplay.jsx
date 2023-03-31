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
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';

const SITHistoryItem = ({ sitItem }) => (
  <dl>
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

const SITHistoryHeader = ({ sitItem }) => (
  <div className={styles.sitHistoryHeader}>
    Total days of SIT approved: {sitItem.approvedDays}{' '}
    <span>updated on {formatDateFromIso(sitItem.decisionDate, 'DD MMM YYYY')} </span>
  </div>
);
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

  const sitEndDate = moment().utc().add(sitStatus.totalDaysRemaining, 'days').format('DD MMM YYYY');

  const sitHistory = React.useMemo(
    () => sitExtensions.filter((sitItem) => sitItem.status !== SIT_EXTENSION_STATUS.PENDING),
    [sitExtensions],
  );
  // Currently active SIT
  const currentLocation = sitStatus.location === LOCATION_TYPES.ORIGIN ? 'origin SIT' : 'destination SIT';

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
            <Restricted to={permissionTypes.createSITExtension}>
              <p>
                <Button type="button" onClick={() => showReviewSITExtension(true)} unstyled>
                  View request
                </Button>
              </p>
            </Restricted>
          ) : (
            <Restricted to={permissionTypes.updateSITExtension}>
              <Button
                type="button"
                onClick={() => showSubmitSITExtension(true)}
                unstyled
                className={styles.submitSITEXtensionLink}
              >
                Edit
              </Button>
            </Restricted>
          ))}
      </div>

      <DataTable
        columnHeaders={['Total days of SIT approved', 'Total days used', 'Total days remaining']}
        dataRow={[shipment.sitDaysAllowance, sitStatus.totalSITDaysUsed, sitStatus.totalDaysRemaining]}
      />
      <p>Current location: {currentLocation}</p>
      <DataTable
        columnHeaders={[`SIT start date`, 'SIT authorized end date']}
        dataRow={[currentDateEnteredSit, sitEndDate]}
      />
      <DataTable columnHeaders={['Total days in destination SIT']} dataRow={[currentDaysInSit]} />
      {sitStatus.pastSITServiceItems && (
        <DataTable columnHeaders={['Previously used SIT']} dataRow={[previousDaysUsed]} />
      )}
      {sitExtensions && sitHistory.length > 0 && (
        <>
          <p>SIT history</p>
          {sitHistory.map((sitItem) => (
            <DataTable
              key={sitItem.id}
              columnHeaders={[<SITHistoryHeader sitItem={sitItem} />]}
              dataRow={[<SITHistoryItem sitItem={sitItem} />]}
            />
          ))}
        </>
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
