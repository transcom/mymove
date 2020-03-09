import React, { Fragment } from 'react';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get, omit, cloneDeep, isEmpty } from 'lodash';
import { reduxForm, getFormValues, FormSection } from 'redux-form';

import { renderStatusIcon, convertDollarsToCents } from 'shared/utils';
import { formatDate, formatCents } from 'shared/formatters';
import { PanelSwaggerField, editablePanelify } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { selectMoveDocument, updateMoveDocument } from 'shared/Entities/modules/moveDocuments';
import { selectPPMForMove } from 'shared/Entities/modules/ppms';
import { isMovingExpenseDocument } from 'shared/Entities/modules/movingExpenseDocuments';
import { MOVE_DOC_TYPE, WEIGHT_TICKET_SET_TYPE } from 'shared/constants';

import ExpenseDocumentForm from 'scenes/Office/DocumentViewer/ExpenseDocumentForm';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';

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
        <h3>
          {renderStatusIcon(moveDocument.status)}
          {moveDocument.title}
        </h3>
        <p className="uploaded-at">Uploaded {formatDate(get(moveDocument, 'document.uploads.0.created_at'))}</p>
        <PanelSwaggerField title="Document title" fieldName="title" required {...moveDocFieldProps} />

        <PanelSwaggerField title="Document type" fieldName="move_document_type" required {...moveDocFieldProps} />
        {isExpenseDocument && moveDocument.moving_expense_type && (
          <PanelSwaggerField fieldName="moving_expense_type" {...moveDocFieldProps} />
        )}
        {isExpenseDocument && get(moveDocument, 'requested_amount_cents') && (
          <PanelSwaggerField fieldName="requested_amount_cents" {...moveDocFieldProps} />
        )}
        {isExpenseDocument && get(moveDocument, 'payment_method') && (
          <PanelSwaggerField fieldName="payment_method" {...moveDocFieldProps} />
        )}
        {isWeightTicketDocument && (
          <>
            <PanelSwaggerField
              title="Weight ticket set type"
              fieldName="weight_ticket_set_type"
              required
              {...moveDocFieldProps}
            />

            {isWeightTicketTypeBoxTruck && (
              <PanelSwaggerField
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
                <PanelSwaggerField title="Vehicle make" fieldName="vehicle_make" required {...moveDocFieldProps} />
                <PanelSwaggerField title="Vehicle model" fieldName="vehicle_model" required {...moveDocFieldProps} />
              </>
            )}
            <PanelSwaggerField title="Empty weight" fieldName="empty_weight" required {...moveDocFieldProps} />
            <PanelSwaggerField title="Full weight" fieldName="full_weight" required {...moveDocFieldProps} />
          </>
        )}
        {isStorageExpenseDocument && (
          <>
            <PanelSwaggerField title="Start date" fieldName="storage_start_date" required {...moveDocFieldProps} />
            <PanelSwaggerField title="End date" fieldName="storage_end_date" required {...moveDocFieldProps} />
          </>
        )}
        <PanelSwaggerField title="Document status" fieldName="status" required {...moveDocFieldProps} />
        <PanelSwaggerField title="Notes" fieldName="notes" {...moveDocFieldProps} />
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

