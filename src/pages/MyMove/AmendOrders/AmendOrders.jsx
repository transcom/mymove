import { React, createRef, useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';
import { connect } from 'react-redux';
import { useHistory } from 'react-router-dom';

import Hint from 'components/Hint';
import SectionWrapper from 'components/Customer/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import ScrollToTop from 'components/ScrollToTop';
import FileUpload from 'components/FileUpload/FileUpload';
import { UploadsShape, OrdersShape } from 'types/customerShapes';
import {
  getOrdersForServiceMember,
  createAmendedOrdersUploadForDocument,
  deleteUpload,
  getResponseError,
  submitAmendedOrders,
} from 'services/internalApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import scrollToTop from 'shared/scrollToTop';
import {
  selectCurrentOrders,
  selectServiceMemberFromLoggedInUser,
  selectUploadsForCurrentAmendedOrders,
} from 'store/entities/selectors';
import { updateOrders as updateOrdersAction } from 'store/entities/actions';
import { generalRoutes } from 'constants/routes';

export const AmendOrders = ({ uploads, updateOrders, serviceMemberId, currentOrders }) => {
  const [isLoading, setLoading] = useState(true);
  const filePondEl = createRef();
  const history = useHistory();
  const [serverError, setServerError] = useState(null);

  const handleDelete = (uploadId) => {
    return deleteUpload(uploadId).then(() => {
      // TODO Temporarily using the original uploaded orders, will change to use amended orders once that is available
      getOrdersForServiceMember(serviceMemberId).then((response) => {
        updateOrders(response);
      });
    });
  };
  const handleUpload = (file) => {
    // TODO Temporarily using the original uploaded orders, will change to use amended orders once that is available
    const documentId = currentOrders?.uploaded_amended_orders?.id;
    return createAmendedOrdersUploadForDocument(file, documentId);
  };
  const handleUploadComplete = () => {
    // TODO Temporarily using the original uploaded orders, will change to use amended orders once that is available
    getOrdersForServiceMember(serviceMemberId).then((response) => {
      updateOrders(response);
    });
  };

  const onChange = () => {
    filePondEl.current?.removeFiles();
    handleUploadComplete();
  };

  const handleSave = () => {
    return submitAmendedOrders(currentOrders?.moves[0])
      .then(() => {
        history.push(generalRoutes.HOME_PATH);
      })
      .catch((e) => {
        // TODO - error handling - below is rudimentary error handling to approximate existing UX
        // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
        const { response } = e;
        const errorMessage = getResponseError(response, 'failed to save amended orders due to server error');
        setServerError(errorMessage);

        scrollToTop();
      });
  };
  const handleCancel = () => {
    // TODO (After MB-8336 is complete) Delete amended orders files before navigating away
    history.push(generalRoutes.HOME_PATH);
  };

  useEffect(() => {
    getOrdersForServiceMember(serviceMemberId).then((response) => {
      updateOrders(response);
      setLoading(false);
    });
  }, [updateOrders, serviceMemberId]);

  if (isLoading) return <LoadingPlaceholder />;

  return (
    <GridContainer>
      <ScrollToTop otherDep={serverError} />

      {serverError && (
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Alert type="error" heading="An error occurred">
              {serverError}
            </Alert>
          </Grid>
        </Grid>
      )}

      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <h1>Orders</h1>
          <p>
            Upload any amended orders here. The office will update your move info to match the new orders. Talk directly
            with your movers to coordinate any changes.
          </p>
        </Grid>
      </Grid>
      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <SectionWrapper>
            <h5>Upload orders</h5>
            <Hint>PDF, JPG, or PNG only. Maximum file size 25MB. Each page must be clear and legible</Hint>
            {uploads && uploads.length > 0 && (
              <>
                <br />
                <UploadsTable uploads={uploads} onDelete={handleDelete} />
              </>
            )}
            <div className="uploader-box">
              <FileUpload
                ref={filePondEl}
                createUpload={handleUpload}
                onChange={onChange}
                labelIdle={'Drag files here or <span class="filepond--label-action">choose from folder</span>'}
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
  currentOrders: OrdersShape,
  uploads: UploadsShape,
};

AmendOrders.defaultProps = {
  uploads: [],
  currentOrders: {},
};

function mapStateToProps(state) {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const serviceMemberId = serviceMember?.id;
  const currentOrders = selectCurrentOrders(state);

  const props = {
    serviceMemberId,
    currentOrders,
    uploads: selectUploadsForCurrentAmendedOrders(state),
  };

  return props;
}

const mapDispatchToProps = {
  // TODO we might need a new action to handle updating amended orders
  updateOrders: updateOrdersAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(AmendOrders);
