import React from 'react';
// import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get } from 'lodash';
import { reduxForm, getFormValues, isValid, FormSection } from 'redux-form';

import editablePanel from '../editablePanel';
import { updateMoveDocumentInfo } from './ducks.js';
import { renderStatusIcon } from 'shared/utils';
import { PanelSwaggerField, PanelField } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import '../office.css';

const DocumentDetailDisplay = props => {
  const moveDoc = props.moveDocument;
  const moveDocFieldProps = {
    values: props.moveDocument,
    schema: props.moveDocSchema,
  };
  return (
    <React.Fragment>
      <div>
        <span className="panel-subhead">
          {renderStatusIcon(moveDoc.status)}
          {moveDoc.title}
        </span>
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
          <SwaggerField fieldName="move_document_type" swagger={schema} />
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

function mapStateToProps(state) {
  return {
    // reduxForm
    initialValues: {
      // TODO: Update to use move doc retrieved from UI store
      moveDocument: get(state, 'moveDocuments.moveDocuments.0', {}),
    },

    moveDocSchema: get(
      state,
      'swagger.spec.definitions.UpdateMoveDocumentPayload',
      {},
    ),

    hasError: false,
    errorMessage: state.office.error,
    isUpdating: false,
    // TODO: Get the appropriate move doc
    moveDocument: get(state, 'moveDocuments.moveDocuments.0', {}),

    // editablePanel
    formIsValid: isValid(formName)(state),
    getUpdateArgs: function() {
      let values = getFormValues(formName)(state);
      return [
        get(state, 'office.officeMove.id'),
        // TODO: get real movedoc
        get(state, 'moveDocuments.moveDocuments.0.id'),
        values.moveDocument,
      ];
    },
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: updateMoveDocumentInfo,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(
  DocumentDetailPanel,
);