const DocumentDetailEdit = ({ formValues, moveDocSchema }) => {
  const isExpenseDocument = get(formValues.moveDocument, 'move_document_type') === MOVE_DOC_TYPE.EXPENSE;
  const isWeightTicketDocument = get(formValues.moveDocument, 'move_document_type') === MOVE_DOC_TYPE.WEIGHT_TICKET_SET;
  const isStorageExpenseDocument =
    get(formValues.moveDocument, 'move_document_type') === 'EXPENSE' &&
    get(formValues.moveDocument, 'moving_expense_type') === 'STORAGE';
  const isWeightTicketTypeCarOrTrailer =
    isWeightTicketDocument &&
    (formValues.moveDocument.weight_ticket_set_type === WEIGHT_TICKET_SET_TYPE.CAR ||
      formValues.moveDocument.weight_ticket_set_type === WEIGHT_TICKET_SET_TYPE.CAR_TRAILER);
  const isWeightTicketTypeBoxTruck =
    isWeightTicketDocument && formValues.moveDocument.weight_ticket_set_type === WEIGHT_TICKET_SET_TYPE.BOX_TRUCK;
  const isWeightTicketTypeProGear =
    isWeightTicketDocument && formValues.moveDocument.weight_ticket_set_type === WEIGHT_TICKET_SET_TYPE.PRO_GEAR;

  return isEmpty(formValues.moveDocument) ? (
    <LoadingPlaceholder />
  ) : (
    <Fragment>
      <div>
        <FormSection name="moveDocument">
          <SwaggerField fieldName="title" swagger={moveDocSchema} required />
          <SwaggerField fieldName="move_document_type" swagger={moveDocSchema} required />
          {isExpenseDocument && <ExpenseDocumentForm moveDocSchema={moveDocSchema} />}
          {isWeightTicketDocument && (
            <>
              <div className="field-with-units">
                <SwaggerField
                  className="short-field"
                  fieldName="weight_ticket_set_type"
                  swagger={moveDocSchema}
                  required
                />
              </div>
              {isWeightTicketTypeBoxTruck && (
                <div className="field-with-units">
                  <SwaggerField className="short-field" fieldName="vehicle_nickname" swagger={moveDocSchema} required />
                </div>
              )}
              {isWeightTicketTypeProGear && (
                <div className="field-with-units">
                  <SwaggerField
                    className="short-field"
                    fieldName="vehicle_nickname"
                    title="Pro-gear type (ex. 'My Pro-gear', 'Spouse Pro-Gear', 'Both')"
                    swagger={moveDocSchema}
                    required
                  />
                </div>
              )}
              {isWeightTicketTypeCarOrTrailer && (
                <>
                  <div className="field-with-units">
                    <SwaggerField className="short-field" fieldName="vehicle_make" swagger={moveDocSchema} required />
                  </div>
                  <div className="field-with-units">
                    <SwaggerField className="short-field" fieldName="vehicle_model" swagger={moveDocSchema} required />
                  </div>
                </>
              )}

              <div className="field-with-units">
                <SwaggerField className="short-field" fieldName="empty_weight" swagger={moveDocSchema} required /> lbs
              </div>
              <div className="field-with-units">
                <SwaggerField className="short-field" fieldName="full_weight" swagger={moveDocSchema} required /> lbs
              </div>
            </>
          )}
          {isStorageExpenseDocument && (
            <>
              <SwaggerField title="Start date" fieldName="storage_start_date" required swagger={moveDocSchema} />
              <SwaggerField title="End date" fieldName="storage_end_date" required swagger={moveDocSchema} />
            </>
          )}
          <SwaggerField fieldName="status" swagger={moveDocSchema} required />
          <SwaggerField fieldName="notes" swagger={moveDocSchema} />
        </FormSection>
      </div>
    </Fragment>
  );
};

DocumentDetailEdit.propTypes = {
  isExpenseDocument: bool.isRequired,
  isWeightTicketDocument: bool.isRequired,
  moveDocSchema: shape({
    properties: object.isRequired,
    required: arrayOf(string).isRequired,
    type: string.isRequired,
  }).isRequired,
};

const formName = 'move_document_viewer';

let DocumentDetailPanel = editablePanelify(DocumentDetailDisplay, DocumentDetailEdit);

DocumentDetailPanel = reduxForm({ form: formName })(DocumentDetailPanel);

function mapStateToProps(state, props) {
  const { moveId, moveDocumentId } = props;
  const moveDocument = selectMoveDocument(state, moveDocumentId);
  const isExpenseDocument = isMovingExpenseDocument(moveDocument);
  const isWeightTicketDocument = get(moveDocument, 'move_document_type') === 'WEIGHT_TICKET_SET';
  const isStorageExpenseDocument =
    get(moveDocument, 'move_document_type') === 'EXPENSE' && get(moveDocument, 'moving_expense_type') === 'STORAGE';
  // Convert cents to collars - make a deep clone copy to not modify moveDocument itself
  const initialMoveDocument = cloneDeep(moveDocument);
  const requested_amount = get(initialMoveDocument, 'requested_amount_cents');
  if (requested_amount) {
    initialMoveDocument.requested_amount_cents = formatCents(requested_amount);
  }

  return {
    // reduxForm
    initialValues: {
      moveDocument: initialMoveDocument,
    },
    isExpenseDocument,
    isWeightTicketDocument,
    isStorageExpenseDocument,
    formValues: getFormValues(formName)(state),
    moveDocSchema: get(state, 'swaggerInternal.spec.definitions.MoveDocumentPayload', {}),
    hasError: false,
    isUpdating: false,
    moveDocument,

    // editablePanelify
    getUpdateArgs: function() {
      // Make a copy of values to not modify moveDocument
      let values = cloneDeep(getFormValues(formName)(state));
      values.moveDocument.personally_procured_move_id = selectPPMForMove(state, props.moveId).id;
      if (
        get(values.moveDocument, 'move_document_type', '') !== 'EXPENSE' &&
        get(values.moveDocument, 'payment_method', false)
      ) {
        values.moveDocument = omit(values.moveDocument, ['payment_method', 'requested_amount_cents']);
      }
      if (get(values.moveDocument, 'move_document_type', '') === 'EXPENSE') {
        values.moveDocument.requested_amount_cents = convertDollarsToCents(values.moveDocument.requested_amount_cents);
      }
      return [moveId, moveDocumentId, values.moveDocument];
    },
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: updateMoveDocument,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(DocumentDetailPanel);
