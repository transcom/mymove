import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { GridContainer, Alert } from '@trussworks/react-uswds';
import { connect } from 'react-redux';
import { useHistory } from 'react-router-dom';

import ScrollToTop from 'components/ScrollToTop';
import ServiceInfoForm from 'components/Customer/ServiceInfoForm/ServiceInfoForm';
import { patchServiceMember, getResponseError } from 'services/internalApi';
import { updateServiceMember as updateServiceMemberAction } from 'store/entities/actions';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import {
  selectServiceMemberFromLoggedInUser,
  selectCurrentOrders,
  selectMoveIsInDraft,
} from 'store/entities/selectors';
import { generalRoutes, customerRoutes } from 'constants/routes';
import { OrdersShape, ServiceMemberShape } from 'types/customerShapes';
import { formatWeight } from 'utils/formatters';

export const EditServiceInfo = ({
  serviceMember,
  currentOrders,
  updateServiceMember,
  setFlashMessage,
  moveIsInDraft,
}) => {
  const history = useHistory();
  const [serverError, setServerError] = useState(null);

  useEffect(() => {
    if (!moveIsInDraft) {
      // Redirect to the home page
      history.push(generalRoutes.HOME_PATH);
    }
  }, [moveIsInDraft, history]);

  const initialValues = {
    first_name: serviceMember?.first_name || '',
    middle_name: serviceMember?.middle_name || '',
    last_name: serviceMember?.last_name || '',
    suffix: serviceMember?.suffix || '',
    affiliation: serviceMember?.affiliation || '',
    edipi: serviceMember?.edipi || '',
    rank: currentOrders?.grade || '',
    current_location: currentOrders?.origin_duty_location || {},
  };

  const handleSubmit = (values) => {
    const entitlementCouldChange = values.rank !== currentOrders.grade;

    const payload = {
      id: serviceMember.id,
      first_name: values.first_name,
      middle_name: values.middle_name,
      last_name: values.last_name,
      suffix: values.suffix,
      affiliation: values.affiliation,
      edipi: values.edipi,
      rank: values.rank,
      current_location_id: values.current_location.id,
    };

    return patchServiceMember(payload)
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

        history.push(customerRoutes.PROFILE_PATH);
      })
      .catch((e) => {
        // TODO - error handling - below is rudimentary error handling to approximate existing UX
        // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
        const { response } = e;
        const errorMessage = getResponseError(response, 'failed to update service member due to server error');
        setServerError(errorMessage);
      });
  };

  const handleCancel = () => {
    history.push(customerRoutes.PROFILE_PATH);
  };

  return (
    <GridContainer>
      <ScrollToTop />
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
  setFlashMessage: setFlashMessageAction,
};

const mapStateToProps = (state) => ({
  serviceMember: selectServiceMemberFromLoggedInUser(state),
  currentOrders: selectCurrentOrders(state) || {},
  moveIsInDraft: selectMoveIsInDraft(state),
});

export default connect(mapStateToProps, mapDispatchToProps)(EditServiceInfo);
