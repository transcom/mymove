import { isEmpty } from 'lodash';

import { MTOAgentType } from 'shared/constants';
import { formatSwaggerDate } from 'shared/formatters';

function formatAgentForDisplay(agent) {
  const agentCopy = { ...agent };
  // handle the diff between expected FE and BE phone format
  Object.keys(agentCopy).forEach((key) => {
    /* eslint-disable security/detect-object-injection */
    if (key === 'phone') {
      const phoneNum = agentCopy[key];
      // will be in format xxxxxxxxxx
      agentCopy[key] = phoneNum.split('-').join('');
    }
  });
  return agentCopy;
}

function formatAgentForAPI(agent) {
  const agentCopy = { ...agent };
  Object.keys(agentCopy).forEach((key) => {
    const sanitizedKey = `${key}`;
    /* eslint-disable security/detect-object-injection */
    if (agentCopy[sanitizedKey] === '') {
      delete agentCopy[sanitizedKey];
    } else if (sanitizedKey === 'phone') {
      const phoneNum = agentCopy[sanitizedKey];
      // will be in format xxx-xxx-xxxx
      agentCopy[sanitizedKey] = `${phoneNum.slice(0, 3)}-${phoneNum.slice(3, 6)}-${phoneNum.slice(6, 10)}`;
    }
    /* eslint-enable security/detect-object-injection */
  });
  return agentCopy;
}

function formatAddressForAPI(address) {
  const formattedAddress = address;

  if (formattedAddress.state) {
    formattedAddress.state = formattedAddress.state?.toUpperCase();
    delete formattedAddress.id;
    return formattedAddress;
  }

  return undefined;
}

const emptyAgentShape = {
  firstName: '',
  lastName: '',
  email: '',
  phone: '',
};

const emptyAddressShape = {
  street_address_1: '',
  street_address_2: '',
  city: '',
  state: '',
  postal_code: '',
};

/**
 * formatMtoShipmentForDisplay converts mtoShipment data from the format API calls expect to the template format
 * @param {*} mtoShipment - (see MtoShipmentShape)
 */
export function formatMtoShipmentForDisplay({
  agents,
  shipmentType,
  requestedPickupDate,
  pickupAddress,
  requestedDeliveryDate,
  destinationAddress,
  customerRemarks,
  moveTaskOrderID,
}) {
  const displayValues = {
    shipmentType,
    moveTaskOrderID,
    customerRemarks: customerRemarks || '',
    pickup: {
      requestedDate: '',
      address: { ...emptyAddressShape },
      agent: { ...emptyAgentShape },
    },
    delivery: {
      requestedDate: '',
      address: { ...emptyAddressShape },
      agent: { ...emptyAgentShape },
    },
    hasDeliveryAddress: 'no',
  };

  if (agents) {
    const receivingAgent = agents.find((agent) => agent.agentType === 'RECEIVING_AGENT');
    const releasingAgent = agents.find((agent) => agent.agentType === 'RELEASING_AGENT');

    if (receivingAgent) {
      const formattedAgent = formatAgentForDisplay(receivingAgent);
      if (Object.keys(formattedAgent).length) {
        displayValues.delivery.agent = { ...emptyAgentShape, ...formattedAgent };
      }
    }
    if (releasingAgent) {
      const formattedAgent = formatAgentForDisplay(releasingAgent);
      if (Object.keys(formattedAgent).length) {
        displayValues.pickup.agent = { ...emptyAgentShape, ...formattedAgent };
      }
    }
  }

  if (pickupAddress) {
    displayValues.pickup.address = { ...emptyAddressShape, ...pickupAddress };
  }

  if (requestedPickupDate) {
    displayValues.pickup.requestedDate = requestedPickupDate;
  }

  if (destinationAddress) {
    displayValues.delivery.address = { ...emptyAddressShape, ...destinationAddress };
    displayValues.hasDeliveryAddress = 'yes';
  }

  if (requestedDeliveryDate) {
    displayValues.delivery.requestedDate = requestedDeliveryDate;
  }

  return displayValues;
}

/**
 * formatMtoShipmentForAPI converts mtoShipment data from the template format to the format API calls expect
 * @param {*} param - unnamed object representing various mtoShipment data parts
 */
export function formatMtoShipmentForAPI({ moveId, shipmentType, pickup, delivery, customerRemarks }) {
  const formattedMtoShipment = {
    moveTaskOrderID: moveId,
    shipmentType,
    customerRemarks,
    agents: [],
  };

  if (pickup?.requestedDate && pickup.requestedDate !== '') {
    formattedMtoShipment.requestedPickupDate = formatSwaggerDate(pickup.requestedDate);
    formattedMtoShipment.pickupAddress = formatAddressForAPI(pickup.address);

    if (pickup.agent) {
      const formattedAgent = formatAgentForAPI(pickup.agent);
      if (!isEmpty(formattedAgent)) {
        formattedMtoShipment.agents.push({ ...formattedAgent, agentType: MTOAgentType.RELEASING });
      }
    }
  }

  if (delivery?.requestedDate && delivery.requestedDate !== '') {
    formattedMtoShipment.requestedDeliveryDate = formatSwaggerDate(delivery.requestedDate);

    if (delivery.address) {
      formattedMtoShipment.destinationAddress = formatAddressForAPI(delivery.address);
    }

    if (delivery.agent) {
      const formattedAgent = formatAgentForAPI(delivery.agent);
      if (!isEmpty(formattedAgent)) {
        formattedMtoShipment.agents.push({ ...formattedAgent, agentType: MTOAgentType.RECEIVING });
      }
    }
  }

  if (!formattedMtoShipment.agents?.length) {
    formattedMtoShipment.agents = undefined;
  }

  return formattedMtoShipment;
}

export default { formatMtoShipmentForAPI, formatMtoShipmentForDisplay };
