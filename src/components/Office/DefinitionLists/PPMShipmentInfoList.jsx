import React, { useEffect, useState } from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';

import shipmentDefinitionListsStyles from './ShipmentDefinitionLists.module.scss';

import styles from 'styles/descriptionList.module.scss';
import { formatDate } from 'shared/dates';
import AsyncPacketDownloadLink from 'shared/AsyncPacketDownloadLink/AsyncPacketDownloadLink';
import { ShipmentShape } from 'types/shipment';
import { formatAddress } from 'utils/shipmentDisplay';
import { formatCentsTruncateWhole, formatWeight } from 'utils/formatters';
import { setFlagStyles, setDisplayFlags, getDisplayFlags, fieldValidationShape } from 'utils/displayFlags';
import { ADVANCE_STATUSES } from 'constants/ppms';
import { ppmShipmentStatuses } from 'constants/shipments';
import affiliation from 'content/serviceMemberAgencies';
import { permissionTypes } from 'constants/permissions';
import Restricted from 'components/Restricted/Restricted';
import { downloadPPMAOAPacket, downloadPPMPaymentPacket } from 'services/ghcApi';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

const PPMShipmentInfoList = ({
  className,
  shipment,
  warnIfMissing,
  errorIfMissing,
  showWhenCollapsed,
  isExpanded,
  isForEvaluationReport,
  onErrorModalToggle,
}) => {
  const {
    hasRequestedAdvance,
    advanceAmountRequested,
    advanceStatus,
    status,
    destinationAddress,
    estimatedIncentive,
    estimatedWeight,
    expectedDepartureDate,
    actualMoveDate,
    pickupAddress,
    proGearWeight,
    secondaryDestinationAddress,
    secondaryPickupAddress,
    tertiaryDestinationAddress,
    tertiaryPickupAddress,
    sitExpected,
    spouseProGearWeight,
  } = shipment.ppmShipment || {};

  const { closeoutOffice, agency } = shipment;
  const ppmShipmentInfo = { ...shipment.ppmShipment, ...shipment };
  let closeoutDisplay;

  switch (agency) {
    case affiliation.MARINES:
      closeoutDisplay = 'TVCB';
      break;
    case affiliation.NAVY:
      closeoutDisplay = 'NAVY';
      break;
    case affiliation.COAST_GUARD:
      closeoutDisplay = 'USCG';
      break;
    default:
      closeoutDisplay = closeoutOffice || '-';
      break;
  }
  setFlagStyles({
    row: styles.row,
    warning: shipmentDefinitionListsStyles.warning,
    missingInfoError: shipmentDefinitionListsStyles.missingInfoError,
  });
  setDisplayFlags(errorIfMissing, warnIfMissing, showWhenCollapsed, null, ppmShipmentInfo);

  const [isTertiaryAddressEnabled, setIsTertiaryAddressEnabled] = useState(false);
  useEffect(() => {
    const fetchData = async () => {
      setIsTertiaryAddressEnabled(await isBooleanFlagEnabled('third_address_available'));
    };
    if (!isForEvaluationReport) fetchData();
  }, [isForEvaluationReport]);

  const showElement = (elementFlags) => {
    return (isExpanded || elementFlags.alwaysShow) && !elementFlags.hideRow;
  };

  const expectedDepartureDateElementFlags = getDisplayFlags('expectedDepartureDate');
  const expectedDepartureDateElement = (
    <div className={expectedDepartureDateElementFlags.classes}>
      <dt>Estimated Departure date</dt>
      <dd data-testid="expectedDepartureDate">
        {(expectedDepartureDate && formatDate(expectedDepartureDate, 'DD MMM YYYY')) || '—'}
      </dd>
    </div>
  );

  const actualDepartureDateElementFlags = getDisplayFlags('actualMoveDate');
  const actualDepartureDateElement = (
    <div className={actualDepartureDateElementFlags.classes}>
      <dt>Actual Departure date</dt>
      <dd data-testid="actualDepartureDate">{(actualMoveDate && formatDate(actualMoveDate, 'DD MMM YYYY')) || '—'}</dd>
    </div>
  );

  const pickupAddressElementFlags = getDisplayFlags('pickupAddress');
  const pickupAddressElement = (
    <div className={pickupAddressElementFlags.classes}>
      <dt>Pickup address</dt>
      <dd data-testid="pickupAddress">{pickupAddress ? formatAddress(pickupAddress) : '-'}</dd>
    </div>
  );

  const secondaryPickupAddressElementFlags = getDisplayFlags('secondaryPickupAddress');
  const secondaryPickupAddressElement = (
    <div className={secondaryPickupAddressElementFlags.classes}>
      <dt>Second pickup address</dt>
      <dd data-testid="secondaryPickupAddress">
        {secondaryPickupAddress ? formatAddress(secondaryPickupAddress) : '—'}
      </dd>
    </div>
  );

  const tertiaryPickupAddressElementFlags = getDisplayFlags('tertiaryPickupAddress');
  const tertiaryPickupAddressElement = (
    <div className={tertiaryPickupAddressElementFlags.classes}>
      <dt>Third pickup address</dt>
      <dd data-testid="tertiaryPickupAddress">{tertiaryPickupAddress ? formatAddress(tertiaryPickupAddress) : '—'}</dd>
    </div>
  );

  const destinationAddressElementFlags = getDisplayFlags('destinationAddress');
  const destinationAddressElement = (
    <div className={destinationAddressElementFlags.classes}>
      <dt>Destination address</dt>
      <dd data-testid="destinationAddress">{destinationAddress ? formatAddress(destinationAddress) : '-'}</dd>
    </div>
  );

  const secondaryDestinationAddressElementFlags = getDisplayFlags('secondaryDestinationAddress');
  const secondaryDestinationAddressElement = (
    <div className={secondaryDestinationAddressElementFlags.classes}>
      <dt>Second destination address</dt>
      <dd data-testid="secondaryDestinationAddress">
        {secondaryDestinationAddress ? formatAddress(secondaryDestinationAddress) : '—'}
      </dd>
    </div>
  );

  const tertiaryDestinationAddressElementFlags = getDisplayFlags('tertiaryDestinationAddress');
  const tertiaryDestinationAddressElement = (
    <div className={tertiaryDestinationAddressElementFlags.classes}>
      <dt>Third destination address</dt>
      <dd data-testid="tertiaryDestinationAddress">
        {tertiaryDestinationAddress ? formatAddress(tertiaryDestinationAddress) : '—'}
      </dd>
    </div>
  );

  const closeoutOfficeElementFlags = getDisplayFlags('closeoutOffice');
  const closeoutOfficeElement = (
    <div className={closeoutOfficeElementFlags.classes}>
      <dt>Closeout office</dt>
      <dd data-testid="closeout">{closeoutDisplay}</dd>
    </div>
  );

  const sitPlannedElementFlags = getDisplayFlags('sitPlanned');
  const sitPlannedElement = (
    <div className={sitPlannedElementFlags.classes}>
      <dt>SIT planned?</dt>
      <dd data-testid="sitPlanned">{sitExpected ? 'Yes' : 'No'}</dd>
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

  const advanceStatusElementFlags = getDisplayFlags('advanceStatus');
  const advanceStatusElement = (
    <div className={advanceStatusElementFlags.classes}>
      <dt>Advance request status</dt>
      <dd data-testid="advanceRequestStatus">
        {ADVANCE_STATUSES[advanceStatus] ? ADVANCE_STATUSES[advanceStatus].displayValue : `Review required`}
      </dd>
    </div>
  );

  const aoaPacketElement = (
    <div>
      <dt>AOA Packet</dt>
      <dd data-testid="aoaPacketDownload">
        <p className={styles.downloadLink}>
          <AsyncPacketDownloadLink
            id={shipment?.ppmShipment?.id}
            label="Download AOA Paperwork (PDF)"
            asyncRetrieval={downloadPPMAOAPacket}
            onFailure={onErrorModalToggle}
          />
        </p>
      </dd>
    </div>
  );

  const paymentPacketElement = (
    <div>
      <dt>Payment Packet</dt>
      <dd data-testid="paymentPacketDownload">
        <p className={styles.downloadLink}>
          <AsyncPacketDownloadLink
            id={shipment?.ppmShipment?.id}
            label="Download Payment Packet (PDF)"
            asyncRetrieval={downloadPPMPaymentPacket}
            onFailure={onErrorModalToggle}
          />
        </p>
      </dd>
    </div>
  );

  const counselorRemarksElement = (
    <div className={styles.row}>
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
      {!actualMoveDate && expectedDepartureDateElement}
      {actualMoveDate && actualDepartureDateElement}
      {pickupAddressElement}
      {secondaryPickupAddressElement}
      {isTertiaryAddressEnabled ? tertiaryPickupAddressElement : null}
      {destinationAddressElement}
      {secondaryDestinationAddressElement}
      {isTertiaryAddressEnabled ? tertiaryDestinationAddressElement : null}
      <Restricted to={permissionTypes.viewCloseoutOffice}>{closeoutOfficeElement}</Restricted>
      {sitPlannedElement}
      {estimatedWeightElement}
      {showElement(proGearWeightElementFlags) && proGearWeightElement}
      {showElement(spouseProGearElementFlags) && spouseProGearElement}
      {showElement(estimatedIncentiveElementFlags) && estimatedIncentiveElement}
      {hasRequestedAdvanceElement}
      {hasRequestedAdvance === true && advanceStatusElement}
      {(advanceStatus === ADVANCE_STATUSES.APPROVED.apiValue || advanceStatus === ADVANCE_STATUSES.EDITED.apiValue) &&
        aoaPacketElement}
      {status === ppmShipmentStatuses.CLOSEOUT_COMPLETE && paymentPacketElement}
      {counselorRemarksElement}
    </dl>
  );

  const evaluationReportDetails = (
    <div className={shipmentDefinitionListsStyles.sideBySideContainer}>
      <div className={shipmentDefinitionListsStyles.sidebySideItem}>
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
          {isExpanded && pickupAddressElement}
          {isExpanded && expectedDepartureDateElement}
        </dl>
      </div>
      <div className={shipmentDefinitionListsStyles.sidebySideItem}>
        <dl
          className={classNames(
            shipmentDefinitionListsStyles.evaluationReportDL,
            styles.descriptionList,
            styles.tableDisplay,
            styles.compact,
            className,
          )}
          data-testid="shipment-info-list-right"
        >
          {isExpanded && destinationAddressElement}
        </dl>
      </div>
    </div>
  );

  return <div>{isForEvaluationReport ? evaluationReportDetails : defaultDetails}</div>;
};

PPMShipmentInfoList.propTypes = {
  className: PropTypes.string,
  shipment: ShipmentShape.isRequired,
  warnIfMissing: PropTypes.arrayOf(fieldValidationShape),
  errorIfMissing: PropTypes.arrayOf(fieldValidationShape),
  showWhenCollapsed: PropTypes.arrayOf(PropTypes.string),
  isExpanded: PropTypes.bool,
  isForEvaluationReport: PropTypes.bool,
  onErrorModalToggle: PropTypes.func,
};

PPMShipmentInfoList.defaultProps = {
  className: '',
  warnIfMissing: [],
  errorIfMissing: [],
  showWhenCollapsed: [],
  isExpanded: false,
  isForEvaluationReport: false,
  onErrorModalToggle: undefined,
};

export default PPMShipmentInfoList;
