import { get, map } from 'lodash';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { getFormValues, reduxForm } from 'redux-form';
import PropTypes from 'prop-types';

import Alert from 'shared/Alert';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import Uploader from 'shared/Uploader';
import { documentSizeLimitMsg } from 'shared/constants';
import ExpenseDocumentForm from 'scenes/Office/DocumentViewer/ExpenseDocumentForm';

import './DocumentUploader.css';
import withRouter from 'utils/routing';

export class DocumentUploader extends Component {
  static propTypes = {
    form: PropTypes.string.isRequired,
    initialValues: PropTypes.object.isRequired,
    genericMoveDocSchema: PropTypes.object.isRequired,
    moveDocSchema: PropTypes.object,
    onSubmit: PropTypes.func.isRequired,
    isPublic: PropTypes.bool,
  };

  state = {
    newUploads: [],
    uploaderIsIdle: true,
    moveDocumentCreateError: false,
  };

  componentDidUpdate() {
    const {
      initialValues,
      router: { location, navigate },
    } = this.props;
    // Clear query string after initial values are set
    if (initialValues && get(location, 'search', false)) {
      navigate(location.pathname, { replace: true });
    }
  }

  onSubmit = () => {
    const { formValues, reset } = this.props;
    const uploadIds = map(this.state.newUploads, 'id');
    this.setState({
      moveDocumentCreateError: null,
    });
    this.props
      .onSubmit(uploadIds, formValues)
      .then(() => {
        reset();
        this.uploader.clearFiles();
      })
      .catch((err) => {
        this.setState({
          moveDocumentCreateError: err,
        });
      });

    // todo: we don't want to do this until the details view is working,
    // we may not want to do it at all if users are going to upload several documents at a time
    // .then(response => {
    //   if (!response.error) {
    //     const moveDocumentId = response.payload.id;
    //     this.props.push(`/moves/${moveId}/documents/${moveDocumentId}`);
    //   }
    // });
  };

  onAddFile = () => {
    this.setState({
      uploaderIsIdle: false,
    });
  };

  onChange = (newUploads, uploaderIsIdle) => {
    this.setState({
      newUploads,
      uploaderIsIdle,
    });
  };

  render() {
    const { handleSubmit, moveDocSchema, genericMoveDocSchema, formValues, isPublic } = this.props;
    const isExpenseDocument = get(this.props, 'formValues.move_document_type', false) === 'EXPENSE';
    const hasFormFilled = formValues && formValues.move_document_type;
    const hasFiles = this.state.newUploads.length;
    const isValid = hasFormFilled && hasFiles && this.state.uploaderIsIdle;
    return (
      <Fragment>
        {this.state.moveDocumentCreateError && (
          <div className="grid-row">
            <div className="grid-col-12 error-message">
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
              <p>{documentSizeLimitMsg}</p>
            </div>
            <Uploader
              isPublic={isPublic}
              onRef={(ref) => (this.uploader = ref)}
              onChange={this.onChange}
              onAddFile={this.onAddFile}
            />
            <div className="hint">(Each page must be clear and legible)</div>
          </div>
          <button className="usa-button submit" disabled={!isValid}>
            Save
          </button>
        </form>
      </Fragment>
    );
  }
}

// form name comes from prop 'form' that is passed in when used
DocumentUploader = reduxForm()(DocumentUploader);

const mapStateToProps = (state, props) => ({
  formValues: getFormValues(props.form)(state),
});

export default withRouter(connect(mapStateToProps)(DocumentUploader));
