import React from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';

import styles from 'styles/descriptionList.module.scss';
import { formatDate } from 'shared/dates';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { ShipmentShape } from 'types/shipment';
import { formatAddress, formatAgent } from 'utils/shipmentDisplay';

const ShipmentInfoList = ({ className, shipment, shipmentType }) => {
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
    storageFacilityAddress,
    serviceOrderNumber,
    tacType,
    sacType,
  } = shipment;

  return (
    <dl
      className={classNames(styles.descriptionList, styles.tableDisplay, styles.compact, className)}
      data-testid="shipment-info-list"
    >
      {primeActualWeight && (
        <div className={styles.row}>
          <dt>Shipment weight</dt>
          <dd data-testid="primeActualWeight">{primeActualWeight || '—'}</dd>
        </div>
      )}
      {storageFacility && (
        <div className={styles.row}>
          <dt>Storage facility info</dt>
          <dd data-testid="storageFacilityName">{storageFacility.facilityName}</dd>
        </div>
      )}
      {serviceOrderNumber && (
        <div className={styles.row}>
          <dt>Service order #</dt>
          <dd data-testid="serviceOrderNumber">{serviceOrderNumber || '—'}</dd>
        </div>
      )}
      {storageFacilityAddress && (
        <div className={styles.row}>
          <dt>Storage facility address</dt>
          <dd>
            {formatAddress(storageFacilityAddress)}
            {storageFacility && storageFacility.lotNumber && (
              <>
                <br /> Lot #{storageFacility.lotNumber}
              </>
            )}
          </dd>
        </div>
      )}
      {requestedPickupDate && (
        <div className={styles.row}>
          <dt>Requested move date</dt>
          <dd>{formatDate(requestedPickupDate, 'DD MMM YYYY')}</dd>
        </div>
      )}
      {pickupAddress && (
        <div className={styles.row}>
          <dt>Origin address</dt>
          <dd>{pickupAddress && formatAddress(pickupAddress)}</dd>
        </div>
      )}
      {secondaryPickupAddress && (
        <div className={styles.row}>
          <dt>Second pickup address</dt>
          <dd>{formatAddress(secondaryPickupAddress)}</dd>
        </div>
      )}
      {requestedDeliveryDate && (
        <div className={styles.row}>
          <dt>Preferred delivery date</dt>
          <dd>{formatDate(requestedDeliveryDate, 'DD MMM YYYY')}</dd>
        </div>
      )}
      {destinationAddress && (
        <div className={styles.row}>
          <dt>{shipmentType === SHIPMENT_OPTIONS.NTSR ? 'Delivery address' : 'Destination address'}</dt>
          <dd data-testid="destinationAddress">{formatAddress(destinationAddress)}</dd>
        </div>
      )}
      {secondaryDeliveryAddress && (
        <div className={styles.row}>
          <dt>{shipmentType === SHIPMENT_OPTIONS.NTSR ? 'Second delivery address' : 'Second destination address'}</dt>
          <dd>{formatAddress(secondaryDeliveryAddress)}</dd>
        </div>
      )}
      {agents &&
        agents.map((agent) => (
          <div className={styles.row} key={`${agent.agentType}-${agent.email}`}>
            <dt>{agent.agentType === 'RELEASING_AGENT' ? 'Releasing agent' : 'Receiving agent'}</dt>
            <dd>{formatAgent(agent)}</dd>
          </div>
        ))}
      <div className={styles.row}>
        <dt>Counselor remarks</dt>
        <dd data-testid="counselorRemarks">{counselorRemarks || '—'}</dd>
      </div>
      {customerRemarks && (
        <div className={styles.row}>
          <dt>Customer remarks</dt>
          <dd data-testid="customerRemarks">{customerRemarks || '—'}</dd>
        </div>
      )}
      {tacType && (
        <div className={styles.row}>
          <dt>TAC</dt>
          <dd data-testid="tacType">{tacType || '—'}</dd>
        </div>
      )}
      {sacType && (
        <div className={styles.row}>
          <dt>SAC</dt>
          <dd data-testid="sacType">{tacType || '—'}</dd>
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
};

ShipmentInfoList.defaultProps = {
  className: '',
  shipmentType: SHIPMENT_OPTIONS.HHG,
};

export default ShipmentInfoList;
