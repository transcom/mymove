import React, { useState } from 'react';
import { connect } from 'react-redux';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';
import { useNavigate } from 'react-router';

import NotificationScrollToTop from 'components/NotificationScrollToTop';
import OrdersInfoForm from 'components/Customer/OrdersInfoForm/OrdersInfoForm';
import { getServiceMember, createOrders, getResponseError } from 'services/internalApi';
import {
  updateOrders as updateOrdersAction,
  updateServiceMember as updateServiceMemberAction,
} from 'store/entities/actions';
import { withContext } from 'shared/AppContext';
import { formatDateForSwagger } from 'shared/dates';
import { formatYesNoAPIValue, dropdownInputOptions } from 'utils/formatters';
import { ORDERS_TYPE_OPTIONS } from 'constants/orders';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import { generalRoutes } from 'constants/routes';
import withRouter from 'utils/routing';

const AddOrders = ({ context, serviceMemberId, updateServiceMember, updateOrders }) => {
  const [serverError, setServerError] = useState('');
  const navigate = useNavigate();

  const handleBack = () => {
    navigate(generalRoutes.HOME_PATH);
  };

  const submitOrders = async (values) => {
    const pendingValues = {
      ...values,
      service_member_id: serviceMemberId,
      new_duty_location_id: values.new_duty_location.id,
      has_dependents: formatYesNoAPIValue(values.has_dependents),
      report_by_date: formatDateForSwagger(values.report_by_date),
      issue_date: formatDateForSwagger(values.issue_date),
      grade: values.grade,
      origin_duty_location_id: values.origin_duty_location.id,
      spouse_has_pro_gear: false,
    };

    try {
      const createdOrders = await createOrders(pendingValues);
      const newOrderId = createdOrders.id;
      updateOrders(createdOrders);
      const updatedServiceMember = await getServiceMember(serviceMemberId);
      updateServiceMember(updatedServiceMember);
      navigate(`/orders/upload/${newOrderId}`);
    } catch (error) {
      const { response } = error;
      const errorMessage = getResponseError(response, 'failed to update/create orders due to server error');
      setServerError(errorMessage);
    }
  };

  const initialValues = {
    orders_type: '',
    issue_date: '',
    report_by_date: '',
    has_dependents: '',
    new_duty_location: '',
    grade: '',
    origin_duty_location: '',
  };

  // Only allow PCS unless feature flag is on
  const showAllOrdersTypes = context.flags?.allOrdersTypes;
  const allowedOrdersTypes = showAllOrdersTypes
    ? ORDERS_TYPE_OPTIONS
    : { PERMANENT_CHANGE_OF_STATION: ORDERS_TYPE_OPTIONS.PERMANENT_CHANGE_OF_STATION };

  const ordersTypeOptions = dropdownInputOptions(allowedOrdersTypes);

  return (
    <GridContainer data-testid="main-container">
      <NotificationScrollToTop dependency={serverError} />

      {serverError && (
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Alert type="error" headingLevel="h4" heading="An error occurred">
              {serverError}
            </Alert>
          </Grid>
        </Grid>
      )}

      <Grid row data-testid="orders-form-container">
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <OrdersInfoForm
            ordersTypeOptions={ordersTypeOptions}
            initialValues={initialValues}
            onSubmit={submitOrders}
            onBack={handleBack}
          />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

const mapStateToProps = (state) => {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);

  return {
    serviceMemberId: serviceMember?.id,
  };
};

const mapDispatchToProps = {
  updateOrders: updateOrdersAction,
  updateServiceMember: updateServiceMemberAction,
};

export default withContext(withRouter(connect(mapStateToProps, mapDispatchToProps)(AddOrders)));
