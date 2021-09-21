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
import { formatDateFromIso } from 'shared/formatters';

const ShipmentSITExtensions = (props) => {
  const { sitExtensions, handleReviewSITExtension } = props;
  const [isReviewSITExtensionModalVisible, setisReviewSITExtensionModalVisible] = useState(false);
  const reviewSITExtensionSubmit = (sitExtensionID, formValues) => {
    setisReviewSITExtensionModalVisible(false);
    handleReviewSITExtension(sitExtensionID, formValues);
  };

  const pendingSITExtension = sitExtensions.find((se) => se.status === SIT_EXTENSION_STATUS.PENDING);
  const showModal = isReviewSITExtensionModalVisible && pendingSITExtension !== undefined;

  const mappedSITExtensionList = sitExtensions.map((sitExt) => {
    overallTotalDaysAuthorized += sitExt.approvedDays;
    return (
      <dl key={sitExt.id}>
        {sitExt.approvedDays > 0 && (
          <div>
            <dt>{sitExt.approvedDays} days added</dt>
            <dd>on {formatDateFromIso(sitExt.decisionDate, 'DD MMM YYYY')}</dd>
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

  // Display overall total days
  const overallTotalDaysUsed = moment().diff(earliestSITDate, 'days');
  const overallTotalDaysAuthorizedAndUsed = (
    <>
      <p>{overallTotalDaysAuthorized} authorized</p>
      <p>{overallTotalDaysUsed} used</p>
    </>
  );

  const overallTotalDaysRemaining = overallTotalDaysAuthorized - overallTotalDaysUsed;
  const overallEndDate = moment().add(overallTotalDaysRemaining, 'days').format('DD MMM YYYY');
  const overallDaysRemainingAndEndDate = (
    <>
      <p>{overallTotalDaysRemaining} remaining</p>
      <p>Ends {overallEndDate}</p>
    </>
  );

  // Currently active SIT
  const currentLocation = sitStatus.location === LOCATION_TYPES.ORIGIN ? 'origin' : 'destination';
  const currentDaysInSit = <p>{sitStatus.totalSITDaysUsed}</p>;
  const currentDateEnteredSit = <p>{moment(sitStatus.sitEntryDate).format('DD MMM YYYY')}</p>;

  // Previous SIT calculations and date ranges
  const previousDaysUsed = sitStatus.pastSITServiceItems?.map((pastSITItem) => {
    const sitDaysUsed = moment(pastSITItem.sitDepartureDate).diff(pastSITItem.sitEntryDate, 'days');
    const location = pastSITItem.reServiceCode === SERVICE_ITEM_CODES.DOPSIT ? 'origin' : 'destination';

    return (
      <p key={pastSITItem.id}>
        {sitDaysUsed} days at {location} ({moment(pastSITItem.sitEntryDate).format('DD MMM YYYY')} -{' '}
        {moment(pastSITItem.sitDepartureDate).format('DD MMM YYYY')})
      </p>
    );
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
        dataRow={[overallTotalDaysAuthorizedAndUsed, overallDaysRemainingAndEndDate]}
      />
      <p>Current location: {currentLocation}</p>
      <DataTable
        columnHeaders={[`Days in ${currentLocation} SIT`, 'Date entered SIT']}
        dataRow={[currentDaysInSit, currentDateEnteredSit]}
      />
      {sitStatus.pastSITServiceItems?.length > 0 && (
        <DataTable columnHeaders={['Previously used SIT']} dataRow={[previousDaysUsed]} />
      )}
      <DataTable columnHeaders={['SIT extensions']} dataRow={[mappedSITExtensionList]} />
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
  sitExtensions: PropTypes.arrayOf(SITExtensionShape).isRequired,
  handleReviewSITExtension: PropTypes.func.isRequired,
  sitStatus: SitStatusShape.isRequired,
  shipment: ShipmentShape.isRequired,
};

ShipmentSITExtensions.defaultProps = {
  sitExtensions: [],
};

export default ShipmentSITExtensions;
