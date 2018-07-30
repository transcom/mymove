import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get } from 'lodash';
import { reduxForm, getFormValues, isValid, FormSection } from 'redux-form';

import editablePanel from '../editablePanel';
import { renderStatusIcon } from 'shared/utils';
import { formatDate } from 'shared/formatters';
import { PanelSwaggerField, PanelField } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import {
  selectMoveDocument,
  updateMoveDocument,
} from 'shared/Entities/modules/moveDocuments';

import '../office.css';

const DocumentDetailDisplay = props => {
  const moveDoc = props.moveDocument;
  const schema = moveDoc.moving_expense_type
    ? props.movingExpenseSchema
    : props.moveDocSchema;
  const moveDocFieldProps = {
    values: props.moveDocument,
    schema: schema,
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
        {moveDoc.moving_expense_type && (
          <PanelSwaggerField
            fieldName="moving_expense_type"
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
          <PanelField title="Notes" className="missing">
            Missing
          </PanelField>
        )}
      </div>
    </React.Fragment>
  );
};

const DocumentDetailEdit = props => {
  const schema = props.moveDocSchema;

  return (
    <React.Fragment>
      <div>
        <FormSection name="moveDocument">
          <SwaggerField fieldName="title" swagger={schema} required />
          <SwaggerField
            fieldName="move_document_type"
            swagger={schema}
            required
          />
          <SwaggerField fieldName="status" swagger={schema} required />
          <SwaggerField fieldName="notes" swagger={schema} />
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
  const moveDocument = selectMoveDocument(state, moveDocumentId);

  return {
    // reduxForm
    initialValues: {
      moveDocument: moveDocument,
    },

    moveDocSchema: get(
      state,
      'swagger.spec.definitions.UpdateGenericMoveDocumentPayload',
      {},
    ),
    movingExpenseSchema: get(
      state,
      'swagger.spec.definitions.UpdateMovingExpenseDocumentPayload',
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
