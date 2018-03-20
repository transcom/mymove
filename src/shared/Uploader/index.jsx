// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { createDocument } from './ducks';
import UploadConfirmation from 'shared/Uploader/UploadConfirmation';

export class Uploader extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Upload Document';
  }

  constructor(props) {
    super(props);
    this.uploadFile = this.uploadFile.bind(this);
  }

  uploadFile() {
    this.props.createDocument(this.fileInput.files[0]);
  }

  render() {
    const { confirmationText } = this.props;
    return (
      <div className="usa-grid">
        <input
          type="file"
          ref={input => {
            this.fileInput = input;
          }}
        />
        <button onClick={this.uploadFile}>Upload Now</button>
        <UploadConfirmation confirmationText={confirmationText} />
      </div>
    );
  }
}

Uploader.propTypes = {
  createDocument: PropTypes.func.isRequired,
  hasSubmitError: PropTypes.bool.isRequired,
  hasSubmitSuccess: PropTypes.bool.isRequired,
  confirmationText: PropTypes.string.isRequired,
};

function mapStateToProps(state) {
  return state.document;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createDocument }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Uploader);
