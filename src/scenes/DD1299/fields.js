import yaml from 'js-yaml'; //todo: yarn add

export const getUiSchema = () => yaml.safeLoad(uiYaml);

const uiYaml = `
order:
  - date_prepared
  - shipment_number
  - name_of_preparing_office
  - dest_office_name
  - origin_office
  - service_member_information
  - item_information
  - orders_information
  - pickup_information
  - destination_information
  - extra_address
  - scheduled_dates
  - remarks
  - other_move_information
  - certification_of_shipment_responsibilities
  - cert_in_lieu_of_signature
groups:
  origin_office:
    title: To (Responsible Origin Personal Property Shipping Office)
    fields:
      - origin_office_address_name
      - origin_office_address
  service_member_information:
    title: Member Or Employee Information
    fields:
      - service_member_first_name
      - service_member_middle_initial
      - service_member_last_name
      - service_member_rank
      - service_member_ssn
      - service_member_agency
  item_information:
    title: Request Action Be Taken to transport or store the following
    fields:
    - household_goods
    - mobile_home
  household_goods:
    title: Household Goods Unaccompanied Baggage Items No Of Containers Enter Quantity Estimate
    fields:
      - hhg_total_pounds
      - hhg_progear_pounds
      - hhg_valuable_items_cartons
  mobile_home:
    title: Mobile Home Information Enter Dimensions In Feet And Inches
    fields:
      - mobile_home_serial_number
      - mobile_home_length
      - mobile_home_width
      - mobile_home_height
      - mobile_home_type_expando
      - mobile_home_services_requested
  orders_information:
    title: This Shipment Storage Is Required Incident To The Following Change Of Station Orders
    fields:
      - station_orders_type
      - station_orders_issued_by
      - station_orders_new_assignment
      - station_orders_date
      - station_orders_number
      - station_orders_paragraph_number
      - station_orders_in_transit_telephone
      - in_transit_address
  pickup_information:
    title: Pickup Origin Information
    fields:
      - pickup_address
      - pickup_telephone
  destination_information:
    title: Destination information
    fields:
      - dest_address
      - agent_to_receive_hhg
  scheduled_dates:
    title: Scheduled Date for
    fields:
      - pack_scheduled_date
      - pickup_scheduled_date
      - delivery_scheduled_date
  other_move_information:
    title: I Certify That No Other Shipments And Or Nontemporary Storage Have Been Made Under These Orders Except As Indicated Below If None Indicate None
    fields:
      - other_move_to
      - other_move_net_pounds
      - other_move_progear_pounds
  certification_of_shipment_responsibilities:
    title: Certification Of Shipment Responsibilities Storage Conditions I Certify That I Have Read And Understand My Shipping Responsibilities And Storage Conditions Printed On The Back Side Of This Form
    fields:
      - service_member_signature
      - date_signed
      - contractor_address
      - contractor_name
  cert_in_lieu_of_signature:
    title: Certificate In Lieu Of Signature On This Form Is Required When Regulations So Authorize. Property Is Baggagehousehold goods, mobile home, and/or professional books, papers and equipment authorized to be shipped at government expense.
    fields:
      - nonavailability_of_signature_reason
      - certified_by_signature
      - title
    `;
