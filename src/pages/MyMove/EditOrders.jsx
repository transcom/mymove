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
  selectEntitlementsForLoggedInUser,
} from 'store/entities/selectors';
import EditOrdersForm from 'components/Customer/EditOrdersForm/EditOrdersForm';
import { OrdersShape } from 'types/customerShapes';
import { formatYesNoInputValue, formatYesNoAPIValue } from 'utils/formatters';
import { ORDERS_TYPE_OPTIONS } from 'constants/orders';
import { dropdownInputOptions } from 'shared/formatters';
import { EntitlementShape, ExistingUploadsShape } from 'types';
import { DutyStationShape } from 'types/dutyStation';

export const EditOrders = ({
  currentOrders,
  serviceMemberId,
  updateOrders,
  existingUploads,
  moveIsApproved,
  spouseHasProGear,
  currentStation,
  setFlashMessage,
  entitlement,
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
    new_duty_station: currentOrders?.new_duty_station || null,
    uploaded_orders: existingUploads || [],
  };

  // Only allow PCS unless feature flag is on
  const showAllOrdersTypes = context.flags?.allOrdersTypes;
  const allowedOrdersTypes = showAllOrdersTypes
    ? ORDERS_TYPE_OPTIONS
    : { PERMANENT_CHANGE_OF_STATION: ORDERS_TYPE_OPTIONS.PERMANENT_CHANGE_OF_STATION };
  const ordersTypeOptions = dropdownInputOptions(allowedOrdersTypes);

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
    let entitlementCouldChange = false;

    const fromFormSpouseHasProGear = (fieldValues.has_dependents && fieldValues.spouse_has_pro_gear) || false;

    if (fieldValues.has_dependents !== currentOrders.has_dependents || fromFormSpouseHasProGear !== spouseHasProGear) {
      entitlementCouldChange = true;
    }

    const newDutyStationId = fieldValues.new_duty_station.id;
    return patchOrders({
      ...fieldValues,
      has_dependents: formatYesNoAPIValue(fieldValues.has_dependents),
      new_duty_station_id: newDutyStationId,
      spouse_has_pro_gear: fromFormSpouseHasProGear,
    })
      .then((response) => {
        updateOrders(response);

        if (entitlementCouldChange) {
          setFlashMessage(
            'EDIT_ORDERS_SUCCESS',
            'info',
            `Your weight entitlement is now ${entitlement.sum.toLocaleString()} lbs.`,
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
                currentStation={currentStation}
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
  serviceMemberId: PropTypes.string.isRequired,
  setFlashMessage: PropTypes.func.isRequired,
  updateOrders: PropTypes.func.isRequired,
  currentOrders: OrdersShape.isRequired,
  entitlement: EntitlementShape.isRequired,
  existingUploads: ExistingUploadsShape,
  spouseHasProGear: PropTypes.bool,
  context: PropTypes.shape({
    flags: PropTypes.shape({
      allOrdersTypes: PropTypes.bool,
    }).isRequired,
  }).isRequired,
  currentStation: DutyStationShape,
};

EditOrders.defaultProps = {
  existingUploads: [],
  spouseHasProGear: false,
  currentStation: {},
};

function mapStateToProps(state) {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const serviceMemberId = serviceMember?.id;
  const currentOrders = selectCurrentOrders(state) || {};
  const uploads = selectUploadsForCurrentOrders(state);

  return {
    currentOrders,
    serviceMemberId,
    existingUploads: uploads,
    moveIsApproved: selectMoveIsApproved(state),
    isPpm: selectHasCurrentPPM(state),
    entitlement: selectEntitlementsForLoggedInUser(state),
    currentStation: serviceMember?.current_station || {},
  };
}

const mapDispatchToProps = {
  updateOrders: updateOrdersAction,
  setFlashMessage: setFlashMessageAction,
};

export default withContext(connect(mapStateToProps, mapDispatchToProps)(EditOrders));
