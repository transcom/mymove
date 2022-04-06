import React, { forwardRef } from 'react';
import PropTypes from 'prop-types';
import isMobile from 'is-mobile';
import { FilePond, registerPlugin } from 'react-filepond';
import 'filepond-polyfill/dist/filepond-polyfill';
import 'filepond/dist/filepond.min.css';
import FilepondPluginFileValidateType from 'filepond-plugin-file-validate-type';
import FilePondPluginFileValidateSize from 'filepond-plugin-file-validate-size';
import FilepondPluginImageExifOrientation from 'filepond-plugin-image-exif-orientation';
import FilePondImagePreview from 'filepond-plugin-image-preview';
import 'filepond-plugin-image-preview/dist/filepond-plugin-image-preview.css';

import 'shared/Uploader/index.css';
import { createUpload as createUploadApi, deleteUpload } from 'services/internalApi';

registerPlugin(
  FilepondPluginFileValidateType,
  FilePondPluginFileValidateSize,
  FilepondPluginImageExifOrientation,
  FilePondImagePreview,
);

const FileUpload = forwardRef(({ name, createUpload, onChange, labelIdle, labelIdleMobile, onAddFile }, ref) => {
  const handleOnChange = () => {
    if (onChange) onChange();
  };

  const processFile = (fieldName, file, metadata, load, error, progress, abort) => {
    createUpload(file)
      .then((response) => {
        load(response.id);
      })
      .catch(error);

    // TODO - in order to handle abort, we need to pass an AbortController to SwaggerRequest, implement this as a future story (it's not working as-is)
    // https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#browser
    return { abort };
  };

  const revertFile = (uploadId, load, error) => {
    deleteUpload(uploadId)
      .then(() => {
        load();
        handleOnChange();
      })
      .catch(error);
  };

  const handleProcessFile = () => {
    handleOnChange();
  };

  const serverConfig = {
    url: '/internal',
    process: processFile,
    fetch: null,
    revert: revertFile,
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
    labelIdle: isMobile() ? labelIdleMobile : labelIdle,
    acceptedFileTypes: [
      'image/jpeg',
      'image/png',
      'application/pdf',
      'application/vnd.ms-excel',
      'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
      'application/octet-stream',
      'application/zip',
    ],
    maxFileSize: '25MB',
    credits: false,
  };

  /* eslint-disable react/jsx-props-no-spreading */
  return (
    <FilePond ref={ref} {...filePondProps} name={name} onprocessfile={handleProcessFile} onaddfilestart={onAddFile} />
  );
  /* eslint-enable react/jsx-props-no-spreading */
});

FileUpload.propTypes = {
  name: PropTypes.string,
  createUpload: PropTypes.func,
  onChange: PropTypes.func,
  onAddFile: PropTypes.func,
  // FilePond instance props
  labelIdle: PropTypes.string,
  labelIdleMobile: PropTypes.string,
};

FileUpload.defaultProps = {
  name: 'file',
  createUpload: createUploadApi,
  onChange: undefined,
  onAddFile: undefined,
  labelIdle: 'Drag & drop or <span class="filepond--label-action">click to upload</span>',
  labelIdleMobile: '<span class="filepond--label-action">Upload</span>',
};

export default FileUpload;
