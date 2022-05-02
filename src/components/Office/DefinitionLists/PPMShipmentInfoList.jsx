import React from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';

import shipmentDefinitionListsStyles from './ShipmentDefinitionLists.module.scss';

import styles from 'styles/descriptionList.module.scss';
import { formatDate } from 'shared/dates';
import { ShipmentShape } from 'types/shipment';
import { setFlagStyles, setDisplayFlags, getDisplayFlags } from 'utils/displayFlags';

const PPMShipmentInfoList = ({ className, shipment, warnIfMissing, errorIfMissing }) => {
  const {
    advanceRequested,
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
  } = shipment;

  setFlagStyles({
    row: styles.row,
    warning: shipmentDefinitionListsStyles.warning,
    missingInfoError: shipmentDefinitionListsStyles.missingInfoError,
  });
  setDisplayFlags(errorIfMissing, warnIfMissing, null, null, shipment);

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
      <dd data-testid="estimatedWeight">{estimatedWeight}</dd>
    </div>
  );

  const proGearWeightElementFlags = getDisplayFlags('proGearWeight');
  const proGearWeightElement = (
    <div className={proGearWeightElementFlags.classes}>
      <dt>Pro-gear</dt>
      <dd data-testid="proGearWeight">{proGearWeight}</dd>
    </div>
  );

  const spouseProGearElementFlags = getDisplayFlags('spouseProGear');
  const spouseProGearElement = (
    <div className={spouseProGearElementFlags.classes}>
      <dt>Spouse pro-gear</dt>
      <dd data-testid="spouseProGear">{spouseProGearWeight}</dd>
    </div>
  );

  const estimatedIncentiveElementFlags = getDisplayFlags('estimatedIncentive');
  const estimatedIncentiveElement = (
    <div className={estimatedIncentiveElementFlags.classes}>
      <dt>Estimated Incentive</dt>
      <dd data-testid="estimatedIncentive">{estimatedIncentive}</dd>
    </div>
  );

  const advanceRequestedElementFlags = getDisplayFlags('advanceRequested');
  const advanceRequestedElement = (
    <div className={advanceRequestedElementFlags.classes}>
      <dt>Advance requested?</dt>
      <dd data-testid="advanceRequested">{advanceRequested ? 'yes' : 'no'}</dd>
    </div>
  );

  // const counselorRemarksElementFlags = getDisplayFlags('counselorRemarks');
  // const counselorRemarksElement = (
  //   <div className={counselorRemarksElementFlags.classes}>
  //     <dt>Counselor remarks</dt>
  //     <dd data-testid="counselorRemarks">{counselorRemarks || '—'}</dd>
  //   </div>
  // );

  return (
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
      {secondOriginZIPElement}
      {destinationZIPElement}
      {secondDestinationZIPElement}
      {sitPlannedElement}
      {estimatedWeightElement}
      {proGearWeightElement}
      {spouseProGearElement}
      {estimatedIncentiveElement}
      {advanceRequestedElement}
      {/* {counselorRemarksElement} */}
    </dl>
  );
};

PPMShipmentInfoList.propTypes = {
  className: PropTypes.string,
  shipment: ShipmentShape.isRequired,
  warnIfMissing: PropTypes.arrayOf(PropTypes.string),
  errorIfMissing: PropTypes.arrayOf(PropTypes.string),
};

PPMShipmentInfoList.defaultProps = {
  className: '',
  warnIfMissing: [],
  errorIfMissing: [],
};

export default PPMShipmentInfoList;
