import { React, createRef, useEffect, useState } from 'react';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';
import { useNavigate, useParams } from 'react-router-dom';
import { connect } from 'react-redux';
import { generatePath } from 'react-router';

import NotificationScrollToTop from 'components/NotificationScrollToTop';
import SectionWrapper from 'components/Customer/SectionWrapper';
import Hint from 'components/Hint';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import FileUpload from 'components/FileUpload/FileUpload';
import { createUploadForAdditionalDocuments, getMove } from 'services/internalApi';
import { selectCurrentMove } from 'store/entities/selectors';
import { customerRoutes } from 'constants/routes';
import { updateMove } from 'store/entities/actions';

const AdditionalDocuments = ({ locator }) => {
  const { moveId, additionalDocumentsId } = useParams();
  const filePondEl = createRef();
  const navigate = useNavigate();
  const [errorMessage, setErrorMessage] = useState(null);

  const uploads = [];

  // useEffect(() => {
  //   // if no existing documents create a new one
  //   if (!additionalDocumentsId) {
  //     createAdditionalDocuments(locator)
  //       .then((res) => {
  //         const path = generatePath(customerRoutes.UPLOAD_ADDITIONAL_DOCUMENTS_EDIT_PATH, {
  //           additionalDocumentsId: res.id,
  //         });
  //         navigate(path, { replace: true });
  //       })
  //       .catch(() => {});
  //   }
  // });

  const handleDelete = async (documentId) => {
    // return deleteAdditionalDocuments(documentId).then(() => {
    //   getMove(moveId).then((res) => {
    //     updateMove(res);
    //   });
    // });
  };

  const handleUpload = async (file) => {
    return createUploadForAdditionalDocuments(file, '40b13561-349d-4e71-a051-eb2b9d69bdd4');
  };

  // const handleUploadComplete = (err) => {
  //   if (err) {
  //     setErrorMessage('Encountered error when completing file upload');
  //   }
  // };

  const onChange = () => {
    // filePondEl.current?.removeFiles();
    // handleUploadComplete();
  };

  const handleSave = () => {
    navigate(-1);
  };

  const handleCancel = () => {
    navigate(-1);
  };

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
              <UploadsTable uploads={uploads} onDelete={handleDelete} showDeleteButton={false} />
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
            <WizardNavigation editMode disableNext={false} onNextClick={handleSave} onCancelClick={handleCancel} />
          </SectionWrapper>
        </Grid>
      </Grid>
    </GridContainer>
  );
};

const mapStateToProps = (state) => {
  const { locator } = selectCurrentMove(state);

  const props = { locator };

  return props;
};

export default connect(mapStateToProps)(AdditionalDocuments);
