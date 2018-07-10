import { get, map } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import { bindActionCreators } from 'redux';

import { getFormValues, reduxForm } from 'redux-form';

import { createMoveDocument } from '../ducks.js';
import Uploader from 'shared/Uploader';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import './DocumentUploader.css';

const moveDocumentFormName = 'move_document_upload';

export class DocumentUploader extends Component {
  constructor(props) {
    super(props);

    this.state = {
      newUploads: [],
    };

    this.onChange = this.onChange.bind(this);
    this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit() {
    const { formValues } = this.props;
    const uploadIds = map(this.state.newUploads, 'id');
    const moveId = this.props.match.params.moveId;
    this.props
      .createMoveDocument(
        moveId,
        uploadIds,
        formValues.title,
        formValues.move_document_type,
        'AWAITING_REVIEW',
        formValues.notes,
      )
      .then(response => {
        const moveDocumentId = response.payload.id;
        this.props.push(`/moves/${moveId}/documents/${moveDocumentId}`);
      });
  }

  onChange(files) {
    this.setState({
      newUploads: files,
    });
  }

  render() {
    const { handleSubmit, schema, formValues } = this.props;
    const hasFormFilled = formValues && formValues.move_document_type;
    const hasFiles = this.state.newUploads.length;
    const isValid = hasFormFilled && hasFiles;
    return (
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
            <Uploader onChange={this.onChange} />
            <div className="hint">(Each page must be clear and legible)</div>
          </div>
        </div>
        <button className="submit" disabled={!isValid}>
          Save
        </button>
      </form>
    );
  }
}

DocumentUploader = reduxForm({
  form: moveDocumentFormName,
})(DocumentUploader);

DocumentUploader.propTypes = {};

function mapStateToProps(state) {
  const props = {
    formValues: getFormValues(moveDocumentFormName)(state),
    schema: get(
      state,
      'swagger.spec.definitions.CreateMoveDocumentPayload',
      {},
    ),
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
