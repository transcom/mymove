import { formatSwaggerDate } from 'shared/formatters';

function formatAgent(agent) {
  const agentCopy = { ...agent };
  Object.keys(agentCopy).forEach((key) => {
    /* eslint-disable security/detect-object-injection */
    if (agentCopy[key] === '') {
      delete agentCopy[key];
    } else if (key === 'phone') {
      const phoneNum = agentCopy[key];
      // will be in format xxx-xxx-xxxx
      agentCopy[key] = `${phoneNum.slice(0, 3)}-${phoneNum.slice(3, 6)}-${phoneNum.slice(6, 10)}`;
    }
    /* eslint-enable security/detect-object-injection */
  });
  return agentCopy;
}

function formatAddress(address) {
  return {
    street_address_1: address.street_address_1,
    street_address_2: address.street_address_2,
    city: address.city,
    state: address.state.toUpperCase(),
    postal_code: address.postal_code,
    country: address.country,
  }
}

export function formatMtoShipment({ moveId, shipmentType, pickup, delivery, customerRemarks }) {
  const formattedMtoShipment = {
    moveTaskOrderID: moveId,
    shipmentType: shipmentType,
    customerRemarks,
    agents: [],
  };
  
  if (pickup) {
    formattedMtoShipment.requestedPickupDate = formatSwaggerDate(pickup.requestedDate);
    formattedMtoShipment.pickupAddress = formatAddress(pickup.address);

    if (pickup.agent) {
      const formattedAgent = formatAgent(pickup.agent);
      if (!isEmpty(formattedAgent)) {
        formattedMtoShipment.agents.push({ ...formattedAgent, agentType: MTOAgentType.RELEASING });
      }
    }
  }
  
  if (delivery) {
    formattedMtoShipment.requestedDeliveryDate = formatSwaggerDate(delivery.requestedDate);
    formattedMtoShipment.destinationAddress = formatAddress(delivery.address);   

    if (delivery.agent) {
      const formattedAgent = formatAgent(delivery.agent);
      if (!isEmpty(formattedAgent)) {
        formattedMtoShipment.agents.push({ ...formattedAgent, agentType: MTOAgentType.RECEIVING });
      }
    }
  }
   
  return formattedMtoShipment;
}
