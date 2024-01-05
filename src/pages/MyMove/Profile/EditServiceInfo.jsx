import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { GridContainer, Alert } from '@trussworks/react-uswds';
import { connect } from 'react-redux';
import { useNavigate } from 'react-router-dom';

import ServiceInfoForm from 'components/Customer/ServiceInfoForm/ServiceInfoForm';
import { patchServiceMember, patchOrders, getResponseError } from 'services/internalApi';
import {
  updateServiceMember as updateServiceMemberAction,
  updateOrders as updateOrdersAction,
} from 'store/entities/actions';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import {
  selectServiceMemberFromLoggedInUser,
  selectCurrentOrders,
  selectMoveIsInDraft,
} from 'store/entities/selectors';
import { generalRoutes, customerRoutes } from 'constants/routes';
import { OrdersShape, ServiceMemberShape } from 'types/customerShapes';
import { formatWeight } from 'utils/formatters';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { formatDateForSwagger } from 'shared/dates';

export const EditServiceInfo = ({
  serviceMember,
  currentOrders,
  updateServiceMember,
  updateOrders,
  setFlashMessage,
  moveIsInDraft,
}) => {
  const navigate = useNavigate();
  const [serverError, setServerError] = useState(null);

  useEffect(() => {
    if (!moveIsInDraft) {
      // Redirect to the home page
      navigate(generalRoutes.HOME_PATH);
    }
  }, [moveIsInDraft, navigate]);

  const initialValues = {
    first_name: serviceMember?.first_name || '',
    middle_name: serviceMember?.middle_name || '',
    last_name: serviceMember?.last_name || '',
    suffix: serviceMember?.suffix || '',
    affiliation: serviceMember?.affiliation || '',
    edipi: serviceMember?.edipi || '',
    grade: currentOrders?.grade || '',
    current_location: currentOrders?.origin_duty_location || {},
    orders_type: currentOrders?.orders_type || '',
    departmentIndicator: currentOrders?.department_indicator,
  };

  const handleSubmit = (values) => {
    const entitlementCouldChange = values.grade !== currentOrders.grade;
    const payload = {
      id: serviceMember.id,
      first_name: values.first_name,
      middle_name: values.middle_name,
      last_name: values.last_name,
      suffix: values.suffix,
      affiliation: values.affiliation,
      edipi: values.edipi,
      rank: values.grade,
      current_location_id: values.current_location.id,
    };

    patchServiceMember(payload)
      .then((response) => {
        updateServiceMember(response);
        if (entitlementCouldChange) {
          const weightAllowance = currentOrders?.has_dependents
            ? response.weight_allotment.total_weight_self_plus_dependents
            : response.weight_allotment.total_weight_self;
          setFlashMessage(
            'EDIT_SERVICE_INFO_SUCCESS',
            'info',
            `Your weight entitlement is now ${formatWeight(weightAllowance)}.`,
            'Your changes have been saved. Note that the entitlement has also changed.',
          );
        } else {
          setFlashMessage('EDIT_SERVICE_INFO_SUCCESS', 'success', '', 'Your changes have been saved.');
        }

        navigate(customerRoutes.PROFILE_PATH);
      })
      .catch((e) => {
        // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
        const { response } = e;
        const errorMessage = getResponseError(response, 'failed to update service member due to server error');
        setServerError(errorMessage);
      });

    const ordersPayload = {
      grade: values.grade,
      origin_duty_location_id: values.current_location.id,
      service_member_id: serviceMember.id,
      id: currentOrders.id,
      new_duty_location_id: currentOrders.new_duty_location.id,
      has_dependents: currentOrders.has_dependents,
      issue_date: formatDateForSwagger(currentOrders.issue_date),
      report_by_date: formatDateForSwagger(currentOrders.report_by_date),
      spouse_has_pro_gear: currentOrders.spouse_has_pro_gear,
      orders_type: currentOrders.orders_type,
    };
    patchOrders(ordersPayload)
      .then((response) => {
        updateOrders(response);
        if (entitlementCouldChange) {
          const weightAllowance = currentOrders?.has_dependents
            ? serviceMember.weight_allotment.total_weight_self_plus_dependents
            : serviceMember.weight_allotment.total_weight_self;
          setFlashMessage(
            'EDIT_SERVICE_INFO_SUCCESS',
            'info',
            `Your weight entitlement is now ${formatWeight(weightAllowance)}.`,
            'Your changes have been saved. Note that the entitlement has also changed.',
          );
        } else {
          setFlashMessage('EDIT_SERVICE_INFO_SUCCESS', 'success', '', 'Your changes have been saved.');
        }

        navigate(customerRoutes.PROFILE_PATH);
      })
      .catch((e) => {
        // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
        const { response } = e;
        const errorMessage = getResponseError(response, 'failed to update orders due to server error');
        setServerError(errorMessage);
      });
  };

  const handleCancel = () => {
    navigate(customerRoutes.PROFILE_PATH);
  };

  return (
    <GridContainer>
      <NotificationScrollToTop dependency={serverError} />
      {serverError && (
        <Alert type="error" headingLevel="h4" heading="An error occurred">
          {serverError}
        </Alert>
      )}
      <ServiceInfoForm
        initialValues={initialValues}
        newDutyLocation={currentOrders?.new_duty_location}
        onSubmit={handleSubmit}
        onCancel={handleCancel}
      />
    </GridContainer>
  );
};

EditServiceInfo.propTypes = {
  updateServiceMember: PropTypes.func.isRequired,
  updateOrders: PropTypes.func.isRequired,
  setFlashMessage: PropTypes.func.isRequired,
  serviceMember: ServiceMemberShape.isRequired,
  currentOrders: OrdersShape.isRequired,
  moveIsInDraft: PropTypes.bool,
};

EditServiceInfo.defaultProps = {
  moveIsInDraft: false,
};

const mapDispatchToProps = {
  updateServiceMember: updateServiceMemberAction,
  updateOrders: updateOrdersAction,
  setFlashMessage: setFlashMessageAction,
};

const mapStateToProps = (state) => ({
  serviceMember: selectServiceMemberFromLoggedInUser(state),
  currentOrders: selectCurrentOrders(state) || {},
  moveIsInDraft: selectMoveIsInDraft(state),
});

export default connect(mapStateToProps, mapDispatchToProps)(EditServiceInfo);
