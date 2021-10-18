import { React, createRef, useState } from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { Grid, Alert } from '@trussworks/react-uswds';
import { useMutation } from 'react-query';

import { createUpload } from 'services/primeApi';
import SectionWrapper from 'components/Customer/SectionWrapper';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import FileUpload from 'components/FileUpload/FileUpload';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';

const UploadPaymentRequest = () => {
  const filePondEl = createRef();
  const history = useHistory();
  const { paymentRequestId } = useParams();
  const [serverError, setServerError] = useState(null);
  // Despite this being named plurarly, only one file is allowed to be uploaded at a time
  // since the endpoint being called only allows one upload at a time.
  const [filesToUpload, setFilesToUpload] = useState([]);

  const handleDelete = () => {
    setFilesToUpload([]);
  };

  const [mutateUploadPaymentRequestDocument] = useMutation(createUpload, {
    onSuccess: () => {
      // TODO - show flash message?
      history.push(`/`);
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      setServerError(errorMsg);
    },
  });

  const handleUpload = (file) => {
    setFilesToUpload([
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

  const handleSave = () => {
    mutateUploadPaymentRequestDocument({ paymentRequestID: paymentRequestId, file: filesToUpload[0].file });
  };

  const handleCancel = () => {
    history.push('/');
  };

  return (
    <>
      {serverError && (
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Alert type="error" heading="An error occurred">
              {serverError}
            </Alert>
          </Grid>
        </Grid>
      )}

      <SectionWrapper>
        <div>
          <h2>Upload Payment Request Document</h2>
          <FileUpload
            ref={filePondEl}
            createUpload={handleUpload}
            onChange={onChange}
            labelIdle={
              'Drag & drop or <span class="filepond--label-action">click to upload a payment request document</span>'
            }
          />
        </div>
        <UploadsTable uploads={filesToUpload} onDelete={handleDelete} />
        <WizardNavigation editMode disableNext={false} onNextClick={handleSave} onCancelClick={handleCancel} />
      </SectionWrapper>
    </>
  );
};

export default UploadPaymentRequest;
