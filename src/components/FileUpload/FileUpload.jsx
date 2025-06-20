import React, { forwardRef, useRef } from 'react';
import { func, string, arrayOf, bool, number } from 'prop-types';
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
import { UPLOAD_SCAN_STATUS } from 'shared/constants';
import { createUpload as createUploadApi, deleteUpload, waitForAvScan } from 'services/internalApi';

registerPlugin(
  FilepondPluginFileValidateType,
  FilePondPluginFileValidateSize,
  FilepondPluginImageExifOrientation,
  FilePondImagePreview,
);

// Helper func to take in a root html element and search it for
// a FilePond listing by original file upload name
const getFilePondItemByFilename = (root, originalName) => {
  if (!root) return null;
  return (
    Array.from(root.querySelectorAll('.filepond--item')).find((item) => {
      const nameEl = item.querySelector('.filepond--file-info-main');
      return nameEl && nameEl.textContent.trim() === originalName;
    }) || null
  );
};

const FileUpload = forwardRef(
  (
    {
      name,
      createUpload,
      className,
      onChange,
      labelIdle,
      labelIdleMobile,
      onAddFile,
      acceptedFileTypes,
      allowMultiple,
      maxParralelUploads,
      fileValidateTypeLabelExpectedTypes,
      labelFileTypeNotAllowed,
      required,
    },
    refFromParent,
  ) => {
    // This helps prevent faulty err signals from the sse
    const internalRef = useRef(null);
    const pondRef = refFromParent ?? internalRef;

    const handleOnChange = (err, file) => {
      if (onChange) onChange(err, file);
    };

    const processFile = (fieldName, file, metadata, load, error) => {
      // Setup abort controller to enable future enhancements.
      // As of right now it is not utilized, it should be possible for client-side only
      const controller = new AbortController();
      const { signal } = controller;

      createUpload(file, { signal })
        .then((response) => {
          // eslint-disable-next-line no-underscore-dangle
          const rootEl = pondRef.current._pond.element; // Grab FilePond HTML
          const itemEl = getFilePondItemByFilename(rootEl, file.name); // Grab our file's HTML entry
          if (itemEl) {
            // Manually adjust status as we have successfully uploaded
            const fileProcessingLabel = itemEl.querySelector('.filepond--file-status-main');
            const subLabelWhenProcessing = itemEl.querySelector('.filepond--file-status-sub');
            if (fileProcessingLabel) fileProcessingLabel.textContent = 'Scanning'; // FilePond doesn't offer a state/api to change it mid-processing
            if (subLabelWhenProcessing) subLabelWhenProcessing.textContent = ''; // Change from FilePond's click to abort -> nothing. You can't abort after it's uploaded
          }
          return waitForAvScan(response.id, { signal }).then(() => response);
        })
        .then((response) => {
          // This only triggers on an AV clean response
          load(response.id); // Make FilePond store the server id
        })
        .catch((err) => {
          if (err.name === 'AbortError') return; // controller close
          if (err.message === UPLOAD_SCAN_STATUS.THREATS_FOUND || err.message === UPLOAD_SCAN_STATUS.LEGACY_INFECTED) {
            pondRef.current?.setOptions({
              labelFileProcessing: 'File failed virus scan',
            });
            error('File failed virus scan');
          } else {
            error(err);
          }
        });

      // TODO - in order to handle abort, we need to pass an AbortController to SwaggerRequest, implement this as a future story (it's not working as-is)
      // https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#browser
      return {
        abort: () => {
          // Send the abort signal (Not yet functional)
          controller.abort();
        },
      };
    };

    const revertFile = (uploadId, load, error) => {
      deleteUpload(uploadId)
        .then(() => {
          load();
          handleOnChange();
        })
        .catch(error);
    };

    const handleProcessFile = (err, file) => {
      handleOnChange(err, file);
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
     * component and passed through (like labelIdle).
     * Note that FilePond does not support api/state changing of labels
     * if we want the labels to change mid-processing we must adjust the HTML manually
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
        required={required}
        ref={pondRef}
        {...filePondProps}
        className={className}
        allowMultiple={allowMultiple}
        maxParralelUploads={maxParralelUploads}
        acceptedFileTypes={acceptedFileTypes}
        name={name}
        onprocessfile={handleProcessFile}
        onaddfilestart={onAddFile}
        labelFileTypeNotAllowed={labelFileTypeNotAllowed}
        fileValidateTypeLabelExpectedTypes={fileValidateTypeLabelExpectedTypes}
      />
    );
  },
);

FileUpload.propTypes = {
  required: bool,
  name: string,
  className: string,
  createUpload: func,
  onChange: func,
  onAddFile: func,
  acceptedFileTypes: arrayOf(string),
  allowMultiple: bool,
  maxParralelUploads: number,
  labelFileTypeNotAllowed: string,
  fileValidateTypeLabelExpectedTypes: string,
  // FilePond instance props
  labelIdle: string,
  labelIdleMobile: string,
};

FileUpload.defaultProps = {
  required: false,
  name: 'file',
  className: null,
  createUpload: createUploadApi,
  onChange: undefined,
  onAddFile: undefined,
  allowMultiple: true,
  maxParralelUploads: 2,
  labelFileTypeNotAllowed: 'File of invalid type',
  fileValidateTypeLabelExpectedTypes: 'Expects {allButLastType} or {lastType}',
  acceptedFileTypes: ['image/jpeg', 'image/png', 'application/pdf'],
  labelIdle: 'Drag & drop or <span class="filepond--label-action">click to upload</span>',
  labelIdleMobile: '<span class="filepond--label-action">Upload</span>',
};

export default FileUpload;
