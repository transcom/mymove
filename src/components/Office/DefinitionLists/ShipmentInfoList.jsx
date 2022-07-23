import React from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';

import shipmentDefinitionListsStyles from './ShipmentDefinitionLists.module.scss';

import styles from 'styles/descriptionList.module.scss';
import { formatDate } from 'shared/dates';
import { ShipmentShape } from 'types/shipment';
import { formatAddress, formatAgent } from 'utils/shipmentDisplay';
import { setFlagStyles, setDisplayFlags, getDisplayFlags, getMissingOrDash } from 'utils/displayFlags';

const ShipmentInfoList = ({ className, shipment, warnIfMissing, errorIfMissing, showWhenCollapsed, isExpanded }) => {
  const {
    requestedPickupDate,
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

  const requestedPickupDateElementFlags = getDisplayFlags('requestedPickupDate');
  const requestedPickupDateElement = (
    <div className={requestedPickupDateElementFlags.classes}>
      <dt>Requested pickup date</dt>
      <dd data-testid="requestedPickupDate">{requestedPickupDate && formatDate(requestedPickupDate, 'DD MMM YYYY')}</dd>
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

  return (
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
      {destinationAddressElement}
      {showElement(destinationTypeFlags) && displayDestinationType && destinationTypeElement}
      {showElement(secondaryDeliveryAddressElementFlags) && secondaryDeliveryAddressElement}
      {showElement(agentsElementFlags) && receivingAgentElement}
      {counselorRemarksElement}
      {customerRemarksElement}
    </dl>
  );
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
