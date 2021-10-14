import { React, createRef, useState } from 'react';
// import { useHistory } from 'react-router-dom';
// import { Grid, Alert } from '@trussworks/react-uswds';

import SectionWrapper from 'components/Customer/SectionWrapper';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import FileUpload from 'components/FileUpload/FileUpload';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';

const UploadPaymentRequest = () => {
  const filePondEl = createRef();
  // const history = useHistory();
  // const [serverError, setServerError] = useState(null);
  const [uploadedFiles, setUploadedFiles] = useState([]);

  const handleDelete = (uploadId) => {
    setUploadedFiles(uploadedFiles.filter((file) => file.id !== uploadId));
  };

  const handleUpload = (file) => {
    // console.log(file);
    setUploadedFiles([...uploadedFiles, { file, id: file.lastModified }]);
    return Promise.resolve();
  };

  const onChange = () => {
    filePondEl.current?.removeFiles();
  };

  const handleSave = () => {};

  const handleCancel = () => {};

  return (
    <>
      {/*
      {serverError && (
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Alert type="error" heading="An error occurred">
              {serverError}
            </Alert>
          </Grid>
        </Grid>
      )}
      */}

      <SectionWrapper>
        <UploadsTable uploads={uploadedFiles} onDelete={handleDelete} />
        <div>
          <h2>Upload Payment Request Documents</h2>
          <FileUpload
            ref={filePondEl}
            createUpload={handleUpload}
            onChange={onChange}
            labelIdle={
              'Drag & drop or <span class="filepond--label-action">click to upload payment request documents</span>'
            }
          />
        </div>
        <WizardNavigation editMode disableNext={false} onNextClick={handleSave} onCancelClick={handleCancel} />
      </SectionWrapper>
    </>
  );
};

export default UploadPaymentRequest;
