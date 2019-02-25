import { includes, get, map } from 'lodash';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { push, replace } from 'react-router-redux';
import { bindActionCreators } from 'redux';
import { getFormValues, reduxForm } from 'redux-form';
import PropTypes from 'prop-types';

import Alert from 'shared/Alert';
import { selectPPMForMove } from 'shared/Entities/modules/ppms';
import { createMoveDocument } from 'shared/Entities/modules/moveDocuments';
import { createMovingExpenseDocument } from 'shared/Entities/modules/movingExpenseDocuments';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import Uploader from 'shared/Uploader';
import ExpenseDocumentForm from './ExpenseDocumentForm';
import { convertDollarsToCents } from 'shared/utils';

import './DocumentUploader.css';

const moveDocumentFormName = 'move_document_upload';

export class DocumentUploader extends Component {
  componentDidUpdate() {
    const { initialValues, location } = this.props;
    // Clear query string after initial values are set
    if (initialValues && get(location, 'search', false)) {
      this.props.replace(this.props.location.pathname);
    }
  }

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
    if (get(formValues, 'move_document_type', false) === 'EXPENSE') {
      formValues.requested_amount_cents = convertDollarsToCents(formValues.requested_amount_cents);
      this.props
        .createMovingExpenseDocument(
          moveId,
          currentPpm.id,
          uploadIds,
          formValues.title,
          formValues.moving_expense_type,
          formValues.move_document_type,
          formValues.requested_amount_cents,
          formValues.payment_method,
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
    } else {
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
    const { handleSubmit, moveDocSchema, genericMoveDocSchema, formValues } = this.props;
    const isExpenseDocument = get(this.props, 'formValues.move_document_type', false) === 'EXPENSE';
    const hasFormFilled = formValues && formValues.move_document_type;
    const hasFiles = this.state.newUploads.length;
    const isValid = hasFormFilled && hasFiles && this.state.uploaderIsIdle;
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
          <SwaggerField title="Document type" fieldName="move_document_type" swagger={genericMoveDocSchema} required />
          <div className="uploader-box">
            <SwaggerField title="Document title" fieldName="title" swagger={genericMoveDocSchema} required />
            {isExpenseDocument && <ExpenseDocumentForm moveDocSchema={moveDocSchema} />}
            <SwaggerField title="Notes" fieldName="notes" swagger={genericMoveDocSchema} />
            <div>
              <h4>Attach PDF or image</h4>
              <p>Upload a PDF or take a picture of each page and upload the images.</p>
            </div>
            <Uploader onRef={ref => (this.uploader = ref)} onChange={this.onChange} onAddFile={this.onAddFile} />
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
  moveDocumentType: PropTypes.string,
  location: PropTypes.object,
};

function mapStateToProps(state, props) {
  const docTypes = get(state, 'swaggerInternal.spec.definitions.MoveDocumentType.enum', []);

  let initialValues = {};
  // Verify the provided doc type against the schema
  if (includes(docTypes, props.moveDocumentType)) {
    initialValues.move_document_type = props.moveDocumentType;
  }

  const newProps = {
    initialValues: initialValues,
    formValues: getFormValues(moveDocumentFormName)(state),
    genericMoveDocSchema: get(state, 'swaggerInternal.spec.definitions.CreateGenericMoveDocumentPayload', {}),
    moveDocSchema: get(state, 'swaggerInternal.spec.definitions.MoveDocumentPayload', {}),
    currentPpm: selectPPMForMove(state, props.moveId) || get(state, 'ppm.currentPpm'),
  };
  return newProps;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      createMoveDocument,
      createMovingExpenseDocument,
      push,
      replace,
    },
    dispatch,
  );
}
export default connect(mapStateToProps, mapDispatchToProps)(DocumentUploader);
