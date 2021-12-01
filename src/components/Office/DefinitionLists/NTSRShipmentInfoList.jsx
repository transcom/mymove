import React from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';

import styles from './OfficeDefinitionLists.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';
import { formatDate } from 'shared/dates';
import { ShipmentShape } from 'types/shipment';
import { formatAddress, formatAgent } from 'utils/shipmentDisplay';

const NTSRShipmentInfoList = ({ className, shipment, isExpanded }) => {
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

  const storageFacilityAddressElement = (
    <div className={classNames(descriptionListStyles.row, { [styles.missingInfoError]: !storageFacility })}>
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

  const primeActualWeightElement = (
    <div className={classNames(descriptionListStyles.row, { [styles.warning]: !primeActualWeight })}>
      <dt>Shipment weight</dt>
      <dd data-testid="primeActualWeight">{primeActualWeight ? String(primeActualWeight).concat(' lbs') : '—'}</dd>
    </div>
  );

  const classNameValue = classNames(descriptionListStyles.row, {
    [styles.missingInfoError]: !storageFacility,
  });
  const storageFacilityInfoElement = (
    <div className={classNameValue}>
      <dt>Storage facility info</dt>
      <dd data-testid="storageFacilityName">
        {storageFacility && storageFacility.facilityName ? storageFacility.facilityName : 'Missing'}
      </dd>
    </div>
  );

  const serviceOrderNumberElement = (
    <div className={classNames(descriptionListStyles.row, { [styles.warning]: !serviceOrderNumber })}>
      <dt>Service order #</dt>
      <dd data-testid="serviceOrderNumber">{serviceOrderNumber || '—'}</dd>
    </div>
  );

  const requestedDeliveryDateElement = (
    <div className={descriptionListStyles.row}>
      <dt>Preferred delivery date</dt>
      <dd>{(requestedDeliveryDate && formatDate(requestedDeliveryDate, 'DD MMM YYYY')) || '—'}</dd>
    </div>
  );

  const destinationAddressElement = (
    <div className={descriptionListStyles.row}>
      <dt>Delivery address</dt>
      <dd data-testid="destinationAddress">{formatAddress(destinationAddress)}</dd>
    </div>
  );

  const secondaryDeliveryAddressElement = (
    <div className={descriptionListStyles.row}>
      <dt>Second delivery address</dt>
      <dd data-testid="secondaryDeliveryAddress">
        {secondaryDeliveryAddress ? formatAddress(secondaryDeliveryAddress) : '—'}
      </dd>
    </div>
  );
  const tacElement = (
    <div className={classNames(descriptionListStyles.row, { [styles.warning]: !tacType })}>
      <dt>TAC</dt>
      <dd data-testid="tacType">{tacType && tac ? String(tac).concat(' (', tacType, ')') : '—'}</dd>
    </div>
  );

  const sacElement = (
    <div className={classNames(descriptionListStyles.row, { [styles.warning]: !sacType })}>
      <dt>SAC</dt>
      <dd data-testid="sacType">{sacType && sac ? String(sac).concat(' (', sacType, ')') : '—'}</dd>
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
      {isExpanded && primeActualWeightElement}
      {isExpanded && storageFacilityInfoElement}
      {!isExpanded && requestedDeliveryDateElement}
      {isExpanded && serviceOrderNumberElement}
      {(isExpanded || storageFacility) && storageFacilityAddressElement}
      {isExpanded && requestedDeliveryDateElement}
      {destinationAddressElement}
      {isExpanded && secondaryDeliveryAddressElement}

      {!isExpanded && !storageFacility && storageFacilityInfoElement}
      {!isExpanded && !storageFacility && storageFacilityAddressElement}

      {isExpanded &&
        agents &&
        agents.map((agent) => (
          <div className={descriptionListStyles.row} key={`${agent.agentType}-${agent.email}`}>
            <dt>{agent.agentType === 'RELEASING_AGENT' ? 'Releasing agent' : 'Receiving agent'}</dt>
            <dd data-testid="agent">{formatAgent(agent)}</dd>
          </div>
        ))}

      {isExpanded && (
        <div className={descriptionListStyles.row}>
          <dt>Customer remarks</dt>
          <dd data-testid="customerRemarks">{customerRemarks || '—'}</dd>
        </div>
      )}

      <div className={classNames(descriptionListStyles.row, { [styles.warning]: !counselorRemarks })}>
        <dt>Counselor remarks</dt>
        <dd data-testid="counselorRemarks">{counselorRemarks || '—'}</dd>
      </div>

      {isExpanded && tacElement}
      {isExpanded && sacElement}
    </dl>
  );
};

NTSRShipmentInfoList.propTypes = {
  className: PropTypes.string,
  shipment: ShipmentShape.isRequired,
  isExpanded: PropTypes.bool,
};

NTSRShipmentInfoList.defaultProps = {
  className: '',
  isExpanded: false,
};

export default NTSRShipmentInfoList;
