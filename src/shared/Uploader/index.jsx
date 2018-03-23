// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { createDocument } from './ducks';
import Alert from 'shared/Alert';

export class Uploader extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Upload Document';
  }

  constructor(props) {
    super(props);
    this.uploadFile = this.uploadFile.bind(this);
  }

  uploadFile() {
    this.props.createDocument(this.fileInput.files[0], this.moveIdInput.value);
  }

  render() {
    const { hasErrored, hasSucceeded } = this.props;
    return (
      <div className="usa-grid">
        Enter Move ID:{' '}
        <input
          type="text"
          ref={input => {
            this.moveIdInput = input;
          }}
        />
        <input
          type="file"
          ref={input => {
            this.fileInput = input;
          }}
        />
        <button onClick={this.uploadFile}>Upload Now</button>
        {hasErrored && (
          <Alert type="error" heading="Submission Error">
            Something went wrong with your upload
          </Alert>
        )}
        {hasSucceeded && (
          <Alert type="success" heading="Submission Successful">
            Your document was successfully uploaded.
          </Alert>
        )}
      </div>
    );
  }
}

Uploader.propTypes = {
  createDocument: PropTypes.func.isRequired,
  hasErrored: PropTypes.bool.isRequired,
  hasSucceeded: PropTypes.bool.isRequired,
};

function mapStateToProps(state) {
  return state.upload;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createDocument }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Uploader);
