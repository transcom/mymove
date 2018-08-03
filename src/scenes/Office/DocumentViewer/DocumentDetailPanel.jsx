import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get, omit } from 'lodash';
import { reduxForm, getFormValues, isValid, FormSection } from 'redux-form';

import editablePanel from '../editablePanel';
import { renderStatusIcon } from 'shared/utils';
import { formatDate } from 'shared/formatters';
// import { formatCents } from 'shared/formatters';
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
  const reimbursementFieldProps = {
    values: get(moveDoc, 'reimbursement', {}),
    schema: props.reimbursementSchema,
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
          get(moveDoc, 'reimbursement.requested_amount') && (
            <PanelSwaggerField
              fieldName="requested_amount"
              {...reimbursementFieldProps}
            />
          )}
        {isExpenseDocument &&
          get(moveDoc, 'reimbursement.method_of_receipt') && (
            <PanelSwaggerField
              title="Payment Method"
              fieldName="method_of_receipt"
              {...reimbursementFieldProps}
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
  const { formValues, moveDocSchema, reimbursementSchema } = props;

  // Should the amount be modified here?
  // let moveDoc = props.moveDocument;
  // if (get(moveDoc, 'reimbursement.requested_amount')) {
  //   let requested_amount = get(moveDoc, 'reimbursement.requested_amount')
  //   let requestedAmount = parseFloat(requested_amount) / 100;
  //   moveDoc.reimbursement.requested_amount = requestedAmount
  //   console.log("edit panel amount", moveDoc.reimbursement.requested_amount = requestedAmount)
  // }
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
            <ExpenseDocumentForm
              moveDocSchema={moveDocSchema}
              reimbursementSchema={reimbursementSchema}
            />
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
  // Don't pass 0-value reimbursement values to update endpoint
  if (
    get(moveDocument, 'reimbursement.id', '') ===
    '00000000-0000-0000-0000-000000000000'
  ) {
    moveDocument = omit(moveDocument, 'reimbursement');
  }
  // This modifies the moveDocument model, and so the value gets divided by 100 2 twice.
  // How
  // let money = {};
  // if (get(moveDocument, 'reimbursement.requested_amount')) {
  //   let requested_amount = get(moveDocument, 'reimbursement.requested_amount')
  //   let requestedAmount = formatCents(requested_amount);
  //   money['reimbursement'] = {}
  //   money['reimbursement']['requested_amount'] = requestedAmount;
  //   money['reimbursement']['method_of_receipt'] = moveDocument.reimbursement.method_of_receipt
  //   console.log("amount", money);
  // }
  return {
    // reduxForm
    // Populates the edit panel, NOT the display panel
    initialValues: {
      moveDocument: moveDocument,
      // money: money,
    },
    formValues: getFormValues(formName)(state),
    moveDocSchema: get(
      state,
      'swagger.spec.definitions.MoveDocumentPayload',
      {},
    ),
    reimbursementSchema: get(
      state,
      'swagger.spec.definitions.Reimbursement',
      {},
    ),
    hasError: false,
    errorMessage: state.office.error,
    isUpdating: false,
    moveDocument: moveDocument,

    // editablePanel
    formIsValid: isValid(formName)(state),
    getUpdateArgs: function() {
      let values = getFormValues(formName)(state);
      values.moveDocument.personally_procured_move_id = get(
        state.office,
        'officePPMs.0.id',
      );
      if (
        get(values.moveDocument, 'move_document_type', '') !== 'EXPENSE' &&
        get(values.moveDocument, 'reimbursement', false)
      ) {
        values.moveDocument = omit(values.moveDocument, 'reimbursement');
      }
      if (get(values.moveDocument, 'move_document_type', '') === 'EXPENSE') {
        values.moveDocument.reimbursement.requested_amount = parseFloat(
          values.moveDocument.reimbursement.requested_amount,
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
