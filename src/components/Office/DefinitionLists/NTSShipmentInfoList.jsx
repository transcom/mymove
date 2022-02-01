import React from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';

import shipmentDefinitionListsStyles from './ShipmentDefinitionLists.module.scss';

import styles from 'styles/descriptionList.module.scss';
import { formatDate } from 'shared/dates';
import { ShipmentShape } from 'types/shipment';
import { formatAddress, formatAgent, formatAccountingCode } from 'utils/shipmentDisplay';

const NTSShipmentInfoList = ({
  className,
  shipment,
  isExpanded,
  warnIfMissing,
  errorIfMissing,
  showWhenCollapsed,
  neverShow,
}) => {
  const {
    pickupAddress,
    secondaryPickupAddress,
    agents,
    counselorRemarks,
    customerRemarks,
    requestedPickupDate,
    storageFacility,
    serviceOrderNumber,
    tacType,
    sacType,
    tac,
    sac,
    usesExternalVendor,
  } = shipment;

  function getFlags(fieldname) {
    let alwaysShow = false;
    let classes = styles.row;
    // Hide row will override any always show that is set.
    let hideRow = false;

    if (errorIfMissing.includes(fieldname) && !shipment[fieldname]) {
      alwaysShow = true;
      classes = classNames(styles.row, shipmentDefinitionListsStyles.missingInfoError);
      return {
        alwaysShow,
        classes,
      };
    }
    if (warnIfMissing.includes(fieldname) && !shipment[fieldname]) {
      alwaysShow = true;
      classes = classNames(styles.row, shipmentDefinitionListsStyles.warning);
      return {
        alwaysShow,
        classes,
      };
    }
    if (showWhenCollapsed.includes(fieldname)) {
      alwaysShow = true;
    }

    if (neverShow.includes(fieldname)) {
      hideRow = true;
    }

    return {
      hideRow,
      alwaysShow,
      classes,
    };
  }

  const getMissingOrDash = (fieldName) => {
    return errorIfMissing.includes(fieldName) ? 'Missing' : '—';
  };

  const showElement = (elementFlags) => {
    return (isExpanded || elementFlags.alwaysShow) && !elementFlags.hideRow;
  };

  const usesExternalVendorElementFlags = getFlags('usesExternalVendor');
  const usesExternalVendorElement = (
    <div className={usesExternalVendorElementFlags.classes}>
      <dt>Vendor</dt>
      <dd data-testid="usesExternalVendor">{usesExternalVendor ? 'External vendor' : 'GHC prime contractor'}</dd>
    </div>
  );

  const storageFacilityAddressElementFlags = getFlags('storageFacility');
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

  const storageFacilityInfoElementFlags = getFlags('storageFacility');
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

  const serviceOrderNumberElementFlags = getFlags('serviceOrderNumber');
  const serviceOrderNumberElement = (
    <div className={serviceOrderNumberElementFlags.classes}>
      <dt>Service order #</dt>
      <dd data-testid="serviceOrderNumber">{serviceOrderNumber || 'Missing'}</dd>
    </div>
  );

  const requestedPickupDateElementFlags = getFlags('requestedPickupDate');
  const requestedPickupDateElement = (
    <div className={requestedPickupDateElementFlags.classes}>
      <dt>Preferred pickup date</dt>
      <dd data-testid="requestedPickupDate">
        {(requestedPickupDate && formatDate(requestedPickupDate, 'DD MMM YYYY')) || '—'}
      </dd>
    </div>
  );

  const pickupAddressElementFlags = getFlags('pickupAddress');
  const pickupAddressElement = (
    <div className={pickupAddressElementFlags.classes}>
      <dt>Pickup address</dt>
      <dd data-testid="pickupAddress">{formatAddress(pickupAddress)}</dd>
    </div>
  );

  const secondaryPickupAddressElementFlags = getFlags('secondaryPickupAddress');
  const secondaryPickupAddressElement = (
    <div className={secondaryPickupAddressElementFlags.classes}>
      <dt>Second pickup address</dt>
      <dd data-testid="secondaryPickupAddress">
        {secondaryPickupAddress ? formatAddress(secondaryPickupAddress) : '—'}
      </dd>
    </div>
  );

  const tacElementFlags = getFlags('tacType');
  const tacElement = (
    <div className={tacElementFlags.classes}>
      <dt>TAC</dt>
      <dd data-testid="tacType">{tacType && tac ? formatAccountingCode(tac, tacType) : getMissingOrDash('tacType')}</dd>
    </div>
  );

  const sacElementFlags = getFlags('sacType');
  const sacElement = (
    <div className={sacElementFlags.classes}>
      <dt>SAC</dt>
      <dd data-testid="sacType">{sacType && sac ? formatAccountingCode(sac, sacType) : '—'}</dd>
    </div>
  );

  const agentsElementFlags = getFlags('agents');
  const agentsElement = agents
    ? agents.map((agent) => (
        <div className={agentsElementFlags.classes} key={`${agent.agentType}-${agent.email}`}>
          <dt>{agent.agentType === 'RELEASING_AGENT' ? 'Releasing agent' : 'Receiving agent'}</dt>
          <dd data-testid="agent">{formatAgent(agent)}</dd>
        </div>
      ))
    : null;

  const counselorRemarksElementFlags = getFlags('counselorRemarks');
  const counselorRemarksElement = (
    <div className={counselorRemarksElementFlags.classes}>
      <dt>Counselor remarks</dt>
      <dd data-testid="counselorRemarks">{counselorRemarks || '—'}</dd>
    </div>
  );

  const customerRemarksElementFlags = getFlags('customerRemarks');
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
};

NTSShipmentInfoList.propTypes = {
  className: PropTypes.string,
  shipment: ShipmentShape.isRequired,
  isExpanded: PropTypes.bool,
  warnIfMissing: PropTypes.arrayOf(PropTypes.string),
  errorIfMissing: PropTypes.arrayOf(PropTypes.string),
  showWhenCollapsed: PropTypes.arrayOf(PropTypes.string),
  // Never show is added as an option since NTSShipmentInfoList is used by both the TOO
  // and services counselor and show different things.
  neverShow: PropTypes.arrayOf(PropTypes.string),
};

NTSShipmentInfoList.defaultProps = {
  className: '',
  isExpanded: false,
  warnIfMissing: [],
  errorIfMissing: [],
  showWhenCollapsed: [],
  neverShow: [],
};

export default NTSShipmentInfoList;
