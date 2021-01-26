/* eslint-disable react/destructuring-assignment */
import React, { Component } from 'react';
import PropTypes from 'prop-types';
import isMobile from 'is-mobile';
import { get, reject } from 'lodash';
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
  constructor(props) {
    super(props);

    this.state = {
      files: [],
    };
  }

  handleProcessFile = () => {
    if (this.props.onChange) {
      this.props.onChange(this.state.files, this.isIdle());
    }
  };

  handleAddFileStart = () => {
    if (this.props.onAddFile) {
      this.props.onAddFile();
    }
  };

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

  isIdle() {
    return this.pond?.status === Status.IDLE || false;
  }

  render() {
    const { labelIdle, files } = this.props;

    const serverConfig = {
      url: '/',
      process: this.processFile,
      revert: this.revertFile,
    };

    /**
     * Default FilePond instance props
     * If these need to be overwritten, they can be exposed as a prop on this
     * component and passed through (like labelIdle)
     */
    const filePondProps = {
      allowMultiple: true,
      server: serverConfig,
      iconUndo: this.pond?.iconRemove,
      imagePreviewMaxHeight: 100,
      labelIdle: isMobile() ? '<span class="filepond--label-action">Upload</span>' : labelIdle,
      labelTapToUndo: 'tap to delete',
      acceptedFileTypes: ['image/jpeg', 'image/png', 'application/pdf'],
      maxFileSize: '25MB',
    };

    /* eslint-disable react/jsx-props-no-spreading */
    return (
      <div>
        <FilePond
          ref={(ref) => {
            this.pond = ref;
          }}
          {...filePondProps}
          onprocessfile={this.handleProcessFile}
          onaddfilestart={this.handleAddFileStart}
          files={files.map((f) => ({
            source: f.id,
            options: {
              type: 'local',
              file: {
                name: f.filename,
                size: f.bytes,
                type: f.content_type,
              },
            },
          }))}
        />
      </div>
    );
    /* eslint-enable react/jsx-props-no-spreading */
  }
}

FileUpload.propTypes = {
  onChange: PropTypes.func,
  createUpload: PropTypes.func.isRequired,
  deleteUpload: PropTypes.func,
  onAddFile: PropTypes.func,
  // FilePond instance props
  labelIdle: PropTypes.string,
  // eslint-disable-next-line react/forbid-prop-types
  files: PropTypes.array,
};

FileUpload.defaultProps = {
  onChange: undefined,
  deleteUpload: undefined,
  onAddFile: undefined,
  labelIdle: 'Drag & drop or <span class="filepond--label-action">click to upload</span>',
  files: [],
};

export default FileUpload;
