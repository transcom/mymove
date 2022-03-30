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
            changedValues: [
              { columnName: 'status', columnValue: 'APPROVED' },
              { columnName: 'updated_at', columnValue: '2022-03-08T21:33:38.596072' },
            ],
            clientQuery:
              'UPDATE "moves" AS moves SET "available_to_prime_at" = $1, "billable_weights_reviewed_at" = $2, "cancel_reason" = $3, "contractor_id" = $4, "excess_weight_acknowledged_at" = $5, "excess_weight_qualified_at" = $6, "excess_weight_upload_id" = $7, "financial_review_flag" = $8, "financial_review_flag_set_at" = $9, "financial_review_remarks" = $10, "locator" = $11, "orders_id" = $12, "ppm_estimated_weight" = $13, "ppm_type" = $14, "reference_id" = $15, "selected_move_type" = $16, "service_counseling_completed_at" = $17, "show" = $18, "status" = $19, "submitted_at" = $20, "tio_remarks" = $21, "updated_at" = $22 WHERE moves.id = $23',
            eventName: 'updateOrder',
            id: '7ce7c1ac-a1d7-4caf-858c-09674a00f273',
            objectId: 'abe92574-53a8-4026-a75c-45ff9eea9bc6',
            oldValues: [
              { columnName: 'show', columnValue: 'true' },
              { columnName: 'created_at', columnValue: '2022-02-28T20:31:35.901803' },
              { columnName: 'updated_at', columnValue: '2022-02-28T20:31:35.901803' },
              { columnName: 'id', columnValue: 'abe92574-53a8-4026-a75c-45ff9eea9bc6' },
              { columnName: 'status', columnValue: 'APPROVALS REQUESTED' },
              { columnName: 'locator', columnValue: 'DVRS0N' },
              { columnName: 'orders_id', columnValue: '18b3725c-529c-4add-811b-7345ece8847f' },
              { columnName: 'ppm_type', columnValue: 'PARTIAL' },
              { columnName: 'reference_id', columnValue: '5037-3728' },
              { columnName: 'contractor_id', columnValue: '5db13bb4-6d29-4bdb-bc81-262f4513ecf6' },
              { columnName: 'selected_move_type', columnValue: 'PPM' },
              { columnName: 'available_to_prime_at', columnValue: '2022-02-28T20:31:35.807801+00:00' },
              { columnName: 'financial_review_flag', columnValue: 'false' },
            ],
            relId: 16592,
            schemaName: 'public',
            tableName: 'moves',
            transactionId: 26993,
          },
          {
            action: 'UPDATE',
            actionTstampClk: '2022-03-08T18:28:58.271Z',
            actionTstampStm: '2022-03-08T18:28:58.220Z',
            actionTstampTx: '2022-03-08T18:28:58.152Z',
            changedValues: [
              { columnName: 'new_duty_location_id', columnValue: '2d5ada83-e09a-47f8-8de6-83ec51694a86' },
              { columnName: 'origin_duty_location_id', columnValue: 'fe8b5aa0-a7af-49fa-a705-ee4035ac546b' },
            ],
            clientQuery:
              'UPDATE "orders" AS orders SET "amended_orders_acknowledged_at" = $1, "department_indicator" = $2, "entitlement_id" = $3, "grade" = $4, "has_dependents" = $5, "issue_date" = $6, "new_duty_location_id" = $7, "nts_sac" = $8, "nts_tac" = $9, "orders_number" = $10, "orders_type" = $11, "orders_type_detail" = $12, "origin_duty_location_id" = $13, "report_by_date" = $14, "sac" = $15, "service_member_id" = $16, "spouse_has_pro_gear" = $17, "status" = $18, "tac" = $19, "updated_at" = $20, "uploaded_amended_orders_id" = $21, "uploaded_orders_id" = $22 WHERE orders.id = $23',
            id: '34752aeb-f658-4afa-b1c0-dcdcb5fb3a73',
            objectId: '18b3725c-529c-4add-811b-7345ece8847f',
            oldValues: [
              { columnName: 'status', columnValue: 'DRAFT' },
              { columnName: 'issue_date', columnValue: '2018-03-15' },
              { columnName: 'orders_number', columnValue: 'ORDER3' },
              { columnName: 'origin_duty_location_id', columnValue: 'fe8b5aa0-a7af-49fa-a705-ee4035ac546b' },
              { columnName: 'id', columnValue: '18b3725c-529c-4add-811b-7345ece8847f' },
              { columnName: 'orders_type', columnValue: 'PERMANENT_CHANGE_OF_STATION' },
              { columnName: 'has_dependents', columnValue: 'false' },
              { columnName: 'service_member_id', columnValue: '64109c3e-68cc-4069-ae5c-7f2460733e7c' },
              { columnName: 'uploaded_orders_id', columnValue: 'a18801a0-93f6-4464-bc9d-905bf6109490' },
              { columnName: 'new_duty_location_id', columnValue: '2d5ada83-e09a-47f8-8de6-83ec51694a86' },
              { columnName: 'department_indicator', columnValue: 'AIR_FORCE' },
              { columnName: 'tac', columnValue: 'F8E1' },
              { columnName: 'grade', columnValue: 'E_1' },
              { columnName: 'updated_at', columnValue: '2022-02-28T20:31:35.893474' },
              { columnName: 'entitlement_id', columnValue: '8aff2987-3766-45e1-99ee-7ee4cced4cc3' },
              { columnName: 'report_by_date', columnValue: '2018-08-01' },
              { columnName: 'orders_type_detail', columnValue: 'HHG_PERMITTED' },
              { columnName: 'created_at', columnValue: '2022-02-28T20:31:35.893474' },
              { columnName: 'spouse_has_pro_gear', columnValue: 'false' },
            ],
            relId: 16879,
            schemaName: 'public',
            tableName: 'orders',
            transactionId: 26950,
          },
        ],
        id: 'abe92574-53a8-4026-a75c-45ff9eea9bc6',
        locator: testMoveLocator,
        referenceId: '5037-3728',
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
    expect(screen.getByTestId('move-history-event-0')).toHaveTextContent('Updated orders');
    expect(screen.getByTestId('move-history-details-0')).toBeInTheDocument();
    expect(screen.getByTestId('move-history-modified-by-0')).toBeInTheDocument();

    expect(screen.getByTestId('move-history-date-time-1')).toHaveTextContent('08 Mar 22 18:28');
    expect(screen.getByTestId('move-history-event-1')).toBeInTheDocument();
    expect(screen.getByTestId('move-history-details-1')).toBeInTheDocument();
    expect(screen.getByTestId('move-history-modified-by-1')).toBeInTheDocument();
  });
});
