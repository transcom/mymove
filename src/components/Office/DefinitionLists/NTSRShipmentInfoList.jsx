import React, { useEffect, useState } from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';

import shipmentDefinitionListsStyles from './ShipmentDefinitionLists.module.scss';

import styles from 'styles/descriptionList.module.scss';
import { formatDate } from 'shared/dates';
import { ShipmentShape } from 'types/shipment';
import { formatAddress, formatAgent, formatAccountingCode } from 'utils/shipmentDisplay';
import { formatWeight } from 'utils/formatters';
import {
  setFlagStyles,
  setDisplayFlags,
  getDisplayFlags,
  getMissingOrDash,
  fieldValidationShape,
} from 'utils/displayFlags';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

const NTSRShipmentInfoList = ({
  className,
  shipment,
  isExpanded,
  warnIfMissing,
  errorIfMissing,
  showWhenCollapsed,
  isForEvaluationReport,
}) => {
  const {
    destinationAddress,
    destinationType,
    displayDestinationType,
    secondaryDeliveryAddress,
    tertiaryDeliveryAddress,
    mtoAgents,
    counselorRemarks,
    customerRemarks,
    ntsRecordedWeight,
    requestedDeliveryDate,
    scheduledPickupDate,
    actualPickupDate,
    scheduledDeliveryDate,
    requiredDeliveryDate,
    actualDeliveryDate,
    storageFacility,
    serviceOrderNumber,
    tacType,
    sacType,
    tac,
    sac,
  } = shipment;

  const receivingAgent = mtoAgents ? mtoAgents.find((agent) => agent.agentType === 'RECEIVING_AGENT') : false;

  setFlagStyles({
    row: styles.row,
    warning: shipmentDefinitionListsStyles.warning,
    missingInfoError: shipmentDefinitionListsStyles.missingInfoError,
  });
  setDisplayFlags(errorIfMissing, warnIfMissing, showWhenCollapsed, null, shipment);

  const [isTertiaryAddressEnabled, setIsTertiaryAddressEnabled] = useState(false);
  useEffect(() => {
    const fetchData = async () => {
      isBooleanFlagEnabled('third_address_available').then((enabled) => {
        setIsTertiaryAddressEnabled(enabled);
      });
    };
    fetchData();
  }, []);

  const showElement = (elementFlags) => {
    return (isExpanded || elementFlags.alwaysShow) && !elementFlags.hideRow;
  };

  const storageFacilityAddressElementFlags = getDisplayFlags('storageFacility');
  const storageFacilityAddressElement = (
    <div className={storageFacilityAddressElementFlags.classes}>
      <dt>Storage facility address</dt>
      <dd data-testid="storageFacilityAddress">
        {storageFacility ? formatAddress(storageFacility.address) : 'Missing'}
        {storageFacility && storageFacility.lotNumber && isExpanded && (
          <>
            <br /> Lot #{storageFacility.lotNumber}
          </>
        )}
      </dd>
    </div>
  );

  const ntsRecordedWeightElementFlags = getDisplayFlags('ntsRecordedWeight');
  const ntsRecordedWeightElement = (
    <div className={ntsRecordedWeightElementFlags.classes}>
      <dt>Previously recorded weight</dt>
      <dd data-testid="ntsRecordedWeight">
        {ntsRecordedWeight ? formatWeight(ntsRecordedWeight) : getMissingOrDash('ntsRecordedWeight')}
      </dd>
    </div>
  );

  const storageFacilityInfoElementFlags = getDisplayFlags('storageFacility');
  const storageFacilityInfoElement = (
    <div className={storageFacilityInfoElementFlags.classes}>
      <dt>Storage facility info</dt>
      <dd data-testid="storageFacilityName">
        {storageFacility && storageFacility.facilityName ? storageFacility.facilityName : 'Missing'}
      </dd>
    </div>
  );

  const serviceOrderNumberElementFlags = getDisplayFlags('serviceOrderNumber');
  const serviceOrderNumberElement = (
    <div className={serviceOrderNumberElementFlags.classes}>
      <dt>Service order #</dt>
      <dd data-testid="serviceOrderNumber">{serviceOrderNumber || getMissingOrDash('serviceOrderNumber')}</dd>
    </div>
  );

  const requestedDeliveryDateElementFlags = getDisplayFlags('requestedDeliveryDate');
  const requestedDeliveryDateElement = (
    <div className={requestedDeliveryDateElementFlags.classes}>
      <dt>Requested delivery date</dt>
      <dd data-testid="requestedDeliveryDate">
        {(requestedDeliveryDate && formatDate(requestedDeliveryDate, 'DD MMM YYYY')) || '—'}
      </dd>
    </div>
  );

  const scheduledDeliveryDateElementFlags = getDisplayFlags('scheduledDeliveryDate');
  const scheduledDeliveryDateElement = (
    <div className={scheduledDeliveryDateElementFlags.classes}>
      <dt>Scheduled delivery date</dt>
      <dd data-testid="scheduledDeliveryDate">
        {(scheduledDeliveryDate && formatDate(scheduledDeliveryDate, 'DD MMM YYYY')) ||
          getMissingOrDash(scheduledDeliveryDate)}
      </dd>
    </div>
  );

  const requiredDeliveryDateElementFlags = getDisplayFlags('requiredDeliveryDate');
  const requiredDeliveryDateElement = (
    <div className={requiredDeliveryDateElementFlags.classes}>
      <dt>Required delivery date</dt>
      <dd data-testid="requiredDeliveryDate">
        {(requiredDeliveryDate && formatDate(requiredDeliveryDate, 'DD MMM YYYY')) ||
          getMissingOrDash(requiredDeliveryDate)}
      </dd>
    </div>
  );

  const actualDeliveryDateElementFlags = getDisplayFlags('actualDeliveryDate');
  const actualDeliveryDateElement = (
    <div className={actualDeliveryDateElementFlags.classes}>
      <dt>Actual delivery date</dt>
      <dd data-testid="actualDeliveryDate">
        {(actualDeliveryDate && formatDate(actualDeliveryDate, 'DD MMM YYYY')) || getMissingOrDash(actualDeliveryDate)}
      </dd>
    </div>
  );

  const scheduledPickupDateElementFlags = getDisplayFlags('scheduledPickupDate');
  const scheduledPickupDateElement = (
    <div className={scheduledPickupDateElementFlags.classes}>
      <dt>Scheduled pickup date</dt>
      <dd data-testid="scheduledPickupDate">
        {(scheduledPickupDate && formatDate(scheduledPickupDate, 'DD MMM YYYY')) ||
          getMissingOrDash('scheduledPickupDate')}
      </dd>
    </div>
  );

  const actualPickupDateElementFlags = getDisplayFlags('actualPickupDate');
  const actualPickupDateElement = (
    <div className={actualPickupDateElementFlags.classes}>
      <dt>Actual pickup date</dt>
      <dd data-testid="actualPickupDate">
        {(actualPickupDate && formatDate(actualPickupDate, 'DD MMM YYYY')) || getMissingOrDash('actualPickupDate')}
      </dd>
    </div>
  );

  const storageFacilityContactInfoElementFlags = getDisplayFlags('storageFacilityContactInfo');
  const storageFacilityContactInfoElement = (
    <div className={storageFacilityContactInfoElementFlags.classes}>
      <dt>Storage information</dt>
      <dd data-testid="storageFacilityName">
        {storageFacility && storageFacility.phone ? storageFacility.phone : '—'}
        <br />
        {storageFacility && storageFacility.email ? storageFacility.email : '—'}
      </dd>
    </div>
  );

  const destinationAddressElementFlags = getDisplayFlags('destinationAddress');
  const destinationAddressElement = (
    <div className={destinationAddressElementFlags.classes}>
      <dt>Delivery address</dt>
      <dd data-testid="destinationAddress">{destinationAddress ? formatAddress(destinationAddress) : '—'}</dd>
    </div>
  );

  const destinationTypeFlags = getDisplayFlags('destinationType');
  const destinationTypeElement = (
    <div className={destinationTypeFlags.classes}>
      <dt>Destination type</dt>
      <dd data-testid="destinationType">{destinationType || getMissingOrDash('destinationType')}</dd>
    </div>
  );

  const secondaryDeliveryAddressElementFlags = getDisplayFlags('secondaryDeliveryAddress');
  const secondaryDeliveryAddressElement = (
    <div className={secondaryDeliveryAddressElementFlags.classes}>
      <dt>Second delivery address</dt>
      <dd data-testid="secondaryDeliveryAddress">
        {secondaryDeliveryAddress ? formatAddress(secondaryDeliveryAddress) : '—'}
      </dd>
    </div>
  );

  const tertiaryDeliveryAddressElementFlags = getDisplayFlags('tertiaryDeliveryAddress');
  const tertiaryDeliveryAddressElement = (
    <div className={tertiaryDeliveryAddressElementFlags.classes}>
      <dt>Third delivery address</dt>
      <dd data-testid="tertiaryDeliveryAddress">
        {tertiaryDeliveryAddress ? formatAddress(tertiaryDeliveryAddress) : '—'}
      </dd>
    </div>
  );

  const tacElementFlags = getDisplayFlags('tacType');
  const tacElement = (
    <div className={tacElementFlags.classes}>
      <dt>TAC</dt>
      <dd data-testid="tacType">{tacType && tac ? formatAccountingCode(tac, tacType) : getMissingOrDash('tacType')}</dd>
    </div>
  );

  const sacElementFlags = getDisplayFlags('sacType');
  const sacElement = (
    <div className={sacElementFlags.classes}>
      <dt>SAC</dt>
      <dd data-testid="sacType">{sacType && sac ? formatAccountingCode(sac, sacType) : '—'}</dd>
    </div>
  );

  const receivingAgentFlags = getDisplayFlags('receivingAgent');
  const receivingAgentElement = (
    <div className={receivingAgentFlags.classes} key={`${receivingAgent.agentType}-${receivingAgent.email}`}>
      <dt>Receiving agent</dt>
      <dd data-testid="receivingAgent">{receivingAgent ? formatAgent(receivingAgent) : '—'}</dd>
    </div>
  );

  const counselorRemarksElementFlags = getDisplayFlags('counselorRemarks');
  const counselorRemarksElement = (
    <div className={counselorRemarksElementFlags.classes}>
      <dt>Counselor remarks</dt>
      <dd data-testid="counselorRemarks">{counselorRemarks || '—'}</dd>
    </div>
  );

  const customerRemarksElement = (
    <div className={styles.row}>
      <dt>Customer remarks</dt>
      <dd data-testid="customerRemarks">{customerRemarks || '—'}</dd>
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
      data-testid="nts-release-shipment-info-list"
    >
      {showElement(ntsRecordedWeightElementFlags) && ntsRecordedWeightElement}
      {showElement(storageFacilityInfoElementFlags) && storageFacilityInfoElement}
      {showElement(serviceOrderNumberElementFlags) && serviceOrderNumberElement}
      {storageFacilityAddressElement}
      {requestedDeliveryDateElement}
      {destinationAddressElement}
      {displayDestinationType && destinationTypeElement}
      {isExpanded && secondaryDeliveryAddressElement}
      {isExpanded && isTertiaryAddressEnabled ? tertiaryDeliveryAddressElement : null}
      {showElement(receivingAgentFlags) && receivingAgentElement}
      {isExpanded && customerRemarksElement}
      {showElement(counselorRemarksElementFlags) && counselorRemarksElement}
      {showElement(tacElementFlags) && tacElement}
      {showElement(sacElementFlags) && sacElement}
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
          data-testid="shipment-info-list"
        >
          {isExpanded && scheduledPickupDateElement}
          {isExpanded && actualPickupDateElement}
          {isExpanded && storageFacilityContactInfoElement}
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
          data-testid="shipment-info-list"
        >
          {isExpanded && scheduledDeliveryDateElement}
          {isExpanded && requiredDeliveryDateElement}
          {isExpanded && actualDeliveryDateElement}
          {isExpanded && receivingAgentElement}
        </dl>
      </div>
    </div>
  );

  return <div>{isForEvaluationReport ? evaluationReportDetails : defaultDetails}</div>;
};

NTSRShipmentInfoList.propTypes = {
  className: PropTypes.string,
  shipment: ShipmentShape.isRequired,
  isExpanded: PropTypes.bool,
  warnIfMissing: PropTypes.arrayOf(fieldValidationShape),
  errorIfMissing: PropTypes.arrayOf(fieldValidationShape),
  showWhenCollapsed: PropTypes.arrayOf(PropTypes.string),
  isForEvaluationReport: PropTypes.bool,
};

NTSRShipmentInfoList.defaultProps = {
  className: '',
  isExpanded: false,
  warnIfMissing: [],
  errorIfMissing: [],
  showWhenCollapsed: [],
  isForEvaluationReport: false,
};

export default NTSRShipmentInfoList;
