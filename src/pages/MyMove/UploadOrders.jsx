import React, { useEffect, useRef, useState } from 'react';
import { connect } from 'react-redux';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { generatePath, useNavigate, useParams } from 'react-router';

import { isBooleanFlagEnabled } from '../../utils/featureFlags';

import './UploadOrders.css';

import FileUpload from 'components/FileUpload/FileUpload';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import { documentSizeLimitMsg } from 'shared/constants';
import { createUploadForDocument, deleteUpload, getAllMoves, getOrders } from 'services/internalApi';
import { updateOrders as updateOrdersAction, updateAllMoves as updateAllMovesAction } from 'store/entities/actions';
import { selectOrdersForLoggedInUser, selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { customerRoutes } from 'constants/routes';
import formStyles from 'styles/form.module.scss';
import { withContext } from 'shared/AppContext';

const UploadOrders = ({ orders, updateOrders, updateAllMoves, serviceMemberId }) => {
  const filePondEl = useRef();
  const navigate = useNavigate();
  const { orderId } = useParams();
  const currentOrders = orders.find((order) => order.id === orderId);
  const uploads = currentOrders?.uploaded_orders?.uploads || [];
  const [multiMove, setMultiMove] = useState(false);

  const handleUploadFile = (file) => {
    const documentId = currentOrders?.uploaded_orders?.id;

    const now = new Date();
    const timestamp =
      now.getFullYear().toString() +
      (now.getMonth() + 1).toString().padStart(2, '0') +
      now.getDate().toString().padStart(2, '0') +
      now.getHours().toString().padStart(2, '0') +
      now.getMinutes().toString().padStart(2, '0') +
      now.getSeconds().toString().padStart(2, '0');

    // Create a new filename with the timestamp prepended
    const newFileName = `${file.name}-${timestamp}`;

    // Create and return a new File object with the new filename
    const newFile = new File([file], newFileName, { type: file.type });

    return createUploadForDocument(newFile, documentId);
  };

  const handleUploadComplete = async () => {
    filePondEl.current?.removeFiles();
    return getOrders(orderId).then((response) => {
      updateOrders(response);
    });
  };

  const handleDeleteFile = async (uploadId) => {
    return deleteUpload(uploadId, orderId).then(() => {
      getOrders(orderId).then((response) => {
        updateOrders(response);
      });
    });
  };

  const onChange = () => {
    filePondEl.current?.removeFiles();
    handleUploadComplete();
  };

  useEffect(() => {
    const fetchData = async () => {
      await getOrders(orderId).then((response) => {
        updateOrders(response);
      });
      await getAllMoves(serviceMemberId).then((response) => {
        updateAllMoves(response);
      });
      isBooleanFlagEnabled('multi_move').then((enabled) => {
        setMultiMove(enabled);
      });
    };
    fetchData();
  }, [updateOrders, orderId, serviceMemberId, updateAllMoves]);

  if (!currentOrders || !uploads) return <LoadingPlaceholder />;

  const isValid = !!uploads.length;

  const handleBack = () => {
    const moveId = currentOrders.moves[0];
    if (multiMove) {
      navigate(generatePath(customerRoutes.MOVE_HOME_PATH, { moveId }));
    } else {
      navigate(customerRoutes.MOVE_HOME_PAGE);
    }
  };
  const handleNext = () => {
    const moveId = currentOrders.moves[0];
    navigate(generatePath(customerRoutes.MOVE_HOME_PATH, { moveId }));
  };

  return (
    <GridContainer>
      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }} data-testid="upload-orders-container">
          <h1>Upload your orders</h1>
          <p>In order to schedule your move, we need to have a complete copy of your orders.</p>
          <p>You can upload a PDF, or you can take a picture of each page and upload the images.</p>
          <p>{documentSizeLimitMsg}</p>

          {uploads?.length > 0 && (
            <>
              <br />
              <UploadsTable uploads={uploads} onDelete={handleDeleteFile} />
            </>
          )}

          <div className="uploader-box">
            <FileUpload
              ref={filePondEl}
              createUpload={handleUploadFile}
              onChange={onChange}
              labelIdle={'Drag & drop or <span class="filepond--label-action">click to upload orders</span>'}
            />
            <div className="hint">(Each page must be clear and legible.)</div>
          </div>

          <div className={formStyles.formActions}>
            <WizardNavigation onBackClick={handleBack} disableNext={!isValid} onNextClick={handleNext} />
          </div>
        </Grid>
      </Grid>
    </GridContainer>
  );
};

const mapStateToProps = (state) => {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const serviceMemberId = serviceMember.id;
  const orders = selectOrdersForLoggedInUser(state);

  return {
    serviceMemberId,
    orders,
  };
};

const mapDispatchToProps = {
  updateOrders: updateOrdersAction,
  updateAllMoves: updateAllMovesAction,
};

export default withContext(connect(mapStateToProps, mapDispatchToProps)(UploadOrders));
