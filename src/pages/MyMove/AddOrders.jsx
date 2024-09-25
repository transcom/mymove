import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';
import { generatePath, useNavigate } from 'react-router-dom';

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
import { selectCanAddOrders, selectMoveId, selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import { customerRoutes, generalRoutes } from 'constants/routes';
import withRouter from 'utils/routing';
import { setCanAddOrders as setCanAddOrdersAction, setMoveId as setMoveIdAction } from 'store/general/actions';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';

const AddOrders = ({
  context,
  serviceMemberId,
  updateServiceMember,
  updateOrders,
  canAddOrders,
  setCanAddOrders,
  moveId,
  setMoveId,
}) => {
  const [serverError, setServerError] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [hasSubmitted, setHasSubmitted] = useState(false);
  const navigate = useNavigate();

  // if the user did NOT come from the create a move button, we want to redirect them to their current move
  // this is an effort to override clicking the browser back button from the upload route
  useEffect(() => {
    const redirectUser = async () => {
      if (!canAddOrders && !hasSubmitted) {
        const path = moveId ? generatePath(customerRoutes.MOVE_HOME_PATH, { moveId }) : generalRoutes.HOME_PATH;
        navigate(path);
      }
      setIsLoading(false);
    };
    redirectUser();
  }, [canAddOrders, navigate, hasSubmitted, moveId]);

  const handleBack = () => {
    navigate(generalRoutes.HOME_PATH);
  };

  const submitOrders = async (values) => {
    setHasSubmitted(true);
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
      setMoveId(createdOrders?.moves[0].id);
      setCanAddOrders(false);
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

  if (isLoading) {
    return <LoadingPlaceholder />;
  }

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
  const canAddOrders = selectCanAddOrders(state);
  const moveId = selectMoveId(state);

  return {
    serviceMemberId: serviceMember?.id,
    canAddOrders,
    moveId,
  };
};

const mapDispatchToProps = {
  updateOrders: updateOrdersAction,
  updateServiceMember: updateServiceMemberAction,
  setCanAddOrders: setCanAddOrdersAction,
  setMoveId: setMoveIdAction,
};

export default withContext(withRouter(connect(mapStateToProps, mapDispatchToProps)(AddOrders)));
