import { isEmpty } from 'lodash';

import { MTOAgentType } from 'shared/constants';
import { formatSwaggerDate } from 'shared/formatters';

function formatAgent(agent) {
  const agentCopy = { ...agent };
  Object.keys(agentCopy).forEach((key) => {
    if (agentCopy[`${key}`] === '') {
      delete agentCopy[`${key}`];
    } else if (`${key}` === 'phone') {
      const phoneNum = agentCopy[`${key}`];
      // will be in format xxx-xxx-xxxx
      agentCopy[`${key}`] = `${phoneNum.slice(0, 3)}-${phoneNum.slice(3, 6)}-${phoneNum.slice(6, 10)}`;
    }
  });
  return agentCopy;
}

function formatAddress(address) {
  const formattedAddress = address;

  if (formattedAddress.state) {
    formattedAddress.state = formattedAddress.state?.toUpperCase();
    return formattedAddress;
  }

  return undefined;
}

/**
 * formatMtoShipment converts mtoShipment data from the template format to the format API calls expect
 * @param {*} param -  unnamed object representing various mtoShipment data parts
 */
export function formatMtoShipment({ moveId, shipmentType, pickup, delivery, customerRemarks }) {
  const formattedMtoShipment = {
    moveTaskOrderID: moveId,
    shipmentType,
    customerRemarks,
    agents: [],
  };

  if (pickup?.requestedDate) {
    formattedMtoShipment.requestedPickupDate = formatSwaggerDate(pickup.requestedDate);
    formattedMtoShipment.pickupAddress = formatAddress(pickup.address);

    if (pickup.agent) {
      const formattedAgent = formatAgent(pickup.agent);
      if (!isEmpty(formattedAgent)) {
        formattedMtoShipment.agents.push({ ...formattedAgent, agentType: MTOAgentType.RELEASING });
      }
    }
  }

  if (delivery?.requestedDate) {
    formattedMtoShipment.requestedDeliveryDate = formatSwaggerDate(delivery.requestedDate);

    if (delivery.address) {
      formattedMtoShipment.destinationAddress = formatAddress(delivery.address);
    }

    if (delivery.agent) {
      const formattedAgent = formatAgent(delivery.agent);
      if (!isEmpty(formattedAgent)) {
        formattedMtoShipment.agents.push({ ...formattedAgent, agentType: MTOAgentType.RECEIVING });
      }
    }
  }

  if (!formatMtoShipment.agents?.length) {
    formatMtoShipment.agents = undefined;
  }

  return formattedMtoShipment;
}

export default formatMtoShipment;
