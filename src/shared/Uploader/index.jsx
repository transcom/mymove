// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';
import { connect } from 'react-redux';
import 'filepond-polyfill/dist/filepond-polyfill.js';
import { FilePond, registerPlugin } from 'react-filepond';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { CreateUpload, DeleteUpload } from 'shared/api.js';
import isMobile from 'is-mobile';
import { concat, reject } from 'lodash';

import 'filepond/dist/filepond.min.css';
import './index.css';

import FilepondPluginFileValidateType from 'filepond-plugin-file-validate-type';
import FilepondPluginImageExifOrientation from 'filepond-plugin-image-exif-orientation';
import FilePondImagePreview from 'filepond-plugin-image-preview';
import 'filepond-plugin-image-preview/dist/filepond-plugin-image-preview.css';

registerPlugin(FilepondPluginFileValidateType);
registerPlugin(FilepondPluginImageExifOrientation);
registerPlugin(FilePondImagePreview);

export class Uploader extends Component {
  constructor(props) {
    super(props);

    this.state = {
      files: [],
    };
  }

  handlePondInit() {
    this.pond._pond.setOptions({
      allowMultiple: true,
      server: {
        url: '/',
        process: this.processFile,
        revert: this.revertFile,
      },
      iconUndo: this.pond._pond.iconRemove,
      imagePreviewMaxHeight: 100,
      labelIdle:
        'Drag & drop or <span class="filepond--label-action">click to upload orders</span>',
      labelTapToUndo: 'tap to delete',
      acceptedFileTypes: ['image/*', 'application/pdf'],
    });

    // Don't mention drag and drop if on mobile device
    if (isMobile()) {
      this.pond._pond.setOptions({
        labelIdle: '<span class="filepond--label-action">Upload</span>',
      });
    }
  }

  processFile = (fieldName, file, metadata, load, error, progress, abort) => {
    const self = this;
    const docID = this.props.document ? this.props.document.id : null;
    CreateUpload(file, docID)
      .then(item => {
        load(item.id);
        const newFiles = concat(self.state.files, item);
        self.setState({
          files: newFiles,
        });
        // Call onChange after the upload completes
        self.pond._pond.onOnce('processfile', e => {
          if (self.props.onChange) {
            self.props.onChange(newFiles);
          }
        });
      })
      .catch(error);

    return { abort };
  };

  revertFile = (uploadId, load, error) => {
    DeleteUpload(uploadId)
      .then(item => {
        load(item);
        const newFiles = reject(
          this.state.files,
          upload => upload.id === uploadId,
        );
        this.setState({
          files: newFiles,
        });
        if (this.props.onChange) {
          this.props.onChange(newFiles);
        }
      })
      .catch(error);
  };

  render() {
    return (
      <div>
        <FilePond
          ref={ref => (this.pond = ref)}
          oninit={() => this.handlePondInit()}
        />
      </div>
    );
  }
}

Uploader.propTypes = {
  document: PropTypes.object,
  onChange: PropTypes.func,
};

function mapStateToProps(state) {
  return {};
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({}, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Uploader);
