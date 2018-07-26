import { get, map } from 'lodash';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import { bindActionCreators } from 'redux';
import { getFormValues, reduxForm } from 'redux-form';
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
    const { formValues, officePPM, moveId, reset } = this.props;
    const uploadIds = map(this.state.newUploads, 'id');
    this.setState({
      moveDocumentCreateError: null,
    });
    this.props
      .createMoveDocument(
        moveId,
        officePPM.id,
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
    const { handleSubmit, schema, formValues } = this.props;
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
          <SwaggerField
            title="Title"
            fieldName="title"
            swagger={schema}
            required
          />
          <SwaggerField
            title="Document type"
            fieldName="move_document_type"
            swagger={schema}
            required
          />
          <SwaggerField title="Notes" fieldName="notes" swagger={schema} />
          <div>
            <h4>Attach PDF or image</h4>
            <p>
              Upload a PDF or take a picture of each page and upload the images.
            </p>
            <div className="uploader-box">
              <Uploader
                onRef={ref => (this.uploader = ref)}
                onChange={this.onChange}
                onAddFile={this.onAddFile}
              />
              <div className="hint">(Each page must be clear and legible)</div>
            </div>
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
    schema: get(
      state,
      'swagger.spec.definitions.CreateMoveDocumentPayload',
      {},
    ),
    moveDocumentCreateError: state.office.moveDocumentCreateError,
    officePPM: get(state.office, 'officePPMs.0', {}),
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
