import React from 'react';
import * as PropTypes from 'prop-types';

import SimpleSection from 'containers/SimpleSection/SimpleSection';
import { AddressShape } from 'types/address';
import { AgentShape } from 'types/agent';
import { formatAgent, formatAddress } from 'utils/shipmentDisplay';

const ShipmentDetailsSidebar = ({ className, agents, secondaryAddresses }) => {
  return (
    <div className={className}>
      {agents &&
        agents.map((agent) => (
          <SimpleSection
            header={agent.agentType === 'RELEASING_AGENT' ? 'Releasing agent' : 'Receiving agent'}
            key={`${agent.agentType}-${agent.email}`}
          >
            <div>{formatAgent(agent)}</div>
          </SimpleSection>
        ))}

      {secondaryAddresses && (
        <SimpleSection header="Secondary addresses">
          {secondaryAddresses.secondaryPickupAddress && (
            <SimpleSection header="Pickup" border={false}>
              <div>{formatAddress(secondaryAddresses.secondaryPickupAddress)}</div>
            </SimpleSection>
          )}

          {secondaryAddresses.secondaryDestinationAddress && (
            <SimpleSection header="Destination" border={false}>
              <div>{formatAddress(secondaryAddresses.secondaryDestinationAddress)}</div>
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
    secondaryDestinationAddress: AddressShape,
  }),
};

ShipmentDetailsSidebar.defaultProps = {
  className: '',
  agents: [],
  secondaryAddresses: {},
};

export default ShipmentDetailsSidebar;
