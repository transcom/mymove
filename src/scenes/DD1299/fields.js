import yaml from 'js-yaml'; //todo: yarn add

export const getFields = () => yaml.safeLoad(fieldsYaml);

const fieldsYaml = `
      date_prepared:
        type: string
        format: date
        example: 2018-01-03
        label: Date Prepared
      shipment_number:
        type: string
        example: 4550
        label: Shipment number
      name_of_preparing_office:
        type: string
        example: pumpernickel office
        label: Name of Preparing Office
      name_of_destination_personal_property_shipping_office:
        type: string
        example: rye office
        label: Name of Destination Personal Property Shipping Office
      to:
        type: group
        label: To (Responsible Origin Personal Property Shipping Office)
        fields:
          origin_office_address_name:
            type: string
            example: Office manager John Dough
            label:  Name
          origin_office_address:
            type: string
            example: '3450 Kneading Way, San Francisco, CA 94104'
            label: Address
      service_member_information:
        type: group
        label: Member or Employee Information
        fields:
          service_member_first_name:
            type: string
            example: Focaccia
            label: First Name
          service_member_middle_initial:
            example: L.
            nullable: true
            label: Middle Initial
          service_member_last_name:
            type: string
            example: Roll
            label: Last Name
          service_member_rank:
            type: string
            example: Commodore
            label: Rank/Grade
          service_member_ssn:
            type: string
            pattern: '^\d{3}-\d{2}-\d{4}$'
            example: 555-555-5555
            label: SSN
          service_member_agency:
            type: string
            example: Air Force
            label: Agency
      hhg_total_pounds:
        type: number
        format: double
        example: 125.25
      hhg_progear_pounds:
        type: number
        format: double
        example: 35.11
        nullable: true
      hhg_valuable_items_cartons:
        type: integer
        example: 3
        nullable: true
      mobile_home_serial_number:
        type: string
        example: 45kljs98kljlkwj5
        nullable: true
      mobile_home_length:
        type: number
        example: 72
        nullable: true
      mobile_home_width:
        type: number
        example: 15.4
        nullable: true
      mobile_home_height:
        type: number
        example: 10
        nullable: true
      mobile_home_type_expando:
        type: string
        example: bathroom and shower unit
        nullable: true
      mobile_home_services_requested:
        type: string
        enum:
        - contents packed
        - mobile home blocked
        - mobile home unblocked
        - stored at origin
        - stored at destination
        nullable: true
      station_orders_type:
        type: string
        enum:
        - permanent
        - temporary
      station_orders_issued_by:
        type: string
        example: Sergeant Naan
      station_orders_new_assignment:
        type: string
        example: ACCOUNTING OPS
      station_orders_date:
        type: string
        format: date
        example: 2018-03-15
      station_orders_number:
        type: string
        example: 98374
      station_orders_paragraph_number:
        type: string
        example: 5
      station_orders_in_transit_telephone:
        type: string
        pattern: '^[2-9]\d{2}-\d{3}-\d{4}$'
        example: 212-666-6666
      in_transit_address_number:
        type: string string
      in_transit_address_street:
        type: string
      in_transit_address_city:
        type: string
      in_transit_address_state:
        type: string
      in_transit_address_zip:
        type: string
      pickup_address_number:
        type: string
      pickup_address_street:
        type: string
      pickup_address_city:
        type: string
      pickup_address_county:
        type: string
      pickup_address_state:
        type: string
      pickup_address_zip:
        type: string
      pickup_address_mobile_court_name:
        type: string
        example: Winnebagel court
      pickup_telephone:
        type: string
        pattern: '^[2-9]\d{2}-\d{3}-\d{4}$'
        example: 212-555-5555
      dest_address_number:
        type: string
      dest_address_street:
        type: string
      dest_address_city:
        type: string
      dest_address_county:
        type: string
      dest_address_state:
        type: string
      dest_address_zip:
        type: string
      dest_address_mobile_court_name:
        type: String
        example: Carraway Court
      agent_to_receive_hhg:
        type: string
      extra_address_number:
        type: string
        nullable: true
      extra_address_street:
        type: string
        nullable: true
      extra_address_city:
        type: string
        nullable: true
      extra_address_county:
        type: string
        nullable: true
      extra_address_state:
        type: string
        nullable: true
      extra_address_zip:
        type: string
        nullable: true
      pack_scheduled_date:
        type: string
        format: date
        example: 2018-03-08
      pickup_scheduled_date:
        type: string
        format: date
        example: 2018-03-09
      delivery_scheduled_date:
        type: string
        format: date
        example: 2018-03-10
      remarks:
        type: string
        example: please be careful with my stuff
      other_move_from:
        type: string
        nullable: true
      other_move_to:
        type: string
        nullable: true
      other_move_net_pounds:
        type: number
        format: double
        example: 4.50
        nullable: true
      other_move_progear_pounds:
        type: number
        format: float
        example: 99.09
        nullable: true
      service_member_signature:
        type: string
        example: Focaccia Roll
      date_signed:
        type: string
        format: date
        example: 2018-01-23
      contractor_address_number:
        type: string
      contractor_address_street:
        type: string
      contractor_address_city:
        type: string
      contractor_address_state:
        type: string
      contractor_address_zip:
        type: string
      contractor_name:
        type: string
        example: Mayflower Transit
        nullable: true
      nonavailability_of_signature_reason:
        type: string
        example: service member not present
        nullable: true
      certified_by_signature:
        type: string
        example: Sally Crumpet
      title:
        type: string
        example: Colonel Crumpet
      created_at:
        type: string
        format: date-time
      updated_at:
        type: string
        format: date-time
        `;
