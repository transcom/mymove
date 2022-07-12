import React, { createRef, useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { useHistory } from 'react-router-dom';

import Alert from 'shared/Alert';
import { withContext } from 'shared/AppContext';
import scrollToTop from 'shared/scrollToTop';
import {
  getResponseError,
  getOrdersForServiceMember,
  patchOrders,
  createUploadForDocument,
  deleteUpload,
} from 'services/internalApi';
import { updateOrders as updateOrdersAction } from 'store/entities/actions';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import {
  selectServiceMemberFromLoggedInUser,
  selectCurrentOrders,
  selectMoveIsApproved,
  selectUploadsForCurrentOrders,
  selectHasCurrentPPM,
} from 'store/entities/selectors';
import EditOrdersForm from 'components/Customer/EditOrdersForm/EditOrdersForm';
import { OrdersShape, ServiceMemberShape } from 'types/customerShapes';
import { formatWeight, formatYesNoInputValue, dropdownInputOptions } from 'utils/formatters';
import { ORDERS_TYPE_OPTIONS } from 'constants/orders';
import { ExistingUploadsShape } from 'types/uploads';
import { formatDateForSwagger } from 'shared/dates';

export const EditOrders = ({
  serviceMember,
  currentOrders,
  updateOrders,
  existingUploads,
  moveIsApproved,
  setFlashMessage,
  context,
}) => {
  const filePondEl = createRef();
  const history = useHistory();
  const [serverError, setServerError] = useState(null);

  const initialValues = {
    orders_type: currentOrders?.orders_type || '',
    issue_date: currentOrders?.issue_date || '',
    report_by_date: currentOrders?.report_by_date || '',
    has_dependents: formatYesNoInputValue(currentOrders?.has_dependents),
    new_duty_location: currentOrders?.new_duty_location || null,
    uploaded_orders: existingUploads || [],
  };

  // Only allow PCS unless feature flag is on
  const showAllOrdersTypes = context.flags?.allOrdersTypes;
  const allowedOrdersTypes = showAllOrdersTypes
    ? ORDERS_TYPE_OPTIONS
    : { PERMANENT_CHANGE_OF_STATION: ORDERS_TYPE_OPTIONS.PERMANENT_CHANGE_OF_STATION };
  const ordersTypeOptions = dropdownInputOptions(allowedOrdersTypes);

  const serviceMemberId = serviceMember.id;

  useEffect(() => {
    getOrdersForServiceMember(serviceMemberId).then((response) => {
      updateOrders(response);
    });
  }, [updateOrders, serviceMemberId]);

  const handleUploadFile = (file) => {
    const documentId = currentOrders?.uploaded_orders?.id;
    return createUploadForDocument(file, documentId);
  };

  const handleUploadComplete = () => {
    filePondEl.current?.removeFiles();
    return getOrdersForServiceMember(serviceMemberId).then((response) => {
      updateOrders(response);
    });
  };

  const handleDeleteFile = (uploadId) => {
    return deleteUpload(uploadId).then(() => {
      getOrdersForServiceMember(serviceMemberId).then((response) => {
        updateOrders(response);
      });
    });
  };

  const submitOrders = (fieldValues) => {
    const hasDependents = fieldValues.has_dependents === 'yes';
    const entitlementCouldChange = hasDependents !== currentOrders.has_dependents;
    const newDutyLocationId = fieldValues.new_duty_location.id;

    return patchOrders({
      ...fieldValues,
      id: currentOrders.id,
      service_member_id: serviceMember.id,
      has_dependents: hasDependents,
      new_duty_location_id: newDutyLocationId,
      issue_date: formatDateForSwagger(fieldValues.issue_date),
      report_by_date: formatDateForSwagger(fieldValues.report_by_date),
      // spouse_has_pro_gear is not updated by this form but is a required value because the endpoint is shared with the
      // ppm office edit orders
      spouse_has_pro_gear: currentOrders.spouse_has_pro_gear,
    })
      .then((response) => {
        updateOrders(response);
        if (entitlementCouldChange) {
          const weightAllowance = hasDependents
            ? serviceMember.weight_allotment.total_weight_self_plus_dependents
            : serviceMember.weight_allotment.total_weight_self;
          setFlashMessage(
            'EDIT_ORDERS_SUCCESS',
            'info',
            `Your weight entitlement is now ${formatWeight(weightAllowance)}.`,
            'Your changes have been saved. Note that the entitlement has also changed.',
          );
        } else {
          setFlashMessage('EDIT_ORDERS_SUCCESS', 'success', '', 'Your changes have been saved.');
        }
        history.goBack();
      })
      .catch((e) => {
        // TODO - error handling - below is rudimentary error handling to approximate existing UX
        // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
        const { response } = e;
        const errorMessage = getResponseError(response, 'failed to update orders due to server error');
        setServerError(errorMessage);
        scrollToTop();
      });
  };

  const handleCancel = () => {
    history.goBack();
  };

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
          {moveIsApproved && (
            <div className="usa-width-one-whole error-message">
              <Alert type="warning" heading="Your move is approved">
                To make a change to your orders, you will need to contact your local PPPO office.
              </Alert>
            </div>
          )}
          {!moveIsApproved && (
            <div className="usa-width-one-whole">
              <EditOrdersForm
                initialValues={initialValues}
                onSubmit={submitOrders}
                filePondEl={filePondEl}
                createUpload={handleUploadFile}
                onUploadComplete={handleUploadComplete}
                onDelete={handleDeleteFile}
                ordersTypeOptions={ordersTypeOptions}
                currentDutyLocation={serviceMember.current_location}
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
  moveIsApproved: PropTypes.bool.isRequired,
  serviceMember: ServiceMemberShape.isRequired,
  setFlashMessage: PropTypes.func.isRequired,
  updateOrders: PropTypes.func.isRequired,
  currentOrders: OrdersShape.isRequired,
  existingUploads: ExistingUploadsShape,
  context: PropTypes.shape({
    flags: PropTypes.shape({
      allOrdersTypes: PropTypes.bool,
    }).isRequired,
  }).isRequired,
};

EditOrders.defaultProps = {
  existingUploads: [],
};

function mapStateToProps(state) {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const currentOrders = selectCurrentOrders(state) || {};
  const uploads = selectUploadsForCurrentOrders(state);

  return {
    serviceMember,
    currentOrders,
    existingUploads: uploads,
    moveIsApproved: selectMoveIsApproved(state),
    isPpm: selectHasCurrentPPM(state),
  };
}

const mapDispatchToProps = {
  updateOrders: updateOrdersAction,
  setFlashMessage: setFlashMessageAction,
};

export default withContext(connect(mapStateToProps, mapDispatchToProps)(EditOrders));
