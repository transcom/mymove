import React, { useEffect, useState } from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';

import shipmentDefinitionListsStyles from './ShipmentDefinitionLists.module.scss';

import styles from 'styles/descriptionList.module.scss';
import { formatDateWithUTC } from 'shared/dates';
import { ShipmentShape } from 'types/shipment';
import { formatAddress, formatAgent } from 'utils/shipmentDisplay';
import { convertInchesToFeetAndInches } from 'utils/formatMtoShipment';
import {
  setFlagStyles,
  setDisplayFlags,
  getDisplayFlags,
  getMissingOrDash,
  fieldValidationShape,
} from 'utils/displayFlags';
import { ADDRESS_UPDATE_STATUS } from 'constants/shipments';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

const ShipmentInfoList = ({
  className,
  shipment,
  warnIfMissing,
  errorIfMissing,
  showWhenCollapsed,
  isExpanded,
  isForEvaluationReport,
}) => {
  const {
    actualPickupDate,
    requestedPickupDate,
    requiredDeliveryDate,
    scheduledDeliveryDate,
    requestedDeliveryDate,
    scheduledPickupDate,
    actualDeliveryDate,
    pickupAddress,
    secondaryPickupAddress,
    tertiaryPickupAddress,
    destinationAddress,
    destinationType,
    displayDestinationType,
    secondaryDeliveryAddress,
    tertiaryDeliveryAddress,
    mtoAgents,
    counselorRemarks,
    customerRemarks,
    deliveryAddressUpdate,
  } = shipment;

  const { year, make, model, lengthInInches, widthInInches, heightInInches } = shipment?.mobileHomeShipment || {};

  setFlagStyles({
    row: styles.row,
    warning: shipmentDefinitionListsStyles.warning,
    missingInfoError: shipmentDefinitionListsStyles.missingInfoError,
  });
  setDisplayFlags(errorIfMissing, warnIfMissing, showWhenCollapsed, null, shipment);

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

  const length = convertInchesToFeetAndInches(lengthInInches);
  const width = convertInchesToFeetAndInches(widthInInches);
  const height = convertInchesToFeetAndInches(heightInInches);

  const formattedDimensions = `${length?.feet}'${length?.inches > 0 ? ` ${length.inches}"` : ''} L x ${width?.feet}'${
    width?.inches > 0 ? ` ${width.inches}"` : ''
  } W x ${height?.feet}'${height?.inches > 0 ? ` ${height.inches}"` : ''} H`;

  const releasingAgent = mtoAgents ? mtoAgents.find((agent) => agent.agentType === 'RELEASING_AGENT') : false;
  const receivingAgent = mtoAgents ? mtoAgents.find((agent) => agent.agentType === 'RECEIVING_AGENT') : false;

  const agentsElementFlags = getDisplayFlags('mtoAgents');
  const releasingAgentElement = !releasingAgent ? (
    <div className={agentsElementFlags.classes}>
      <dt>Releasing agent</dt>
      <dd data-testid="RELEASING_AGENT">—</dd>
    </div>
  ) : (
    <div className={agentsElementFlags.classes} key={`${releasingAgent.agentType}-${releasingAgent.email}`}>
      <dt>Releasing agent</dt>
      <dd data-testid={releasingAgent.agentType}>{formatAgent(releasingAgent)}</dd>
    </div>
  );

  const receivingAgentElement = !receivingAgent ? (
    <div className={agentsElementFlags.classes}>
      <dt>Receiving agent</dt>
      <dd data-testid="RECEIVING_AGENT">—</dd>
    </div>
  ) : (
    <div className={agentsElementFlags.classes} key={`${receivingAgent.agentType}-${receivingAgent.email}`}>
      <dt>Receiving agent</dt>
      <dd data-testid={receivingAgent.agentType}>{formatAgent(receivingAgent)}</dd>
    </div>
  );

  const scheduledPickupDateElementFlags = getDisplayFlags('scheduledPickupDate');
  const scheduledPickupDateElement = (
    <div className={scheduledPickupDateElementFlags.classes}>
      <dt>Scheduled pickup date</dt>
      <dd data-testid="requestedPickupDate">
        {(scheduledPickupDate && formatDateWithUTC(scheduledPickupDate, 'DD MMM YYYY')) ||
          getMissingOrDash('scheduledPickupDate')}
      </dd>
    </div>
  );

  const requestedPickupDateElementFlags = getDisplayFlags('scheduledPickupDate');
  const requestedPickupDateElement = (
    <div className={requestedPickupDateElementFlags.classes}>
      <dt>Requested pickup date</dt>
      <dd data-testid="requestedPickupDate">
        {(requestedPickupDate && formatDateWithUTC(requestedPickupDate, 'DD MMM YYYY')) ||
          getMissingOrDash('requestedPickupDate')}
      </dd>
    </div>
  );

  const actualPickupDateElementFlags = getDisplayFlags('actualPickupDate');
  const actualPickupDateElement = (
    <div className={requestedPickupDateElementFlags.classes}>
      <dt>Actual pickup date</dt>
      <dd data-testid="actualPickupDate">
        {(actualPickupDate && formatDateWithUTC(actualPickupDate, 'DD MMM YYYY')) ||
          getMissingOrDash('actualPickupDate')}
      </dd>
    </div>
  );

  const requestedDeliveryDateElementFlags = getDisplayFlags('requestedDeliveryDate');
  const requestedDeliveryDateElement = (
    <div className={requestedDeliveryDateElementFlags.classes}>
      <dt>Requested delivery date</dt>
      <dd data-testid="requestedDeliveryDate">
        {(requestedDeliveryDate && formatDateWithUTC(requestedDeliveryDate, 'DD MMM YYYY')) ||
          getMissingOrDash('requestedDeliveryDate')}
      </dd>
    </div>
  );

  const scheduledDeliveryDateElementFlags = getDisplayFlags('scheduledDeliveryDate');
  const scheduledDeliveryDateElement = (
    <div className={scheduledDeliveryDateElementFlags.classes}>
      <dt>Scheduled delivery date</dt>
      <dd data-testid="scheduledDeliveryDate">
        {(scheduledDeliveryDate && formatDateWithUTC(scheduledDeliveryDate, 'DD MMM YYYY')) ||
          getMissingOrDash('scheduledDeliveryDate')}
      </dd>
    </div>
  );

  const requiredDeliveryDateElementFlags = getDisplayFlags('requiredDeliveryDate');
  const requiredDeliveryDateElement = (
    <div className={requiredDeliveryDateElementFlags.classes}>
      <dt>Required delivery date</dt>
      <dd data-testid="requiredDeliveryDate">
        {(requiredDeliveryDate && formatDateWithUTC(requiredDeliveryDate, 'DD MMM YYYY')) ||
          getMissingOrDash('requiredDeliveryDate')}
      </dd>
    </div>
  );
  const actualDeliveryDateElementFlags = getDisplayFlags('actualDeliveryDate');
  const actualDeliveryDateElement = (
    <div className={actualDeliveryDateElementFlags.classes}>
      <dt>Actual delivery date</dt>
      <dd data-testid="actualDeliveryDate">
        {(actualDeliveryDate && formatDateWithUTC(actualDeliveryDate, 'DD MMM YYYY')) ||
          getMissingOrDash('actualDeliverDate')}
      </dd>
    </div>
  );

  const pickupAddressElementFlags = getDisplayFlags('pickupAddress');
  const pickupAddressElement = (
    <div className={pickupAddressElementFlags.classes}>
      <dt>Origin address</dt>
      <dd data-testid="pickupAddress">
        {(pickupAddress && formatAddress(pickupAddress)) || getMissingOrDash('pickupAddress')}
      </dd>
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

  const destinationTypeFlags = getDisplayFlags('destinationType');
  const destinationTypeElement = (
    <div className={destinationTypeFlags.classes}>
      <dt>Destination type</dt>
      <dd data-testid="destinationType">{destinationType || getMissingOrDash('destinationType')}</dd>
    </div>
  );

  const destinationAddressElementFlags = getDisplayFlags('destinationAddress');
  const destinationAddressElement = (
    <div className={destinationAddressElementFlags.classes}>
      <dt>Destination address</dt>
      <dd data-testid="destinationAddress">
        {deliveryAddressUpdate?.status === ADDRESS_UPDATE_STATUS.REQUESTED
          ? 'Review required'
          : (destinationAddress && formatAddress(destinationAddress)) || '—'}
      </dd>
    </div>
  );

  const secondaryDeliveryAddressElementFlags = getDisplayFlags('secondaryDeliveryAddress');
  const secondaryDeliveryAddressElement = (
    <div className={secondaryDeliveryAddressElementFlags.classes}>
      <dt>Second destination address</dt>
      <dd data-testid="secondaryDeliveryAddress">
        {secondaryDeliveryAddress ? formatAddress(secondaryDeliveryAddress) : '—'}
      </dd>
    </div>
  );

  const tertiaryDeliveryAddressElementFlags = getDisplayFlags('tertiaryDeliveryAddress');
  const tertiaryDeliveryAddressElement = (
    <div className={tertiaryDeliveryAddressElementFlags.classes}>
      <dt>Third destination address</dt>
      <dd data-testid="tertiaryDeliveryAddress">
        {tertiaryDeliveryAddress ? formatAddress(tertiaryDeliveryAddress) : '—'}
      </dd>
    </div>
  );

  const yearElementFlags = getDisplayFlags('year');
  const yearElement = (
    <div className={yearElementFlags.classes}>
      <dt>Mobile home year</dt>
      <dd data-testid="year">{year}</dd>
    </div>
  );

  const makeElementFlags = getDisplayFlags('make');
  const makeElement = (
    <div className={makeElementFlags.classes}>
      <dt>Mobile home make</dt>
      <dd data-testid="make">{make}</dd>
    </div>
  );

  const modelElementFlags = getDisplayFlags('model');
  const modelElement = (
    <div className={modelElementFlags.classes}>
      <dt>Mobile home model</dt>
      <dd data-testid="model">{model}</dd>
    </div>
  );

  const dimensionsElementFlags = getDisplayFlags('dimensions');
  const dimensionsElement = (
    <div className={dimensionsElementFlags.classes}>
      <dt>Dimensions</dt>
      <dd data-testid="dimensions">{formattedDimensions}</dd>
    </div>
  );

  const counselorRemarksElementFlags = getDisplayFlags('counselorRemarks');
  const counselorRemarksElement = (
    <div className={counselorRemarksElementFlags.classes}>
      <dt>Counselor remarks</dt>
      <dd data-testid="counselorRemarks">{counselorRemarks || '—'}</dd>
    </div>
  );

  const customerRemarksElementFlags = getDisplayFlags('customerRemarks');
  const customerRemarksElement = (
    <div className={customerRemarksElementFlags.classes}>
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
      data-testid="shipment-info-list"
    >
      {requestedPickupDateElement}
      {pickupAddressElement}
      {secondaryPickupAddressElement}
      {isTertiaryAddressEnabled ? tertiaryPickupAddressElement : null}
      {showElement(agentsElementFlags) && releasingAgentElement}
      {showElement(requestedDeliveryDateElementFlags) && requestedDeliveryDateElement}
      {requestedDeliveryDateElement}
      {destinationAddressElement}
      {showElement(destinationTypeFlags) && displayDestinationType && destinationTypeElement}
      {secondaryDeliveryAddressElement}
      {isTertiaryAddressEnabled ? tertiaryDeliveryAddressElement : null}
      {showElement(agentsElementFlags) && receivingAgentElement}
      {yearElement}
      {makeElement}
      {modelElement}
      {dimensionsElement}
      {counselorRemarksElement}
      {customerRemarksElement}
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
          {showElement(scheduledPickupDateElement) && scheduledPickupDateElement}
          {showElement(actualPickupDateElementFlags) && actualPickupDateElement}
          {showElement(requestedDeliveryDateElementFlags) && requestedDeliveryDateElement}
          {showElement(agentsElementFlags) && releasingAgentElement}
        </dl>
      </div>
      <div className={shipmentDefinitionListsStyles.sidebySideItem}>
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
          {showElement(scheduledDeliveryDateElementFlags) && scheduledDeliveryDateElement}
          {showElement(requiredDeliveryDateElementFlags) && requiredDeliveryDateElement}
          {showElement(actualDeliveryDateElementFlags) && actualDeliveryDateElement}
          {showElement(agentsElementFlags) && receivingAgentElement}
        </dl>
      </div>
    </div>
  );

  return <div>{isForEvaluationReport ? evaluationReportDetails : defaultDetails}</div>;
};

ShipmentInfoList.propTypes = {
  className: PropTypes.string,
  shipment: ShipmentShape.isRequired,
  warnIfMissing: PropTypes.arrayOf(fieldValidationShape),
  errorIfMissing: PropTypes.arrayOf(fieldValidationShape),
  showWhenCollapsed: PropTypes.arrayOf(PropTypes.string),
  isExpanded: PropTypes.bool,
};

ShipmentInfoList.defaultProps = {
  className: '',
  warnIfMissing: [],
  errorIfMissing: [],
  showWhenCollapsed: [],
  isExpanded: false,
};

export default ShipmentInfoList;
