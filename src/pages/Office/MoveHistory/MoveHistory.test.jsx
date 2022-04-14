import React from 'react';
import { render, screen } from '@testing-library/react';

import MoveHistory from './MoveHistory';

import { MockProviders } from 'testUtils';

const testMoveLocator = 'DVRS0N';

jest.mock('hooks/queries', () => ({
  useGHCGetMoveHistory: () => {
    return {
      isLoading: false,
      isError: false,
      queueResult: {
        totalCount: 2,
        data: [
          {
            action: 'UPDATE',
            actionTstampClk: '2022-03-09T15:33:38.623Z',
            actionTstampStm: '2022-03-09T15:33:38.622Z',
            actionTstampTx: '2022-03-09T15:33:38.579Z',
            changedValues: { available_to_prime_at: '2022-04-11T19:31:18.482947+00:00', status: 'APPROVED' },
            clientQuery:
              'UPDATE "moves" AS moves SET "available_to_prime_at" = $1, "billable_weights_reviewed_at" = $2, "cancel_reason" = $3, "contractor_id" = $4, "excess_weight_acknowledged_at" = $5, "excess_weight_qualified_at" = $6, "excess_weight_upload_id" = $7, "financial_review_flag" = $8, "financial_review_flag_set_at" = $9, "financial_review_remarks" = $10, "locator" = $11, "orders_id" = $12, "ppm_estimated_weight" = $13, "ppm_type" = $14, "reference_id" = $15, "selected_move_type" = $16, "service_counseling_completed_at" = $17, "show" = $18, "status" = $19, "submitted_at" = $20, "tio_remarks" = $21, "updated_at" = $22 WHERE moves.id = $23',
            eventName: 'updateMoveTaskOrderStatus',
            id: '6f5a4601-edde-4df1-aca9-0d3b58c11b59',
            objectId: '3efc84ca-d5a8-4f5f-b9b8-6deca1188e11',
            oldValues: {
              available_to_prime_at: '',
              billable_weights_reviewed_at: '',
              cancel_reason: '',
              contractor_id: '5db13bb4-6d29-4bdb-bc81-262f4513ecf6',
              excess_weight_acknowledged_at: '',
              excess_weight_qualified_at: '',
              excess_weight_upload_id: '',
              financial_review_flag: '',
              financial_review_flag_set_at: '',
              financial_review_remarks: '',
              id: '3efc84ca-d5a8-4f5f-b9b8-6deca1188e11',
              locator: 'YRHCWH',
              orders_id: 'e76cd7b9-689e-439e-a7d3-553b7eefe413',
              ppm_estimated_weight: '',
              ppm_type: 'PARTIAL',
              reference_id: '1895-7770',
              selected_move_type: 'PPM',
              service_counseling_completed_at: '2022-04-11T19:28:00.832092+00:00',
              show: '',
              status: 'SERVICE COUNSELING COMPLETED',
              submitted_at: '2021-09-11T19:36:00.871143',
              tio_remarks: '',
            },
            relId: 16592,
            schemaName: 'public',
            sessionUserEmail: 'too_tio_role@office.mil',
            sessionUserFirstName: 'Leo',
            sessionUserId: '9bda91d2-7a0c-4de1-ae02-b8cf8b4b858b',
            sessionUserLastName: 'Spaceman',
            sessionUserTelephone: '415-555-1212',
            tableName: 'moves',
            transactionId: 5577,
          },
          {
            action: 'UPDATE',
            actionTstampClk: '2022-03-08T18:28:58.271Z',
            actionTstampStm: '2022-03-08T18:28:58.220Z',
            actionTstampTx: '2022-03-08T18:28:58.152Z',
            changedValues: { approved_date: '2022-04-11', status: 'APPROVED' },
            clientQuery:
              'UPDATE "mto_shipments" AS mto_shipments SET "actual_pickup_date" = $1, "approved_date" = $2, "billable_weight_cap" = $3, "billable_weight_justification" = $4, "counselor_remarks" = $5, "customer_remarks" = $6, "deleted_at" = $7, "destination_address_id" = $8, "destination_address_type" = $9, "distance" = $10, "diversion" = $11, "first_available_delivery_date" = $12, "move_id" = $13, "nts_recorded_weight" = $14, "pickup_address_id" = $15, "prime_actual_weight" = $16, "prime_estimated_weight" = $17, "prime_estimated_weight_recorded_date" = $18, "rejection_reason" = $19, "requested_delivery_date" = $20, "requested_pickup_date" = $21, "required_delivery_date" = $22, "sac_type" = $23, "scheduled_pickup_date" = $24, "secondary_delivery_address_id" = $25, "secondary_pickup_address_id" = $26, "service_order_number" = $27, "shipment_type" = $28, "sit_days_allowance" = $29, "status" = $30, "storage_facility_id" = $31, "tac_type" = $32, "updated_at" = $33, "uses_external_vendor" = $34 WHERE mto_shipments.id = $35',
            eventName: 'approveShipment',
            id: '53605c7b-870b-406b-a5fc-e80ce8600eaa',
            objectId: '1d5e7cc1-0a8e-4d75-82e8-b3b05b9ddf38',
            oldValues: {
              actual_pickup_date: '2020-03-16',
              approved_date: '',
              billable_weight_cap: '',
              billable_weight_justification: '',
              counselor_remarks: '',
              customer_remarks: 'Please treat gently',
              days_in_storage: '',
              deleted_at: '',
              destination_address_id: '89550d5c-17fa-4fce-86d8-0f058a4bf236',
              destination_address_type: '',
              distance: '',
              diversion: '',
              first_available_delivery_date: '',
              id: '1d5e7cc1-0a8e-4d75-82e8-b3b05b9ddf38',
              move_id: '3efc84ca-d5a8-4f5f-b9b8-6deca1188e11',
              nts_recorded_weight: '',
              pickup_address_id: '3f80aee2-1252-4666-b9f8-42ff7691be41',
              prime_actual_weight: '',
              prime_estimated_weight: '',
              prime_estimated_weight_recorded_date: '',
              rejection_reason: '',
              requested_delivery_date: '2021-11-17',
              requested_pickup_date: '2021-11-10',
              required_delivery_date: '',
              sac_type: '',
              scheduled_pickup_date: '2020-03-16',
              secondary_delivery_address_id: '',
              secondary_pickup_address_id: '',
              service_order_number: '',
              shipment_type: 'HHG',
              sit_days_allowance: '',
              status: 'SUBMITTED',
              storage_facility_id: '',
              tac_type: '',
              uses_external_vendor: '',
            },
          },
        ],
      },
    };
  },
}));

describe('MoveHistory', () => {
  render(
    <MockProviders initialEntries={[`/moves/${testMoveLocator}/history`]}>
      <MoveHistory moveCode={testMoveLocator} />,
    </MockProviders>,
  );

  it('renders the different elements of the Move history tab', () => {
    expect(screen.getByText('Move history (2)')).toBeInTheDocument();
    expect(screen.getByRole('table')).toBeInTheDocument();

    expect(screen.getByTestId('move-history-date-time-0')).toHaveTextContent('09 Mar 22 15:33');
    expect(screen.getByTestId('move-history-event-0')).toHaveTextContent('Approved move');
    expect(screen.getByTestId('move-history-details-0')).toBeInTheDocument();
    expect(screen.getByTestId('move-history-modified-by-0')).toBeInTheDocument();

    expect(screen.getByTestId('move-history-date-time-1')).toHaveTextContent('08 Mar 22 18:28');
    expect(screen.getByTestId('move-history-event-1')).toHaveTextContent('Approved shipment');
    expect(screen.getByTestId('move-history-details-1')).toBeInTheDocument();
    expect(screen.getByTestId('move-history-modified-by-1')).toBeInTheDocument();

    expect(screen.getByTestId('pagination')).toBeInTheDocument();
  });
});
