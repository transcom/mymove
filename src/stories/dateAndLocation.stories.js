import React from 'react';

import { storiesOf } from '@storybook/react';
import { action } from '@storybook/addon-actions';
import { withKnobs, text } from '@storybook/addon-knobs';

import { Provider } from 'react-redux';
import { HashRouter as Router } from 'react-router-dom';
import { createStore } from 'redux';
import { appReducer } from '../appReducer';

import DateAndLocation from '../scenes/Moves/Ppm/DateAndLocation';

const initstore = {
  swaggerInternal: {
    spec: {
      definitions: {
        UpdatePersonallyProcuredMovePayload: {
          type: 'object',
          properties: {
            size: {
              type: 'string',
              'x-nullable': true,
              title: 'Size',
              enum: ['S', 'M', 'L'],
              $$ref: '/internal/swagger.yaml#/definitions/TShirtSize',
            },
            original_move_date: {
              type: 'string',
              format: 'date',
              example: '2018-04-26',
              title: 'When do you plan to move?',
              'x-nullable': true,
              'x-always-required': true,
            },
            actual_move_date: {
              type: 'string',
              example: '2018-04-26',
              format: 'date',
              title: 'When did you actually move?',
              'x-nullable': true,
            },
            pickup_postal_code: {
              type: 'string',
              format: 'zip',
              title: 'ZIP/Postal Code',
              example: '90210',
              pattern: '^(\\d{5}([\\-]\\d{4})?)$',
              'x-nullable': true,
              'x-always-required': true,
            },
            has_additional_postal_code: {
              type: 'boolean',
              'x-nullable': true,
              title: 'Do you have stuff at another pickup location?',
              'x-always-required': false,
            },
            additional_pickup_postal_code: {
              type: 'string',
              format: 'zip',
              title: 'ZIP/Postal Code',
              example: '90210',
              pattern: '^(\\d{5}([\\-]\\d{4})?)$',
              'x-nullable': true,
            },
            destination_postal_code: {
              type: 'string',
              format: 'zip',
              title: 'ZIP/Postal Code',
              example: '90210',
              pattern: '^(\\d{5}([\\-]\\d{4})?)$',
              'x-nullable': true,
              'x-always-required': true,
            },
            has_sit: {
              type: 'boolean',
              'x-nullable': true,
              title: 'Are you going to put your stuff in temporary storage before moving into your new home?',
              'x-always-required': false,
            },
            days_in_storage: {
              type: 'integer',
              title: 'How many days do you plan to put your stuff in storage?',
              minimum: 0,
              maximum: 90,
              'x-nullable': true,
            },
            total_sit_cost: {
              type: 'integer',
              title: 'How much does your storage cost?',
              minimum: 0,
              'x-nullable': true,
            },
            estimated_storage_reimbursement: {
              type: 'string',
              title: 'Estimated Storage Reimbursement',
              'x-nullable': true,
            },
            weight_estimate: {
              type: 'integer',
              minimum: 0,
              title: 'Weight Estimate',
              'x-nullable': true,
            },
            net_weight: {
              type: 'integer',
              minimum: 1,
              title: 'Net Weight',
              'x-nullable': true,
            },
            has_requested_advance: {
              type: 'boolean',
              default: false,
              title: 'Would you like an advance of up to 60% of your PPM incentive?',
            },
            advance: {
              type: 'object',
              'x-nullable': true,
              properties: {
                id: {
                  type: 'string',
                  format: 'uuid',
                  example: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
                },
                requested_amount: {
                  type: 'integer',
                  format: 'cents',
                  minimum: 1,
                  title: 'Requested Amount',
                  description: 'unit is cents',
                },
                method_of_receipt: {
                  'x-nullable': true,
                  type: 'string',
                  title: 'Method of Receipt',
                  enum: ['MIL_PAY', 'OTHER_DD', 'GTCC'],
                  'x-display-value': {
                    MIL_PAY: 'MilPay',
                    OTHER_DD: 'Other account',
                    GTCC: 'GTCC',
                  },
                  $$ref: '/internal/swagger.yaml#/definitions/MethodOfReceipt',
                },
                status: {
                  'x-nullable': true,
                  type: 'string',
                  title: 'Reimbursement',
                  enum: ['DRAFT', 'REQUESTED', 'APPROVED', 'REJECTED', 'PAID'],
                  $$ref: '/internal/swagger.yaml#/definitions/ReimbursementStatus',
                },
                requested_date: {
                  'x-nullable': true,
                  type: 'string',
                  example: '2018-04-26',
                  format: 'date',
                  title: 'Requested Date',
                },
              },
              required: ['requested_amount', 'method_of_receipt'],
              $$ref: '/internal/swagger.yaml#/definitions/Reimbursement',
            },
            advance_worksheet: {
              type: 'object',
              properties: {
                id: {
                  type: 'string',
                  format: 'uuid',
                  example: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
                },
                service_member_id: {
                  type: 'string',
                  format: 'uuid',
                  title: 'The service member this document belongs to',
                },
                uploads: {
                  type: 'array',
                  items: {
                    type: 'object',
                    properties: {
                      id: {
                        type: 'string',
                        format: 'uuid',
                        example: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
                      },
                      url: {
                        type: 'string',
                        format: 'uri',
                        example: 'https://uploads.domain.test/dir/c56a4180-65aa-42ec-a945-5fd21dec0538',
                      },
                      filename: {
                        type: 'string',
                        example: 'filename.pdf',
                      },
                      content_type: {
                        type: 'string',
                        format: 'mime-type',
                        example: 'application/pdf',
                      },
                      bytes: {
                        type: 'integer',
                      },
                      created_at: {
                        type: 'string',
                        format: 'date-time',
                      },
                      updated_at: {
                        type: 'string',
                        format: 'date-time',
                      },
                    },
                    required: ['id', 'url', 'filename', 'content_type', 'bytes', 'created_at', 'updated_at'],
                    $$ref: '/internal/swagger.yaml#/definitions/UploadPayload',
                  },
                },
              },
              required: ['id', 'service_member_id', 'uploads'],
              $$ref: '/internal/swagger.yaml#/definitions/DocumentPayload',
            },
          },
        },
      },
    },
  },
  orders: {
    currentOrders: {
      new_duty_station: {
        name: text('new duty station', 'McCllellan Air Force Base'),
        address: {
          postal_code: 11209,
        },
      },
    },
  },
  serviceMember: {
    currentServiceMember: {
      current_station: {
        address: {
          postal_code: 90210,
        },
      },
    },
  },
  ppp_date_and_location: {
    values: {
      pickup_postal_code: '90210',
      origin_duty_station_zip: '50309',
      origin_move_date: '2019-09-12',
    },
  },
};

const store = createStore(appReducer, initstore);

const setupKnobs = () => {
  initstore.orders.currentOrders.new_duty_station.name = text('New Duty Station', 'McClellan Air Force Base');
  initstore.orders.currentOrders.new_duty_station.address.postal_code = text('New Duty Station ZIP', '95652');
};

storiesOf('scenes/Moves/Ppm/DateAndLocation', module)
  .addDecorator(withKnobs)
  .addDecorator(story => (
    <Provider store={store}>
      <Router>{story()}</Router>
    </Provider>
  ))
  .addDecorator(story => <div className="my-move site">{story()}</div>)
  .add('sample', () => (
    <div>
      {setupKnobs()}
      <DateAndLocation
        error={null}
        pages={['DateLoc']}
        pageKey="DateLoc"
        windowWidth="1024"
        match={{ params: { moveId: '42' } }}
        createOrUpdatePpm={action('createOrUpdatePpm')}
      />
    </div>
  ));
