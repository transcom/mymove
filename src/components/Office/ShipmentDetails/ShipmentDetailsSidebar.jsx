import React from 'react';
import { Link } from 'react-router-dom';
import * as PropTypes from 'prop-types';

import SimpleSection from 'containers/SimpleSection/SimpleSection';
import { formatAgent, formatAddress } from 'utils/shipmentDisplay';
import { ShipmentShape } from 'types/shipment';

const ShipmentDetailsSidebar = ({ className, shipment }) => {
  const { agents, secondaryAddresses, serviceOrderNumber, storageFacility } = shipment;

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
          <div>{serviceOrderNumber}</div>
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

      {(secondaryAddresses?.secondaryPickupAddress || secondaryAddresses?.secondaryDeliveryAddress) && (
        <SimpleSection header="Secondary addresses" border>
          {secondaryAddresses?.secondaryPickupAddress && (
            <SimpleSection header="Pickup">
              <div>{formatAddress(secondaryAddresses?.secondaryPickupAddress)}</div>
            </SimpleSection>
          )}

          {secondaryAddresses?.secondaryDeliveryAddress && (
            <SimpleSection header="Destination">
              <div>{formatAddress(secondaryAddresses?.secondaryDeliveryAddress)}</div>
            </SimpleSection>
          )}
        </SimpleSection>
      )}
    </div>
  );
};

ShipmentDetailsSidebar.propTypes = {
  className: PropTypes.string,
  shipment: ShipmentShape,
};

ShipmentDetailsSidebar.defaultProps = {
  className: '',
  shipment: {},
};

export default ShipmentDetailsSidebar;
