import { get, map } from 'lodash';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import { bindActionCreators } from 'redux';
import { getFormValues, reduxForm } from 'redux-form';
import PropTypes from 'prop-types';

import Alert from 'shared/Alert';
import { createMoveDocument } from 'shared/Entities/modules/moveDocuments';
import { createMovingExpenseDocument } from 'shared/Entities/modules/movingExpenseDocuments';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import Uploader from 'shared/Uploader';
import ExpenseDocumentForm from './ExpenseDocumentForm';

import './DocumentUploader.css';

const moveDocumentFormName = 'move_document_upload';

export class DocumentUploader extends Component {
  constructor(props) {
    super(props);

    this.state = {
      newUploads: [],
      uploaderIsIdle: true,
      moveDocumentCreateError: null,
    };

    this.onChange = this.onChange.bind(this);
    this.onSubmit = this.onSubmit.bind(this);
    this.onAddFile = this.onAddFile.bind(this);
  }

  onSubmit() {
    const { formValues, currentPpm, moveId, reset } = this.props;
    const uploadIds = map(this.state.newUploads, 'id');
    this.setState({
      moveDocumentCreateError: null,
    });
    if (get(formValues, 'movingExpenseDocument', false)) {
      formValues.reimbursement.requested_amount = parseFloat(
        formValues.reimbursement.requested_amount,
      );
      this.props
        .createMovingExpenseDocument(
          moveId,
          currentPpm.id,
          uploadIds,
          formValues.title,
          formValues.movingExpenseDocument.moving_expense_type,
          formValues.move_document_type,
          formValues.reimbursement,
          formValues.notes,
        )
        .then(() => {
          reset();
          this.uploader.clearFiles();
        })
        .catch(err => {
          this.setState({
            moveDocumentCreateError: err,
          });
        });
    }
    if (get(formValues, 'movingExpenseDocument', false) === false) {
      this.props
        .createMoveDocument(
          moveId,
          currentPpm.id,
          uploadIds,
          formValues.title,
          formValues.move_document_type,
          formValues.notes,
        )
        .then(() => {
          reset();
          this.uploader.clearFiles();
        })
        .catch(err => {
          this.setState({
            moveDocumentCreateError: err,
          });
        });
    }
    //todo: we don't want to do this until the details view is working,
    // we may not want to do it at all if users are going to upload several documents at a time
    // .then(response => {
    //   if (!response.error) {
    //     const moveDocumentId = response.payload.id;
    //     this.props.push(`/moves/${moveId}/documents/${moveDocumentId}`);
    //   }
    // });
  }

  onAddFile() {
    this.setState({
      uploaderIsIdle: false,
    });
  }

  onChange(newUploads, uploaderIsIdle) {
    this.setState({
      newUploads,
      uploaderIsIdle,
    });
  }

  render() {
    const {
      handleSubmit,
      moveDocSchema,
      movingExpenseSchema,
      reimbursementSchema,
      formValues,
    } = this.props;
    const isExpenseDocument =
      get(this.props, 'formValues.move_document_type', false) === 'EXPENSE';
    const hasFormFilled = formValues && formValues.move_document_type;
    const hasFiles = this.state.newUploads.length;
    const isValid = hasFormFilled && hasFiles && this.state.uploaderIsIdle;
    console.log('formValues', formValues);
    return (
      <Fragment>
        {this.state.moveDocumentCreateError && (
          <div className="usa-grid">
            <div className="usa-width-one-whole error-message">
              <Alert type="error" heading="An error occurred">
                Something went wrong contacting the server.
              </Alert>
            </div>
          </div>
        )}
        <form onSubmit={handleSubmit(this.onSubmit)}>
          <h3>Upload a new document</h3>
          <SwaggerField
            title="Document type"
            fieldName="move_document_type"
            swagger={moveDocSchema}
            required
          />
          <div className="uploader-box">
            <SwaggerField
              title="Document title"
              fieldName="title"
              swagger={moveDocSchema}
              required
            />
            {isExpenseDocument && (
              <ExpenseDocumentForm
                movingExpenseSchema={movingExpenseSchema}
                reimbursementSchema={reimbursementSchema}
              />
            )}
            <SwaggerField
              title="Notes"
              fieldName="notes"
              swagger={moveDocSchema}
            />
            <div>
              <h4>Attach PDF or image</h4>
              <p>
                Upload a PDF or take a picture of each page and upload the
                images.
              </p>
            </div>
            <Uploader
              onRef={ref => (this.uploader = ref)}
              onChange={this.onChange}
              onAddFile={this.onAddFile}
            />
            <div className="hint">(Each page must be clear and legible)</div>
          </div>
          <button className="submit" disabled={!isValid}>
            Save
          </button>
        </form>
      </Fragment>
    );
  }
}

DocumentUploader = reduxForm({
  form: moveDocumentFormName,
})(DocumentUploader);

DocumentUploader.propTypes = {
  moveId: PropTypes.string.isRequired,
};

function mapStateToProps(state) {
  const props = {
    formValues: getFormValues(moveDocumentFormName)(state),
    moveDocSchema: get(
      state,
      'swagger.spec.definitions.CreateGenericMoveDocumentPayload',
      {},
    ),
    movingExpenseSchema: get(
      state,
      'swagger.spec.definitions.CreateMovingExpenseDocumentPayload',
      {},
    ),
    reimbursementSchema: get(
      state,
      'swagger.spec.definitions.Reimbursement',
      {},
    ),
    moveDocumentCreateError: state.office.moveDocumentCreateError,
    currentPpm:
      get(state.office, 'officePPMs.0') || get(state, 'ppm.currentPpm'),
  };
  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      createMoveDocument,
      createMovingExpenseDocument,
      push,
    },
    dispatch,
  );
}
export default connect(mapStateToProps, mapDispatchToProps)(DocumentUploader);
