import React from 'react';
import { Link } from 'react-router-dom';
import * as PropTypes from 'prop-types';

import SimpleSection from 'containers/SimpleSection/SimpleSection';
import { AddressShape } from 'types/address';
import { AgentShape } from 'types/agent';
import { formatAgent, formatAddress } from 'utils/shipmentDisplay';

const ShipmentDetailsSidebar = ({ className, agents, secondaryAddresses, storageFacility, serviceOrderNumber }) => {
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

      {storageFacility && storageFacility.facilityName && (
        <SimpleSection
          key="facility-info-and-address"
          header={
            <>
              Facility info and address
              <Link to="" className="usa-link float-right">
                Edit
              </Link>
            </>
          }
          border
        >
          <div>{formatAddress(storageFacility)}</div>
        </SimpleSection>
      )}

      {serviceOrderNumber && (
        <SimpleSection
          key="service-order-number"
          header={
            <>
              Service order number
              <Link to="" className="usa-link float-right">
                Edit
              </Link>
            </>
          }
          border
        >
          <div>{}</div>
        </SimpleSection>
      )}

      <SimpleSection
        key="accounting-codes"
        header={
          <>
            Accounting codes
            <Link to="" className="usa-link float-right">
              Edit
            </Link>
          </>
        }
        border
      >
        <div>{}</div>
      </SimpleSection>

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
  storageFacility: AddressShape,
  serviceOrderNumber: PropTypes.string,
};

ShipmentDetailsSidebar.defaultProps = {
  className: '',
  agents: [],
  secondaryAddresses: {},
  storageFacility: {},
  serviceOrderNumber: '',
};

export default ShipmentDetailsSidebar;
