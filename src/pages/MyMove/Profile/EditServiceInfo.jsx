import React, { useState } from 'react';
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
  selectEntitlementsForLoggedInUser,
} from 'store/entities/selectors';
import { generalRoutes } from 'constants/routes';
import { OrdersShape, ServiceMemberShape } from 'types/customerShapes';
import { EntitlementShape } from 'types';

export const EditServiceInfo = ({
  serviceMember,
  currentOrders,
  entitlement,
  updateServiceMember,
  setFlashMessage,
}) => {
  const history = useHistory();
  const [serverError, setServerError] = useState(null);

  const initialValues = {
    first_name: serviceMember?.first_name || '',
    middle_name: serviceMember?.middle_name || '',
    last_name: serviceMember?.last_name || '',
    suffix: serviceMember?.suffix || '',
    affiliation: serviceMember?.affiliation || '',
    edipi: serviceMember?.edipi || '',
    rank: currentOrders?.grade || '',
    current_station: currentOrders?.origin_duty_station || {},
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
      current_station_id: values.current_station.id,
    };

    return patchServiceMember(payload)
      .then((response) => {
        updateServiceMember(response);

        if (entitlementCouldChange) {
          setFlashMessage(
            'EDIT_SERVICE_INFO_SUCCESS',
            'info',
            `Your weight entitlement is now ${entitlement.sum.toLocaleString()} lbs.`,
            'Your changes have been saved. Note that the entitlement has also changed.',
          );
        } else {
          setFlashMessage('EDIT_SERVICE_INFO_SUCCESS', 'success', '', 'Your changes have been saved.');
        }

        // TODO - change this to profile path?
        history.push(generalRoutes.HOME_PATH);
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
    // TODO - change this to profile path?
    history.push(generalRoutes.HOME_PATH);
  };

  return (
    <GridContainer>
      <ScrollToTop />
      {serverError && (
        <Alert type="error" heading="An error occurred">
          {serverError}
        </Alert>
      )}
      <ServiceInfoForm
        initialValues={initialValues}
        newDutyStation={currentOrders?.newDutyStation}
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
  entitlement: EntitlementShape.isRequired,
};

const mapDispatchToProps = {
  updateServiceMember: updateServiceMemberAction,
  setFlashMessage: setFlashMessageAction,
};

const mapStateToProps = (state) => ({
  serviceMember: selectServiceMemberFromLoggedInUser(state),
  currentOrders: selectCurrentOrders(state),
  entitlement: selectEntitlementsForLoggedInUser(state),
});

export default connect(mapStateToProps, mapDispatchToProps)(EditServiceInfo);
