import { React, createRef, useState } from 'react';
import { useHistory } from 'react-router-dom';
// import { Grid, Alert } from '@trussworks/react-uswds';

import SectionWrapper from 'components/Customer/SectionWrapper';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import FileUpload from 'components/FileUpload/FileUpload';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';

const UploadPaymentRequest = () => {
  const filePondEl = createRef();
  const history = useHistory();
  // const { paymentRequestId } = useParams();
  // const [serverError, setServerError] = useState(null);
  const [uploadedFiles, setUploadedFiles] = useState([]);

  const handleDelete = (uploadId) => {
    setUploadedFiles(uploadedFiles.filter((file) => file.id !== uploadId));
  };

  const handleUpload = (file) => {
    setUploadedFiles([
      ...uploadedFiles,
      {
        file,
        filename: file.name,
        url: '',
        bytes: file.size,
        created_at: new Date().toISOString(),
        id: `${Date.now()}${file.name}`,
      },
    ]);
    return Promise.resolve();
  };

  const onChange = () => {
    filePondEl.current?.removeFiles();
  };

  const handleSave = () => {};

  const handleCancel = () => {
    history.push('/');
  };

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
        <UploadsTable uploads={uploadedFiles} onDelete={handleDelete} />
        <WizardNavigation editMode disableNext={false} onNextClick={handleSave} onCancelClick={handleCancel} />
      </SectionWrapper>
    </>
  );
};

export default UploadPaymentRequest;
