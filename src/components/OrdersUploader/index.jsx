/* eslint-disable no-underscore-dangle */
/* eslint-disable react/destructuring-assignment */
import React, { Component } from 'react';
import 'filepond-polyfill/dist/filepond-polyfill';
import { FilePond, registerPlugin } from 'react-filepond';
import { FileStatus } from 'filepond';
import PropTypes from 'prop-types';
import isMobile from 'is-mobile';
import { get, reject } from 'lodash';
import 'filepond/dist/filepond.min.css';
import FilepondPluginFileValidateType from 'filepond-plugin-file-validate-type';
import FilePondPluginFileValidateSize from 'filepond-plugin-file-validate-size';
import FilepondPluginImageExifOrientation from 'filepond-plugin-image-exif-orientation';
import FilePondImagePreview from 'filepond-plugin-image-preview';
import 'filepond-plugin-image-preview/dist/filepond-plugin-image-preview.css';

import 'shared/Uploader/index.css';

registerPlugin(FilepondPluginFileValidateType);
registerPlugin(FilePondPluginFileValidateSize);
registerPlugin(FilepondPluginImageExifOrientation);
registerPlugin(FilePondImagePreview);

const idleStatuses = [FileStatus.PROCESSING_COMPLETE, FileStatus.PROCESSING_ERROR];

class OrdersUploader extends Component {
  constructor(props) {
    super(props);

    this.state = {
      files: [],
    };
  }

  componentDidMount() {
    if (this.props.onRef) {
      this.props.onRef(this);
    }
  }

  componentWillUnmount() {
    if (this.props.onRef) {
      this.props.onRef(undefined);
    }
  }

  handlePondInit() {
    // If this component is unloaded quickly, this function can be called after the ref is deleted,
    // so check that the ref still exists before continuing
    if (!this.pond) {
      return;
    }
    this.setPondOptions();

    this.pond._pond.on('processfile', () => {
      if (this.props.onChange) {
        this.props.onChange(this.state.files, this.isIdle());
      }
    });

    this.pond._pond.on('addfilestart', () => {
      if (this.props.onAddFile) {
        this.props.onAddFile();
      }
    });

    // Don't mention drag and drop if on mobile device
    if (isMobile()) {
      this.pond._pond.setOptions({
        labelIdle: '<span class="filepond--label-action">Upload</span>',
      });
    }
  }

  getFiles() {
    return this.state.files;
  }

  setPondOptions() {
    // If this component is unloaded quickly, this function can be called after the ref is deleted,
    // so check that the ref still exists before continuing
    if (!this.pond) {
      return;
    }
    const { options } = this.props;
    const defaultOptions = {
      allowMultiple: true,
      server: {
        url: '/',
        process: this.processFile,
        revert: this.revertFile,
      },
      iconUndo: this.pond._pond.iconRemove,
      imagePreviewMaxHeight: 100,
      labelIdle: 'Drag & drop or <span class="filepond--label-action">click to upload</span>',
      labelTapToUndo: 'tap to delete',
      acceptedFileTypes: ['image/jpeg', 'image/png', 'application/pdf'],
      maxFileSize: '25MB',
    };
    this.pond._pond.setOptions({ ...defaultOptions, ...options });
  }

  processFile = (fieldName, file, metadata, load, error, progress, abort) => {
    const { createUpload } = this.props;
    const self = this;
    createUpload(file)
      .then((item) => {
        const response = get(item, 'response.body', {});
        load(response.id);
        const newFiles = self.state.files.concat(response);
        self.setState({
          files: newFiles,
        });
      })
      .catch(error);

    return { abort };
  };

  revertFile = (uploadId, load, error) => {
    const { onChange, deleteUpload } = this.props;
    deleteUpload(uploadId)
      .then((item) => {
        const response = get(item, 'response', {});
        load(response);
        const getNewFiles = (state) => reject(state.files, (upload) => upload.id === uploadId);
        this.setState(
          (prevState) => ({
            files: getNewFiles(prevState),
          }),
          () => {
            if (onChange) {
              onChange(this.state.files, this.isIdle());
            }
          },
        );
      })
      .catch(error);
  };

  isEmpty() {
    return this.state.files.length === 0;
  }

  // TODO: Remove isIdle function- the onChange where it is called does not actually expect second argument
  isIdle() {
    // If this component is unloaded quickly, this function can be called after the ref is deleted,
    // so check that the ref still exists before continuing
    if (!this.pond) {
      return false;
    }
    // Returns a boolean: is FilePond done with all uploading?
    const existingFiles = this.pond._pond.getFiles();
    return existingFiles.every((f) => idleStatuses.indexOf(f.status) > -1);
  }

  clearFiles() {
    // If this component is unloaded quickly, this function can be called after the ref is deleted,
    // so check that the ref still exists before continuing
    if (!this.pond) {
      return;
    }
    this.pond._pond.removeFiles();

    this.setState({
      files: [],
    });

    if (this.props.onChange) {
      this.props.onChange([], true);
    }
  }

  render() {
    return (
      <div>
        <FilePond
          ref={(ref) => {
            this.pond = ref;
          }}
          oninit={() => this.handlePondInit()}
        />
      </div>
    );
  }
}

OrdersUploader.propTypes = {
  onChange: PropTypes.func,
  createUpload: PropTypes.func.isRequired,
  onRef: PropTypes.func,
  deleteUpload: PropTypes.func,
  onAddFile: PropTypes.func,
  // eslint-disable-next-line react/forbid-prop-types
  options: PropTypes.object,
};

OrdersUploader.defaultProps = {
  onChange: undefined,
  onRef: undefined,
  deleteUpload: undefined,
  onAddFile: undefined,
  options: {},
};

export default OrdersUploader;
