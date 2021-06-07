import React from 'react';
import * as PropTypes from 'prop-types';

import styles from 'components/Office/ShipmentInfoList/ShipmentInfoList.module.scss';
import { formatDate } from 'shared/dates';
import { AddressShape } from 'types/address';
import { AgentShape } from 'types/agent';
import { formatAddress, formatAgent } from 'utils/shipmentDisplay';

const ShipmentInfoList = ({
  className,
  requestedMoveDate,
  originAddress,
  secondPickupAddress,
  destinationAddress,
  secondDestinationAddress,
  agents,
  counselorRemarks,
  customerRemarks,
}) => {
  return (
    <dl className={`${styles.ShipmentInfoList} ${className}`} data-testid="shipment-info-list">
      <div className={styles.row}>
        <dt>Requested move date</dt>
        <dd>{formatDate(requestedMoveDate, 'DD MMM YYYY')}</dd>
      </div>
      <div className={styles.row}>
        <dt>Origin address</dt>
        <dd>{originAddress && formatAddress(originAddress)}</dd>
      </div>
      {secondPickupAddress && (
        <div className={styles.row}>
          <dt>Second pickup address</dt>
          <dd>{formatAddress(secondPickupAddress)}</dd>
        </div>
      )}
      <div className={styles.row}>
        <dt>Destination address</dt>
        <dd data-testid="shipmentDestinationAddress">{formatAddress(destinationAddress)}</dd>
      </div>
      {secondDestinationAddress && (
        <div className={styles.row}>
          <dt>Second destination address</dt>
          <dd data-testid="shipmentSecondDestinationAddress">{formatAddress(secondDestinationAddress)}</dd>
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
  requestedMoveDate: PropTypes.string.isRequired,
  originAddress: AddressShape.isRequired,
  secondPickupAddress: AddressShape,
  destinationAddress: AddressShape,
  secondDestinationAddress: AddressShape,
  agents: PropTypes.arrayOf(AgentShape),
  counselorRemarks: PropTypes.string,
  customerRemarks: PropTypes.string,
};

ShipmentInfoList.defaultProps = {
  className: '',
  secondPickupAddress: null,
  destinationAddress: null,
  secondDestinationAddress: null,
  agents: [],
  counselorRemarks: '',
  customerRemarks: '',
};

export default ShipmentInfoList;
