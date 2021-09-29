import React, { useState } from 'react';
import classnames from 'classnames';
import moment from 'moment';
import { PropTypes } from 'prop-types';
import { Button, Tag } from '@trussworks/react-uswds';

import DataTableWrapper from '../../DataTableWrapper/index';
import DataTable from '../../DataTable/index';
import ReviewSITExtensionsModal from '../ReviewSITExtensionModal/ReviewSITExtensionModal';
import { SITExtensionShape } from '../../../types/sitExtensions';

import styles from './ShipmentSITExtensions.module.scss';

import { sitExtensionReasons, SIT_EXTENSION_STATUS } from 'constants/sitExtensions';
import { formatDateFromIso, formatDate } from 'shared/formatters';
import { utcDateFormat } from 'shared/dates';
import { SERVICE_ITEM_CODES } from 'constants/serviceItems';
import { ShipmentShape } from 'types/shipment';
import { SitStatusShape, LOCATION_TYPES } from 'types/sitStatusShape';

const ShipmentSITExtensions = (props) => {
  const { sitExtensions, sitStatus, shipment, handleReviewSITExtension } = props;
  const { totalSITDaysUsed, totalDaysRemaining } = sitStatus;

  const [isReviewSITExtensionModalVisible, setisReviewSITExtensionModalVisible] = useState(false);
  const reviewSITExtensionSubmit = (sitExtensionID, formValues) => {
    setisReviewSITExtensionModalVisible(false);
    handleReviewSITExtension(sitExtensionID, formValues);
  };

  const pendingSITExtension = sitExtensions.find((se) => se.status === SIT_EXTENSION_STATUS.PENDING);
  const showModal = isReviewSITExtensionModalVisible && pendingSITExtension !== undefined;

  const sitEndDate = `Ends ${moment().utc().add(totalDaysRemaining, 'days').format('DD MMM YYYY')}`;

  const mappedSITExtensionList = sitExtensions.map((sitExt) => {
    return (
      <dl key={sitExt.id}>
        <div>
          <dt>{sitExt.approvedDays} days added</dt>
          <dd>on {formatDateFromIso(sitExt.decisionDate, 'DD MMM YYYY')}</dd>
        </div>
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
      <p>{totalSITDaysUsed} used</p>
    </>
  );

  const daysRemainingAndEndDate = (
    <>
      <p>{totalDaysRemaining} remaining</p>
      <p>{sitEndDate}</p>
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
    <DataTableWrapper className={classnames('maxw-tablet', styles.mtoShipmentSITExtensions)} testID="sitExtensions">
      <div className={styles.title}>
        <p>SIT (STORAGE IN TRANSIT){pendingSITExtension && <Tag>Extension requested</Tag>}</p>
        {pendingSITExtension && (
          <p>
            <Button type="button" onClick={() => setisReviewSITExtensionModalVisible(true)} unstyled>
              View request
            </Button>
          </p>
        )}
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
      {sitStatus.pastSITServiceItems?.length > 0 && (
        <DataTable columnHeaders={['Previously used SIT']} dataRow={[previousDaysUsed]} />
      )}
      {sitExtensions.length > 0 && <DataTable columnHeaders={['SIT extensions']} dataRow={[mappedSITExtensionList]} />}
      {showModal && (
        <ReviewSITExtensionsModal
          onClose={() => setisReviewSITExtensionModalVisible(false)}
          onSubmit={reviewSITExtensionSubmit}
          sitExtension={pendingSITExtension}
        />
      )}
    </DataTableWrapper>
  );
};

ShipmentSITExtensions.propTypes = {
  sitExtensions: PropTypes.arrayOf(SITExtensionShape),
  handleReviewSITExtension: PropTypes.func.isRequired,
  sitStatus: SitStatusShape.isRequired,
  shipment: ShipmentShape.isRequired,
};

ShipmentSITExtensions.defaultProps = {
  sitExtensions: [],
};

export default ShipmentSITExtensions;
