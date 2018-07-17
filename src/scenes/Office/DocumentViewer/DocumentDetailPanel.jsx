import React from 'react';
// import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { compact, get } from 'lodash';
import {
  reduxForm,
  getFormValues,
  isValid,
  FormSection,
  Field,
} from 'redux-form';

import editablePanel from '../editablePanel';
import { updateMoveDocument } from './ducks.js';
import { PanelSwaggerField, PanelField } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faClock from '@fortawesome/fontawesome-free-solid/faClock';

import '../office.css';

const DocumentDetailDisplay = props => {
  console.log('props in detail edit', props);
  const moveDoc = props.moveDocument;
  const moveDocFieldProps = {
    values: props.moveDocument,
    schema: props.moveDocSchema,
  };
  return (
    <React.Fragment>
      <div>
        <span className="panel-subhead">
          <FontAwesomeIcon
            aria-hidden
            className="icon approval-waiting"
            icon={faClock}
            title="Awaiting Review"
          />
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
  console.log('props in detail edit', props);
  return (
    <React.Fragment>
      <div />
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
      'swagger.spec.definitions.MoveDocumentPayload',
      {},
    ),

    hasError: false,
    errorMessage: state.office.error,
    isUpdating: false,
    moveDocuments: get(state, 'moveDocuments.moveDocuments', {}),

    moveDocument: get(state, 'moveDocuments.moveDocuments.0', {}),
    serviceMember: get(state, 'office.officeServiceMember', {}),
    move: get(state, 'office.officeMove', {}),

    // editablePanel
    formIsValid: isValid(formName)(state),
    // getUpdateArgs: function() {
    //   let values = getFormValues(formName)(state);
    //   return [
    //     get(state, 'moveDocument.id'),
    //     values.moveDocument,
    //     get(state, 'office.officeServiceMember.id'),
    //     values.serviceMember,
    //   ];
    // },
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
