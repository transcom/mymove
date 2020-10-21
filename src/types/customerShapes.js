import { arrayOf, bool, func, string, shape } from 'prop-types';

import { AddressShape } from 'types/address';

export const MtoAgentShape = shape({
  firstName: string,
  lastName: string,
  phone: string,
  email: string,
  agentType: string,
});

export const HhgShipmentShape = shape({
  agents: arrayOf(MtoAgentShape),
  customerRemarks: string,
  shipmentType: string,
  requestedPickupDate: string,
  pickupAddress: AddressShape,
  requestedDeliveryDate: string,
  destinationAddress: AddressShape,
});

export const NtsShipmentShape = shape({
  agents: arrayOf(MtoAgentShape),
  customerRemarks: string,
  shipmentType: string,
  requestedPickupDate: string,
  pickupAddress: AddressShape,
});

export const NtsrShipmentShape = shape({
  agents: arrayOf(MtoAgentShape),
  customerRemarks: string,
  shipmentType: string,
  requestedDeliveryDate: string,
  destinationAddress: AddressShape,
});

export const MatchShape = shape({
  isExact: bool.isRequired,
  params: shape({
    moveId: string.isRequired,
  }),
  path: string.isRequired,
  url: string.isRequired,
});

export const HistoryShape = shape({
  goBack: func.isRequired,
  push: func.isRequired,
});

export const PageListShape = arrayOf(string);

export const PageKeyShape = string;

export const WizardPageShape = shape({
  pageList: PageListShape.isRequired,
  pageKey: PageKeyShape.isRequired,
  match: MatchShape.isRequired,
  history: HistoryShape.isRequired,
});

export default {
  MatchShape,
  HistoryShape,
  PageListShape,
  PageKeyShape,
  WizardPageShape,
  MtoAgentShape,
  HhgShipmentShape,
  NtsShipmentShape,
  NtsrShipmentShape,
};
