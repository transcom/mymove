import { arrayOf, bool, func, number, object, shape, string } from 'prop-types';

import { AddressShape } from 'types/address';
import { DutyLocationShape } from 'types/dutyLocation';

export const WeightAllotment = shape({
  total_weight_self: number.isRequired,
  total_weight_self_plus_dependents: number.isRequired,
  pro_gear_weight: number.isRequired,
  pro_gear_weight_spouse: number.isRequired,
});

export const ServiceMemberShape = shape({
  id: string.isRequired,
  affiliation: string,
  edipi: string,
  rank: string,
  first_name: string,
  middle_name: string,
  last_name: string,
  suffix: string,
  telephone: string,
  secondary_telephone: string,
  personal_email: string,
  email_is_preferred: bool,
  phone_is_preferred: bool,
  residential_address: AddressShape,
  backup_mailing_address: AddressShape,
  weight_allotment: WeightAllotment,
});

export const MoveShape = shape({
  id: string,
  locator: string,
  selected_move_type: string,
  status: string,
});

export const UploadShape = shape({
  filename: string,
  content_type: string,
  id: string,
  status: string,
  bytes: number,
  created_at: string,
  updated_at: string,
  url: string,
});

export const UploadsShape = arrayOf(UploadShape);

export const OrdersShape = shape({
  has_dependents: bool,
  id: string,
  issue_date: string,
  moves: arrayOf(string),
  new_duty_location: DutyLocationShape,
  orders_type: string,
  report_by_date: string,
  service_member_id: string,
  spouse_has_pro_gear: bool,
  status: string,
  updated_at: string,
  uploaded_orders: shape({
    id: string,
    uploads: UploadsShape,
  }),
  uploaded_amended_orders: shape({
    id: string,
    uploads: UploadsShape,
  }),
});

export const DocumentShape = shape({});

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

export const PPMShipmentShape = shape({
  pickupPostalCode: string,
  actualPickupPostalCode: string,
  secondaryPickupPostalCode: string,
  destinationPostalCode: string,
  actualDestinationPostalCode: string,
  secondaryDestinationPostalCode: string,
  sitExpected: bool,
  expectedDepartureDate: string,
  actualDepartureDate: string,
  hasProGear: bool,
  proGearWeight: number,
  spouseProGearWeight: number,
  estimatedWeight: number,
  estimatedIncentive: number,
  advance: number,
  advanceRequested: bool,
  hasReceivedAdvance: bool,
  advanceAmountReceived: number,
  status: string,
});

export const MtoShipmentShape = shape({
  agents: arrayOf(MtoAgentShape),
  customerRemarks: string,
  counselorRemarks: string,
  shipmentType: string,
  requestedPickupDate: string,
  pickupAddress: AddressShape,
  requestedDeliveryDate: string,
  destinationAddress: AddressShape,
  secondaryDeliveryAddress: AddressShape,
  secondaryPickupAddress: AddressShape,
  ppmShipment: PPMShipmentShape,
  status: string,
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

export const AdditionalParamsShape = object;

export const WizardPageShape = shape({
  pageList: PageListShape.isRequired,
  pageKey: PageKeyShape.isRequired,
  match: MatchShape.isRequired,
  history: HistoryShape.isRequired,
});

export const BackupContactShape = shape({
  name: string.isRequired,
  telephone: string.isRequired,
  email: string.isRequired,
});

export default {
  MatchShape,
  HistoryShape,
  PageListShape,
  PageKeyShape,
  WizardPageShape,
  MtoShipmentFormValuesShape,
  MtoAgentShape,
  MtoShipmentShape,
  HhgShipmentShape,
  NtsShipmentShape,
  NtsrShipmentShape,
  AdditionalParamsShape,
};
