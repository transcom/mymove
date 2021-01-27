/* eslint-disable react/destructuring-assignment */
import React, { Component } from 'react';
import PropTypes from 'prop-types';
import isMobile from 'is-mobile';
import { FilePond, registerPlugin } from 'react-filepond';
import { Status } from 'filepond';
import 'filepond-polyfill/dist/filepond-polyfill';
import 'filepond/dist/filepond.min.css';
import FilepondPluginFileValidateType from 'filepond-plugin-file-validate-type';
import FilePondPluginFileValidateSize from 'filepond-plugin-file-validate-size';
import FilepondPluginImageExifOrientation from 'filepond-plugin-image-exif-orientation';
import FilePondImagePreview from 'filepond-plugin-image-preview';
import 'filepond-plugin-image-preview/dist/filepond-plugin-image-preview.css';

import 'shared/Uploader/index.css';

registerPlugin(
  FilepondPluginFileValidateType,
  FilePondPluginFileValidateSize,
  FilepondPluginImageExifOrientation,
  FilePondImagePreview,
);

// TODO:
// - forwardRef if necessary (props.onRef)

class FileUpload extends Component {
  processFile = (fieldName, file, metadata, load, error, progress, abort) => {
    const { createUpload } = this.props;

    createUpload(file)
      .then((response) => {
        load(response.id);
      })
      .catch(error);

    // TODO - abort handler?
    return { abort };
  };

  /*
  revertFile = (uploadId, load, error) => {
    // TODO
  };
  */

  handleProcessFile = () => {
    if (this.props.onChange) {
      this.props.onChange(this.pond?.getFiles(), this.isIdle());
    }

    // TODO - make this an option
    this.pond?.removeFiles();
  };

  isIdle() {
    return this.pond?.status === Status.IDLE || false;
  }

  render() {
    const { labelIdle, onAddFile } = this.props;

    const serverConfig = {
      url: '/internal',
      process: this.processFile,
      fetch: null,
      revert: null,
    };

    /**
     * Default FilePond instance props
     * If these need to be overwritten, they can be exposed as a prop on this
     * component and passed through (like labelIdle)
     */
    const filePondProps = {
      allowMultiple: true,
      server: serverConfig,
      imagePreviewMaxHeight: 100,
      labelIdle: isMobile() ? '<span class="filepond--label-action">Upload</span>' : labelIdle,
      acceptedFileTypes: ['image/jpeg', 'image/png', 'application/pdf'],
      maxFileSize: '25MB',
    };

    /* eslint-disable react/jsx-props-no-spreading */
    return (
      <FilePond
        ref={(ref) => {
          this.pond = ref;
        }}
        {...filePondProps}
        name="file"
        onprocessfile={this.handleProcessFile}
        onaddfilestart={onAddFile}
      />
    );
    /* eslint-enable react/jsx-props-no-spreading */
  }
}

FileUpload.propTypes = {
  createUpload: PropTypes.func.isRequired,
  onChange: PropTypes.func,
  onAddFile: PropTypes.func,
  // FilePond instance props
  labelIdle: PropTypes.string,
};

FileUpload.defaultProps = {
  onChange: undefined,
  onAddFile: undefined,
  labelIdle: 'Drag & drop or <span class="filepond--label-action">click to upload</span>',
};

export default FileUpload;
