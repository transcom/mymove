import { React, createRef, useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { connect } from 'react-redux';

import Hint from 'components/Hint';
import SectionWrapper from 'components/Customer/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import ScrollToTop from 'components/ScrollToTop';
import FileUpload from 'components/FileUpload/FileUpload';
import { UploadsShape, OrdersShape } from 'types/customerShapes';
import { getOrdersForServiceMember, createUploadForDocument, deleteUpload } from 'services/internalApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import {
  selectCurrentOrders,
  selectServiceMemberFromLoggedInUser,
  selectUploadsForCurrentOrders,
} from 'store/entities/selectors';
import { updateOrders as updateOrdersAction } from 'store/entities/actions';

export const AmendOrders = ({ uploads, updateOrders, serviceMemberId, currentOrders }) => {
  const [isLoading, setLoading] = useState(true);
  const filePondEl = createRef();

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
    const documentId = currentOrders?.uploaded_orders?.id;
    return createUploadForDocument(file, documentId);
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
    // push(generalRoutes.HOME_PATH);
  };
  const handleCancel = () => {
    // push(generalRoutes.HOME_PATH);
  };

  useEffect(() => {
    getOrdersForServiceMember(serviceMemberId).then((response) => {
      updateOrders(response);
      setLoading(false);
    });
  });

  if (isLoading) return <LoadingPlaceholder />;

  return (
    <GridContainer>
      <ScrollToTop />
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
  // push: PropTypes.func.isRequired,
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
    uploads: selectUploadsForCurrentOrders(state),
  };

  return props;
}

const mapDispatchToProps = {
  // TODO we might need a new action to handle updating amended orders
  updateOrders: updateOrdersAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(AmendOrders);
