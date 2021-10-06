import React from 'react';
import * as PropTypes from 'prop-types';

import SimpleSection from 'containers/SimpleSection/SimpleSection';
import { AddressShape } from 'types/address';
import { AgentShape } from 'types/agent';
import { formatAgent, formatAddress } from 'utils/shipmentDisplay';

const ShipmentDetailsSidebar = ({ className, agents, secondaryAddresses }) => {
  const { secondaryPickupAddress, secondaryDeliveryAddress } = secondaryAddresses;
  return (
    <div className={className}>
      {agents &&
        agents.map((agent) => (
          <SimpleSection
            key={`${agent.agentType}-${agent.email}`}
            header={agent.agentType === 'RELEASING_AGENT' ? 'Releasing agent' : 'Receiving agent'}
            border
          >
            <div>{formatAgent(agent)}</div>
          </SimpleSection>
        ))}

      {(secondaryPickupAddress || secondaryDeliveryAddress) && (
        <SimpleSection header="Secondary addresses" border>
          {secondaryPickupAddress && (
            <SimpleSection header="Pickup">
              <div>{formatAddress(secondaryPickupAddress)}</div>
            </SimpleSection>
          )}

          {secondaryDeliveryAddress && (
            <SimpleSection header="Destination">
              <div>{formatAddress(secondaryDeliveryAddress)}</div>
            </SimpleSection>
          )}
        </SimpleSection>
      )}
    </div>
  );
};

ShipmentDetailsSidebar.propTypes = {
  className: PropTypes.string,
  agents: PropTypes.arrayOf(AgentShape),
  secondaryAddresses: PropTypes.shape({
    secondaryPickupAddress: AddressShape,
    secondaryDeliveryAddress: AddressShape,
  }),
};

ShipmentDetailsSidebar.defaultProps = {
  className: '',
  agents: [],
  secondaryAddresses: {},
};

export default ShipmentDetailsSidebar;
