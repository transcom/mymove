import { get, map } from 'lodash';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import { bindActionCreators } from 'redux';
import { getFormValues, reduxForm, FormSection } from 'redux-form';
import PropTypes from 'prop-types';

import Alert from 'shared/Alert';
import { createMoveDocument } from 'shared/Entities/modules/moveDocuments';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import Uploader from 'shared/Uploader';

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
    const { formValues, moveId, reset } = this.props;
    const uploadIds = map(this.state.newUploads, 'id');
    this.setState({
      moveDocumentCreateError: null,
    });
    // if (get(formValues, 'movingExpenseDocument'), false) {
    //   this.props
    //     .createMovingExpenseDocument(
    //       moveId,
    //       uploadIds,
    //       formValues.title,
    //       formValues.movingExpenseDocument.moving_expense_type,
    //       formValues.move_document_type,
    //       formValues.reimbursement,
    //       formValues.notes,
    //     )
    //     .then(() => {
    //       reset();
    //       this.uploader.clearFiles();
    //     })
    //     .catch(err => {
    //       this.setState({
    //         moveDocumentCreateError: err,
    //       });
    //     });
    // }
    if (get(formValues, 'movingExpenseDocument', false) === false) {
      this.props
        .createMoveDocument(
          moveId,
          uploadIds,
          formValues.title,
          formValues.move_document_type,
          'AWAITING_REVIEW',
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
      get(this.props, 'formValues.move_document_type', false) === 'EXPENSE' ||
      get(this.props, 'formValues.move_document_type', false) ===
        'STORAGE_EXPENSE';
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
          {isExpenseDocument && (
            <Fragment>
              <FormSection name="movingExpenseDocument">
                <SwaggerField
                  title="Document title"
                  fieldName="title"
                  swagger={movingExpenseSchema}
                  required
                />
                <SwaggerField
                  title="Expense type"
                  fieldName="moving_expense_type"
                  swagger={movingExpenseSchema}
                  required
                />
              </FormSection>
              <FormSection name="reimbursement">
                <SwaggerField
                  title="Amount"
                  fieldName="requested_amount"
                  swagger={reimbursementSchema}
                  required
                />
                <SwaggerField
                  title="Method of Payment"
                  fieldName="method_of_receipt"
                  swagger={reimbursementSchema}
                  required
                />
              </FormSection>
            </Fragment>
          )}
          <div className="uploader-box">
            <SwaggerField
              title="Document title"
              fieldName="title"
              swagger={moveDocSchema}
              required
            />
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
      'swagger.spec.definitions.CreateMoveDocumentPayload',
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
  };
  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      createMoveDocument,
      push,
    },
    dispatch,
  );
}
export default connect(mapStateToProps, mapDispatchToProps)(DocumentUploader);
