import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { GridContainer, Alert } from '@trussworks/react-uswds';
import { connect } from 'react-redux';
import { useLocation, useNavigate } from 'react-router-dom';

import ServiceInfoForm from 'components/Customer/ServiceInfoForm/ServiceInfoForm';
import { patchServiceMember, getResponseError } from 'services/internalApi';
import { updateServiceMember as updateServiceMemberAction } from 'store/entities/actions';
import {
  selectServiceMemberFromLoggedInUser,
  selectCurrentOrders,
  selectMoveIsInDraft,
} from 'store/entities/selectors';
import { generalRoutes, customerRoutes } from 'constants/routes';
import { OrdersShape, ServiceMemberShape } from 'types/customerShapes';
import NotificationScrollToTop from 'components/NotificationScrollToTop';

export const EditServiceInfo = ({ serviceMember, currentOrders, updateServiceMember, moveIsInDraft }) => {
  const navigate = useNavigate();
  const [serverError, setServerError] = useState(null);
  const { state } = useLocation();

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
    preferred_name: serviceMember?.preferred_name || '',
    affiliation: serviceMember?.affiliation || '',
    edipi: serviceMember?.edipi || '',
    orders_type: currentOrders?.orders_type || '',
    departmentIndicator: currentOrders?.department_indicator,
    emplid: serviceMember?.emplid || '',
  };

  const handleSubmit = (values) => {
    const payload = {
      id: serviceMember.id,
      first_name: values.first_name,
      middle_name: values.middle_name,
      last_name: values.last_name,
      suffix: values.suffix,
      preferred_name: values.preferred_name,
      affiliation: values.affiliation,
      edipi: values.edipi,
      emplid: values.affiliation === 'COAST_GUARD' ? values.emplid : null,
    };

    patchServiceMember(payload)
      .then((response) => {
        updateServiceMember(response);
        navigate(customerRoutes.PROFILE_PATH, { state });
      })
      .catch((e) => {
        // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
        const { response } = e;
        let errorMessage;
        if (e.response.body.message === 'Unhandled data error encountered') {
          errorMessage = 'This EMPLID is already in use';
        } else {
          errorMessage = getResponseError(response, 'failed to update service member due to server error');
        }
        setServerError(errorMessage);
      });
  };

  const handleCancel = () => {
    navigate(customerRoutes.PROFILE_PATH, { state });
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
  serviceMember: ServiceMemberShape.isRequired,
  currentOrders: OrdersShape.isRequired,
  moveIsInDraft: PropTypes.bool,
};

EditServiceInfo.defaultProps = {
  moveIsInDraft: false,
};

const mapDispatchToProps = {
  updateServiceMember: updateServiceMemberAction,
};

const mapStateToProps = (state) => ({
  serviceMember: selectServiceMemberFromLoggedInUser(state),
  currentOrders: selectCurrentOrders(state) || {},
  moveIsInDraft: selectMoveIsInDraft(state),
});

export default connect(mapStateToProps, mapDispatchToProps)(EditServiceInfo);
