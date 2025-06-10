import { React, createRef, useEffect, useState } from 'react';
import { GridContainer, Grid, Alert, Button } from '@trussworks/react-uswds';
import { useNavigate, useParams } from 'react-router-dom';
import { connect } from 'react-redux';
import classNames from 'classnames';

import { filepondButtonStyle, filepondWrapperStyle } from '../UploadOrders';

import styles from './AdditionalDocuments.module.scss';

import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import Hint from 'components/Hint';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import FileUpload from 'components/FileUpload/FileUpload';
import {
  createUploadForAdditionalDocuments,
  deleteAdditionalDocumentUpload,
  getMove,
  getResponseError,
} from 'services/internalApi';
import { selectMovesForLoggedInUser } from 'store/entities/selectors';
import { updateMove as updateMoveAction } from 'store/entities/actions';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import scrollToTop from 'shared/scrollToTop';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import appendTimestampToFilename from 'utils/fileUpload';
import { primaryButtonStyle } from 'shared/standardUI/Buttons/ButtonUsa';

const desktopFileUploadActionElement = `<div class='${filepondWrapperStyle}'>
    <span>Dragfiles here or </span>
    <button class='${filepondButtonStyle}'>Choose from folder</button>
  </div>`;

const mobileFileUploadActionElement = `<div>
    <button class='${filepondButtonStyle}'>Upload files</span>
  </div>`;

const AdditionalDocuments = ({ moves, updateMove }) => {
  const { moveId } = useParams();
  const filePondEl = createRef();
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(true);
  const [serverError, setServerError] = useState(null);
  const currentlyViewedMove = moves.find((move) => move.id === moveId);
  const uploads = currentlyViewedMove?.additionalDocuments?.uploads;

  const handleDelete = async (uploadId) => {
    return deleteAdditionalDocumentUpload(uploadId, moveId).then(() => {
      getMove(moveId).then((res) => {
        updateMove(res);
      });
    });
  };

  const handleUpload = async (file) => {
    return createUploadForAdditionalDocuments(appendTimestampToFilename(file), moveId);
  };

  const handleUploadComplete = () => {
    getMove(moveId)
      .then((res) => {
        updateMove(res);
      })
      .catch((e) => {
        const { response } = e;
        const error = getResponseError(response, 'failed to upload due to server error');
        setServerError(error);

        scrollToTop();
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
      <NotificationScrollToTop dependency={serverError} />

      {serverError && (
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Alert type="error" headingLevel="h4" heading="An error occurred">
              {serverError}
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
            <>
              <br />
              <UploadsTable uploads={uploads} onDelete={handleDelete} showDownloadLink />
            </>
            <div className="uploader-box">
              <FileUpload
                ref={filePondEl}
                createUpload={handleUpload}
                onChange={onChange}
                labelIdle={desktopFileUploadActionElement}
                labelIdleMobile={mobileFileUploadActionElement}
              />
            </div>
            <Button className={classNames(primaryButtonStyle, styles['navigation-back-btn'])} onClick={handleBack}>
              Back
            </Button>
          </SectionWrapper>
        </Grid>
      </Grid>
    </GridContainer>
  );
};

const mapStateToProps = (state) => {
  const moves = selectMovesForLoggedInUser(state);

  const props = { moves };

  return props;
};

const mDTP = {
  updateMove: updateMoveAction,
};

export default connect(mapStateToProps, mDTP)(AdditionalDocuments);
