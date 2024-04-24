import React, { createRef, useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { useNavigate, useParams } from 'react-router-dom';

import Alert from 'shared/Alert';
import { withContext } from 'shared/AppContext';
import scrollToTop from 'shared/scrollToTop';
import {
  getResponseError,
  patchOrders,
  createUploadForDocument,
  deleteUpload,
  getAllMoves,
  getOrders,
} from 'services/internalApi';
import { updateOrders as updateOrdersAction, updateAllMoves as updateAllMovesAction } from 'store/entities/actions';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import {
  selectServiceMemberFromLoggedInUser,
  selectOrdersForLoggedInUser,
  selectAllMoves,
} from 'store/entities/selectors';
import EditOrdersForm from 'components/Customer/EditOrdersForm/EditOrdersForm';
import { formatWeight, formatYesNoInputValue, dropdownInputOptions } from 'utils/formatters';
import { ORDERS_TYPE_OPTIONS } from 'constants/orders';
import { formatDateForSwagger } from 'shared/dates';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';

const EditOrders = ({
  serviceMemberId,
  serviceMemberMoves,
  updateOrders,
  setFlashMessage,
  context,
  orders,
  updateAllMoves,
}) => {
  const filePondEl = createRef();
  const navigate = useNavigate();
  const { moveId, orderId } = useParams();
  const [serverError, setServerError] = useState('');

  const currentOrder = orders.find((order) => order.moves[0] === moveId);

  const checkIfMoveStatusIsApproved = (status) => {
    return status === 'APPROVED';
  };

  let move;
  let isMoveApproved;
  if (serviceMemberMoves && Object.keys(serviceMemberMoves).length !== 0) {
    const currentMove = serviceMemberMoves.currentMove.find((m) => m.id === moveId);
    const previousMoves = serviceMemberMoves.previousMoves.find((m) => m.id === moveId);
    move = currentMove || previousMoves;
    isMoveApproved = checkIfMoveStatusIsApproved(move.status);
  }

  useEffect(() => {
    const fetchData = async () => {
      getOrders(orderId).then((response) => {
        updateOrders(response);
      });
      getAllMoves(serviceMemberId).then((response) => {
        updateAllMoves(response);
      });
    };
    fetchData();
  }, [updateOrders, serviceMemberId, updateAllMoves, orderId]);

  const initialValues = {
    orders_type: currentOrder?.orders_type || '',
    issue_date: currentOrder?.issue_date || '',
    report_by_date: currentOrder?.report_by_date || '',
    has_dependents: formatYesNoInputValue(currentOrder?.has_dependents),
    new_duty_location: currentOrder?.new_duty_location || null,
    uploaded_orders: currentOrder?.uploaded_orders?.uploads || [],
    move_status: move?.status,
    grade: currentOrder?.grade || null,
    origin_duty_location: currentOrder?.origin_duty_location || {},
  };

  // Only allow PCS unless feature flag is on
  const showAllOrdersTypes = context.flags?.allOrdersTypes;
  const allowedOrdersTypes = showAllOrdersTypes
    ? ORDERS_TYPE_OPTIONS
    : { PERMANENT_CHANGE_OF_STATION: ORDERS_TYPE_OPTIONS.PERMANENT_CHANGE_OF_STATION };
  const ordersTypeOptions = dropdownInputOptions(allowedOrdersTypes);

  const handleUploadFile = (file) => {
    const documentId = currentOrder?.uploaded_orders?.id;
    return createUploadForDocument(file, documentId);
  };

  const handleUploadComplete = async () => {
    filePondEl.current?.removeFiles();
    return getOrders(orderId).then((response) => {
      updateOrders(response);
    });
  };

  const handleDeleteFile = async (uploadId) => {
    return deleteUpload(uploadId, orderId).then(() => {
      return getOrders(orderId).then((response) => {
        updateOrders(response);
      });
    });
  };

  const submitOrders = async (fieldValues) => {
    let hasDependents = false;
    if (fieldValues.has_dependents === 'yes') {
      hasDependents = true;
    }

    const entitlementCouldChange =
      hasDependents !== currentOrder.has_dependents || fieldValues.grade !== currentOrder.grade;
    const newDutyLocationId = fieldValues.new_duty_location.id;
    const newPayGrade = fieldValues.grade;
    const newOriginDutyLocationId = fieldValues.origin_duty_location.id;

    return patchOrders({
      ...fieldValues,
      id: currentOrder.id,
      service_member_id: serviceMemberId,
      has_dependents: hasDependents,
      new_duty_location_id: newDutyLocationId,
      issue_date: formatDateForSwagger(fieldValues.issue_date),
      report_by_date: formatDateForSwagger(fieldValues.report_by_date),
      grade: newPayGrade,
      origin_duty_location_id: newOriginDutyLocationId,
      // spouse_has_pro_gear is not updated by this form but is a required value because the endpoint is shared with the
      // ppm office edit orders
      spouse_has_pro_gear: currentOrder.spouse_has_pro_gear,
    })
      .then((response) => {
        updateOrders(response);
        if (entitlementCouldChange) {
          const weightAllowance = response.authorizedWeight;
          setFlashMessage(
            'EDIT_ORDERS_SUCCESS',
            'info',
            `Your weight entitlement is now ${formatWeight(weightAllowance)}.`,
            'Your changes have been saved. Note that the entitlement has also changed.',
          );
        } else {
          setFlashMessage('EDIT_ORDERS_SUCCESS', 'success', '', 'Your changes have been saved.');
        }
        navigate(-1);
      })
      .catch((e) => {
        // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
        const { response } = e;
        const errorMessage = getResponseError(response, 'failed to update orders due to server error');
        setServerError(errorMessage);
        scrollToTop();
      });
  };

  const handleCancel = () => {
    navigate(-1);
  };

  // early return while api call loads object
  if (!currentOrder) {
    return <LoadingPlaceholder />;
  }

  return (
    <div className="grid-container usa-prose">
      <div className="grid-row">
        <div className="grid-col-12">
          {serverError && (
            <div className="usa-width-one-whole error-message">
              <Alert type="error" heading="An error occurred">
                {serverError}
              </Alert>
            </div>
          )}
          {isMoveApproved && (
            <div className="usa-width-one-whole error-message">
              <Alert type="warning" heading="Your move is approved">
                To make a change to your orders, you will need to contact your local PPPO office.
              </Alert>
            </div>
          )}
          {!isMoveApproved && (
            <div className="usa-width-one-whole" data-testid="edit-orders-form-container">
              <EditOrdersForm
                initialValues={initialValues}
                onSubmit={submitOrders}
                filePondEl={filePondEl}
                createUpload={handleUploadFile}
                onUploadComplete={handleUploadComplete}
                onDelete={handleDeleteFile}
                ordersTypeOptions={ordersTypeOptions}
                currentDutyLocation={currentOrder?.origin_duty_location}
                onCancel={handleCancel}
              />
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

EditOrders.propTypes = {
  setFlashMessage: PropTypes.func.isRequired,
  updateOrders: PropTypes.func.isRequired,
  context: PropTypes.shape({
    flags: PropTypes.shape({
      allOrdersTypes: PropTypes.bool,
    }).isRequired,
  }).isRequired,
};

function mapStateToProps(state) {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const serviceMemberId = serviceMember.id;
  const orders = selectOrdersForLoggedInUser(state);
  const serviceMemberMoves = selectAllMoves(state);

  return {
    serviceMemberId,
    serviceMemberMoves,
    orders,
  };
}

const mapDispatchToProps = {
  updateOrders: updateOrdersAction,
  updateAllMoves: updateAllMovesAction,
  setFlashMessage: setFlashMessageAction,
};

export default withContext(connect(mapStateToProps, mapDispatchToProps)(EditOrders));
