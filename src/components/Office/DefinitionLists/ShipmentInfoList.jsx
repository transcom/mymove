import React from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import shipmentDefinitionListsStyles from './ShipmentDefinitionLists.module.scss';

import styles from 'styles/descriptionList.module.scss';
import { formatDate } from 'shared/dates';
import { ShipmentShape } from 'types/shipment';
import { formatAddress, formatAgent } from 'utils/shipmentDisplay';
import { setFlagStyles, setDisplayFlags, getDisplayFlags, getMissingOrDash } from 'utils/displayFlags';

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
    destinationAddress,
    destinationType,
    displayDestinationType,
    secondaryDeliveryAddress,
    mtoAgents,
    counselorRemarks,
    customerRemarks,
  } = shipment;

  setFlagStyles({
    row: styles.row,
    warning: shipmentDefinitionListsStyles.warning,
    missingInfoError: shipmentDefinitionListsStyles.missingInfoError,
  });
  setDisplayFlags(errorIfMissing, warnIfMissing, showWhenCollapsed, null, shipment);

  const showElement = (elementFlags) => {
    return (isExpanded || elementFlags.alwaysShow) && !elementFlags.hideRow;
  };

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
      <dd data-testid="requestedPickupDate">{scheduledPickupDate && formatDate(scheduledPickupDate, 'DD MMM YYYY')}</dd>
    </div>
  );

  const requestedPickupDateElementFlags = getDisplayFlags('scheduledPickupDate');
  const requestedPickupDateElement = (
    <div className={requestedPickupDateElementFlags.classes}>
      <dt>Requested pickup date</dt>
      <dd data-testid="requestedPickupDate">{requestedPickupDate && formatDate(requestedPickupDate, 'DD MMM YYYY')}</dd>
    </div>
  );

  const actualPickupDateElementFlags = getDisplayFlags('actualPickupDate');
  const actualPickupDateElement = (
    <div className={requestedPickupDateElementFlags.classes}>
      <dt>Actual pickup date</dt>
      <dd data-testid="actualPickupDate">{actualPickupDate && formatDate(actualPickupDate, 'DD MMM YYYY')}</dd>
    </div>
  );

  const requestedDeliveryDateElementFlags = getDisplayFlags('requestedDeliveryDate');
  const requestedDeliveryDateElement = (
    <div className={requestedDeliveryDateElementFlags.classes}>
      <dt>Requested delivery date</dt>
      <dd data-testid="requestedDeliveryDate">
        {(requestedDeliveryDate && formatDate(requestedDeliveryDate, 'DD MMM YYYY')) ||
          getMissingOrDash('requestedDeliveryDate')}
      </dd>
    </div>
  );

  const scheduledDeliveryDateElementFlags = getDisplayFlags('scheduledDeliveryDate');
  const scheduledDeliveryDateElement = (
    <div className={scheduledDeliveryDateElementFlags.classes}>
      <dt>Scheduled delivery date</dt>
      <dd data-testid="scheduledDeliveryDate">
        {(scheduledDeliveryDate && formatDate(scheduledDeliveryDate, 'DD MMM YYYY')) ||
          getMissingOrDash('scheduledDeliveryDate')}
      </dd>
    </div>
  );

  const requiredDeliveryDateElementFlags = getDisplayFlags('requiredDeliveryDate');
  const requiredDeliveryDateElement = (
    <div className={requiredDeliveryDateElementFlags.classes}>
      <dt>Required delivery date</dt>
      <dd data-testid="requiredDeliveryDate">
        {(requiredDeliveryDate && formatDate(requiredDeliveryDate, 'DD MMM YYYY')) ||
          getMissingOrDash('requiredDeliveryDate')}
      </dd>
    </div>
  );
  const actualDeliveryDateElementFlags = getDisplayFlags('actualDeliveryDate');
  const actualDeliveryDateElement = (
    <div className={actualDeliveryDateElementFlags.classes}>
      <dt>Actual delivery date</dt>
      <dd data-testid="actualDeliveryDate">
        {(actualDeliveryDate && formatDate(actualDeliveryDate, 'DD MMM YYYY')) || getMissingOrDash('actualDeliverDate')}
      </dd>
    </div>
  );

  const pickupAddressElementFlags = getDisplayFlags('pickupAddress');
  const pickupAddressElement = (
    <div className={pickupAddressElementFlags.classes}>
      <dt>Origin address</dt>
      <dd data-testid="pickupAddress">{pickupAddress && formatAddress(pickupAddress)}</dd>
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
      <dd data-testid="destinationAddress">{formatAddress(destinationAddress)}</dd>
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
      {showElement(secondaryPickupAddressElementFlags) && secondaryPickupAddressElement}
      {showElement(agentsElementFlags) && releasingAgentElement}
      {showElement(requestedDeliveryDateElementFlags) && requestedDeliveryDateElement}
      {destinationAddressElement}
      {showElement(destinationTypeFlags) && displayDestinationType && destinationTypeElement}
      {showElement(secondaryDeliveryAddressElementFlags) && secondaryDeliveryAddressElement}
      {showElement(agentsElementFlags) && receivingAgentElement}
      {counselorRemarksElement}
      {customerRemarksElement}
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
            {showElement(scheduledPickupDateElement) && scheduledPickupDateElement}
            {showElement(actualPickupDateElementFlags) && actualPickupDateElement}
            {showElement(requestedDeliveryDateElementFlags) && requestedDeliveryDateElement}
            {showElement(agentsElementFlags) && releasingAgentElement}
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
            {showElement(scheduledDeliveryDateElementFlags) && scheduledDeliveryDateElement}
            {showElement(requiredDeliveryDateElementFlags) && requiredDeliveryDateElement}
            {showElement(actualDeliveryDateElementFlags) && actualDeliveryDateElement}
            {showElement(agentsElementFlags) && receivingAgentElement}
          </dl>
        </Grid>
      </Grid>
    </GridContainer>
  );

  return <div>{isForEvaluationReport ? evaluationReportDetails : defaultDetails}</div>;
};

ShipmentInfoList.propTypes = {
  className: PropTypes.string,
  shipment: ShipmentShape.isRequired,
  warnIfMissing: PropTypes.arrayOf(PropTypes.string),
  errorIfMissing: PropTypes.arrayOf(PropTypes.string),
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
