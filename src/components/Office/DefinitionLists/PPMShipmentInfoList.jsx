import React from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import shipmentDefinitionListsStyles from './ShipmentDefinitionLists.module.scss';

import styles from 'styles/descriptionList.module.scss';
import { formatDate } from 'shared/dates';
import { ShipmentShape } from 'types/shipment';
import { formatCentsTruncateWhole, formatWeight } from 'utils/formatters';
import { setFlagStyles, setDisplayFlags, getDisplayFlags } from 'utils/displayFlags';

const PPMShipmentInfoList = ({
  className,
  shipment,
  warnIfMissing,
  errorIfMissing,
  showWhenCollapsed,
  isExpanded,
  isForEvaluationReport,
}) => {
  const {
    hasRequestedAdvance,
    advanceAmountRequested,
    destinationPostalCode,
    estimatedIncentive,
    estimatedWeight,
    expectedDepartureDate,
    pickupPostalCode,
    proGearWeight,
    secondaryDestinationPostalCode,
    secondaryPickupPostalCode,
    sitExpected,
    spouseProGearWeight,
  } = shipment.ppmShipment || {};

  setFlagStyles({
    row: styles.row,
    warning: shipmentDefinitionListsStyles.warning,
  });
  setDisplayFlags(errorIfMissing, warnIfMissing, showWhenCollapsed, null, shipment);

  const showElement = (elementFlags) => {
    return (isExpanded || elementFlags.alwaysShow) && !elementFlags.hideRow;
  };

  const expectedDepartureDateElementFlags = getDisplayFlags('expectedDepartureDate');
  const expectedDepartureDateElement = (
    <div className={expectedDepartureDateElementFlags.classes}>
      <dt>Departure date</dt>
      <dd data-testid="expectedDepartureDate">
        {expectedDepartureDate && formatDate(expectedDepartureDate, 'DD MMM YYYY')}
      </dd>
    </div>
  );

  const originZIPElementFlags = getDisplayFlags('originZIP');
  const originZIPElement = (
    <div className={originZIPElementFlags.classes}>
      <dt>Origin ZIP</dt>
      <dd data-testid="originZIP">{pickupPostalCode}</dd>
    </div>
  );

  const secondOriginZIPElementFlags = getDisplayFlags('secondOriginZIP');
  const secondOriginZIPElement = (
    <div className={secondOriginZIPElementFlags.classes}>
      <dt>Second origin ZIP</dt>
      <dd data-testid="secondOriginZIP">{secondaryPickupPostalCode}</dd>
    </div>
  );

  const destinationZIPElementFlags = getDisplayFlags('DestinationZIP');
  const destinationZIPElement = (
    <div className={destinationZIPElementFlags.classes}>
      <dt>Destination ZIP</dt>
      <dd data-testid="destinationZIP">{destinationPostalCode}</dd>
    </div>
  );

  const secondDestinationZIPElementFlags = getDisplayFlags('secondDestinationZIP');
  const secondDestinationZIPElement = (
    <div className={secondDestinationZIPElementFlags.classes}>
      <dt>Second destination ZIP</dt>
      <dd data-testid="secondDestinationZIP">{secondaryDestinationPostalCode}</dd>
    </div>
  );

  const sitPlannedElementFlags = getDisplayFlags('sitPlanned');
  const sitPlannedElement = (
    <div className={sitPlannedElementFlags.classes}>
      <dt>SIT planned?</dt>
      <dd data-testid="sitPlanned">{sitExpected ? 'yes' : 'no'}</dd>
    </div>
  );

  const estimatedWeightElementFlags = getDisplayFlags('estimatedWeight');
  const estimatedWeightElement = (
    <div className={estimatedWeightElementFlags.classes}>
      <dt>Estimated weight</dt>
      <dd data-testid="estimatedWeight">{formatWeight(estimatedWeight)}</dd>
    </div>
  );

  const proGearWeightElementFlags = getDisplayFlags('proGearWeight');
  const proGearWeightElement = (
    <div className={proGearWeightElementFlags.classes}>
      <dt>Pro-gear</dt>
      <dd data-testid="proGearWeight">{proGearWeight ? `Yes, ${formatWeight(proGearWeight)}` : 'No'}</dd>
    </div>
  );

  const spouseProGearElementFlags = getDisplayFlags('spouseProGear');
  const spouseProGearElement = (
    <div className={spouseProGearElementFlags.classes}>
      <dt>Spouse pro-gear</dt>
      <dd data-testid="spouseProGear">{spouseProGearWeight ? `Yes, ${formatWeight(spouseProGearWeight)}` : 'No'}</dd>
    </div>
  );

  const estimatedIncentiveElementFlags = getDisplayFlags('estimatedIncentive');
  const estimatedIncentiveElement = (
    <div className={estimatedIncentiveElementFlags.classes}>
      <dt>Estimated Incentive</dt>
      <dd data-testid="estimatedIncentive">
        ${estimatedIncentive ? formatCentsTruncateWhole(estimatedIncentive) : '0'}
      </dd>
    </div>
  );

  const hasRequestedAdvanceElementFlags = getDisplayFlags('hasRequestedAdvance');
  const hasRequestedAdvanceElement = (
    <div className={hasRequestedAdvanceElementFlags.classes}>
      <dt>Advance requested?</dt>
      <dd data-testid="hasRequestedAdvance">
        {hasRequestedAdvance ? `Yes, $${formatCentsTruncateWhole(advanceAmountRequested)}` : 'No'}
      </dd>
    </div>
  );

  const counselorRemarksElementFlags = getDisplayFlags('counselorRemarks');
  const counselorRemarksElement = (
    <div className={counselorRemarksElementFlags.classes}>
      <dt>Counselor remarks</dt>
      <dd data-testid="counselorRemarks">{shipment.counselorRemarks || '—'}</dd>
    </div>
  );

  const defaultDetails = (
    <dl
      className={classNames(
        shipmentDefinitionListsStyles.ShipmentDefinitionLists,
        styles.descriptionList,
        styles.tableDisplay,
        styles.compact,
        className,
      )}
      data-testid="ppm-shipment-info-list"
    >
      {expectedDepartureDateElement}
      {originZIPElement}
      {showElement(secondOriginZIPElementFlags) && secondOriginZIPElement}
      {destinationZIPElement}
      {showElement(secondDestinationZIPElementFlags) && secondDestinationZIPElement}
      {sitPlannedElement}
      {estimatedWeightElement}
      {showElement(proGearWeightElementFlags) && proGearWeightElement}
      {showElement(spouseProGearElementFlags) && spouseProGearElement}
      {showElement(estimatedIncentiveElementFlags) && estimatedIncentiveElement}
      {hasRequestedAdvanceElement}
      {counselorRemarksElement}
    </dl>
  );

  const evaluationReportDetails = (
    <GridContainer className={shipmentDefinitionListsStyles.evaluationReportDLContainer}>
      <Grid row className={shipmentDefinitionListsStyles.evaluationReportRow}>
        <Grid col={6}>
          <dl
            className={classNames(
              shipmentDefinitionListsStyles.evaluationReportDL,
              styles.descriptionList,
              styles.tableDisplay,
              styles.compact,
              className,
            )}
            data-testid="shipment-info-list-left"
          >
            {isExpanded && originZIPElement}
            {isExpanded && expectedDepartureDateElement}
          </dl>
        </Grid>
        <Grid col={6}>
          <dl
            className={classNames(
              shipmentDefinitionListsStyles.evaluationReportDL,
              shipmentDefinitionListsStyles.evaluationReportDLRight,
              styles.descriptionList,
              styles.tableDisplay,
              styles.compact,
              className,
            )}
            data-testid="shipment-info-list-right"
          >
            {isExpanded && destinationZIPElement}
          </dl>
        </Grid>
      </Grid>
    </GridContainer>
  );

  return <div>{isForEvaluationReport ? evaluationReportDetails : defaultDetails}</div>;
};

PPMShipmentInfoList.propTypes = {
  className: PropTypes.string,
  shipment: ShipmentShape.isRequired,
  warnIfMissing: PropTypes.arrayOf(PropTypes.string),
  errorIfMissing: PropTypes.arrayOf(PropTypes.string),
  showWhenCollapsed: PropTypes.arrayOf(PropTypes.string),
  isExpanded: PropTypes.bool,
  isForEvaluationReport: PropTypes.bool,
};

PPMShipmentInfoList.defaultProps = {
  className: '',
  warnIfMissing: [],
  errorIfMissing: [],
  showWhenCollapsed: [],
  isExpanded: false,
  isForEvaluationReport: false,
};

export default PPMShipmentInfoList;
