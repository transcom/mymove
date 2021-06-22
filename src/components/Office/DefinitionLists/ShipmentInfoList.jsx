import React from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';

import styles from 'styles/descriptionList.module.scss';
import { formatDate } from 'shared/dates';
import { ShipmentShape } from 'types/shipment';
import { formatAddress, formatAgent } from 'utils/shipmentDisplay';

const ShipmentInfoList = ({ className, shipment }) => {
  const {
    requestedPickupDate,
    pickupAddress,
    secondaryPickupAddress,
    destinationAddress,
    secondaryDeliveryAddress,
    agents,
    counselorRemarks,
    customerRemarks,
  } = shipment;

  return (
    <dl
      className={classNames(styles.descriptionList, styles.tableDisplay, styles.compact, className)}
      data-testid="shipment-info-list"
    >
      <div className={styles.row}>
        <dt>Requested move date</dt>
        <dd>{formatDate(requestedPickupDate, 'DD MMM YYYY')}</dd>
      </div>
      <div className={styles.row}>
        <dt>Origin address</dt>
        <dd>{pickupAddress && formatAddress(pickupAddress)}</dd>
      </div>
      {secondaryPickupAddress && (
        <div className={styles.row}>
          <dt>Second pickup address</dt>
          <dd>{formatAddress(secondaryPickupAddress)}</dd>
        </div>
      )}
      <div className={styles.row}>
        <dt>Destination address</dt>
        <dd data-testid="shipmentDestinationAddress">{formatAddress(destinationAddress)}</dd>
      </div>
      {secondaryDeliveryAddress && (
        <div className={styles.row}>
          <dt>Second destination address</dt>
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
      <div className={styles.row}>
        <dt>Customer remarks</dt>
        <dd data-testid="customerRemarks">{customerRemarks || '—'}</dd>
      </div>
    </dl>
  );
};

ShipmentInfoList.propTypes = {
  className: PropTypes.string,
  shipment: ShipmentShape.isRequired,
};

ShipmentInfoList.defaultProps = {
  className: '',
};

export default ShipmentInfoList;
