import React from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import shipmentDefinitionListsStyles from './ShipmentDefinitionLists.module.scss';

import styles from 'styles/descriptionList.module.scss';
import { formatDate } from 'shared/dates';
import { ShipmentShape } from 'types/shipment';
import { formatAddress, formatAgent, formatAccountingCode } from 'utils/shipmentDisplay';
import { setFlagStyles, setDisplayFlags, getDisplayFlags, getMissingOrDash } from 'utils/displayFlags';

const NTSShipmentInfoList = ({
  className,
  shipment,
  isExpanded,
  warnIfMissing,
  errorIfMissing,
  showWhenCollapsed,
  neverShow,
  isForEvaluationReport,
}) => {
  const {
    pickupAddress,
    secondaryPickupAddress,
    agents,
    mtoAgents,
    counselorRemarks,
    customerRemarks,
    requestedPickupDate,
    storageFacility,
    serviceOrderNumber,
    scheduledPickupDate,
    actualPickupDate,
    scheduledDeliveryDate,
    requiredDeliveryDate,
    actualDeliveryDate,
    tacType,
    sacType,
    tac,
    sac,
    usesExternalVendor,
  } = shipment;

  setFlagStyles({
    row: styles.row,
    warning: shipmentDefinitionListsStyles.warning,
    missingInfoError: shipmentDefinitionListsStyles.missingInfoError,
  });
  // Never show is an option since NTSShipmentInfoList is used by both the TOO
  // and services counselor and show different things.
  setDisplayFlags(errorIfMissing, warnIfMissing, showWhenCollapsed, neverShow, shipment);

  const releasingAgent = mtoAgents ? mtoAgents.find((agent) => agent.agentType === 'RELEASING_AGENT') : false;

  const showElement = (elementFlags) => {
    return (isExpanded || elementFlags.alwaysShow) && !elementFlags.hideRow;
  };

  const usesExternalVendorElementFlags = getDisplayFlags('usesExternalVendor');
  const usesExternalVendorElement = (
    <div className={usesExternalVendorElementFlags.classes}>
      <dt>Vendor</dt>
      <dd data-testid="usesExternalVendor">{usesExternalVendor ? 'External vendor' : 'GHC prime contractor'}</dd>
    </div>
  );

  const storageFacilityAddressElementFlags = getDisplayFlags('storageFacility');
  const storageFacilityAddressElement = (
    <div className={storageFacilityAddressElementFlags.classes}>
      <dt>Storage facility address</dt>
      <dd data-testid="storageFacilityAddress">
        {storageFacility ? formatAddress(storageFacility.address) : getMissingOrDash('storageFacility')}
        {storageFacility && storageFacility.lotNumber && isExpanded && (
          <>
            <br /> Lot #{storageFacility.lotNumber}
          </>
        )}
      </dd>
    </div>
  );

  const storageFacilityInfoElementFlags = getDisplayFlags('storageFacility');
  const storageFacilityInfoElement = (
    <div className={storageFacilityInfoElementFlags.classes}>
      <dt>Storage facility info</dt>
      <dd data-testid="storageFacilityName">
        {storageFacility && storageFacility.facilityName
          ? storageFacility.facilityName
          : getMissingOrDash('storageFacility')}
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

  const serviceOrderNumberElementFlags = getDisplayFlags('serviceOrderNumber');
  const serviceOrderNumberElement = (
    <div className={serviceOrderNumberElementFlags.classes}>
      <dt>Service order #</dt>
      <dd data-testid="serviceOrderNumber">{serviceOrderNumber || 'Missing'}</dd>
    </div>
  );

  const requestedPickupDateElementFlags = getDisplayFlags('requestedPickupDate');
  const requestedPickupDateElement = (
    <div className={requestedPickupDateElementFlags.classes}>
      <dt>Requested pickup date</dt>
      <dd data-testid="requestedPickupDate">
        {(requestedPickupDate && formatDate(requestedPickupDate, 'DD MMM YYYY')) || '—'}
      </dd>
    </div>
  );

  const pickupAddressElementFlags = getDisplayFlags('pickupAddress');
  const pickupAddressElement = (
    <div className={pickupAddressElementFlags.classes}>
      <dt>Pickup address</dt>
      <dd data-testid="pickupAddress">{formatAddress(pickupAddress)}</dd>
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

  const releasingAgentFlags = getDisplayFlags('releasingAgent');
  const releasingAgentElement = !releasingAgent ? (
    <div className={releasingAgentFlags.classes}>
      <dt>Releasing agent</dt>
      <dd data-testid="RELEASING_AGENT">—</dd>
    </div>
  ) : (
    <div className={releasingAgentFlags.classes} key={`${releasingAgent.agentType}-${releasingAgent.email}`}>
      <dt>Releasing agent</dt>
      <dd data-testid={releasingAgent.agentType}>{formatAgent(releasingAgent)}</dd>
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

  const agentsElementFlags = getDisplayFlags('agents');
  const agentsElement = agents
    ? agents.map((agent) => (
        <div className={agentsElementFlags.classes} key={`${agent.agentType}-${agent.email}`}>
          <dt>{agent.agentType === 'RELEASING_AGENT' ? 'Releasing agent' : 'Receiving agent'}</dt>
          <dd data-testid="agent">{formatAgent(agent)}</dd>
        </div>
      ))
    : null;

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
      data-testid="nts-shipment-info-list"
    >
      {showElement(usesExternalVendorElementFlags) && usesExternalVendorElement}
      {requestedPickupDateElement}
      {pickupAddressElement}
      {showElement(secondaryPickupAddressElementFlags) && secondaryPickupAddressElement}
      {showElement(agentsElementFlags) && agentsElement}
      {showElement(storageFacilityInfoElementFlags) && storageFacilityInfoElement}
      {showElement(serviceOrderNumberElementFlags) && serviceOrderNumberElement}
      {showElement(storageFacilityAddressElementFlags) && storageFacilityAddressElement}
      {showElement(customerRemarksElementFlags) && customerRemarksElement}
      {showElement(counselorRemarksElementFlags) && counselorRemarksElement}
      {showElement(tacElementFlags) && tacElement}
      {showElement(sacElementFlags) && sacElement}
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
            data-testid="shipment-info-list"
          >
            {isExpanded && scheduledPickupDateElement}
            {isExpanded && actualPickupDateElement}
            {isExpanded && releasingAgentElement}
          </dl>
        </Grid>
        <Grid col={6}>
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
            {isExpanded && storageFacilityContactInfoElement}
          </dl>
        </Grid>
      </Grid>
    </GridContainer>
  );

  return <div>{isForEvaluationReport ? evaluationReportDetails : defaultDetails}</div>;
};

NTSShipmentInfoList.propTypes = {
  className: PropTypes.string,
  shipment: ShipmentShape.isRequired,
  isExpanded: PropTypes.bool,
  warnIfMissing: PropTypes.arrayOf(PropTypes.string),
  errorIfMissing: PropTypes.arrayOf(PropTypes.string),
  showWhenCollapsed: PropTypes.arrayOf(PropTypes.string),
  neverShow: PropTypes.arrayOf(PropTypes.string),
  isForEvaluationReport: PropTypes.bool,
};

NTSShipmentInfoList.defaultProps = {
  className: '',
  isExpanded: false,
  warnIfMissing: [],
  errorIfMissing: [],
  showWhenCollapsed: [],
  neverShow: [],
  isForEvaluationReport: false,
};

export default NTSShipmentInfoList;
