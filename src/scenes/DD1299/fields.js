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
        label: Shipment Number
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
        label: Member Or Employee Information
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
            pattern: '^\\d{3}-\\d{2}-\\d{4}$'
            example: 555-555-5555
            label: SSN
          service_member_agency:
            type: string
            example: Air Force
            label: Agency
      item_information:
        type: group
        label: Request Action Be Taken to transport or store the following
        fields:
          household_goods:
            type: group
            label: HOUSEHOLD GOODS/UNACCOMPANIED BAGGAGE/ITEMS/NO. OF CONTAINERS (Enter quantity estimate)
            fields:
              hhg_total_pounds:
                type: number
                format: double
                example: 125.25
                label: Pounds
              hhg_progear_pounds:
                type: number
                format: double
                example: 35.11
                nullable: true
                label: Pounds of Professional Books, Papers, and Equipment (PBP&E) (Enter "none" if not applicable)
              hhg_valuable_items_cartons:
                type: integer
                example: 3
                nullable: true
                label: EXPENSIVE AND VALUABLE ITEMS (Number of cartons)
          mobile_home:
            type: group
            label: MOBILE HOME INFORMATION (Enter dimensions in feet and inches)
            fields:
              mobile_home_serial_number:
                type: string
                example: 45kljs98kljlkwj5
                nullable: true
                label: SERIAL NUMBER
              mobile_home_length:
                type: number
                example: 72
                nullable: true
                label: LENGTH
              mobile_home_width:
                type: number
                example: 15.4
                nullable: true
                label: WIDTH
              mobile_home_height:
                type: number
                example: 10
                nullable: true
                label: HEIGHT
              mobile_home_type_expando:
                type: string
                example: bathroom and shower unit
                nullable: true
                label: TYPE EXPANDO
              mobile_home_services_requested:
                type: string
                enum:
                - contents packed
                - mobile home blocked
                - mobile home unblocked
                - stored at origin
                - stored at destination
                nullable: true
      orders_information:
        type: group
        label: This Shipment Storage Is Required Incident To The Following Change Of Station Orders
        fields:
            station_orders_type:
              type: string
              label: Type Orders
              enum:
              - permanent
              - temporary
            station_orders_issued_by:
              type: string
              example: Sergeant Naan
              label: Issued by
            station_orders_new_assignment:
              type: string
              example: ACCOUNTING OPS
              label: New duty Assignment
            station_orders_date:
              type: string
              format: date
              example: 2018-03-15
              label: Date Of Orders
            station_orders_number:
              type: string
              example: 98374
              label: Orders Number
            station_orders_paragraph_number:
              type: string
              example: 5
              label: Paragraph Number
            station_orders_in_transit_telephone:
              type: string
              pattern: '^[2-9]\\d{2}-\\d{3}-\\d{4}$'
              example: 212-666-6666
              label: In Transit Telephone No
            in_transit_address:
              type: string
              label: In Transit Address
      pickup_information:
        type: group
        label: Pickup Origin Information
        fields:
          pickup_address:
            type: string
            label: Address (Street, Apartment Number, City, County, State, ZIP Code)(If a mobile home park, include mobile home court name)
          pickup_telephone:
            type: string
            pattern: '^[2-9]\\d{2}-\\d{3}-\\d{4}$'
            example: 212-555-5555
            label: Telephone Number (Include Area Code)
      destination_information:
        type: group
        label: Destination information
        fields:
          dest_address:
            type: string
            label: Address (Street, Apartment Number, City, County, State, ZIP Code)(If a mobile home park, include mobile home court name)
          agent_to_receive_hhg:
            type: string
            label: Agent Designated To Receive Property
      extra_address:
        type: string
        nullable: true
        label: Extra Pickup/Delivery Address (if applicable)
      scheduled_dates:
        type: group
        label: Scheduled Date for
        fields:
          pack_scheduled_date:
            type: string
            format: date
            example: 2018-03-08
            label: Pack
          pickup_scheduled_date:
            type: string
            format: date
            example: 2018-03-09
            label: Pickup
          delivery_scheduled_date:
            type: string
            format: date
            example: 2018-03-10
            label: Delivery
      remarks:
        type: string
        example: please be careful with my stuff
        label: Remarks
      other_move_information:
        type: group
        label: I Certify That No Other Shipments And Or Nontemporary Storage Have Been Made Under These Orders Except As Indicated Below If None Indicate None
        fields:
          other_move_from:
            type: string
            nullable: true
            label: From
          other_move_to:
            type: string
            nullable: true
            label: To
          other_move_net_pounds:
            type: number
            format: double
            example: 4.50
            nullable: true
            label: Net Pounds (Actual Or Estimated)
          other_move_progear_pounds:
            type: number
            format: float
            example: 99.09
            nullable: true
            label: Pounds Of PBP&E
      certification_of_shipment_responsibilities:
        type: group
        label: Certification Of Shipment Responsibilities Storage Conditions I Certify That I Have Read And Understand My Shipping Responsibilities And Storage Conditions Printed On The Back Side Of This Form
        fields:
          service_member_signature:
            type: string
            example: Focaccia Roll
            label: Signature Of Member/Employee
          date_signed:
            type: string
            format: date
            example: 2018-01-23
            label: Date Signed
          contractor_address:
            type: string
            label: Address Of Contractor (Street Suite No City State Zip Code)
          contractor_name:
            type: string
            example: Mayflower Transit
            nullable: true
            label: Name of Contractor
      cert_in_lieu_of_signature:
        type: group
        label: Certificate In Lieu Of Signature On This Form Is Required When Regulations So Authorize. Property Is Baggagehousehold goods, mobile home, and/or professional books, papers and equipment authorized to be shipped at government expense.
        fields:
          nonavailability_of_signature_reason:
            type: string
            example: service member not present
            nullable: true
            label: Reason For Nonavailability Of Signature
          certified_by_signature:
            type: string
            example: Sally Crumpet
            label: Certified By Signature
          title:
            type: string
            example: Colonel Crumpet
            label: Title
        `;
