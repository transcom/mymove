import { arrayOf, bool, func, string, shape, object } from 'prop-types';

import { AddressShape } from 'types/address';

export const MtoDisplayOptionsShape = shape({
  schema: object.isRequired,
  showPickupFields: bool.isRequired,
  showDeliveryFields: bool.isRequired,
  displayName: string.isRequired,
});

export const MtoAgentShape = shape({
  firstName: string,
  lastName: string,
  phone: string,
  email: string,
  agentType: string,
});

const placeShape = shape({
  requestedDate: string,
  address: AddressShape,
  agent: MtoAgentShape,
});

export const MtoShipmentFormValuesShape = shape({
  pickup: placeShape,
  delivery: placeShape,
  customerRemarks: string,
});

export const MtoShipmentShape = shape({
  agents: arrayOf(MtoAgentShape),
  customerRemarks: string,
  shipmentType: string,
  requestedPickupDate: string,
  pickupAddress: AddressShape,
  requestedDeliveryDate: string,
  destinationAddress: AddressShape,
});

export const HhgShipmentShape = MtoShipmentShape;

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
  MtoShipmentFormValuesShape,
  MtoDisplayOptionsShape,
  MtoAgentShape,
  MtoShipmentShape,
  HhgShipmentShape,
  NtsShipmentShape,
  NtsrShipmentShape,
};
