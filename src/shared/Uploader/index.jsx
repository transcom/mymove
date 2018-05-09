// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';
import { connect } from 'react-redux';
import 'filepond-polyfill/dist/filepond-polyfill.js';
import { FilePond, registerPlugin } from 'react-filepond';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { CreateUpload, DeleteUpload } from 'shared/api.js';

import 'filepond/dist/filepond.min.css';
import './index.css';

// Register the image preview plugin
import FilePondImagePreview from 'filepond-plugin-image-preview';
import 'filepond-plugin-image-preview/dist/filepond-plugin-image-preview.css';
registerPlugin(FilePondImagePreview);

export class Uploader extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Upload Document';
  }

  handlePondInit() {
    this.pond._pond.setOptions({
      allowMultiple: true,
      server: {
        url: '/',
        process: this.processFile,
        revert: this.revertFile,
      },
    });
  }

  processFile = (fieldName, file, metadata, load, error, progress, abort) => {
    CreateUpload(file, this.props.document.id)
      .then(item => load(item.id))
      .catch(error);

    return { abort };
  };

  revertFile = (uploadId, load, error) => {
    DeleteUpload(uploadId)
      .then(load)
      .catch(error);
  };

  render() {
    return (
      <div className="usa-grid">
        <FilePond
          ref={ref => (this.pond = ref)}
          oninit={() => this.handlePondInit()}
        />
      </div>
    );
  }
}

Uploader.propTypes = {
  document: PropTypes.object.isRequired,
};

function mapStateToProps(state) {
  return {};
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({}, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Uploader);
