// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { createDocument } from './ducks';

export class Uploader extends Component {
  constructor(props) {
    super(props);
    this.uploadFile = this.uploadFile.bind(this);
  }

  uploadFile() {
    this.props.createDocument(this.fileInput.files[0]);
  }

  render() {
    return (
      <div className="uploader">
        <input
          type="file"
          ref={input => {
            this.fileInput = input;
          }}
        />
        <button onClick={this.uploadFile}>Upload Now</button>
      </div>
    );
  }
}

Uploader.propTypes = {
  createDocument: PropTypes.func.isRequired,
  hasSubmitError: PropTypes.bool.isRequired,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapStateToProps(state) {
  return state.document;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createDocument }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Uploader);
