import { arrayOf, bool, number, shape, string } from 'prop-types';

import { TransportationOfficeShape } from './user';

import { AddressShape } from 'types/address';
import { DutyLocationShape } from 'types/dutyLocation';

export const Entitlements = shape({
  proGear: number.isRequired,
  proGearSpouse: number.isRequired,
});

export const ServiceMemberShape = shape({
  id: string.isRequired,
  affiliation: string,
  edipi: string,
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
});

export const MoveShape = shape({
  id: string,
  locator: string,
  selected_move_type: string,
  status: string,
  closeout_office: TransportationOfficeShape,
});

export const UploadShape = shape({
  filename: string,
  contentType: string,
  id: string,
  status: string,
  bytes: number,
  createdAt: string,
  updatedAt: string,
  url: string,
});

export const UploadsShape = arrayOf(UploadShape);

export const OrdersShape = shape({
  has_dependents: bool,
  id: string,
  issue_date: string,
  moves: arrayOf(string),
  origin_duty_location: DutyLocationShape,
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
  grade: string,
  authorizedWeight: number,
  entitlement: Entitlements,
});

export const BackupContactShape = shape({
  name: string.isRequired,
  telephone: string.isRequired,
  email: string.isRequired,
});
