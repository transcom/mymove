import { React, createRef, useEffect, useState } from 'react';
import { GridContainer, Grid, Alert, Button } from '@trussworks/react-uswds';
import { useNavigate, useParams } from 'react-router-dom';
import { connect } from 'react-redux';
import { generatePath } from 'react-router';

import NotificationScrollToTop from 'components/NotificationScrollToTop';
import SectionWrapper from 'components/Customer/SectionWrapper';
import Hint from 'components/Hint';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import FileUpload from 'components/FileUpload/FileUpload';
import { createUploadForAdditionalDocuments, deleteAdditionalDocumentUpload, getMove } from 'services/internalApi';
import { selectCurrentMove } from 'store/entities/selectors';
import { customerRoutes } from 'constants/routes';
import { updateMove as updateMoveAction } from 'store/entities/actions';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';

const AdditionalDocuments = ({ move, updateMove }) => {
  const [isLoading, setIsLoading] = useState(true);
  const moveId = move?.id;
  const filePondEl = createRef();
  const navigate = useNavigate();
  const [errorMessage, setErrorMessage] = useState(null);
  const uploads = move?.additionalDocuments?.uploads;

  const handleDelete = async (uploadId) => {
    return deleteAdditionalDocumentUpload(uploadId, moveId).then(() => {
      getMove(moveId).then((res) => {
        updateMove(res);
      });
    });
  };

  const handleUpload = async (file) => {
    return createUploadForAdditionalDocuments(file, moveId);
  };

  const handleUploadComplete = () => {
    getMove(moveId).then((res) => {
      updateMove(res);
    });
  };

  const onChange = () => {
    filePondEl.current?.removeFiles();
    handleUploadComplete();
  };

  const handleBack = () => {
    navigate(-1);
  };

  useEffect(() => {
    getMove(moveId).then((res) => {
      updateMove(res);
      setIsLoading(false);
    });
  }, [updateMove, moveId]);

  if (isLoading) return <LoadingPlaceholder />;

  const warningMessage =
    'Documents uploaded here will not amend a customers move. Please upload new orders/amendments via the "Upload documents" link next to the Orders section of the customers move.';

  return (
    <GridContainer>
      <NotificationScrollToTop dependency={errorMessage} />

      {errorMessage && (
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Alert type="error" headingLevel="h4" heading="An error occurred">
              {errorMessage}
            </Alert>
          </Grid>
        </Grid>
      )}

      <Grid row data-testid="info-container">
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <h1>Additional Documents</h1>
          <p>Upload any additional documentation that may help your services counselor complete your request.</p>
          <Alert type="info" headingLevel="" heading="">
            {warningMessage}
          </Alert>
        </Grid>
      </Grid>
      <Grid row data-testid="upload-info-container">
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <SectionWrapper>
            <h5>Upload documents</h5>
            <Hint>PDF, JPG, or PNG only. Maximum file size 25MB. Each page must be clear and legible</Hint>
            {/* {uploads?.length > 0 && ( */}
            <>
              <br />
              <UploadsTable uploads={uploads} onDelete={handleDelete} />
            </>
            {/* )} */}
            <div className="uploader-box">
              <FileUpload
                ref={filePondEl}
                createUpload={handleUpload}
                onChange={onChange}
                labelIdle={`Drag files here or <span class="filepond--label-action">choose from folder</span>`}
                labelIdleMobile={`<span class="filepond--label-action">Upload files</span>`}
              />
            </div>
            <Button onClick={handleBack}>Back</Button>
            {/* <WizardNavigation editMode disableNext={false} onBackClick={handleBack} /> */}
          </SectionWrapper>
        </Grid>
      </Grid>
    </GridContainer>
  );
};

const mapStateToProps = (state) => {
  const move = selectCurrentMove(state);

  const props = { move };

  return props;
};

const mDTP = {
  updateMove: updateMoveAction,
};

export default connect(mapStateToProps, mDTP)(AdditionalDocuments);
