import { React, createRef, useState } from 'react';
import { useNavigate, useParams, generatePath } from 'react-router-dom';
import { Grid, Alert } from '@trussworks/react-uswds';
import { useMutation } from '@tanstack/react-query';
import { func } from 'prop-types';
import { connect } from 'react-redux';

import styles from './UploadServiceRequestDocuments.module.scss';

import formStyles from 'styles/form.module.scss';
import { createServiceRequestDocumentUpload } from 'services/primeApi';
import { primeSimulatorRoutes } from 'constants/routes';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import FileUpload from 'components/FileUpload/FileUpload';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';

const UploadServiceRequest = ({ setFlashMessage }) => {
  const { moveCodeOrID } = useParams();
  const filePondEl = createRef();
  const navigate = useNavigate();
  const { mtoServiceItemId } = useParams();
  const [serverError, setServerError] = useState(null);
  // Despite this being named plurarly, only one file is allowed to be uploaded at a time
  // since the endpoint being called only allows one upload at a time.
  const [filesToUpload, setFilesToUpload] = useState([]);
  const [uploadSuccess, setUploadSuccess] = useState(false);

  const handleDelete = () => {
    setFilesToUpload([]);
  };

  const { mutate: mutateUploadServiceRequestDocument } = useMutation(createServiceRequestDocumentUpload, {
    onSuccess: () => {
      // TODO - show flash message?
      setUploadSuccess(true);

      setFlashMessage(`MSG_UPLOAD_DOC_SUCCESS${moveCodeOrID}`, 'success', 'Successfully uploaded document', '', true);

      navigate(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
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
    mutateUploadServiceRequestDocument({ mtoServiceItemID: mtoServiceItemId, file: filesToUpload[0].file });
  };

  const handleCancel = () => {
    navigate(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
  };

  return (
    <>
      {serverError && (
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Alert type="error" headingLevel="h4" heading="An error occurred">
              {serverError}
            </Alert>
          </Grid>
        </Grid>
      )}

      {uploadSuccess && (
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Alert headingLevel="h4" type="success">
              Upload saved successfully
            </Alert>
          </Grid>
        </Grid>
      )}

      <div>
        <SectionWrapper className={styles.container}>
          <div>
            <h2>Upload Service Request Document</h2>
            <FileUpload
              ref={filePondEl}
              createUpload={handleUpload}
              onChange={onChange}
              labelIdle='Drag & drop or <span class="filepond--label-action">click to upload a service request document</span>'
            />
          </div>
          <UploadsTable uploads={filesToUpload} onDelete={handleDelete} />
          <div className={formStyles.formActions}>
            <WizardNavigation editMode disableNext={false} onNextClick={handleSave} onCancelClick={handleCancel} />
          </div>
        </SectionWrapper>
      </div>
    </>
  );
};

UploadServiceRequest.propTypes = {
  setFlashMessage: func.isRequired,
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
};

export default connect(() => ({}), mapDispatchToProps)(UploadServiceRequest);
