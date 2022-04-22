import React, { Fragment } from 'react';
import PropTypes from 'prop-types';

import { get } from 'lodash';

import { renderStatusIcon } from 'shared/utils';
import { formatDate } from 'utils/formatters';
import { PanelSwaggerField } from 'shared/EditablePanel';
import { WEIGHT_TICKET_SET_TYPE } from 'shared/constants';

const DocumentDetailDisplay = ({
  isExpenseDocument,
  isWeightTicketDocument,
  moveDocument,
  moveDocSchema,
  isStorageExpenseDocument,
}) => {
  const moveDocFieldProps = {
    values: moveDocument,
    schema: moveDocSchema,
  };
  const isWeightTicketTypeCarOrTrailer =
    isWeightTicketDocument &&
    (moveDocument.weight_ticket_set_type === WEIGHT_TICKET_SET_TYPE.CAR ||
      moveDocument.weight_ticket_set_type === WEIGHT_TICKET_SET_TYPE.CAR_TRAILER);
  const isWeightTicketTypeBoxTruck =
    isWeightTicketDocument && moveDocument.weight_ticket_set_type === WEIGHT_TICKET_SET_TYPE.BOX_TRUCK;
  const isWeightTicketTypeProGear =
    isWeightTicketDocument && moveDocument.weight_ticket_set_type === WEIGHT_TICKET_SET_TYPE.PRO_GEAR;
  return (
    <Fragment>
      <div>
        <h3 data-testid="panel-subhead">
          {renderStatusIcon(moveDocument.status)}
          {moveDocument.title}
        </h3>
        <p className="uploaded-at" data-testid="uploaded-at">
          Uploaded {formatDate(get(moveDocument, 'document.uploads.0.created_at'))}
        </p>
        <PanelSwaggerField
          data-testid="document-title"
          title="Document title"
          fieldName="title"
          required
          {...moveDocFieldProps}
        />
        <PanelSwaggerField
          data-testid="move-document-type"
          title="move-document-type"
          fieldName="move_document_type"
          required
          {...moveDocFieldProps}
        />
        {isExpenseDocument && moveDocument.moving_expense_type && (
          <PanelSwaggerField data-testid="moving-expense-type" fieldName="moving_expense_type" {...moveDocFieldProps} />
        )}
        {isExpenseDocument && get(moveDocument, 'requested_amount_cents') && (
          <PanelSwaggerField
            data-testid="requested-amount-cents"
            fieldName="requested_amount_cents"
            {...moveDocFieldProps}
          />
        )}
        {isExpenseDocument && get(moveDocument, 'payment_method') && (
          <PanelSwaggerField data-testid="payment-method" fieldName="payment_method" {...moveDocFieldProps} />
        )}
        {isWeightTicketDocument && (
          <>
            <PanelSwaggerField
              data-testid="weight-ticket-set-type"
              title="Weight ticket set type"
              fieldName="weight_ticket_set_type"
              required
              {...moveDocFieldProps}
            />

            {isWeightTicketTypeBoxTruck && (
              <PanelSwaggerField
                data-testid="vehicle-nickname"
                title="Vehicle nickname"
                fieldName="vehicle_nickname"
                required
                {...moveDocFieldProps}
              />
            )}
            {isWeightTicketTypeProGear && (
              <PanelSwaggerField title="Pro-gear type" fieldName="vehicle_nickname" required {...moveDocFieldProps} />
            )}
            {isWeightTicketTypeCarOrTrailer && (
              <>
                <PanelSwaggerField
                  data-testid="vehicle-make"
                  title="Vehicle make"
                  fieldName="vehicle_make"
                  required
                  {...moveDocFieldProps}
                />
                <PanelSwaggerField
                  data-testid="vehicle-model"
                  title="Vehicle model"
                  fieldName="vehicle_model"
                  required
                  {...moveDocFieldProps}
                />
              </>
            )}
            <PanelSwaggerField
              data-testid="empty-weight"
              title="Empty weight"
              fieldName="empty_weight"
              required
              {...moveDocFieldProps}
            />
            <PanelSwaggerField
              data-testid="full-weight"
              title="Full weight"
              fieldName="full_weight"
              required
              {...moveDocFieldProps}
            />
          </>
        )}
        {isStorageExpenseDocument && (
          <>
            <PanelSwaggerField
              data-testid="storage-start-date"
              title="Start date"
              fieldName="storage_start_date"
              required
              {...moveDocFieldProps}
            />
            <PanelSwaggerField
              data-testid="storage-end-date"
              title="End date"
              fieldName="storage_end_date"
              required
              {...moveDocFieldProps}
            />
          </>
        )}
        <PanelSwaggerField
          data-testid="status"
          title="Document status"
          fieldName="status"
          required
          {...moveDocFieldProps}
        />
        <PanelSwaggerField data-testid="notes" title="Notes" fieldName="notes" {...moveDocFieldProps} />
      </div>
    </Fragment>
  );
};

const { bool, object, shape, string, number, arrayOf } = PropTypes;

DocumentDetailDisplay.propTypes = {
  isExpenseDocument: bool.isRequired,
  isWeightTicketDocument: bool.isRequired,
  moveDocSchema: shape({
    properties: object.isRequired,
    required: arrayOf(string).isRequired,
    type: string.isRequired,
  }).isRequired,
  moveDocument: shape({
    document: shape({
      id: string.isRequired,
      service_member_id: string.isRequired,
      uploads: arrayOf(
        shape({
          byes: number,
          content_type: string.isRequired,
          created_at: string.isRequired,
          filename: string.isRequired,
          id: string.isRequired,
          update_at: string,
          url: string.isRequired,
        }),
      ).isRequired,
    }),
    id: string.isRequired,
    move_document_type: string.isRequired,
    move_id: string.isRequired,
    notes: string,
    personally_procured_move_id: string,
    status: string.isRequired,
    title: string.isRequired,
  }).isRequired,
};

export default DocumentDetailDisplay;
