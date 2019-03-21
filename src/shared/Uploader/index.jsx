// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';
import { connect } from 'react-redux';
import 'filepond-polyfill/dist/filepond-polyfill.js';
import { FilePond, registerPlugin } from 'react-filepond';
import { FileStatus } from 'filepond';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { CreateUpload, DeleteUpload } from 'shared/api.js';
import isMobile from 'is-mobile';
import { concat, reject, every, includes } from 'lodash';

import 'filepond/dist/filepond.min.css';
import './index.css';

import FilepondPluginFileValidateType from 'filepond-plugin-file-validate-type';
import FilepondPluginImageExifOrientation from 'filepond-plugin-image-exif-orientation';
import FilePondImagePreview from 'filepond-plugin-image-preview';
import 'filepond-plugin-image-preview/dist/filepond-plugin-image-preview.css';

registerPlugin(FilepondPluginFileValidateType);
registerPlugin(FilepondPluginImageExifOrientation);
registerPlugin(FilePondImagePreview);

const idleStatuses = [FileStatus.PROCESSING_COMPLETE, FileStatus.PROCESSING_ERROR];

export class Uploader extends Component {
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

  clearFiles() {
    this.pond._pond.removeFiles();

    this.setState({
      files: [],
    });

    if (this.props.onChange) {
      this.props.onChange([], true);
    }
  }

  isIdle() {
    // Returns a boolean: is FilePond done with all uploading?
    const existingFiles = this.pond._pond.getFiles();
    const isIdle = every(existingFiles, f => {
      return includes(idleStatuses, f.status);
    });

    return isIdle;
  }

  handlePondInit() {
    // If this component is unloaded quickly, this function can be called after the ref is deleted,
    // so check that the ref still exists before continuing
    if (!this.pond) {
      return;
    }

    const { labelIdle } = this.props;
    this.pond._pond.setOptions({
      allowMultiple: true,
      server: {
        url: '/',
        process: this.processFile,
        revert: this.revertFile,
      },
      iconUndo: this.pond._pond.iconRemove,
      imagePreviewMaxHeight: 100,
      labelIdle: labelIdle || 'Drag & drop or <span class="filepond--label-action">click to upload</span>',
      labelTapToUndo: 'tap to delete',
      acceptedFileTypes: ['image/jpeg', 'image/png', 'application/pdf'],
    });

    this.pond._pond.on('processfile', e => {
      if (this.props.onChange) {
        this.props.onChange(this.state.files, this.isIdle());
      }
    });

    this.pond._pond.on('addfilestart', e => {
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

  processFile = (fieldName, file, metadata, load, error, progress, abort) => {
    const { document, isPublic } = this.props;
    const self = this;
    const docID = document ? document.id : null;
    CreateUpload(file, docID, isPublic)
      .then(item => {
        load(item.id);
        const newFiles = concat(self.state.files, item);
        self.setState({
          files: newFiles,
        });
      })
      .catch(error);

    return { abort };
  };

  revertFile = (uploadId, load, error) => {
    const { onChange, isPublic } = this.props;
    DeleteUpload(uploadId, isPublic)
      .then(item => {
        load(item);
        const newFiles = reject(this.state.files, upload => upload.id === uploadId);
        this.setState({
          files: newFiles,
        });
        if (onChange) {
          onChange(newFiles, this.isIdle());
        }
      })
      .catch(error);
  };

  render() {
    return (
      <div>
        <FilePond ref={ref => (this.pond = ref)} oninit={() => this.handlePondInit()} />
      </div>
    );
  }
}

Uploader.propTypes = {
  document: PropTypes.object,
  onChange: PropTypes.func,
  labelIdle: PropTypes.string,
  isPublic: PropTypes.bool,
};

function mapStateToProps(state) {
  return {};
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({}, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Uploader);
