import React from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';

import styles from './OfficeDefinitionLists.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';
import { formatDate } from 'shared/dates';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { ShipmentShape } from 'types/shipment';
import { formatAddress, formatAgent } from 'utils/shipmentDisplay';

const ShipmentInfoList = ({ className, shipment, shipmentType, isExpanded }) => {
  const {
    requestedPickupDate,
    pickupAddress,
    secondaryPickupAddress,
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
  } = shipment;

  // if (shipmentType === SHIPMENT_OPTIONS.NTSR) {
  //   if (isExpanded) {
  //     let listItems = [storageFacilityAddressElement]
  //   } else {
  //
  //   }
  //
  // }
  // const present = [];
  // const requiredButMissing = [];

  const storageFacilityAddressElement =
    shipmentType === SHIPMENT_OPTIONS.NTSR ? (
      <div className={classNames((descriptionListStyles.row, { [styles.missingInfoError]: !storageFacility }))}>
        <dt>Storage facility address</dt>
        <dd>
          {storageFacility ? formatAddress(storageFacility.address) : 'Missing'}
          {storageFacility && storageFacility.lotNumber && (
            <>
              <br /> Lot #{storageFacility.lotNumber}
            </>
          )}
        </dd>
      </div>
    ) : null;

  // if (storageFacilityAddress) {
  //   present.push(storageFacilityAddressElement);
  // } else {
  //   requiredButMissing.push(storageFacilityAddressElement);
  // }

  const primeActualWeightElement =
    shipmentType === SHIPMENT_OPTIONS.NTSR ? (
      <div className={classNames((descriptionListStyles.row, { [styles.missingInfoError]: !primeActualWeight }))}>
        <dt>Shipment weight</dt>
        <dd data-testid="primeActualWeight">{primeActualWeight || '—'}</dd>
      </div>
    ) : null;

  const storageFacilityInfoElement =
    shipmentType === SHIPMENT_OPTIONS.NTSR ? (
      <div className={classNames((descriptionListStyles.row, { [styles.missingInfoError]: !storageFacility }))}>
        <dt>Storage facility info</dt>
        <dd data-testid="storageFacilityName">
          {storageFacility && storageFacility.facilityName ? storageFacility.facilityName : 'Missing'}
        </dd>
      </div>
    ) : null;

  const serviceOrderNumberElement =
    shipmentType === SHIPMENT_OPTIONS.NTSR ? (
      <div className={descriptionListStyles.row}>
        <dt>Service order #</dt>
        <dd data-testid="serviceOrderNumber">{serviceOrderNumber || '—'}</dd>
      </div>
    ) : null;

  const requestedDeliveryDateElement =
    shipmentType === SHIPMENT_OPTIONS.NTSR ? (
      <div className={descriptionListStyles.row}>
        <dt>Preferred delivery date</dt>
        <dd>{(requestedDeliveryDate && formatDate(requestedDeliveryDate, 'DD MMM YYYY')) || '—'}</dd>
      </div>
    ) : null;

  const requestedPickupDateElement =
    shipmentType === SHIPMENT_OPTIONS.NTS || shipmentType === SHIPMENT_OPTIONS.HHG ? (
      <div className={descriptionListStyles.row}>
        <dt>{shipmentType === SHIPMENT_OPTIONS.NTS ? 'Preferred pickup date' : 'Requested move date'}</dt>
        <dd>{requestedPickupDate && formatDate(requestedPickupDate, 'DD MMM YYYY')}</dd>
      </div>
    ) : null;

  const pickupAddressElement =
    shipmentType === SHIPMENT_OPTIONS.NTS || shipmentType === SHIPMENT_OPTIONS.HHG ? (
      <div className={descriptionListStyles.row}>
        <dt>{shipmentType === SHIPMENT_OPTIONS.NTS ? 'Origin address' : 'Current address'}</dt>
        <dd>{pickupAddress && formatAddress(pickupAddress)}</dd>
      </div>
    ) : null;

  const secondaryPickupAddressElement =
    shipmentType === SHIPMENT_OPTIONS.NTS || shipmentType === SHIPMENT_OPTIONS.HHG ? (
      <div className={descriptionListStyles.row}>
        <dt>Second pickup address</dt>
        <dd>{secondaryPickupAddress && formatAddress(secondaryPickupAddress)}</dd>
      </div>
    ) : null;

  const destinationAddressElement =
    shipmentType === SHIPMENT_OPTIONS.NTSR || shipmentType === SHIPMENT_OPTIONS.HHG ? (
      <div className={descriptionListStyles.row}>
        <dt>{shipmentType === SHIPMENT_OPTIONS.NTSR ? 'Delivery address' : 'Destination address'}</dt>
        <dd data-testid="destinationAddress">{formatAddress(destinationAddress)}</dd>
      </div>
    ) : null;

  const secondaryDeliveryAddressElement =
    shipmentType === SHIPMENT_OPTIONS.NTSR || shipmentType === SHIPMENT_OPTIONS.HHG ? (
      <div className={descriptionListStyles.row}>
        <dt>{shipmentType === SHIPMENT_OPTIONS.NTSR ? 'Second delivery address' : 'Second destination address'}</dt>
        <dd data-testid="secondaryDeliveryAddress">
          {secondaryDeliveryAddress ? formatAddress(secondaryDeliveryAddress) : '—'}
        </dd>
      </div>
    ) : null;

  return (
    <dl
      className={classNames(
        descriptionListStyles.descriptionList,
        descriptionListStyles.tableDisplay,
        descriptionListStyles.compact,
        className,
      )}
      data-testid="shipment-info-list"
    >
      {requestedPickupDateElement}
      {isExpanded && primeActualWeightElement}
      {pickupAddress && pickupAddressElement}
      {isExpanded && secondaryPickupAddress && secondaryPickupAddressElement}
      {isExpanded && storageFacilityInfoElement}
      {!isExpanded && requestedDeliveryDate && requestedDeliveryDateElement}
      {isExpanded && serviceOrderNumberElement}
      {(isExpanded || storageFacility) && storageFacilityAddressElement}
      {isExpanded && requestedDeliveryDate && requestedDeliveryDateElement}
      {destinationAddressElement}
      {isExpanded && secondaryDeliveryAddressElement}

      {!isExpanded && !storageFacility && storageFacilityInfoElement}
      {!isExpanded && !storageFacility && storageFacilityAddressElement}

      {isExpanded &&
        agents &&
        agents.map((agent) => (
          <div className={descriptionListStyles.row} key={`${agent.agentType}-${agent.email}`}>
            <dt>{agent.agentType === 'RELEASING_AGENT' ? 'Releasing agent' : 'Receiving agent'}</dt>
            <dd>{formatAgent(agent)}</dd>
          </div>
        ))}
      {isExpanded && (
        <div className={descriptionListStyles.row}>
          <dt>Customer remarks</dt>
          <dd data-testid="customerRemarks">{customerRemarks || '—'}</dd>
        </div>
      )}

      <div className={descriptionListStyles.row}>
        <dt>Counselor remarks</dt>
        <dd data-testid="counselorRemarks">{counselorRemarks || '—'}</dd>
      </div>
      {isExpanded && (
        <div className={descriptionListStyles.row}>
          <dt>TAC</dt>
          <dd data-testid="tacType">{tacType || '—'}</dd>
        </div>
      )}
      {isExpanded && (
        <div className={descriptionListStyles.row}>
          <dt>SAC</dt>
          <dd data-testid="sacType">{sacType || '—'}</dd>
        </div>
      )}
    </dl>
  );
};

ShipmentInfoList.propTypes = {
  className: PropTypes.string,
  shipment: ShipmentShape.isRequired,
  shipmentType: PropTypes.oneOf([
    SHIPMENT_OPTIONS.HHG,
    SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC,
    SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
    SHIPMENT_OPTIONS.NTS,
    SHIPMENT_OPTIONS.NTSR,
  ]),
  isExpanded: PropTypes.bool,
};

ShipmentInfoList.defaultProps = {
  className: '',
  shipmentType: SHIPMENT_OPTIONS.HHG,
  isExpanded: false,
};

export default ShipmentInfoList;
