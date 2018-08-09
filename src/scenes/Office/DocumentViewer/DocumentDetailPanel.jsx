import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get, omit } from 'lodash';
import { reduxForm, getFormValues, isValid, FormSection } from 'redux-form';

import editablePanel from '../editablePanel';
import { renderStatusIcon, convertDollarsToCents } from 'shared/utils';
import { formatDate, formatCents } from 'shared/formatters';
import { PanelSwaggerField, PanelField } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import ExpenseDocumentForm from './ExpenseDocumentForm';
import {
  selectMoveDocument,
  updateMoveDocument,
} from 'shared/Entities/modules/moveDocuments';
import { isMovingExpenseDocument } from 'shared/Entities/modules/movingExpenseDocuments';

import '../office.css';

const DocumentDetailDisplay = props => {
  const moveDoc = props.moveDocument;
  const isExpenseDocument = isMovingExpenseDocument(moveDoc);
  const moveDocFieldProps = {
    values: moveDoc,
    schema: props.moveDocSchema,
  };
  return (
    <React.Fragment>
      <div>
        <span className="panel-subhead">
          {renderStatusIcon(moveDoc.status)}
          {moveDoc.title}
        </span>
        <p className="uploaded-at">
          Uploaded {formatDate(get(moveDoc, 'document.uploads.0.created_at'))}
        </p>
        {moveDoc.title ? (
          <PanelSwaggerField fieldName="title" {...moveDocFieldProps} />
        ) : (
          <PanelField title="Document Title" className="missing">
            Missing
          </PanelField>
        )}
        {moveDoc.move_document_type ? (
          <PanelSwaggerField
            fieldName="move_document_type"
            {...moveDocFieldProps}
          />
        ) : (
          <PanelField title="Document Type" className="missing">
            Missing
          </PanelField>
        )}
        {isExpenseDocument &&
          moveDoc.moving_expense_type && (
            <PanelSwaggerField
              fieldName="moving_expense_type"
              {...moveDocFieldProps}
            />
          )}
        {isExpenseDocument &&
          get(moveDoc, 'requested_amount_cents') && (
            <PanelSwaggerField
              fieldName="requested_amount_cents"
              {...moveDocFieldProps}
            />
          )}
        {isExpenseDocument &&
          get(moveDoc, 'payment_method') && (
            <PanelSwaggerField
              fieldName="payment_method"
              {...moveDocFieldProps}
            />
          )}
        {moveDoc.status ? (
          <PanelSwaggerField fieldName="status" {...moveDocFieldProps} />
        ) : (
          <PanelField title="Document Status" className="missing">
            Missing
          </PanelField>
        )}
        {moveDoc.notes ? (
          <PanelSwaggerField fieldName="notes" {...moveDocFieldProps} />
        ) : (
          <PanelField title="Notes" />
        )}
      </div>
    </React.Fragment>
  );
};

const DocumentDetailEdit = props => {
  const { formValues, moveDocSchema } = props;
  const isExpenseDocument =
    get(formValues, 'moveDocument.move_document_type', '') === 'EXPENSE';
  return (
    <React.Fragment>
      <div>
        <FormSection name="moveDocument">
          <SwaggerField fieldName="title" swagger={moveDocSchema} required />
          <SwaggerField
            fieldName="move_document_type"
            swagger={moveDocSchema}
            required
          />
          {isExpenseDocument && (
            <ExpenseDocumentForm moveDocSchema={moveDocSchema} />
          )}
          <SwaggerField fieldName="status" swagger={moveDocSchema} required />
          <SwaggerField fieldName="notes" swagger={moveDocSchema} />
        </FormSection>
      </div>
    </React.Fragment>
  );
};

const formName = 'move_document_viewer';

let DocumentDetailPanel = editablePanel(
  DocumentDetailDisplay,
  DocumentDetailEdit,
);
DocumentDetailPanel = reduxForm({ form: formName })(DocumentDetailPanel);

function mapStateToProps(state, props) {
  const moveDocumentId = props.moveDocumentId;
  let moveDocument = selectMoveDocument(state, moveDocumentId);
  // Convert cents to collars - make a deep clone copy to not modify moveDocument itself
  let initialMoveDocument = JSON.parse(JSON.stringify(moveDocument));
  let requested_amount = get(initialMoveDocument, 'requested_amount_cents');
  if (requested_amount) {
    initialMoveDocument.requested_amount_cents = formatCents(requested_amount);
  }

  return {
    // reduxForm
    initialValues: {
      moveDocument: initialMoveDocument,
    },
    formValues: getFormValues(formName)(state),
    moveDocSchema: get(
      state,
      'swagger.spec.definitions.MoveDocumentPayload',
      {},
    ),
    hasError: false,
    errorMessage: state.office.error,
    isUpdating: false,
    moveDocument: moveDocument,

    // editablePanel
    formIsValid: isValid(formName)(state),
    getUpdateArgs: function() {
      // Make a copy of values to not modify moveDocument
      let values = JSON.parse(JSON.stringify(getFormValues(formName)(state)));
      values.moveDocument.personally_procured_move_id = get(
        state.office,
        'officePPMs.0.id',
      );
      if (
        get(values.moveDocument, 'move_document_type', '') !== 'EXPENSE' &&
        get(values.moveDocument, 'payment_method', false)
      ) {
        values.moveDocument = omit(values.moveDocument, [
          'payment_method',
          'requested_amount_cents',
        ]);
      }
      if (get(values.moveDocument, 'move_document_type', '') === 'EXPENSE') {
        values.moveDocument.requested_amount_cents = parseFloat(
          values.moveDocument.requested_amount_cents,
        );
      }
      let requested_amount = get(
        values.moveDocument,
        'requested_amount_cents',
        '',
      );
      if (requested_amount) {
        values.moveDocument.requested_amount_cents = convertDollarsToCents(
          requested_amount,
        );
      }
      return [
        get(state, 'office.officeMove.id'),
        get(moveDocument, 'id'),
        values.moveDocument,
      ];
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

export default connect(mapStateToProps, mapDispatchToProps)(
  DocumentDetailPanel,
);
