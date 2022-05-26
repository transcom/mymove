import { arrayOf, bool, func, number, shape, string } from 'prop-types';

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

export const BackupContactShape = shape({
  name: string.isRequired,
  telephone: string.isRequired,
  email: string.isRequired,
});

export default {
  MatchShape,
  HistoryShape,
};
