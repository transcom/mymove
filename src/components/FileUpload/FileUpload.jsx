import React, { forwardRef } from 'react';
import { func, string, arrayOf, bool, int } from 'prop-types';
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

const FileUpload = forwardRef(
  (
    {
      name,
      createUpload,
      onChange,
      labelIdle,
      labelIdleMobile,
      onAddFile,
      acceptedFileTypes,
      allowMultiple,
      maxParralelUploads,
    },
    ref,
  ) => {
    const handleOnChange = () => {
      if (onChange) onChange();
    };

    const processFile = (fieldName, file, metadata, load, error, progress, abort) => {
      createUpload(file)
        .then((response) => {
          load(response);
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
      server: serverConfig,
      imagePreviewMaxHeight: 100,
      labelIdle: isMobile() ? labelIdleMobile : labelIdle,
      maxFileSize: '25MB',
      credits: false,
    };

    /* eslint-disable react/jsx-props-no-spreading */
    return (
      <FilePond
        ref={ref}
        {...filePondProps}
        allowMultiple={allowMultiple}
        maxParralelUploads={maxParralelUploads}
        acceptedFileTypes={acceptedFileTypes}
        name={name}
        onprocessfile={handleProcessFile}
        onaddfilestart={onAddFile}
      />
    );
    /* eslint-enable react/jsx-props-no-spreading */
  },
);

FileUpload.propTypes = {
  name: string,
  createUpload: func,
  onChange: func,
  onAddFile: func,
  acceptedFileTypes: arrayOf(string),
  allowMultiple: bool,
  maxParralelUploads: int,
  // FilePond instance props
  labelIdle: string,
  labelIdleMobile: string,
};

FileUpload.defaultProps = {
  name: 'file',
  createUpload: createUploadApi,
  onChange: undefined,
  onAddFile: undefined,
  allowMultiple: true,
  maxParralelUploads: 2,
  acceptedFileTypes: ['image/jpeg', 'image/png', 'application/pdf'],
  labelIdle: 'Drag & drop or <span class="filepond--label-action">click to upload</span>',
  labelIdleMobile: '<span class="filepond--label-action">Upload</span>',
};

export default FileUpload;
