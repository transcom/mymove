import { React, createRef, useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';
import { connect } from 'react-redux';
import { generatePath, useNavigate, useParams } from 'react-router-dom';

import styles from './AmendOrders.module.scss';

import Hint from 'components/Hint';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import FileUpload from 'components/FileUpload/FileUpload';
import {
  createUploadForAmendedOrdersDocument,
  deleteUpload,
  getResponseError,
  submitAmendedOrders,
  getOrders,
} from 'services/internalApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import scrollToTop from 'shared/scrollToTop';
import {
  selectOrdersForLoggedInUser,
  selectServiceMemberFromLoggedInUser,
  selectUploadsForCurrentAmendedOrders,
} from 'store/entities/selectors';
import { updateOrders as updateOrdersAction } from 'store/entities/actions';
import { customerRoutes } from 'constants/routes';
import appendTimestampToFilename from 'utils/fileUpload';

export const AmendOrders = ({ updateOrders, serviceMemberId, orders }) => {
  const [isLoading, setLoading] = useState(true);
  const filePondEl = createRef();
  const navigate = useNavigate();
  const { orderId } = useParams();
  const [serverError, setServerError] = useState(null);
  const currentOrders = orders.find((order) => order.id === orderId);
  const uploads = currentOrders?.uploaded_amended_orders?.uploads;

  const handleDelete = async (uploadId) => {
    return deleteUpload(uploadId, orderId).then(() => {
      getOrders(orderId).then((response) => {
        updateOrders(response);
      });
    });
  };

  const handleUpload = (file) => {
    return createUploadForAmendedOrdersDocument(appendTimestampToFilename(file), orderId);
  };

  const handleUploadComplete = () => {
    getOrders(orderId).then((response) => {
      updateOrders(response);
    });
  };

  const onChange = () => {
    filePondEl.current?.removeFiles();
    handleUploadComplete();
  };

  const handleSave = async () => {
    return submitAmendedOrders(currentOrders?.moves[0])
      .then(() => {
        const moveId = currentOrders?.moves[0];
        navigate(generatePath(customerRoutes.MOVE_HOME_PATH, { moveId }));
      })
      .catch((e) => {
        const { response } = e;
        const errorMessage = getResponseError(response, 'failed to save amended orders due to server error');
        setServerError(errorMessage);

        scrollToTop();
      });
  };
  const handleCancel = () => {
    navigate(-1);
  };

  useEffect(() => {
    getOrders(orderId).then((response) => {
      updateOrders(response);
      setLoading(false);
    });
  }, [updateOrders, serviceMemberId, orderId]);

  if (isLoading) return <LoadingPlaceholder />;

  const additionalText = uploads && uploads.length > 0 ? 'additional ' : '';

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
          <h1>Orders</h1>
          <p>
            Upload any amended orders here. The office will update your move info to match the new orders. Talk directly
            with your movers to coordinate any changes.
          </p>
        </Grid>
      </Grid>
      <Grid row data-testid="upload-info-container">
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <SectionWrapper>
            <h5 className={styles.uploadOrdersHeader}>Upload orders</h5>
            <Hint>PDF, JPG, or PNG only. Maximum file size 25MB. Each page must be clear and legible</Hint>
            {uploads?.length > 0 && (
              <>
                <br />
                <UploadsTable uploads={uploads} onDelete={handleDelete} showDeleteButton={false} showDownloadLink />
              </>
            )}
            <div className="uploader-box">
              <FileUpload
                ref={filePondEl}
                createUpload={handleUpload}
                onChange={onChange}
                labelIdle={`Drag ${additionalText}files here or <span class="filepond--label-action">choose from folder</span>`}
                labelIdleMobile={`<span class="filepond--label-action">Upload ${additionalText}files</span>`}
              />
            </div>
            <WizardNavigation editMode disableNext={false} onNextClick={handleSave} onCancelClick={handleCancel} />
          </SectionWrapper>
        </Grid>
      </Grid>
    </GridContainer>
  );
};

AmendOrders.propTypes = {
  serviceMemberId: PropTypes.string.isRequired,
  updateOrders: PropTypes.func.isRequired,
};

function mapStateToProps(state) {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const serviceMemberId = serviceMember?.id;
  const orders = selectOrdersForLoggedInUser(state);

  const props = {
    serviceMemberId,
    orders,
    uploads: selectUploadsForCurrentAmendedOrders(state),
  };

  return props;
}

const mapDispatchToProps = {
  updateOrders: updateOrdersAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(AmendOrders);
