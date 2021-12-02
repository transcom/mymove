import React from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';

import styles from './OfficeDefinitionLists.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';
import { formatDate } from 'shared/dates';
import { ShipmentShape } from 'types/shipment';
import { formatAddress, formatAgent } from 'utils/shipmentDisplay';

const NTSRShipmentInfoList = ({
  className,
  shipment,
  isExpanded,
  warnIfMissing,
  errorIfMissing,
  showWhenCollapsed,
}) => {
  const {
    destinationAddress,
    secondaryDeliveryAddress,
    agents,
    counselorRemarks,
    customerRemarks,
    primeActualWeight,
    requestedDeliveryDate,
    storageFacility,
    serviceOrderNumber,
    tacType,
    sacType,
    tac,
    sac,
  } = shipment;

  function getFlags(fieldname) {
    let alwaysShow = false;
    let classes = descriptionListStyles.row;

    if (errorIfMissing.includes(fieldname) && !shipment[fieldname]) {
      alwaysShow = true;
      classes = classNames(descriptionListStyles.row, styles.missingInfoError);
      return {
        alwaysShow,
        classes,
      };
    }
    if (warnIfMissing.includes(fieldname) && !shipment[fieldname]) {
      alwaysShow = true;
      classes = classNames(descriptionListStyles.row, styles.warning);
      return {
        alwaysShow,
        classes,
      };
    }
    if (showWhenCollapsed.includes(fieldname)) {
      alwaysShow = true;
    }
    return {
      alwaysShow,
      classes,
    };
  }
  const storageFacilityAddressElementFlags = getFlags('storageFacility');
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

  const primeActualWeightElementFlags = getFlags('primeActualWeight');
  const primeActualWeightElement = (
    <div className={primeActualWeightElementFlags.classes}>
      <dt>Shipment weight</dt>
      <dd data-testid="primeActualWeight">{primeActualWeight ? String(primeActualWeight).concat(' lbs') : '—'}</dd>
    </div>
  );

  const storageFacilityInfoElementFlags = getFlags('storageFacility');
  const storageFacilityInfoElement = (
    <div className={storageFacilityInfoElementFlags.classes}>
      <dt>Storage facility info</dt>
      <dd data-testid="storageFacilityName">
        {storageFacility && storageFacility.facilityName ? storageFacility.facilityName : 'Missing'}
      </dd>
    </div>
  );

  const serviceOrderNumberElementFlags = getFlags('serviceOrderNumber');
  const serviceOrderNumberElement = (
    <div className={serviceOrderNumberElementFlags.classes}>
      <dt>Service order #</dt>
      <dd data-testid="serviceOrderNumber">{serviceOrderNumber || '—'}</dd>
    </div>
  );

  const requestedDeliveryDateElementFlags = getFlags('requestedDeliveryDate');
  const requestedDeliveryDateElement = (
    <div className={requestedDeliveryDateElementFlags.classes}>
      <dt>Preferred delivery date</dt>
      <dd>{(requestedDeliveryDate && formatDate(requestedDeliveryDate, 'DD MMM YYYY')) || '—'}</dd>
    </div>
  );

  const destinationAddressElementFlags = getFlags('destinationAddress');
  const destinationAddressElement = (
    <div className={destinationAddressElementFlags.classes}>
      <dt>Delivery address</dt>
      <dd data-testid="destinationAddress">{formatAddress(destinationAddress)}</dd>
    </div>
  );

  const secondaryDeliveryAddressElement = (
    <div className={descriptionListStyles.rows}>
      <dt>Second delivery address</dt>
      <dd data-testid="secondaryDeliveryAddress">
        {secondaryDeliveryAddress ? formatAddress(secondaryDeliveryAddress) : '—'}
      </dd>
    </div>
  );

  const tacElementFlags = getFlags('tacType');
  const tacElement = (
    <div className={tacElementFlags.classes}>
      <dt>TAC</dt>
      <dd data-testid="tacType">{tacType && tac ? String(tac).concat(' (', tacType, ')') : '—'}</dd>
    </div>
  );

  const sacElementFlags = getFlags('sacType');
  const sacElement = (
    <div className={sacElementFlags.classes}>
      <dt>SAC</dt>
      <dd data-testid="sacType">{sacType && sac ? String(sac).concat(' (', sacType, ')') : '—'}</dd>
    </div>
  );

  const agentsElement = agents
    ? agents.map((agent) => (
        <div className={descriptionListStyles.row} key={`${agent.agentType}-${agent.email}`}>
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

  const customerRemarksElement = (
    <div className={descriptionListStyles.row}>
      <dt>Customer remarks</dt>
      <dd data-testid="customerRemarks">{customerRemarks || '—'}</dd>
    </div>
  );
  return (
    <dl
      className={classNames(
        descriptionListStyles.descriptionList,
        descriptionListStyles.tableDisplay,
        descriptionListStyles.compact,
        className,
        styles.OfficeDefinitionLists,
        styles.ShipmentInfoList,
      )}
      data-testid="shipment-info-list"
    >
      {(isExpanded || primeActualWeightElementFlags.alwaysShow) && primeActualWeightElement}
      {(isExpanded || storageFacilityInfoElementFlags.alwaysShow) && storageFacilityInfoElement}
      {!isExpanded && requestedDeliveryDateElement}
      {(isExpanded || serviceOrderNumberElementFlags.alwaysShow) && serviceOrderNumberElement}
      {storageFacilityAddressElement}
      {isExpanded && requestedDeliveryDateElement}
      {destinationAddressElement}
      {isExpanded && secondaryDeliveryAddressElement}
      {isExpanded && agentsElement}
      {isExpanded && customerRemarksElement}
      {(isExpanded || counselorRemarksElementFlags.alwaysShow) && counselorRemarksElement}
      {(isExpanded || tacElementFlags.alwaysShow) && tacElement}
      {(isExpanded || sacElementFlags.alwaysShow) && sacElement}
    </dl>
  );
};

NTSRShipmentInfoList.propTypes = {
  className: PropTypes.string,
  shipment: ShipmentShape.isRequired,
  isExpanded: PropTypes.bool,
  warnIfMissing: PropTypes.arrayOf(PropTypes.string),
  errorIfMissing: PropTypes.arrayOf(PropTypes.string),
  showWhenCollapsed: PropTypes.arrayOf(PropTypes.string),
};

NTSRShipmentInfoList.defaultProps = {
  className: '',
  isExpanded: false,
  warnIfMissing: [],
  errorIfMissing: [],
  showWhenCollapsed: [],
};

export default NTSRShipmentInfoList;
